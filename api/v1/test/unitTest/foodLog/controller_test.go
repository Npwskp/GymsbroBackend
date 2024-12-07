package foodlog_test

import (
	"github.com/stretchr/testify/mock"

	foodlog "github.com/Npwskp/GymsbroBackend/api/v1/nutrition/foodLog"
)

// Mock service
type MockFoodLogService struct {
	mock.Mock
}

func (m *MockFoodLogService) CreateFoodLog(dto *foodlog.CreateFoodLogDto, userid string) (*foodlog.FoodLog, error) {
	args := m.Called(dto, userid)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*foodlog.FoodLog), args.Error(1)
}

func (m *MockFoodLogService) GetFoodLog(id string, userid string) (*foodlog.FoodLog, error) {
	args := m.Called(id, userid)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*foodlog.FoodLog), args.Error(1)
}

func (m *MockFoodLogService) GetFoodLogByUser(userid string) ([]*foodlog.FoodLog, error) {
	args := m.Called(userid)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*foodlog.FoodLog), args.Error(1)
}

func (m *MockFoodLogService) GetFoodLogByUserDate(userid string, date string) (*foodlog.FoodLog, error) {
	args := m.Called(userid, date)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*foodlog.FoodLog), args.Error(1)
}

func (m *MockFoodLogService) DeleteFoodLog(id string, userid string) error {
	args := m.Called(id, userid)
	return args.Error(0)
}

func (m *MockFoodLogService) UpdateFoodLog(doc *foodlog.UpdateFoodLogDto, id string, userid string) (*foodlog.FoodLog, error) {
	args := m.Called(doc, id, userid)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*foodlog.FoodLog), args.Error(1)
}
