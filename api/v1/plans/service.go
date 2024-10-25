package plans

import (
	"context"
	"errors"

	"github.com/Npwskp/GymsbroBackend/api/v1/function"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Plan struct {
	ID         primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	UserID     string             `json:"userid" validate:"required" bson:"userid"`
	TypeOfPlan string             `json:"typeofplan" validate:"required"`
	DayOfWeek  string             `json:"dayofweek" validate:"required"`
	Exercise   []string           `json:"exercise"`
}

type PlanService struct {
	DB *mongo.Database
}

type IPlanService interface {
	CreatePlan(plan *CreatePlanDto) (*Plan, error)
	GetAllPlans() ([]*Plan, error)
	GetPlan(id string) (*Plan, error)
	GetAllPlanByUser(user_id string) ([]*Plan, error)
	GetPlanByUserDay(user_id string, day string) (*Plan, error)
	DeletePlan(id string) error
	DeleteByUserDay(userid string, day string) error
	UpdatePlan(doc *UpdatePlanDto, id string) (*Plan, error)
	UpdatePlanByUserDay(doc *UpdatePlanDto, userid string, day string) (*Plan, error)
}

func (ps *PlanService) CreatePlan(plan *CreatePlanDto) (*Plan, error) {
	find := bson.D{{Key: "userid", Value: plan.UserID}, {Key: "dayofweek", Value: plan.DayOfWeek}}
	res, err := ps.DB.Collection("plans").CountDocuments(context.Background(), find)
	if err != nil {
		return nil, err
	} else if res > 0 {
		return nil, errors.New("Plan already exist")
	}
	if plan.Exercise == nil {
		plan.Exercise = []string{}
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
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	filter := bson.D{{Key: "_id", Value: oid}}
	plan := &Plan{}
	if err := ps.DB.Collection("plans").FindOne(context.Background(), filter).Decode(plan); err != nil {
		return nil, err
	}
	return plan, nil
}

func (ps *PlanService) GetAllPlanByUser(user_id string) ([]*Plan, error) {
	filter := bson.D{{Key: "userid", Value: user_id}}
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

func (ps *PlanService) GetPlanByUserDay(user_id string, day string) (*Plan, error) {
	filter := bson.D{{Key: "userid", Value: user_id}, {Key: "dayofweek", Value: day}}
	plan := &Plan{}
	if err := ps.DB.Collection("plans").FindOne(context.Background(), filter).Decode(plan); err != nil {
		return nil, err
	}
	return plan, nil
}

func (ps *PlanService) DeletePlan(id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	filter := bson.D{{Key: "_id", Value: oid}}
	if _, err := ps.DB.Collection("plans").DeleteOne(context.Background(), filter); err != nil {
		return err
	}
	return nil
}

func (ps *PlanService) DeleteByUserDay(userid string, day string) error {
	filter := bson.D{{Key: "userid", Value: userid}, {Key: "dayofweek", Value: day}}
	if _, err := ps.DB.Collection("plans").DeleteMany(context.Background(), filter); err != nil {
		return err
	}
	return nil
}

func (ps *PlanService) UpdatePlan(doc *UpdatePlanDto, id string) (*Plan, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	filter := bson.D{{Key: "_id", Value: oid}}
	plan, err := ps.GetPlan(id)
	if err != nil {
		return nil, err
	}
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "typeofplan", Value: function.Coalesce(doc.TypeOfPlan, plan.TypeOfPlan)},
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

func (ps *PlanService) UpdatePlanByUserDay(doc *UpdatePlanDto, userid string, day string) (*Plan, error) {
	filter := bson.D{{Key: "userid", Value: userid}, {Key: "dayofweek", Value: day}}
	plan, err := ps.GetPlanByUserDay(userid, day)
	if err != nil {
		return nil, err
	}
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "typeofplan", Value: function.Coalesce(doc.TypeOfPlan, plan.TypeOfPlan)},
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
