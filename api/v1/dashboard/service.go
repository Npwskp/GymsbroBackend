package dashboard

import (
	"context"
	"encoding/json"
	"errors"
	"math"
	"os"
	"time"

	authEnums "github.com/Npwskp/GymsbroBackend/api/v1/auth/enums"
	dashboardEnums "github.com/Npwskp/GymsbroBackend/api/v1/dashboard/enums"
	dashboardFunctions "github.com/Npwskp/GymsbroBackend/api/v1/dashboard/functions"
	"github.com/Npwskp/GymsbroBackend/api/v1/user"
	"github.com/Npwskp/GymsbroBackend/api/v1/workout/exercise"
	exerciseEnums "github.com/Npwskp/GymsbroBackend/api/v1/workout/exercise/enums"
	"github.com/Npwskp/GymsbroBackend/api/v1/workout/exerciseLog"
	"github.com/Npwskp/GymsbroBackend/api/v1/workout/workoutSession"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type DashboardService struct {
	DB *mongo.Database
}

type IDashboardService interface {
	GetDashboard(userId string) (*DashboardResponse, error)
	GetUserStrengthStandards(userId string) (*UserStrengthStandards, error)
	GetRepMax(userId string, exerciseId string) (*RepMaxResponse, error)
}

func getTimeOfDay(t time.Time) string {
	hour := t.Hour()
	switch {
	case hour >= 5 && hour < 12:
		return "Morning"
	case hour >= 12 && hour < 17:
		return "Afternoon"
	case hour >= 17 && hour < 22:
		return "Evening"
	default:
		return "Night"
	}
}

func calculateMovingAverage(values []int, window int) []float64 {
	result := make([]float64, len(values))
	for i := range values {
		count := 0
		sum := 0
		for j := max(0, i-window+1); j <= i; j++ {
			sum += values[j]
			count++
		}
		result[i] = float64(sum) / float64(count)
	}
	return result
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (ds *DashboardService) GetDashboard(userId string) (*DashboardResponse, error) {
	// Get workout sessions
	sessionFilter := bson.D{{Key: "userid", Value: userId}}
	cursor, err := ds.DB.Collection("workoutSessions").Find(context.Background(), sessionFilter)
	if err != nil {
		return nil, err
	}

	var sessions []workoutSession.WorkoutSession
	if err := cursor.All(context.Background(), &sessions); err != nil {
		return nil, err
	}

	// Get exercise logs
	logFilter := bson.D{{Key: "userid", Value: userId}}
	logCursor, err := ds.DB.Collection("exerciseLogs").Find(context.Background(), logFilter)
	if err != nil {
		return nil, err
	}

	var exerciseLogs []exerciseLog.ExerciseLog
	if err := logCursor.All(context.Background(), &exerciseLogs); err != nil {
		return nil, err
	}

	// Initialize response
	response := &DashboardResponse{
		FrequencyGraph: FrequencyGraphData{
			Labels: make([]string, 30),
			Values: make([]int, 30),
		},
		Analysis: WorkoutAnalysis{
			TotalWorkouts:  len(sessions),
			TotalExercises: len(exerciseLogs),
		},
	}

	// Prepare frequency graph data
	now := time.Now()
	dailyCount := make(map[string]int)
	dayFrequency := make(map[string]int)  // For tracking most active day
	timeFrequency := make(map[string]int) // For tracking most active time
	var totalVolume float64
	var currentStreak, bestStreak, streak int

	// Process last 30 days for graph
	for i := 0; i < 30; i++ {
		date := now.AddDate(0, 0, -i)
		dateStr := date.Format("2006-01-02")
		response.FrequencyGraph.Labels[29-i] = dateStr
	}

	// Process sessions
	for _, session := range sessions {
		dateStr := session.StartTime.Format("2006-01-02")
		dayOfWeek := session.StartTime.Weekday().String()
		timeOfDay := getTimeOfDay(session.StartTime)

		dailyCount[dateStr]++
		dayFrequency[dayOfWeek]++
		timeFrequency[timeOfDay]++
		totalVolume += session.TotalVolume

		// Count recent activity
		daysAgo := int(now.Sub(session.StartTime).Hours() / 24)
		if daysAgo < 7 {
			response.Analysis.LastWeekCount++
		}
		if daysAgo < 30 {
			response.Analysis.LastMonthCount++
		}
	}

	// Process exercise logs
	for _, log := range exerciseLogs {
		dateStr := log.CreatedAt.Format("2006-01-02")
		dayOfWeek := log.CreatedAt.Weekday().String()
		timeOfDay := getTimeOfDay(log.CreatedAt)

		dailyCount[dateStr]++
		dayFrequency[dayOfWeek]++
		timeFrequency[timeOfDay]++
		totalVolume += log.TotalVolume
	}

	// Fill in frequency graph values
	for i, dateStr := range response.FrequencyGraph.Labels {
		response.FrequencyGraph.Values[i] = dailyCount[dateStr]

		// Calculate streaks
		if dailyCount[dateStr] > 0 {
			streak++
			if streak > bestStreak {
				bestStreak = streak
			}
			// If this is today or yesterday and we have activity, update current streak
			if i >= 28 { // last two days
				currentStreak = streak
			}
		} else {
			streak = 0
		}
	}

	// Calculate trend line (7-day moving average)
	response.FrequencyGraph.TrendLine = calculateMovingAverage(response.FrequencyGraph.Values, 7)

	// Find most active day and time
	var maxDayCount, maxTimeCount int
	for day, count := range dayFrequency {
		if count > maxDayCount {
			maxDayCount = count
			response.Analysis.MostActiveDay = day
		}
	}
	for timeOfDay, count := range timeFrequency {
		if count > maxTimeCount {
			maxTimeCount = count
			response.Analysis.MostActiveTime = timeOfDay
		}
	}

	// Calculate weekly average
	firstActivity := now
	if len(sessions) > 0 && len(exerciseLogs) > 0 {
		firstSessionDate := sessions[0].StartTime
		firstLogDate := exerciseLogs[0].CreatedAt
		if firstSessionDate.Before(firstLogDate) {
			firstActivity = firstSessionDate
		} else {
			firstActivity = firstLogDate
		}
	} else if len(sessions) > 0 {
		firstActivity = sessions[0].StartTime
	} else if len(exerciseLogs) > 0 {
		firstActivity = exerciseLogs[0].CreatedAt
	}

	weeks := max(1, int(now.Sub(firstActivity).Hours()/(24*7)))
	response.Analysis.AveragePerWeek = float64(len(sessions)+len(exerciseLogs)) / float64(weeks)
	response.Analysis.TotalVolume = totalVolume
	response.Analysis.BestStreak = bestStreak
	response.Analysis.CurrentStreak = currentStreak

	return response, nil
}

func (ds *DashboardService) GetUserStrengthStandards(userId string) (*UserStrengthStandards, error) {
	var userObj user.User
	err := ds.DB.Collection("users").FindOne(context.Background(), bson.D{{Key: "userid", Value: userId}}).Decode(&userObj)
	if err != nil {
		return nil, err
	}

	userBodyWeight := userObj.Weight
	userGender := userObj.Gender

	if userGender == authEnums.GenderMale {
		if userBodyWeight < 50 || userBodyWeight > 140 {
			return nil, errors.New("bodyweight out of range of strength standards processing")
		}
	} else if userGender == authEnums.GenderFemale {
		if userBodyWeight < 45 || userBodyWeight > 120 {
			return nil, errors.New("bodyweight out of range of strength standards processing")
		}
	}

	// Create $or conditions from ConsiderExercises
	orConditions := make([]bson.D, len(dashboardEnums.ConsiderExercises))
	for i, exerciseEquip := range dashboardEnums.ConsiderExercises {
		orConditions[i] = bson.D{
			{Key: "name", Value: exerciseEquip.Exercise},
			{Key: "equipment", Value: exerciseEquip.Equipment},
		}
	}

	exercise_pipeline := []bson.D{
		{{Key: "$match", Value: bson.D{
			{Key: "$or", Value: orConditions},
		}}},
	}

	exercise_cursor, err := ds.DB.Collection("exercises").Aggregate(context.Background(), exercise_pipeline)
	if err != nil {
		return nil, err
	}
	defer exercise_cursor.Close(context.Background())

	var exercises []exercise.Exercise
	if err := exercise_cursor.All(context.Background(), &exercises); err != nil {
		return nil, err
	}

	if len(exercises) != len(dashboardEnums.ConsiderExercises) {
		return nil, errors.New("consider exercise count mismatch")
	}

	// Map queried exercises to ConsiderExercises
	exerciseMap := make(map[string]exercise.Exercise)
	for _, ex := range exercises {
		exerciseMap[ex.ID.Hex()] = ex
	}

	// Now you can use exerciseMap to look up exercise info by ID
	exerciseIds := make([]string, 0)
	for _, exercise := range exercises {
		exerciseIds = append(exerciseIds, exercise.ID.Hex())
	}

	// Find Latest Exercise Logs for each exercise
	pipeline := []bson.D{
		{{Key: "$match", Value: bson.D{
			{Key: "userid", Value: userId},
			{Key: "exerciseid", Value: bson.D{{Key: "$in", Value: exerciseIds}}},
		}}},
		{{Key: "$sort", Value: bson.D{
			{Key: "exerciseid", Value: 1},
			{Key: "datetime", Value: -1},
		}}},
		{{Key: "$group", Value: bson.D{
			{Key: "exerciseid", Value: "$exerciseid"},
			{Key: "doc", Value: bson.D{{Key: "$first", Value: "$$ROOT"}}},
		}}},
		{{Key: "$replaceRoot", Value: bson.D{{Key: "newRoot", Value: "$doc"}}}},
	}

	cursor, err := ds.DB.Collection("exerciseLogs").Aggregate(context.Background(), pipeline)
	if err != nil {
		return nil, err
	}

	var latestLogs []exerciseLog.ExerciseLog
	if err := cursor.All(context.Background(), &latestLogs); err != nil {
		return nil, err
	}

	// Create a map of exercise name to latest log for easy lookup
	latestLogsMap := make(map[string]exerciseLog.ExerciseLog)
	for _, log := range latestLogs {
		latestLogsMap[log.ExerciseID] = log
	}

	// Replace the inline struct with the one from enums
	var strengthStandards struct {
		Male   dashboardEnums.StrengthStandards `json:"Male"`
		Female dashboardEnums.StrengthStandards `json:"Female"`
	}

	// Read and parse the JSON file
	jsonFile, err := os.ReadFile("api/v1/dashboard/json/strengthStandard.json")
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(jsonFile, &strengthStandards); err != nil {
		return nil, err
	}

	// Calculate strength standards for each exercise
	userStrengthStandards := make([]UserStrengthStandardPerExercise, 0)
	muscleGroupStrengths := make(map[exerciseEnums.TargetMuscle][]float64)

	for exerciseId, considerExercise := range exerciseMap {
		latestLog, exists := latestLogsMap[exerciseId]

		if !exists {
			continue
		}

		exerciseName := considerExercise.Name
		exerciseEquipment := considerExercise.Equipment

		// Find the closest bodyweight standards
		var standards dashboardEnums.StrengthStandards

		if userGender == "Male" {
			standards = strengthStandards.Male
		} else {
			standards = strengthStandards.Female
		}

		exerciseStandards, exists := standards[exerciseName]
		if !exists {
			continue
		}

		// Find closest bodyweight bracket
		closestStandardWeight := math.Floor(userBodyWeight/5) * 5
		var closestStandard dashboardEnums.StrengthStandard

		for _, standard := range exerciseStandards {
			if standard.Bodyweight == closestStandardWeight {
				closestStandard = standard
				break
			}
		}

		// Find the set with maximum weight
		var maxWeight float64
		var maxReps int
		for _, set := range latestLog.Sets {
			if set.Weight*float64(set.Reps) > maxWeight*float64(maxReps) {
				maxWeight = set.Weight
				maxReps = set.Reps
			}
		}

		// Calculate 1RM using Brzycki formula if more than 1 rep
		if maxReps > 1 {
			maxWeight, err = dashboardFunctions.CalculateOneRepMax(maxWeight, float64(maxReps))
			if err != nil {
				return nil, err
			}
		}

		// Calculate relative strength (as percentage of bodyweight)
		relativeStrength := maxWeight / userBodyWeight

		// Calculate strength level based on standards
		var strengthLevel dashboardEnums.StrengthType
		var score float64

		switch {
		case maxWeight <= closestStandard.Standards.Beginner:
			strengthLevel = dashboardEnums.StrengthTypeBeginner
			score = (maxWeight / closestStandard.Standards.Beginner) * dashboardEnums.BeginnerScore
		case maxWeight <= closestStandard.Standards.Novice:
			strengthLevel = dashboardEnums.StrengthTypeNovice
			score = dashboardEnums.BeginnerScore + ((maxWeight-closestStandard.Standards.Beginner)/(closestStandard.Standards.Novice-closestStandard.Standards.Beginner))*(dashboardEnums.NoviceScore-dashboardEnums.BeginnerScore)
		case maxWeight <= closestStandard.Standards.Intermediate:
			strengthLevel = dashboardEnums.StrengthTypeIntermediate
			score = dashboardEnums.NoviceScore + ((maxWeight-closestStandard.Standards.Novice)/(closestStandard.Standards.Intermediate-closestStandard.Standards.Novice))*(dashboardEnums.IntermediateScore-dashboardEnums.NoviceScore)
		case maxWeight <= closestStandard.Standards.Advanced:
			strengthLevel = dashboardEnums.StrengthTypeAdvanced
			score = dashboardEnums.IntermediateScore + ((maxWeight-closestStandard.Standards.Intermediate)/(closestStandard.Standards.Advanced-closestStandard.Standards.Intermediate))*(dashboardEnums.AdvancedScore-dashboardEnums.IntermediateScore)
		case maxWeight <= closestStandard.Standards.Elite:
			strengthLevel = dashboardEnums.StrengthTypeElite
			score = dashboardEnums.AdvancedScore + ((maxWeight-closestStandard.Standards.Advanced)/(closestStandard.Standards.Elite-closestStandard.Standards.Advanced))*(dashboardEnums.EliteScore-dashboardEnums.AdvancedScore)
		default:
			strengthLevel = dashboardEnums.StrengthTypeElite
			score = dashboardEnums.EliteScore
		}

		// Add to muscle group calculations
		for _, muscleGroup := range considerExercise.TargetMuscle {
			muscleGroupStrengths[muscleGroup] = append(muscleGroupStrengths[muscleGroup], score)
		}

		// Add to exercise standards
		userStrengthStandards = append(userStrengthStandards, UserStrengthStandardPerExercise{
			Exercise:         exerciseName,
			Equipment:        exerciseEquipment,
			RepMax:           maxWeight,
			RelativeStrength: relativeStrength,
			StrengthLevel:    strengthLevel,
			Score:            score,
			LastPerformed:    latestLog.DateTime,
		})
	}

	// Calculate average strength per muscle group
	muscleGroupStandards := make([]UserStrengthStandardPerMuscleGroup, 0)
	for muscleGroup, scores := range muscleGroupStrengths {
		var totalScore float64
		for _, score := range scores {
			totalScore += score
		}
		avgScore := totalScore / float64(len(scores))

		muscleGroupStandards = append(muscleGroupStandards, UserStrengthStandardPerMuscleGroup{
			TargetMuscle:  muscleGroup,
			StrengthLevel: dashboardEnums.ClassifyStrength(avgScore),
			Score:         avgScore,
		})
	}

	return &UserStrengthStandards{
		ExerciseStandards:    userStrengthStandards,
		MuscleGroupStrengths: muscleGroupStandards,
	}, nil
}

func (ds *DashboardService) GetRepMax(userId string, exerciseId string) (*RepMaxResponse, error) {
	// Create an aggregation pipeline to calculate 1RM at database level
	pipeline := []bson.D{
		{{Key: "$match", Value: bson.D{
			{Key: "userid", Value: userId},
			{Key: "exerciseid", Value: exerciseId},
		}}},
		{{Key: "$unwind", Value: "$sets"}},
		{{Key: "$match", Value: bson.D{
			{Key: "sets.weight", Value: bson.D{{Key: "$gt", Value: 0}}},
			{Key: "sets.reps", Value: bson.D{{Key: "$gt", Value: 0}}},
		}}},
		{{Key: "$addFields", Value: bson.D{
			{Key: "oneRM", Value: bson.D{
				{Key: "$multiply", Value: bson.A{
					"$sets.weight",
					bson.D{{Key: "$divide", Value: bson.A{
						36.0,
						bson.D{{Key: "$subtract", Value: bson.A{37.0, "$sets.reps"}}},
					}}},
				}},
			}},
		}}},
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: nil},
			{Key: "bestOneRM", Value: bson.D{{Key: "$max", Value: "$oneRM"}}},
			{Key: "lastUpdated", Value: bson.D{{Key: "$max", Value: "$datetime"}}},
		}}},
	}

	cursor, err := ds.DB.Collection("exerciseLogs").Aggregate(context.Background(), pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	// Get the result
	var results []struct {
		BestOneRM   float64   `bson:"bestOneRM"`
		LastUpdated time.Time `bson:"lastUpdated"`
	}

	if err := cursor.All(context.Background(), &results); err != nil {
		return nil, err
	}

	if len(results) == 0 || results[0].BestOneRM == 0 {
		return nil, errors.New("no valid sets found for rep max calculation")
	}

	bestOneRM := math.Round(results[0].BestOneRM*100) / 100

	// Calculate other rep maxes using the best one rep max
	eightRM, err := dashboardFunctions.EstimateRepMax(bestOneRM, 8)
	if err != nil {
		return nil, err
	}

	twelveRM, err := dashboardFunctions.EstimateRepMax(bestOneRM, 12)
	if err != nil {
		return nil, err
	}

	return &RepMaxResponse{
		OneRepMax:    bestOneRM,
		EightRepMax:  eightRM,
		TwelveRepMax: twelveRM,
		LastUpdated:  results[0].LastUpdated,
	}, nil
}
