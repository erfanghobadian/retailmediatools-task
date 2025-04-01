package utils

import (
	"sweng-task/internal/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAvgConversionRateStrategy_Calculate(t *testing.T) {
	strategy := AvgConversionRateStrategy{}
	maxBid := 2.0

	tests := []struct {
		name          string
		global        model.EventCounts
		placement     model.EventCounts
		item          model.EventCounts
		itemPlacement model.EventCounts
		expectedMin   float64
		expectedMax   float64
	}{
		{
			name:          "High CVR -> max bid",
			global:        model.EventCounts{Impressions: 200, Conversions: 10},
			placement:     model.EventCounts{Impressions: 200, Conversions: 10},
			item:          model.EventCounts{Impressions: 200, Conversions: 30},
			itemPlacement: model.EventCounts{Impressions: 200, Conversions: 30},
			expectedMin:   1.9,
			expectedMax:   2.0,
		},
		{
			name:          "Low CVR -> min bid",
			global:        model.EventCounts{Impressions: 200, Conversions: 20},
			placement:     model.EventCounts{Impressions: 200, Conversions: 20},
			item:          model.EventCounts{Impressions: 200, Conversions: 1},
			itemPlacement: model.EventCounts{Impressions: 200, Conversions: 1},
			expectedMin:   0.59,
			expectedMax:   0.61,
		},
		{
			name:          "Medium CVR -> mid bid",
			global:        model.EventCounts{Impressions: 200, Conversions: 20},
			placement:     model.EventCounts{Impressions: 200, Conversions: 20},
			item:          model.EventCounts{Impressions: 200, Conversions: 12},
			itemPlacement: model.EventCounts{Impressions: 200, Conversions: 12},
			expectedMin:   0.6,
			expectedMax:   2.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bid := strategy.Calculate(
				maxBid,
				tt.global,
				tt.placement,
				tt.item,
				tt.itemPlacement,
			)
			assert.GreaterOrEqual(t, bid, tt.expectedMin)
			assert.LessOrEqual(t, bid, tt.expectedMax)
		})
	}
}
