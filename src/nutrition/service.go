package nutrition

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Nutrition struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	UserID    string             `json:"userid" validate:"required" bson:"userid"`
	Carb      float64            `json:"carb" default:"0"`
	Protein   float64            `json:"protein" default:"0"`
	Fat       float64            `json:"fat" default:"0"`
	Calories  float64            `json:"calories" default:"0"`
	CreatedAt time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
}

type NutritionService struct {
	DB *mongo.Database
}

type INutritionService interface {
	CreateNutrition(nutrition *CreateNutritionDto) (*Nutrition, error)
	GetAllNutritions() ([]*Nutrition, error)
	GetNutrition(id string) (*Nutrition, error)
	GetNutritionByUser(userid string) ([]*Nutrition, error)
	DeleteNutrition(id string) error
	UpdateNutrition(doc *UpdateNutritionDto, id string) (*Nutrition, error)
}

func (ns *NutritionService) CreateNutrition(nutrition *CreateNutritionDto) (*Nutrition, error) {
	nutrition.CreatedAt = time.Now()
	localLocation, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		return nil, err
	}
	nutrition.CreatedAt = nutrition.CreatedAt.In(localLocation)
	res, err := ns.GetNutritionByUser(nutrition.UserID)
	if err != nil {
		return nil, err
	}
	if len(res) > 0 {
		for _, v := range res {
			tmp := v.CreatedAt.In(localLocation)
			if tmp.Format("2006-01-02") == nutrition.CreatedAt.Format("2006-01-02") {
				return nil, errors.New("Nutrition already exist, please update it")
			}
		}
	}

	result, err := ns.DB.Collection("nutrition").InsertOne(context.Background(), nutrition)
	if err != nil {
		return nil, err
	}
	filter := bson.D{{Key: "_id", Value: result.InsertedID}}
	createdRecord := ns.DB.Collection("nutrition").FindOne(context.Background(), filter)
	createdNutrition := &Nutrition{}
	if err := createdRecord.Decode(createdNutrition); err != nil {
		return nil, err
	}
	return createdNutrition, nil
}

func (ns *NutritionService) GetAllNutritions() ([]*Nutrition, error) {
	cursor, err := ns.DB.Collection("nutrition").Find(context.Background(), bson.D{})
	if err != nil {
		return nil, err
	}
	var nutritions []*Nutrition
	if err := cursor.All(context.Background(), &nutritions); err != nil {
		return nil, err
	}
	return nutritions, nil
}

func (ns *NutritionService) GetNutrition(id string) (*Nutrition, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	filter := bson.D{{Key: "_id", Value: oid}}
	nutrition := &Nutrition{}
	if err := ns.DB.Collection("nutrition").FindOne(context.Background(), filter).Decode(nutrition); err != nil {
		return nil, err
	}
	return nutrition, nil
}

func (ns *NutritionService) GetNutritionByUser(userid string) ([]*Nutrition, error) {
	filter := bson.M{"userid": bson.M{"$in": []string{userid}}}
	cursor, err := ns.DB.Collection("nutrition").Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	var nutritions []*Nutrition
	if err := cursor.All(context.Background(), &nutritions); err != nil {
		return nil, err
	}
	return nutritions, nil
}

func (ns *NutritionService) DeleteNutrition(id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	filter := bson.D{{Key: "_id", Value: oid}}
	if _, err := ns.DB.Collection("nutrition").DeleteOne(context.Background(), filter); err != nil {
		return err
	}
	return nil
}

func (ns *NutritionService) UpdateNutrition(doc *UpdateNutritionDto, id string) (*Nutrition, error) {
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
	if _, err := ns.DB.Collection("nutrition").UpdateOne(context.Background(), filter, update); err != nil {
		return nil, err
	}
	return ns.GetNutrition(id)
}
