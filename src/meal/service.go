package meal

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Meal struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	UserID    string             `json:"userid" validate:"required" bson:"userid"`
	Carb      float64            `json:"carb" default:"0"`
	Protein   float64            `json:"protein" default:"0"`
	Fat       float64            `json:"fat" default:"0"`
	Calories  float64            `json:"calories" default:"0"`
	CreatedAt time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty" default:"null"`
}

type Nutrient struct {
	Carb          float64         `json:"carb"`
	Protein       float64         `json:"protein"`
	Fat           float64         `json:"fat"`
	MicroNutrient []MicroNutrient `json:"microNutrient"`
}

type MicroNutrient struct {
	Name  string  `json:"name"`
	Value float64 `json:"value"`
	Unit  string  `json:"unit"`
}

type MealService struct {
	DB *mongo.Database
}

type IMealService interface {
	CreateMeal(meal *CreateMealDto) (*Meal, error)
	GetAllMeals() ([]*Meal, error)
	GetMeal(id string) (*Meal, error)
	GetMealByUser(userid string) ([]*Meal, error)
	GetMealByUserDate(userid string, start int, end int) ([]*Meal, error)
	DeleteMeal(id string) error
	UpdateMeal(doc *UpdateMealDto, id string) (*Meal, error)
}

func (ns *MealService) CreateMeal(meal *CreateMealDto) (*Meal, error) {
	meal.CreatedAt = time.Now()
	localLocation, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		return nil, err
	}
	res, err := ns.GetMealByUser(meal.UserID)
	if err != nil {
		return nil, err
	}
	if len(res) > 0 {
		for _, v := range res {
			tmp := v.CreatedAt.In(localLocation)
			if tmp.Format("2006-01-02") == meal.CreatedAt.In(localLocation).Format("2006-01-02") {
				return nil, errors.New("Meal already exist, please update instead")
			}
		}
	}

	result, err := ns.DB.Collection("meal").InsertOne(context.Background(), meal)
	if err != nil {
		return nil, err
	}
	filter := bson.D{{Key: "_id", Value: result.InsertedID}}
	createdRecord := ns.DB.Collection("meal").FindOne(context.Background(), filter)
	createdMeal := &Meal{}
	if err := createdRecord.Decode(createdMeal); err != nil {
		return nil, err
	}
	return createdMeal, nil
}

func (ns *MealService) GetAllMeals() ([]*Meal, error) {
	cursor, err := ns.DB.Collection("meal").Find(context.Background(), bson.D{})
	if err != nil {
		return nil, err
	}
	var Meals []*Meal
	if err := cursor.All(context.Background(), &Meals); err != nil {
		return nil, err
	}
	return Meals, nil
}

func (ns *MealService) GetMeal(id string) (*Meal, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	filter := bson.D{{Key: "_id", Value: oid}}
	meal := &Meal{}
	if err := ns.DB.Collection("meal").FindOne(context.Background(), filter).Decode(meal); err != nil {
		return nil, err
	}
	return meal, nil
}

func (ns *MealService) GetMealByUser(userid string) ([]*Meal, error) {
	filter := bson.M{"userid": bson.M{"$in": []string{userid}}}
	cursor, err := ns.DB.Collection("meal").Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	var Meals []*Meal
	if err := cursor.All(context.Background(), &Meals); err != nil {
		return nil, err
	}
	return Meals, nil
}

func (ns *MealService) GetMealByUserDate(userid string, start int, end int) ([]*Meal, error) {
	start_time := time.Unix(int64(start), 0).Format("2006-01-02T15:04:05Z")
	end_time := time.Unix(int64(end), 0).Format("2006-01-02T15:04:05Z")
	filter := bson.M{"userid": bson.M{"$in": []string{userid}}, "created_at": bson.M{"$gte": start_time, "$lte": end_time}}
	cursor, err := ns.DB.Collection("meal").Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	var Meals []*Meal
	if err := cursor.All(context.Background(), &Meals); err != nil {
		return nil, err
	}
	return Meals, nil
}

func (ns *MealService) DeleteMeal(id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	filter := bson.D{{Key: "_id", Value: oid}}
	if _, err := ns.DB.Collection("meal").DeleteOne(context.Background(), filter); err != nil {
		return err
	}
	return nil
}

func (ns *MealService) UpdateMeal(doc *UpdateMealDto, id string) (*Meal, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	filter := bson.D{{Key: "_id", Value: oid}}
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "carb", Value: doc.Carb},
			{Key: "protein", Value: doc.Protein},
			{Key: "fat", Value: doc.Fat},
			{Key: "calories", Value: doc.Calories},
		}},
	}
	if _, err := ns.DB.Collection("meal").UpdateOne(context.Background(), filter, update); err != nil {
		return nil, err
	}
	return ns.GetMeal(id)
}
