package transform

import (
	"errors"
	"testing"

	"github.com/opplieam/bb-transform/.jetgen/postgres/public/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGenerateDataset(t *testing.T) {
	type mockBehavior func(storer *MockCategoryStorer)
	tests := []struct {
		name         string
		mockBehavior mockBehavior
		cfg          Config
		wantErr      bool
	}{
		{
			name: "Success",
			cfg: Config{
				Version:       "v1",
				Shuffle:       false,
				TrainRatio:    60,
				ValidateRatio: 20,
				TestRatio:     20,
			},
			mockBehavior: func(storer *MockCategoryStorer) {
				storer.EXPECT().OriginalCategory().Return(Category{}, nil)
				storer.EXPECT().MatchedCategory().Return([]model.MatchCategory{}, nil)
				storer.EXPECT().CleanUp("v1").Return(nil)
				storer.EXPECT().InsertDataset(mock.AnythingOfType("[]model.CategoryDataset")).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "Success with shuffle",
			cfg: Config{
				Version:       "v1",
				Shuffle:       true,
				TrainRatio:    60,
				ValidateRatio: 20,
				TestRatio:     20,
			},
			mockBehavior: func(storer *MockCategoryStorer) {
				storer.EXPECT().OriginalCategory().Return(Category{
					1: {Name: "Cat1", Path: "/cat1"},
					2: {Name: "Cat2", Path: "/cat2"},
				}, nil)
				storer.EXPECT().MatchedCategory().Return([]model.MatchCategory{
					{MatchID: func() *int32 { i := int32(1); return &i }()},
					{MatchID: func() *int32 { i := int32(2); return &i }()},
				}, nil)
				storer.EXPECT().CleanUp("v1").Return(nil)
				storer.EXPECT().InsertDataset(mock.AnythingOfType("[]model.CategoryDataset")).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "Failed to get original category",
			cfg: Config{
				Version: "v1",
			},
			mockBehavior: func(storer *MockCategoryStorer) {
				storer.EXPECT().OriginalCategory().Return(nil, errors.New("original category error"))
			},
			wantErr: true,
		},
		{
			name: "Failed to get matched category",
			cfg: Config{
				Version: "v1",
			},
			mockBehavior: func(storer *MockCategoryStorer) {
				storer.EXPECT().OriginalCategory().Return(Category{}, nil)
				storer.EXPECT().MatchedCategory().Return(nil, errors.New("matched category error"))
			},
			wantErr: true,
		},
		{
			name: "Failed to clean up",
			cfg: Config{
				Version: "v1",
			},
			mockBehavior: func(storer *MockCategoryStorer) {
				storer.EXPECT().OriginalCategory().Return(Category{}, nil)
				storer.EXPECT().MatchedCategory().Return([]model.MatchCategory{}, nil)
				storer.EXPECT().CleanUp("v1").Return(errors.New("cleanup error"))
			},
			wantErr: true,
		},
		{
			name: "Failed to insert dataset",
			cfg: Config{
				Version: "v1",
			},
			mockBehavior: func(storer *MockCategoryStorer) {
				storer.EXPECT().OriginalCategory().Return(Category{}, nil)
				storer.EXPECT().MatchedCategory().Return([]model.MatchCategory{}, nil)
				storer.EXPECT().CleanUp("v1").Return(nil)
				storer.EXPECT().InsertDataset(mock.AnythingOfType("[]model.CategoryDataset")).Return(errors.New("insert dataset error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStorer := NewMockCategoryStorer(t)
			tt.mockBehavior(mockStorer)

			tr := NewTransform(mockStorer, tt.cfg)
			err := tr.GenerateDataset()

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
