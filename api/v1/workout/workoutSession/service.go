package workoutSession

import (
	"context"
	"errors"
	"time"

	"github.com/Npwskp/GymsbroBackend/api/v1/workout/exerciseLog"
	"github.com/Npwskp/GymsbroBackend/api/v1/workout/workout"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type WorkoutSessionService struct {
	DB *mongo.Database
}

type IWorkoutSessionService interface {
	StartSession(dto *CreateWorkoutSessionDto, userId string) (*WorkoutSession, error)
	EndSession(id string, userId string) (*WorkoutSession, error)
	UpdateSession(id string, dto *UpdateWorkoutSessionDto, userId string) (*WorkoutSession, error)
	GetSession(id string, userId string) (*WorkoutSession, error)
	GetUserSessions(userId string) ([]*WorkoutSession, error)
	GetOnGoingSession(userId string) (*WorkoutSession, error)
	DeleteSession(id string, userId string) error
	LogSession(dto *LoggedSessionDto, userId string) (*WorkoutSession, error)
}

func (s *WorkoutSessionService) StartSession(dto *CreateWorkoutSessionDto, userId string) (*WorkoutSession, error) {
	var exercises []SessionExercise

	ongoingSession, err := s.GetOnGoingSession(userId)
	if err != nil {
		if err != mongo.ErrNoDocuments {
			return nil, err
		}
	}

	if ongoingSession != nil {
		return nil, errors.New("user already has an ongoing session")
	}

	if dto.Type == PlannedSession && dto.WorkoutID != "" {
		// If starting from a plan, fetch and copy the workout exercises
		workout := &workout.Workout{}
		workoutOid, err := primitive.ObjectIDFromHex(dto.WorkoutID)
		if err != nil {
			return nil, err
		}

		err = s.DB.Collection("workout").FindOne(context.Background(), bson.D{
			{Key: "_id", Value: workoutOid},
			{Key: "userid", Value: userId},
		}).Decode(workout)
		if err != nil {
			return nil, err
		}

		// Convert workout exercises to session exercises
		for i, ex := range workout.Exercises {
			exercises = append(exercises, SessionExercise{
				ExerciseID: ex.ExerciseID,
				Order:      i,
			})
		}
	}

	session := &WorkoutSession{
		UserID:      userId,
		WorkoutID:   dto.WorkoutID,
		Type:        dto.Type,
		StartTime:   time.Now(),
		Status:      StatusInProgress,
		TotalVolume: 0,
		Duration:    0,
		Exercises:   exercises,
		Notes:       dto.Notes,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	result, err := s.DB.Collection("workoutSessions").InsertOne(context.Background(), session)
	if err != nil {
		return nil, err
	}

	filter := bson.D{{Key: "_id", Value: result.InsertedID}}
	createdSession := &WorkoutSession{}
	err = s.DB.Collection("workoutSessions").FindOne(context.Background(), filter).Decode(createdSession)
	if err != nil {
		return nil, err
	}

	return createdSession, nil
}

func (s *WorkoutSessionService) EndSession(id string, userId string) (*WorkoutSession, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.D{
		{Key: "_id", Value: oid},
		{Key: "userid", Value: userId},
	}

	session := &WorkoutSession{}
	if err := s.DB.Collection("workoutSessions").FindOne(context.Background(), filter).Decode(session); err != nil {
		return nil, err
	}

	if session.Status != StatusInProgress {
		return nil, errors.New("session is not in progress")
	}

	endTime := time.Now()
	duration := int(endTime.Sub(session.StartTime).Seconds())

	// Calculate total volume from exercise logs
	var totalVolume float64
	for _, exercise := range session.Exercises {
		if exercise.ExerciseLogID != "" {
			logOid, err := primitive.ObjectIDFromHex(exercise.ExerciseLogID)
			if err != nil {
				continue
			}

			log := &exerciseLog.ExerciseLog{}
			err = s.DB.Collection("exerciseLogs").FindOne(context.Background(), bson.D{
				{Key: "_id", Value: logOid},
			}).Decode(log)
			if err != nil {
				continue
			}

			totalVolume += log.TotalVolume
		}
	}

	update := bson.D{{Key: "$set", Value: bson.D{
		{Key: "end_time", Value: endTime},
		{Key: "status", Value: StatusCompleted},
		{Key: "duration", Value: duration},
		{Key: "total_volume", Value: totalVolume},
		{Key: "updated_at", Value: time.Now()},
	}}}

	after := options.After
	opts := options.FindOneAndUpdate().SetReturnDocument(after)

	result := &WorkoutSession{}
	err = s.DB.Collection("workoutSessions").FindOneAndUpdate(context.Background(), filter, update, opts).Decode(result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *WorkoutSessionService) UpdateSession(id string, dto *UpdateWorkoutSessionDto, userId string) (*WorkoutSession, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.D{
		{Key: "_id", Value: oid},
		{Key: "userid", Value: userId},
	}

	update := bson.D{{Key: "$set", Value: bson.D{
		{Key: "status", Value: dto.Status},
		{Key: "exercises", Value: dto.Exercises},
		{Key: "notes", Value: dto.Notes},
		{Key: "updatedAt", Value: time.Now()},
	}}}

	after := options.After
	opts := options.FindOneAndUpdate().SetReturnDocument(after)

	result := &WorkoutSession{}
	err = s.DB.Collection("workoutSessions").FindOneAndUpdate(context.Background(), filter, update, opts).Decode(result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *WorkoutSessionService) GetSession(id string, userId string) (*WorkoutSession, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.D{
		{Key: "_id", Value: oid},
		{Key: "userid", Value: userId},
	}

	session := &WorkoutSession{}
	if err := s.DB.Collection("workoutSessions").FindOne(context.Background(), filter).Decode(session); err != nil {
		return nil, err
	}

	return session, nil
}

func (s *WorkoutSessionService) GetUserSessions(userId string) ([]*WorkoutSession, error) {
	filter := bson.D{{Key: "userid", Value: userId}}
	opts := options.Find().SetSort(bson.D{{Key: "startTime", Value: -1}})

	cursor, err := s.DB.Collection("workoutSessions").Find(context.Background(), filter, opts)
	if err != nil {
		return nil, err
	}

	sessions := make([]*WorkoutSession, 0)
	if err := cursor.All(context.Background(), &sessions); err != nil {
		return nil, err
	}

	return sessions, nil
}

func (s *WorkoutSessionService) GetOnGoingSession(userId string) (*WorkoutSession, error) {
	filter := bson.D{{Key: "userid", Value: userId}, {Key: "status", Value: StatusInProgress}}
	opts := options.FindOne().SetSort(bson.D{{Key: "startTime", Value: -1}})

	session := &WorkoutSession{}
	err := s.DB.Collection("workoutSessions").FindOne(context.Background(), filter, opts).Decode(session)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (s *WorkoutSessionService) DeleteSession(id string, userId string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.D{
		{Key: "_id", Value: oid},
		{Key: "userid", Value: userId},
	}

	result, err := s.DB.Collection("workoutSessions").DeleteOne(context.Background(), filter)
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

func (s *WorkoutSessionService) LogSession(dto *LoggedSessionDto, userId string) (*WorkoutSession, error) {
	session := &WorkoutSession{
		UserID:    userId,
		WorkoutID: dto.WorkoutID,
		Type:      LoggedSession,
		StartTime: dto.StartTime,
		EndTime:   dto.EndTime,
		Status:    dto.Status,
		Exercises: dto.Exercises,
		Notes:     dto.Notes,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Calculate duration in minutes
	duration := int(dto.EndTime.Sub(dto.StartTime).Minutes())
	if duration < 0 {
		return nil, errors.New("end time must be after start time")
	}
	session.Duration = duration

	// Calculate total volume if exercises are provided
	var totalVolume float64
	for _, exercise := range dto.Exercises {
		// Fetch exercise log to get volume
		if exercise.ExerciseLogID != "" {
			logOid, err := primitive.ObjectIDFromHex(exercise.ExerciseLogID)
			if err != nil {
				return nil, err
			}

			exerciseLog := &exerciseLog.ExerciseLog{}
			err = s.DB.Collection("exerciseLog").FindOne(context.Background(), bson.D{
				{Key: "_id", Value: logOid},
				{Key: "userid", Value: userId},
			}).Decode(exerciseLog)
			if err != nil {
				return nil, err
			}

			totalVolume += exerciseLog.TotalVolume
		}
	}
	session.TotalVolume = totalVolume

	result, err := s.DB.Collection("workoutSessions").InsertOne(context.Background(), session)
	if err != nil {
		return nil, err
	}

	createdSession := &WorkoutSession{}
	err = s.DB.Collection("workoutSessions").FindOne(context.Background(), bson.D{{Key: "_id", Value: result.InsertedID}}).Decode(createdSession)
	if err != nil {
		return nil, err
	}

	return createdSession, nil
}
