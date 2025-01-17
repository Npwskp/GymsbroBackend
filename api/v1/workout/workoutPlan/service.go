package workoutPlan

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type WorkoutPlanService struct {
	DB *mongo.Database
}

type IWorkoutPlanService interface {
	CreatePlanByDaysOfWeek(dto *CreatePlanByDaysOfWeekDto, userId string) ([]WorkoutPlan, error)
	CreatePlanByCyclicWorkout(dto *CreatePlanByCyclicWorkoutDto, userId string) ([]WorkoutPlan, error)
	GetWorkoutPlansByUser(userId string) ([]WorkoutPlan, error)
}

func (s *WorkoutPlanService) CreatePlanByDaysOfWeek(dto *CreatePlanByDaysOfWeekDto, userId string) ([]WorkoutPlan, error) {
	currentDate := time.Now()

	// Map to group dates by workoutID
	workoutDates := make(map[string][]time.Time)

	// Calculate total days based on weeks duration
	totalDays := dto.WeeksDuration * 7

	// Iterate through each day and group dates by workout
	for i := 0; i < totalDays; i++ {
		date := currentDate.AddDate(0, 0, i)
		var workoutId string

		switch date.Weekday() {
		case time.Monday:
			workoutId = dto.MondayWorkoutID
		case time.Tuesday:
			workoutId = dto.TuesdayWorkoutID
		case time.Wednesday:
			workoutId = dto.WednesdayWorkoutID
		case time.Thursday:
			workoutId = dto.ThursdayWorkoutID
		case time.Friday:
			workoutId = dto.FridayWorkoutID
		case time.Saturday:
			workoutId = dto.SaturdayWorkoutID
		case time.Sunday:
			workoutId = dto.SundayWorkoutID
		}

		workoutDates[workoutId] = append(workoutDates[workoutId], date)
	}

	// Get existing workout plans
	var existingPlans []WorkoutPlan
	filter := bson.M{
		"userid": userId,
		"workoutid": bson.M{
			"$in": []string{
				dto.MondayWorkoutID,
				dto.TuesdayWorkoutID,
				dto.WednesdayWorkoutID,
				dto.ThursdayWorkoutID,
				dto.FridayWorkoutID,
				dto.SaturdayWorkoutID,
				dto.SundayWorkoutID,
			},
		},
	}

	cursor, err := s.DB.Collection("workoutPlan").Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	if err := cursor.All(context.Background(), &existingPlans); err != nil {
		return nil, err
	}

	// Separate workouts into existing and new
	existingWorkoutIds := make(map[string]bool)
	for _, plan := range existingPlans {
		existingWorkoutIds[plan.WorkoutID] = true
	}

	var updatedPlans []WorkoutPlan
	var newWorkouts []interface{}

	// Handle existing plans: update with new dates
	for _, plan := range existingPlans {
		newDates := workoutDates[plan.WorkoutID]
		// Filter out past dates from existing plan
		var futureDates []time.Time
		for _, date := range plan.Dates {
			if !date.Before(currentDate) {
				futureDates = append(futureDates, date)
			}
		}
		// Combine with new dates
		plan.Dates = append(futureDates, newDates...)
		plan.UpdatedAt = time.Now()

		// Update in database
		_, err := s.DB.Collection("workoutPlan").UpdateOne(
			context.Background(),
			bson.M{"_id": plan.ID},
			bson.M{"$set": bson.M{
				"dates":      plan.Dates,
				"updated_at": plan.UpdatedAt,
			}},
		)
		if err != nil {
			return nil, err
		}
		updatedPlans = append(updatedPlans, plan)
	}

	// Handle new plans: create new entries
	for workoutId, dates := range workoutDates {
		if !existingWorkoutIds[workoutId] {
			newWorkouts = append(newWorkouts, WorkoutPlan{
				UserID:    userId,
				WorkoutID: workoutId,
				Dates:     dates,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			})
		}
	}

	// Insert new workout plans if any
	if len(newWorkouts) > 0 {
		result, err := s.DB.Collection("workoutPlan").InsertMany(context.Background(), newWorkouts)
		if err != nil {
			return nil, err
		}

		// Fetch newly inserted documents
		newFilter := bson.M{"_id": bson.M{"$in": result.InsertedIDs}}
		cursor, err := s.DB.Collection("workoutPlan").Find(context.Background(), newFilter)
		if err != nil {
			return nil, err
		}
		defer cursor.Close(context.Background())

		var newPlans []WorkoutPlan
		if err := cursor.All(context.Background(), &newPlans); err != nil {
			return nil, err
		}
		updatedPlans = append(updatedPlans, newPlans...)
	}

	return updatedPlans, nil
}

func (s *WorkoutPlanService) CreatePlanByCyclicWorkout(dto *CreatePlanByCyclicWorkoutDto, userId string) ([]WorkoutPlan, error) {
	currentDate := time.Now()

	// Map to group dates by workoutID
	workoutDates := make(map[string][]time.Time)

	// Calculate total days based on weeks duration
	totalDays := dto.WeeksDuration * 7

	// Iterate through each day and assign workouts cyclically
	for i := 0; i < totalDays; i++ {
		date := currentDate.AddDate(0, 0, i)
		// Get workout ID from the cycle using modulo operation
		workoutId := dto.WorkoutIDs[i%len(dto.WorkoutIDs)]
		workoutDates[workoutId] = append(workoutDates[workoutId], date)
	}

	// Get existing workout plans
	var existingPlans []WorkoutPlan
	filter := bson.M{
		"userid": userId,
		"workoutid": bson.M{
			"$in": dto.WorkoutIDs,
		},
	}

	cursor, err := s.DB.Collection("workoutPlan").Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	if err := cursor.All(context.Background(), &existingPlans); err != nil {
		return nil, err
	}

	// Separate workouts into existing and new
	existingWorkoutIds := make(map[string]bool)
	for _, plan := range existingPlans {
		existingWorkoutIds[plan.WorkoutID] = true
	}

	var updatedPlans []WorkoutPlan
	var newWorkouts []interface{}

	// Handle existing plans: update with new dates
	for _, plan := range existingPlans {
		newDates := workoutDates[plan.WorkoutID]
		// Filter out past dates from existing plan
		var futureDates []time.Time
		for _, date := range plan.Dates {
			if !date.Before(currentDate) {
				futureDates = append(futureDates, date)
			}
		}
		// Combine with new dates
		plan.Dates = append(futureDates, newDates...)
		plan.UpdatedAt = time.Now()

		// Update in database
		_, err := s.DB.Collection("workoutPlan").UpdateOne(
			context.Background(),
			bson.M{"_id": plan.ID},
			bson.M{"$set": bson.M{
				"dates":      plan.Dates,
				"updated_at": plan.UpdatedAt,
			}},
		)
		if err != nil {
			return nil, err
		}
		updatedPlans = append(updatedPlans, plan)
	}

	// Handle new plans: create new entries
	for workoutId, dates := range workoutDates {
		if !existingWorkoutIds[workoutId] {
			newWorkouts = append(newWorkouts, WorkoutPlan{
				UserID:    userId,
				WorkoutID: workoutId,
				Dates:     dates,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			})
		}
	}

	// Insert new workout plans if any
	if len(newWorkouts) > 0 {
		result, err := s.DB.Collection("workoutPlan").InsertMany(context.Background(), newWorkouts)
		if err != nil {
			return nil, err
		}

		// Fetch newly inserted documents
		newFilter := bson.M{"_id": bson.M{"$in": result.InsertedIDs}}
		cursor, err := s.DB.Collection("workoutPlan").Find(context.Background(), newFilter)
		if err != nil {
			return nil, err
		}
		defer cursor.Close(context.Background())

		var newPlans []WorkoutPlan
		if err := cursor.All(context.Background(), &newPlans); err != nil {
			return nil, err
		}
		updatedPlans = append(updatedPlans, newPlans...)
	}

	return updatedPlans, nil
}

func (s *WorkoutPlanService) GetWorkoutPlansByUser(userId string) ([]WorkoutPlan, error) {
	var workoutPlans []WorkoutPlan

	filter := bson.M{"userid": userId}
	cursor, err := s.DB.Collection("workoutPlan").Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	if err := cursor.All(context.Background(), &workoutPlans); err != nil {
		return nil, err
	}

	return workoutPlans, nil
}
