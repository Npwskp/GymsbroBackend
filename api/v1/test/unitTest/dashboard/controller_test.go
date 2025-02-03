package dashboard_test

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/Npwskp/GymsbroBackend/api/v1/dashboard"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock service
type MockDashboardService struct {
	mock.Mock
}

func (m *MockDashboardService) GetDashboard(userId string, startDate, endDate time.Time) (*dashboard.DashboardResponse, error) {
	args := m.Called(userId, startDate, endDate)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dashboard.DashboardResponse), args.Error(1)
}

func (m *MockDashboardService) GetUserStrengthStandards(userId string) (*dashboard.UserStrengthStandards, error) {
	args := m.Called(userId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dashboard.UserStrengthStandards), args.Error(1)
}

func (m *MockDashboardService) GetRepMax(userId string, exerciseId string, useLatest bool) (*dashboard.RepMaxResponse, error) {
	args := m.Called(userId, exerciseId, useLatest)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dashboard.RepMaxResponse), args.Error(1)
}

func (m *MockDashboardService) GetNutritionSummary(userId string, startDate, endDate time.Time) (*dashboard.NutritionSummaryResponse, error) {
	args := m.Called(userId, startDate, endDate)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dashboard.NutritionSummaryResponse), args.Error(1)
}

func (m *MockDashboardService) GetBodyCompositionAnalysis(userId string, startDate, endDate time.Time) (*dashboard.BodyCompositionAnalysisResponse, error) {
	args := m.Called(userId, startDate, endDate)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dashboard.BodyCompositionAnalysisResponse), args.Error(1)
}

// TestMiddleware sets up the test context with a mock user
func testMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		claims := jwt.MapClaims{
			"sub": c.Get("userid", ""),
		}
		token := &jwt.Token{
			Claims: claims,
		}
		c.Locals("user", token)
		return c.Next()
	}
}

// Test setup helper
func setupTest() (*fiber.App, *MockDashboardService) {
	app := fiber.New()
	api := app.Group("/api/v1", testMiddleware())

	mockService := new(MockDashboardService)
	controller := &dashboard.DashboardController{
		Instance: api,
		Service:  mockService,
	}
	controller.Handle()
	return app, mockService
}

func TestGetDashboardHandler(t *testing.T) {
	app, mockService := setupTest()

	t.Run("Successfully get dashboard", func(t *testing.T) {
		startDate := time.Now().AddDate(0, 0, -7).UTC().Truncate(time.Second)
		endDate := time.Now().UTC().Truncate(time.Second)

		expectedResponse := &dashboard.DashboardResponse{
			FrequencyGraph: dashboard.FrequencyGraphData{
				Labels:    []string{"2024-01-01", "2024-01-02"},
				Values:    []int{2, 3},
				TrendLine: []float64{2.5, 2.5},
			},
			Analysis: dashboard.WorkoutAnalysis{
				TotalWorkouts:          5,
				TotalExercises:         15,
				TotalVolume:            5000,
				AverageWorkoutDuration: 3600,
			},
		}

		mockService.On("GetDashboard", "test_user",
			mock.MatchedBy(func(t time.Time) bool { return t.UTC().Truncate(time.Second).Equal(startDate) }),
			mock.MatchedBy(func(t time.Time) bool { return t.UTC().Truncate(time.Second).Equal(endDate) }),
		).Return(expectedResponse, nil)

		url := fmt.Sprintf("/api/v1/dashboard?startDate=%s&endDate=%s",
			url.QueryEscape(startDate.Format("2006-01-02 15:04:05")),
			url.QueryEscape(endDate.Format("2006-01-02 15:04:05")))

		req := httptest.NewRequest("GET", url, nil)
		req.Header.Set("userid", "test_user")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var result dashboard.DashboardResponse
		json.NewDecoder(resp.Body).Decode(&result)
		assert.Equal(t, expectedResponse.Analysis.TotalWorkouts, result.Analysis.TotalWorkouts)
		assert.Equal(t, expectedResponse.Analysis.TotalExercises, result.Analysis.TotalExercises)
	})

	t.Run("Invalid date format", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/dashboard?startDate=invalid&endDate=invalid", nil)
		req.Header.Set("userid", "test_user")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	})
}

func TestGetUserStrengthStandardsHandler(t *testing.T) {
	app, mockService := setupTest()

	t.Run("Successfully get strength standards", func(t *testing.T) {
		expectedResponse := &dashboard.UserStrengthStandards{
			ExerciseStandards: []dashboard.UserStrengthStandardPerExercise{
				{
					Exercise:         "Bench Press",
					RepMax:           100,
					RelativeStrength: 1.25,
					Score:            75,
				},
			},
			MuscleGroupStrengths: []dashboard.UserStrengthStandardPerMuscleGroup{
				{
					Score: 80,
				},
			},
		}

		mockService.On("GetUserStrengthStandards", "test_user").Return(expectedResponse, nil)

		req := httptest.NewRequest("GET", "/api/v1/dashboard/strength-standards", nil)
		req.Header.Set("userid", "test_user")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var result dashboard.UserStrengthStandards
		json.NewDecoder(resp.Body).Decode(&result)
		assert.Equal(t, len(expectedResponse.ExerciseStandards), len(result.ExerciseStandards))
		assert.Equal(t, len(expectedResponse.MuscleGroupStrengths), len(result.MuscleGroupStrengths))
	})
}

func TestGetRepMaxHandler(t *testing.T) {
	app, mockService := setupTest()

	t.Run("Successfully get rep max", func(t *testing.T) {
		exerciseId := "test_exercise"
		expectedResponse := &dashboard.RepMaxResponse{
			OneRepMax:    100,
			EightRepMax:  80,
			TwelveRepMax: 70,
			LastUpdated:  time.Now(),
		}

		mockService.On("GetRepMax", "test_user", exerciseId, false).Return(expectedResponse, nil)

		req := httptest.NewRequest("GET", "/api/v1/dashboard/rep-max/"+exerciseId, nil)
		req.Header.Set("userid", "test_user")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var result dashboard.RepMaxResponse
		json.NewDecoder(resp.Body).Decode(&result)
		assert.Equal(t, expectedResponse.OneRepMax, result.OneRepMax)
		assert.Equal(t, expectedResponse.EightRepMax, result.EightRepMax)
		assert.Equal(t, expectedResponse.TwelveRepMax, result.TwelveRepMax)
	})
}

func TestGetNutritionSummaryHandler(t *testing.T) {
	app, mockService := setupTest()

	t.Run("Successfully get nutrition summary", func(t *testing.T) {
		startDate := time.Now().AddDate(0, 0, -7).UTC().Truncate(time.Second)
		endDate := time.Now().UTC().Truncate(time.Second)

		expectedResponse := &dashboard.NutritionSummaryResponse{
			DailySummaries: []dashboard.DailyNutritionSummary{
				{
					Date:          "2024-01-01",
					TotalCalories: 2000,
					TotalProtein:  150,
					TotalCarbs:    200,
					TotalFat:      70,
				},
			},
			AverageCalories: 2000,
			AverageProtein:  150,
			AverageCarbs:    200,
			AverageFat:      70,
		}

		mockService.On("GetNutritionSummary", "test_user",
			mock.MatchedBy(func(t time.Time) bool { return t.UTC().Truncate(time.Second).Equal(startDate) }),
			mock.MatchedBy(func(t time.Time) bool { return t.UTC().Truncate(time.Second).Equal(endDate) }),
		).Return(expectedResponse, nil)

		url := fmt.Sprintf("/api/v1/dashboard/nutrition-summary?startDate=%s&endDate=%s",
			url.QueryEscape(startDate.Format("2006-01-02 15:04:05")),
			url.QueryEscape(endDate.Format("2006-01-02 15:04:05")))

		req := httptest.NewRequest("GET", url, nil)
		req.Header.Set("userid", "test_user")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var result dashboard.NutritionSummaryResponse
		json.NewDecoder(resp.Body).Decode(&result)
		assert.Equal(t, len(expectedResponse.DailySummaries), len(result.DailySummaries))
		assert.Equal(t, expectedResponse.AverageCalories, result.AverageCalories)
	})

	t.Run("Invalid date format", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/dashboard/nutrition-summary?startDate=invalid&endDate=invalid", nil)
		req.Header.Set("userid", "test_user")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	})
}

func TestGetBodyCompositionAnalysisHandler(t *testing.T) {
	app, mockService := setupTest()

	t.Run("Successfully get body composition analysis", func(t *testing.T) {
		startDate := time.Now().AddDate(0, 0, -7).UTC().Truncate(time.Second)
		endDate := time.Now().UTC().Truncate(time.Second)

		expectedResponse := &dashboard.BodyCompositionAnalysisResponse{
			Labels: []string{"2024-01-01", "2024-01-02"},
			Data: []dashboard.DailyBodyCompositionSummary{
				{
					Weight:             70,
					BMI:                22.5,
					BodyFatMass:        12,
					BodyFatPercentage:  17,
					SkeletalMuscleMass: 35,
					ExtracellularWater: 20,
					ECWRatio:           0.38,
				},
			},
			Changes: []dashboard.DailyBodyCompositionSummary{
				{
					Weight:             -1,
					BMI:                -0.3,
					BodyFatMass:        -0.5,
					BodyFatPercentage:  -0.5,
					SkeletalMuscleMass: 0.2,
					ExtracellularWater: 0,
					ECWRatio:           0,
				},
			},
		}

		mockService.On("GetBodyCompositionAnalysis", "test_user",
			mock.MatchedBy(func(t time.Time) bool { return t.UTC().Truncate(time.Second).Equal(startDate) }),
			mock.MatchedBy(func(t time.Time) bool { return t.UTC().Truncate(time.Second).Equal(endDate) }),
		).Return(expectedResponse, nil)

		url := fmt.Sprintf("/api/v1/dashboard/body-composition?startDate=%s&endDate=%s",
			url.QueryEscape(startDate.Format("2006-01-02 15:04:05")),
			url.QueryEscape(endDate.Format("2006-01-02 15:04:05")))

		req := httptest.NewRequest("GET", url, nil)
		req.Header.Set("userid", "test_user")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var result dashboard.BodyCompositionAnalysisResponse
		json.NewDecoder(resp.Body).Decode(&result)
		assert.Equal(t, len(expectedResponse.Labels), len(result.Labels))
		assert.Equal(t, len(expectedResponse.Data), len(result.Data))
		assert.Equal(t, len(expectedResponse.Changes), len(result.Changes))
	})

	t.Run("Invalid date format", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/dashboard/body-composition?startDate=invalid&endDate=invalid", nil)
		req.Header.Set("userid", "test_user")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	})
}
