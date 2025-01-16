package dashboard

import (
	"context"
	"time"

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

func (s *DashboardService) GetDashboard(userId string) (*DashboardResponse, error) {
	// Get workout sessions
	sessionFilter := bson.D{{Key: "userid", Value: userId}}
	cursor, err := s.DB.Collection("workoutSessions").Find(context.Background(), sessionFilter)
	if err != nil {
		return nil, err
	}

	var sessions []workoutSession.WorkoutSession
	if err := cursor.All(context.Background(), &sessions); err != nil {
		return nil, err
	}

	// Get exercise logs
	logFilter := bson.D{{Key: "userid", Value: userId}}
	logCursor, err := s.DB.Collection("exerciseLogs").Find(context.Background(), logFilter)
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
