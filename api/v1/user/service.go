package user

import (
	"context"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"time"

	authEnums "github.com/Npwskp/GymsbroBackend/api/v1/auth/enums"
	"github.com/Npwskp/GymsbroBackend/api/v1/function"
	minio "github.com/Npwskp/GymsbroBackend/api/v1/storage"
	userFitnessPreferenceEnums "github.com/Npwskp/GymsbroBackend/api/v1/user/enums"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserService struct {
	DB           *mongo.Database
	MinioService minio.MinioService
}

const (
	UserPictureBucketName = "user-profile-image"
)

type IUserService interface {
	CreateUser(user *CreateUserDto) (*User, error)
	GetAllUsers() ([]*User, error)
	GetUser(id string) (*User, error)
	GetUserByEmail(email string) (*User, error)
	GetUserByOAuthID(oauthid string) (*User, error)
	GetUserEnergyConsumePlan(id string) (*userFitnessPreferenceEnums.EnergyConsumptionPlan, error)
	DeleteUser(id string) error
	UpdateUsernamePassword(doc *UpdateUsernamePasswordDto, id string) (*User, error)
	UpdateBody(doc *UpdateBodyDto, id string) (*User, error)
	UpdateFirstLoginStatus(id string) error
	UpdateUserPicture(c *fiber.Ctx, id string, file io.Reader, filename string, contentType string) (*User, error)
}

func (us *UserService) CreateUser(user *CreateUserDto) (*User, error) {
	model := CreateUserModel(user)

	// Calculate BMR if possible
	calculateAndUpdateBMRAndBMI(model)

	find := bson.D{{Key: "email", Value: user.Email}}
	check, err := us.DB.Collection("users").CountDocuments(context.Background(), find)
	if err != nil {
		return nil, err
	}
	if check > 0 {
		return nil, errors.New("email have been used")
	}
	result, err := us.DB.Collection("users").InsertOne(context.Background(), model)
	if err != nil {
		return nil, err
	}
	filter := bson.D{{Key: "_id", Value: result.InsertedID}}
	createdRecord := us.DB.Collection("users").FindOne(context.Background(), filter)
	createdUser := &User{}
	if err := createdRecord.Decode(createdUser); err != nil {
		return nil, err
	}
	return createdUser, nil
}

func (us *UserService) GetAllUsers() ([]*User, error) {
	cursor, err := us.DB.Collection("users").Find(context.Background(), bson.D{})
	if err != nil {
		return nil, err
	}
	users := make([]*User, 0)
	if err := cursor.All(context.Background(), &users); err != nil {
		return nil, err
	}
	return users, nil
}

func (us *UserService) GetUser(id string) (*User, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	filter := bson.D{{Key: "_id", Value: oid}}
	user := &User{}
	if err := us.DB.Collection("users").FindOne(context.Background(), filter).Decode(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (us *UserService) GetUserByEmail(email string) (*User, error) {
	filter := bson.D{{Key: "email", Value: email}}
	user := &User{}
	if err := us.DB.Collection("users").FindOne(context.Background(), filter).Decode(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (us *UserService) GetUserByOAuthID(oauthid string) (*User, error) {
	filter := bson.D{{Key: "oauth_id", Value: oauthid}}
	user := &User{}
	if err := us.DB.Collection("users").FindOne(context.Background(), filter).Decode(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (us *UserService) GetUserEnergyConsumePlan(id string) (*userFitnessPreferenceEnums.EnergyConsumptionPlan, error) {
	user, err := us.GetUser(id)
	if err != nil {
		return nil, err
	}

	if err := validateUserForEnergyPlan(user); err != nil {
		return nil, err
	}

	return userFitnessPreferenceEnums.GetUserEnergyConsumePlan(user.Weight, user.Height, user.Age, user.Gender, user.NutritionInfo.ActivityLevel, user.NutritionInfo.Goal)
}

func (us *UserService) DeleteUser(id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	filter := bson.D{{Key: "_id", Value: oid}}
	if _, err := us.DB.Collection("users").DeleteOne(context.Background(), filter); err != nil {
		return err
	}
	return nil
}

func (us *UserService) UpdateUsernamePassword(doc *UpdateUsernamePasswordDto, id string) (*User, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	filter := bson.D{{Key: "_id", Value: oid}}
	user, err := us.GetUser(id)
	if err != nil {
		return nil, err
	}
	if strings.Compare(user.Password, doc.Password) != 0 {
		return nil, errors.New("password is not correct")
	}
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "username", Value: function.Coalesce(doc.Username, user.Username)},
			{Key: "password", Value: function.Coalesce(doc.NewPassword, user.Password)},
			{Key: "updated_at", Value: time.Now()},
		}},
	}

	// Perform the update
	result, err := us.DB.Collection("users").UpdateOne(context.Background(), filter, update)
	if err != nil {
		return nil, err
	}

	// Check if any document was modified
	if result.ModifiedCount == 0 {
		return nil, errors.New("no user found for the given ID")
	}

	// Retrieve the updated document
	filter = bson.D{{Key: "_id", Value: oid}}
	UpdatedUser := &User{}
	updatedRecord := us.DB.Collection("users").FindOne(context.Background(), filter)
	if err := updatedRecord.Decode(&UpdatedUser); err != nil {
		return nil, err
	}

	return UpdatedUser, nil
}

func (us *UserService) UpdateBody(doc *UpdateBodyDto, id string) (*User, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.D{{Key: "_id", Value: oid}}
	user, err := us.GetUser(id)
	if err != nil {
		return nil, err
	}

	new_weight := function.Coalesce(doc.Weight, user.Weight).(float64)
	new_height := function.Coalesce(doc.Height, user.Height).(float64)
	new_age := function.Coalesce(doc.Age, user.Age).(int)
	new_gender := function.Coalesce(doc.Gender, user.Gender).(authEnums.GenderType)

	new_BMR := userFitnessPreferenceEnums.CalculateBMR(new_weight, new_height, new_age, new_gender)
	new_BMI := userFitnessPreferenceEnums.CalculateBMI(new_weight, new_height)

	doc.NutritionInfo.BMR = new_BMR
	doc.BodyComposition.BMI = new_BMI

	new_macronutrients := &userFitnessPreferenceEnums.Macronutrients{
		CarbPreference: function.Coalesce(doc.Macronutrients.CarbPreference, user.Macronutrients.CarbPreference).(userFitnessPreferenceEnums.CarbPreferenceType),
		Calories:       function.Coalesce(doc.Macronutrients.Calories, user.Macronutrients.Calories).(float64),
		Protein:        function.Coalesce(doc.Macronutrients.Protein, user.Macronutrients.Protein).(float64),
		Fat:            function.Coalesce(doc.Macronutrients.Fat, user.Macronutrients.Fat).(float64),
		Carbs:          function.Coalesce(doc.Macronutrients.Carbs, user.Macronutrients.Carbs).(float64),
	}

	new_nutrition_info := &userFitnessPreferenceEnums.NutritionInfo{
		BMR:           new_BMR,
		ActivityLevel: function.Coalesce(doc.NutritionInfo.ActivityLevel, user.NutritionInfo.ActivityLevel).(userFitnessPreferenceEnums.ActivityLevelType),
		Goal:          function.Coalesce(doc.NutritionInfo.Goal, user.NutritionInfo.Goal).(userFitnessPreferenceEnums.GoalType),
	}

	new_body_composition := &userFitnessPreferenceEnums.BodyCompositionInfo{
		BMI:                function.Coalesce(doc.BodyComposition.BMI, user.BodyComposition.BMI).(float64),
		BodyFatMass:        function.Coalesce(doc.BodyComposition.BodyFatMass, user.BodyComposition.BodyFatMass).(float64),
		BodyFatPercentage:  function.Coalesce(doc.BodyComposition.BodyFatPercentage, user.BodyComposition.BodyFatPercentage).(float64),
		SkeletalMuscleMass: function.Coalesce(doc.BodyComposition.SkeletalMuscleMass, user.BodyComposition.SkeletalMuscleMass).(float64),
		ExtracellularWater: function.Coalesce(doc.BodyComposition.ExtracellularWater, user.BodyComposition.ExtracellularWater).(float64),
		ECWRatio:           function.Coalesce(doc.BodyComposition.ECWRatio, user.BodyComposition.ECWRatio).(float64),
	}

	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "weight", Value: function.Coalesce(doc.Weight, user.Weight)},
			{Key: "height", Value: function.Coalesce(doc.Height, user.Height)},
			{Key: "age", Value: function.Coalesce(doc.Age, user.Age)},
			{Key: "gender", Value: function.Coalesce(doc.Gender, user.Gender)},
			{Key: "nutrition_info", Value: new_nutrition_info},
			{Key: "body_composition", Value: new_body_composition},
			{Key: "macronutrients", Value: new_macronutrients},
			{Key: "updated_at", Value: time.Now()},
		}},
	}

	result, err := us.DB.Collection("users").UpdateOne(context.Background(), filter, update)
	if err != nil {
		return nil, err
	}

	if result.ModifiedCount == 0 {
		return nil, errors.New("no user found for the given ID")
	}

	filter = bson.D{{Key: "_id", Value: oid}}
	UpdatedUser := &User{}
	updatedRecord := us.DB.Collection("users").FindOne(context.Background(), filter)
	if err := updatedRecord.Decode(&UpdatedUser); err != nil {
		return nil, err
	}

	return UpdatedUser, nil
}

func (us *UserService) UpdateFirstLoginStatus(id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.D{{Key: "_id", Value: oid}}
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "is_first_login", Value: false},
			{Key: "updated_at", Value: time.Now()},
		}},
	}

	result, err := us.DB.Collection("users").UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}

	if result.ModifiedCount == 0 {
		return errors.New("no user found for the given ID")
	}

	return nil
}

func (us *UserService) UpdateUserPicture(c *fiber.Ctx, id string, file io.Reader, filename string, contentType string) (*User, error) {
	// Get user first to verify existence and get current picture URL
	user, err := us.GetUser(id)
	if err != nil {
		return nil, err
	}

	oldPictureURL := user.Picture

	ext := strings.ToLower(filepath.Ext(filename))
	// Generate unique filename using user ID, timestamp, and a random component
	timestamp := time.Now().UnixNano()
	objectName := fmt.Sprintf("users/%s/profile_%d%s", id, timestamp, ext)

	// Upload to MinIO
	err = us.MinioService.UploadFile(c.Context(), file, UserPictureBucketName, objectName, contentType)
	if err != nil {
		return nil, fmt.Errorf("failed to upload picture: %v", err)
	}

	// Get the URL of the uploaded file
	url, err := us.MinioService.GetFileURL(c.Context(), UserPictureBucketName, objectName)
	if err != nil {
		return nil, fmt.Errorf("failed to get picture URL: %v", err)
	}

	// Update user's picture URL in database
	oid, _ := primitive.ObjectIDFromHex(id)
	filter := bson.D{{Key: "_id", Value: oid}}
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "picture", Value: url},
			{Key: "updated_at", Value: time.Now()},
		}},
	}

	result, err := us.DB.Collection("users").UpdateOne(context.Background(), filter, update)
	if err != nil {
		return nil, err
	}

	if result.ModifiedCount == 0 {
		return nil, errors.New("no user found for the given ID")
	}

	// Delete old picture after successful upload and update
	if oldPictureURL != "" {
		baseURL := strings.Split(oldPictureURL, "?")[0]
		urlParts := strings.Split(baseURL, us.MinioService.GetFullBucketName(UserPictureBucketName)+"/")
		if len(urlParts) > 1 {
			oldObjectName := urlParts[1]
			if err := us.MinioService.DeleteFile(c.Context(), UserPictureBucketName, oldObjectName); err != nil {
				fmt.Printf("Warning: Failed to delete old profile picture: %v\n", err)
			}
		}
	}

	// Get updated user
	updatedUser := &User{}
	if err := us.DB.Collection("users").FindOne(context.Background(), filter).Decode(updatedUser); err != nil {
		return nil, err
	}

	return updatedUser, nil
}

func validateUserForEnergyPlan(user *User) error {
	var missingFields []string
	if user.Weight == 0 {
		missingFields = append(missingFields, "Weight")
	}
	if user.Height == 0 {
		missingFields = append(missingFields, "Height")
	}
	if user.Age == 0 {
		missingFields = append(missingFields, "Age")
	}
	if user.Gender == "" {
		missingFields = append(missingFields, "Gender")
	}
	if user.NutritionInfo.ActivityLevel == "" {
		missingFields = append(missingFields, "ActivityLevel")
	}
	if user.NutritionInfo.Goal == "" {
		missingFields = append(missingFields, "Goal")
	}
	if len(missingFields) > 0 {
		return errors.New("missing fields for energy consume plan calculation: " + strings.Join(missingFields, ", "))
	}
	return nil
}

func canCalculateBMR(user *User) bool {
	return user.Weight > 0 &&
		user.Height > 0 &&
		user.Age > 0 &&
		(user.Gender == authEnums.GenderMale || user.Gender == authEnums.GenderFemale)
}

func calculateAndUpdateBMRAndBMI(user *User) {
	if canCalculateBMR(user) {
		user.NutritionInfo.BMR = userFitnessPreferenceEnums.CalculateBMR(
			user.Weight,
			user.Height,
			user.Age,
			user.Gender,
		)
		user.BodyComposition.BMI = userFitnessPreferenceEnums.CalculateBMI(
			user.Weight,
			user.Height,
		)
	}
}
