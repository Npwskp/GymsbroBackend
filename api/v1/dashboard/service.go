package dashboard

import (
	"context"
	"encoding/json"
	"errors"
	"math"
	"os"
	"sort"
	"time"

	authEnums "github.com/Npwskp/GymsbroBackend/api/v1/auth/enums"
	dashboardEnums "github.com/Npwskp/GymsbroBackend/api/v1/dashboard/enums"
	dashboardFunctions "github.com/Npwskp/GymsbroBackend/api/v1/dashboard/functions"
	foodLog "github.com/Npwskp/GymsbroBackend/api/v1/nutrition/foodLog"
	"github.com/Npwskp/GymsbroBackend/api/v1/user"
	userFitnessPreferenceEnums "github.com/Npwskp/GymsbroBackend/api/v1/user/enums"
	"github.com/Npwskp/GymsbroBackend/api/v1/workout/exercise"
	exerciseEnums "github.com/Npwskp/GymsbroBackend/api/v1/workout/exercise/enums"
	"github.com/Npwskp/GymsbroBackend/api/v1/workout/exerciseLog"
	"github.com/Npwskp/GymsbroBackend/api/v1/workout/workoutSession"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DashboardService struct {
	DB *mongo.Database
}

type ExerciseData struct {
	ID           string                    `bson:"_id"`
	RootExercise exercise.Exercise         `bson:"exercise"`
	Logs         []exerciseLog.ExerciseLog `bson:"logs"`
}

type IDashboardService interface {
	GetDashboard(userId string, startDate, endDate time.Time) (*DashboardResponse, error)
	GetUserStrengthStandards(userId string) (*UserStrengthStandards, error)
	GetRepMax(userId string, exerciseId string, useLatest bool) (*RepMaxResponse, error)
	GetNutritionSummary(userid string, startDate, endDate time.Time) (*NutritionSummaryResponse, error)
	GetBodyCompositionAnalysis(userId string, startDate, endDate time.Time) (*BodyCompositionAnalysisResponse, error)
}

func calculateMovingAverageInt(values []int, window int) []float64 {
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

// Helper function to calculate max set volume
func calculateMaxSetVolume(sets []exerciseLog.SetLog) float64 {
	maxSetVolume := 0.0
	for _, set := range sets {
		setVolume := set.Weight * float64(set.Reps)
		if setVolume > maxSetVolume {
			maxSetVolume = setVolume
		}
	}
	return maxSetVolume
}

// Helper function to calculate best 1RM from sets
func calculateBestOneRM(sets []exerciseLog.SetLog) (float64, error) {
	bestOneRM := 0.0
	for _, set := range sets {
		if set.Weight <= 0 || set.Reps <= 0 {
			continue
		}
		oneRM, err := dashboardFunctions.CalculateOneRepMax(set.Weight, float64(set.Reps))
		if err != nil {
			continue
		}
		if oneRM > bestOneRM {
			bestOneRM = oneRM
		}
	}
	if bestOneRM <= 0 {
		return 0, errors.New("no valid sets found for 1RM calculation")
	}
	return bestOneRM, nil
}

func (ds *DashboardService) GetDashboard(userId string, startDate, endDate time.Time) (*DashboardResponse, error) {
	// Get workout sessions with date filter
	sessionFilter := bson.D{
		{Key: "userid", Value: userId},
		{Key: "start_time", Value: bson.D{
			{Key: "$gte", Value: startDate},
			{Key: "$lte", Value: endDate},
		}},
	}
	cursor, err := ds.DB.Collection("workoutSessions").Find(context.Background(), sessionFilter)
	if err != nil {
		return nil, err
	}

	sessions := make([]workoutSession.WorkoutSession, 0)
	if err := cursor.All(context.Background(), &sessions); err != nil {
		return nil, err
	}

	// Get exercise logs with date filter
	logFilter := bson.D{
		{Key: "userid", Value: userId},
		{Key: "datetime", Value: bson.D{
			{Key: "$gte", Value: startDate},
			{Key: "$lte", Value: endDate},
		}},
	}
	logCursor, err := ds.DB.Collection("exerciseLogs").Find(context.Background(), logFilter)
	if err != nil {
		return nil, err
	}

	exerciseLogs := make([]exerciseLog.ExerciseLog, 0)
	if err := logCursor.All(context.Background(), &exerciseLogs); err != nil {
		return nil, err
	}

	// Calculate number of days between start and end date
	days := int(endDate.Sub(startDate).Hours()/24) + 1

	// Initialize response
	response := &DashboardResponse{
		FrequencyGraph: FrequencyGraphData{
			Labels: make([]string, days),
			Values: make([]int, days),
		},
		Analysis: WorkoutAnalysis{
			TotalWorkouts:  len(sessions),
			TotalExercises: len(exerciseLogs),
		},
	}

	// Prepare frequency graph data
	dailyCount := make(map[string]int)
	var totalVolume float64
	var totalDuration float64

	// Process sessions
	for _, session := range sessions {
		dateStr := session.StartTime.Format("2006-01-02")
		dailyCount[dateStr]++
		totalVolume += session.TotalVolume
		totalDuration += float64(session.Duration)
	}

	// Process exercise logs
	for _, log := range exerciseLogs {
		dateStr := log.CreatedAt.Format("2006-01-02")
		dailyCount[dateStr]++
		totalVolume += log.TotalVolume
	}

	// Fill in frequency graph values for each day in the date range
	for i := 0; i < days; i++ {
		date := startDate.AddDate(0, 0, i)
		dateStr := date.Format("2006-01-02")
		response.FrequencyGraph.Labels[i] = dateStr
		response.FrequencyGraph.Values[i] = dailyCount[dateStr]
	}

	// Calculate trend line (7-day moving average)
	response.FrequencyGraph.TrendLine = calculateMovingAverageInt(response.FrequencyGraph.Values, 7)
	response.Analysis.TotalVolume = totalVolume

	// Handle potential division by zero for average workout duration
	if len(sessions) > 0 {
		avgDuration := totalDuration / float64(len(sessions))
		if !math.IsNaN(avgDuration) && !math.IsInf(avgDuration, 0) {
			response.Analysis.AverageWorkoutDuration = avgDuration
		} else {
			response.Analysis.AverageWorkoutDuration = 0
		}
	} else {
		response.Analysis.AverageWorkoutDuration = 0
	}

	exerciseData, err := ds.getExerciseLogsData(userId, startDate, endDate)
	if err != nil {
		return nil, err
	}

	// Get top progress exercises
	topProgress, err := ds.GetTopProgressExercises(exerciseData)
	if err == nil && len(topProgress) > 0 {
		response.TopProgress = topProgress
	} else {
		response.TopProgress = make([]ExerciseProgress, 0)
	}

	// Get top frequency exercises
	topFrequency, err := ds.GetTopFrequencyExercises(exerciseData)
	if err == nil && len(topFrequency) > 0 {
		response.TopFrequency = topFrequency
	} else {
		response.TopFrequency = make([]ExerciseFrequency, 0)
	}

	return response, nil
}

func (ds *DashboardService) GetUserStrengthStandards(userId string) (*UserStrengthStandards, error) {
	var userObj user.User
	objectId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return nil, err
	}
	err = ds.DB.Collection("users").FindOne(context.Background(), bson.D{{Key: "_id", Value: objectId}}).Decode(&userObj)
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

	// Find matching exercises
	var exercises []exercise.Exercise
	cursor, err := ds.DB.Collection("exercises").Find(context.Background(), bson.D{
		{Key: "$or", Value: orConditions},
	})
	if err != nil {
		return nil, err
	}
	if err = cursor.All(context.Background(), &exercises); err != nil {
		return nil, err
	}

	// Create a map of exercise IDs
	exerciseIds := make([]string, len(exercises))
	exerciseMap := make(map[string]exercise.Exercise)
	for i, ex := range exercises {
		exerciseIds[i] = ex.ID.Hex()
		exerciseMap[ex.ID.Hex()] = ex
	}

	// Find latest exercise logs
	pipeline := []bson.D{
		{{Key: "$match", Value: bson.D{
			{Key: "userid", Value: userId},
			{Key: "exerciseid", Value: bson.D{{Key: "$in", Value: exerciseIds}}},
		}}},
		{{Key: "$sort", Value: bson.D{
			{Key: "datetime", Value: -1},
		}}},
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$exerciseid"},
			{Key: "latestLog", Value: bson.D{{Key: "$first", Value: "$$ROOT"}}},
		}}},
	}

	logCursor, err := ds.DB.Collection("exerciseLogs").Aggregate(context.Background(), pipeline)
	if err != nil {
		return nil, err
	}

	var logResults []struct {
		ID        string                  `bson:"_id"`
		LatestLog exerciseLog.ExerciseLog `bson:"latestLog"`
	}
	if err = logCursor.All(context.Background(), &logResults); err != nil {
		return nil, err
	}

	// Create a map of exercise ID to latest log
	latestLogsMap := make(map[string]exerciseLog.ExerciseLog)
	for _, result := range logResults {
		latestLogsMap[result.ID] = result.LatestLog
	}

	// Read strength standards
	var strengthStandards struct {
		Male   dashboardEnums.StrengthStandards `json:"Male"`
		Female dashboardEnums.StrengthStandards `json:"Female"`
	}

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

	for exerciseId, ex := range exerciseMap {
		latestLog, exists := latestLogsMap[exerciseId]
		if !exists {
			continue
		}

		// Find the closest bodyweight standards
		var standards dashboardEnums.StrengthStandards
		if userGender == authEnums.GenderMale {
			standards = strengthStandards.Male
		} else {
			standards = strengthStandards.Female
		}

		exerciseStandards, exists := standards[ex.Name]
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

		// Calculate max weight and reps
		var maxWeight float64
		var maxReps int
		for _, set := range latestLog.Sets {
			if set.Weight*float64(set.Reps) > maxWeight*float64(maxReps) {
				maxWeight = set.Weight
				maxReps = set.Reps
			}
		}

		// Calculate 1RM if more than 1 rep
		if maxReps > 1 {
			maxWeight, err = dashboardFunctions.CalculateOneRepMax(maxWeight, float64(maxReps))
			if err != nil {
				continue
			}
		}

		// Calculate relative strength
		relativeStrength := maxWeight / userBodyWeight

		// Calculate strength level and score
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
		for _, muscleGroup := range ex.TargetMuscle {
			muscleGroupStrengths[muscleGroup] = append(muscleGroupStrengths[muscleGroup], score)
		}

		// Add to exercise standards
		userStrengthStandards = append(userStrengthStandards, UserStrengthStandardPerExercise{
			Exercise:         ex.Name,
			Equipment:        ex.Equipment,
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

func (ds *DashboardService) GetRepMax(userId string, exerciseId string, useLatest bool) (*RepMaxResponse, error) {
	// Create an aggregation pipeline to calculate 1RM at database level
	pipeline := []bson.D{
		{{Key: "$match", Value: bson.D{
			{Key: "userid", Value: userId},
			{Key: "exerciseid", Value: exerciseId},
		}}},
	}

	if useLatest {
		// Add stages to get only the latest exercise log
		pipeline = append(pipeline,
			bson.D{{Key: "$sort", Value: bson.D{{Key: "datetime", Value: -1}}}},
			bson.D{{Key: "$limit", Value: 1}},
		)
	}

	// Add remaining stages
	pipeline = append(pipeline,
		bson.D{{Key: "$unwind", Value: "$sets"}},
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "sets.weight", Value: bson.D{{Key: "$gt", Value: 0}}},
			{Key: "sets.reps", Value: bson.D{{Key: "$gt", Value: 0}}},
		}}},
		bson.D{{Key: "$addFields", Value: bson.D{
			{Key: "oneRM", Value: bson.D{
				{Key: "$divide", Value: bson.A{
					"$sets.weight",
					bson.D{{Key: "$subtract", Value: bson.A{
						1.0278,
						bson.D{{Key: "$multiply", Value: bson.A{0.0278, "$sets.reps"}}},
					}}},
				}},
			}},
		}}},
		bson.D{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: nil},
			{Key: "bestOneRM", Value: bson.D{{Key: "$max", Value: "$oneRM"}}},
			{Key: "lastUpdated", Value: bson.D{{Key: "$max", Value: "$datetime"}}},
		}}},
	)

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

func (ds *DashboardService) GetTopProgressExercises(exerciseData []ExerciseData) ([]ExerciseProgress, error) {
	// Calculate progress for each exercise
	var progressList []ExerciseProgress
	for _, exercise := range exerciseData {
		logs := exercise.Logs
		if len(logs) < 2 {
			continue // Skip if we don't have at least 2 logs
		}

		firstLog := logs[0]
		lastLog := logs[len(logs)-1]

		// Skip if logs are on the same day
		if firstLog.DateTime.Equal(lastLog.DateTime) {
			continue
		}

		// Calculate Volume Progress
		startMaxSetVolume, endMaxSetVolume := calculateMaxSetVolume(firstLog.Sets), calculateMaxSetVolume(lastLog.Sets)
		if startMaxSetVolume <= 0 || endMaxSetVolume <= 0 {
			continue
		}

		volumeProgress := ((endMaxSetVolume - startMaxSetVolume) / startMaxSetVolume) * 100
		if math.IsNaN(volumeProgress) || math.IsInf(volumeProgress, 0) {
			continue
		}

		// Calculate 1RM Progress
		startOneRM, err := calculateBestOneRM(firstLog.Sets)
		if err != nil {
			continue
		}
		endOneRM, err := calculateBestOneRM(lastLog.Sets)
		if err != nil {
			continue
		}

		oneRMProgress := ((endOneRM - startOneRM) / startOneRM) * 100
		if math.IsNaN(oneRMProgress) || math.IsInf(oneRMProgress, 0) {
			continue
		}

		// Use the average of both progress metrics
		averageProgress := (volumeProgress + oneRMProgress) / 2
		if math.IsNaN(averageProgress) || math.IsInf(averageProgress, 0) {
			continue
		}

		progressList = append(progressList, ExerciseProgress{
			ExerciseID:     exercise.ID,
			Exercise:       exercise.RootExercise,
			StartVolume:    startMaxSetVolume,
			EndVolume:      endMaxSetVolume,
			VolumeProgress: math.Round(volumeProgress*100) / 100,
			StartOneRM:     startOneRM,
			EndOneRM:       endOneRM,
			OneRMProgress:  math.Round(oneRMProgress*100) / 100,
			Progress:       math.Round(averageProgress*100) / 100,
			StartDate:      firstLog.DateTime,
			EndDate:        lastLog.DateTime,
		})
	}

	// Sort by average progress in descending order
	sort.Slice(progressList, func(i, j int) bool {
		return progressList[i].Progress > progressList[j].Progress
	})

	return progressList, nil
}

func (ds *DashboardService) GetTopFrequencyExercises(exerciseData []ExerciseData) ([]ExerciseFrequency, error) {
	// Calculate frequency for each exercise
	var frequencyList []ExerciseFrequency
	for _, exercise := range exerciseData {
		frequency := ExerciseFrequency{
			ExerciseID: exercise.ID,
			Exercise:   exercise.RootExercise,
			Frequency:  float64(len(exercise.Logs)),
		}
		frequencyList = append(frequencyList, frequency)
	}

	// Sort by frequency in descending order
	sort.Slice(frequencyList, func(i, j int) bool {
		return frequencyList[i].Frequency > frequencyList[j].Frequency
	})

	return frequencyList, nil
}

func (ds *DashboardService) getExerciseLogsData(userId string, startDate, endDate time.Time) ([]ExerciseData, error) {
	pipeline := []bson.D{
		{{Key: "$match", Value: bson.D{
			{Key: "userid", Value: userId},
			{Key: "datetime", Value: bson.D{
				{Key: "$gte", Value: startDate},
				{Key: "$lte", Value: endDate},
			}},
		}}},
		{{Key: "$addFields", Value: bson.D{
			{Key: "exerciseObjectId", Value: bson.D{
				{Key: "$toObjectId", Value: "$exerciseid"},
			}},
		}}},
		{{Key: "$sort", Value: bson.D{
			{Key: "datetime", Value: 1},
		}}},
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$exerciseid"},
			{Key: "logs", Value: bson.D{{Key: "$push", Value: "$$ROOT"}}},
		}}},
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "exercises"},
			{Key: "localField", Value: "logs.exerciseObjectId"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "exercise"},
		}}},
		{{Key: "$unwind", Value: bson.D{
			{Key: "path", Value: "$exercise"},
			{Key: "preserveNullAndEmptyArrays", Value: true},
		}}},
	}

	cursor, err := ds.DB.Collection("exerciseLogs").Aggregate(context.Background(), pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var exerciseData []ExerciseData
	if err := cursor.All(context.Background(), &exerciseData); err != nil {
		return nil, err
	}

	return exerciseData, nil
}

func (s *DashboardService) GetNutritionSummary(userid string, startDate, endDate time.Time) (*NutritionSummaryResponse, error) {
	foodLogService := &foodLog.FoodLogService{DB: s.DB}

	var dailySummaries []DailyNutritionSummary
	totalCalories, totalProtein, totalCarbs, totalFat := 0.0, 0.0, 0.0, 0.0
	daysCount := 0

	// Iterate through each day in the range
	for currentDate := startDate; !currentDate.After(endDate); currentDate = currentDate.AddDate(0, 0, 1) {
		dateStr := currentDate.Format("2006-01-02")
		nutrients, err := foodLogService.CalculateDailyNutrients(dateStr, userid)
		if err != nil || len(nutrients.Nutrients) == 0 || nutrients.Calories == 0 {
			// Skip days with no data
			continue
		}

		// Extract macronutrients from nutrients array
		var protein, carbs, fat float64
		for _, nutrient := range nutrients.Nutrients {
			switch nutrient.Name {
			case "Protein":
				protein = nutrient.Amount
			case "Carbohydrate, by difference":
				carbs = nutrient.Amount
			case "Total lipid (fat)":
				fat = nutrient.Amount
			}
		}

		summary := DailyNutritionSummary{
			Date:          dateStr,
			TotalCalories: nutrients.Calories,
			TotalProtein:  protein,
			TotalCarbs:    carbs,
			TotalFat:      fat,
		}

		dailySummaries = append(dailySummaries, summary)
		totalCalories += nutrients.Calories
		totalProtein += protein
		totalCarbs += carbs
		totalFat += fat
		daysCount++
	}

	// Calculate averages
	var response NutritionSummaryResponse
	if daysCount > 0 {
		response = NutritionSummaryResponse{
			DailySummaries:  dailySummaries,
			AverageCalories: totalCalories / float64(daysCount),
			AverageProtein:  totalProtein / float64(daysCount),
			AverageCarbs:    totalCarbs / float64(daysCount),
			AverageFat:      totalFat / float64(daysCount),
		}
	} else {
		response = NutritionSummaryResponse{
			DailySummaries: []DailyNutritionSummary{},
		}
	}

	return &response, nil
}

func (ds *DashboardService) GetBodyCompositionAnalysis(userId string, startDate, endDate time.Time) (*BodyCompositionAnalysisResponse, error) {
	// Query body composition logs within date range
	filter := bson.D{
		{Key: "userid", Value: userId},
		{Key: "created_at", Value: bson.D{
			{Key: "$gte", Value: startDate},
			{Key: "$lte", Value: endDate},
		}},
	}

	// Sort by date ascending
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: 1}})
	cursor, err := ds.DB.Collection("bodyCompositionLog").Find(context.Background(), filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var logs []struct {
		CreatedAt       time.Time                                      `bson:"created_at"`
		Weight          float64                                        `bson:"weight"`
		BodyComposition userFitnessPreferenceEnums.BodyCompositionInfo `bson:"body_composition"`
	}

	if err := cursor.All(context.Background(), &logs); err != nil {
		return nil, err
	}

	if len(logs) == 0 {
		return &BodyCompositionAnalysisResponse{
			Labels:  []string{},
			Data:    []DailyBodyCompositionSummary{},
			Changes: []DailyBodyCompositionSummary{},
		}, nil
	}

	response := BodyCompositionAnalysisResponse{
		Labels:  []string{},
		Data:    []DailyBodyCompositionSummary{},
		Changes: []DailyBodyCompositionSummary{},
	}

	for i, log := range logs {
		response.Labels = append(response.Labels, log.CreatedAt.Format("2006-01-02"))
		response.Data = append(response.Data, DailyBodyCompositionSummary{
			Weight:             log.Weight,
			BMI:                log.BodyComposition.BMI,
			BodyFatMass:        log.BodyComposition.BodyFatMass,
			BodyFatPercentage:  log.BodyComposition.BodyFatPercentage,
			SkeletalMuscleMass: log.BodyComposition.SkeletalMuscleMass,
			ExtracellularWater: log.BodyComposition.ExtracellularWater,
			ECWRatio:           log.BodyComposition.ECWRatio,
		})
		if i > 0 {
			response.Changes = append(response.Changes, DailyBodyCompositionSummary{
				Weight:             log.Weight - logs[i-1].Weight,
				BMI:                log.BodyComposition.BMI - logs[i-1].BodyComposition.BMI,
				BodyFatMass:        log.BodyComposition.BodyFatMass - logs[i-1].BodyComposition.BodyFatMass,
				BodyFatPercentage:  log.BodyComposition.BodyFatPercentage - logs[i-1].BodyComposition.BodyFatPercentage,
				SkeletalMuscleMass: log.BodyComposition.SkeletalMuscleMass - logs[i-1].BodyComposition.SkeletalMuscleMass,
				ECWRatio:           log.BodyComposition.ECWRatio - logs[i-1].BodyComposition.ECWRatio,
				ExtracellularWater: log.BodyComposition.ExtracellularWater - logs[i-1].BodyComposition.ExtracellularWater,
			})
		} else {
			response.Changes = append(response.Changes, DailyBodyCompositionSummary{})
		}
	}
	return &response, nil
}
