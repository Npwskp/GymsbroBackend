package plans

import (
	"context"
	"errors"

	"github.com/Npwskp/GymsbroBackend/src/function"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Plan struct {
	ID         string   `json:"id,omitempty" bson:"_id,omitempty"`
	UserID     string   `json:"userId" validate:"required"`
	TypeOfPlan string   `json:"typeOfPlan" validate:"required"`
	DayOfWeek  string   `json:"dayOfWeek" validate:"required"`
	Exercise   []string `json:"exercise" default:"[]"`
}

type PlanService struct {
	DB *mongo.Database
}

type IPlanService interface {
	CreatePlan(plan *CreatePlanDto) (*Plan, error)
	GetAllPlans() ([]*Plan, error)
	GetPlan(id string) (*Plan, error)
	GetAllPlanByUser(user_id string, day string) ([]*Plan, error)
	DeletePlan(id string, day string) error
	UpdatePlan(doc *UpdatePlanDto, id string) (*Plan, error)
}

func (ps *PlanService) CreatePlan(plan *CreatePlanDto) (*Plan, error) {
	find := bson.D{{Key: "userId", Value: plan.UserID}, {Key: "dayOfWeek", Value: plan.DayOfWeek}}
	res := ps.DB.Collection("plans").FindOne(context.Background(), find)
	if res != nil {
		return nil, errors.New("plan already exists")
	}
	result, err := ps.DB.Collection("plans").InsertOne(context.Background(), plan)
	if err != nil {
		return nil, err
	}
	filter := bson.D{{Key: "_id", Value: result.InsertedID}}
	createdRecord := ps.DB.Collection("plans").FindOne(context.Background(), filter)
	createdPlan := &Plan{}
	if err := createdRecord.Decode(createdPlan); err != nil {
		return nil, err
	}
	return createdPlan, nil
}

func (ps *PlanService) GetAllPlans() ([]*Plan, error) {
	cursor, err := ps.DB.Collection("plans").Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	var plans []*Plan
	if err := cursor.All(context.Background(), &plans); err != nil {
		return nil, err
	}
	return plans, nil
}

func (ps *PlanService) GetPlan(id string) (*Plan, error) {
	filter := bson.D{{Key: "_id", Value: id}}
	plan := &Plan{}
	if err := ps.DB.Collection("plans").FindOne(context.Background(), filter).Decode(plan); err != nil {
		return nil, err
	}
	return plan, nil
}

func (ps *PlanService) GetAllPlanByUser(user_id string, day string) ([]*Plan, error) {
	filter := bson.D{{Key: "dayOfWeek", Value: day}, {Key: "userId", Value: user_id}}
	cursor, err := ps.DB.Collection("plans").Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	var plans []*Plan
	if err := cursor.All(context.Background(), &plans); err != nil {
		return nil, err
	}
	return plans, nil
}

func (ps *PlanService) DeletePlan(id string, day string) error {
	filter := bson.D{{Key: "_id", Value: id}, {Key: "dayOfWeek", Value: day}}
	if _, err := ps.DB.Collection("plans").DeleteOne(context.Background(), filter); err != nil {
		return err
	}
	return nil
}

func (ps *PlanService) UpdatePlan(doc *UpdatePlanDto, id string) (*Plan, error) {
	filter := bson.D{{Key: "_id", Value: id}}
	plan, err := ps.GetPlan(id)
	if err != nil {
		return nil, err
	}
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "userId", Value: function.Coalesce(doc.UserID, plan.UserID)},
			{Key: "typeOfPlan", Value: function.Coalesce(doc.TypeOfPlan, plan.TypeOfPlan)},
			{Key: "exercise", Value: function.Coalesce(doc.Exercise, plan.Exercise)},
		}},
	}
	if _, err := ps.DB.Collection("plans").UpdateOne(context.Background(), filter, update); err != nil {
		return nil, err
	}
	updatedRecord := ps.DB.Collection("plans").FindOne(context.Background(), filter)
	updatedPlan := &Plan{}
	if err := updatedRecord.Decode(updatedPlan); err != nil {
		return nil, err
	}
	return updatedPlan, nil
}
