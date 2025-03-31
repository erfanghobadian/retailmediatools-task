package utils

import (
	"sweng-task/internal/model"
)

const (
	MinImpressionThreshold   = 100
	BidFallbackMultiplier    = 0.5
	BidMinMultiplier         = 0.3
	BidHighPerformanceFactor = 2.0
	BidLowPerformanceFactor  = 0.5
)

type BidStrategy interface {
	Calculate(maxBid float64, global, placement, item, itemPlacement model.EventCounts) float64
}

type AvgConversionRateStrategy struct{}

func (s AvgConversionRateStrategy) Calculate(maxBid float64, global, placement, item, itemPlacement model.EventCounts) float64 {
	cvr := calculateRateWithFallbacks([]model.EventCounts{itemPlacement, item, placement, global}, MinImpressionThreshold, func(e model.EventCounts) int { return e.Conversions })
	avgCVR := calculateRateWithFallbacks([]model.EventCounts{placement, global}, MinImpressionThreshold, func(e model.EventCounts) int { return e.Conversions })

	if avgCVR == 0 {
		return maxBid * BidFallbackMultiplier
	}

	switch {
	case cvr >= BidHighPerformanceFactor*avgCVR:
		return maxBid
	case cvr <= BidLowPerformanceFactor*avgCVR:
		return maxBid * BidMinMultiplier
	default:
		minBid := maxBid * BidMinMultiplier
		ratio := (cvr - BidLowPerformanceFactor*avgCVR) / ((BidHighPerformanceFactor - BidLowPerformanceFactor) * avgCVR)
		return minBid + ratio*(maxBid-minBid)
	}
}

type AvgClickThroughRateStrategy struct{}

func (s AvgClickThroughRateStrategy) Calculate(maxBid float64, global, placement, item, itemPlacement model.EventCounts) float64 {
	ctr := calculateRateWithFallbacks([]model.EventCounts{itemPlacement, item, placement, global}, MinImpressionThreshold, func(e model.EventCounts) int { return e.Clicks })
	avgCTR := calculateRateWithFallbacks([]model.EventCounts{placement, global}, MinImpressionThreshold, func(e model.EventCounts) int { return e.Clicks })

	if avgCTR == 0 {
		return maxBid * BidFallbackMultiplier
	}

	switch {
	case ctr >= BidHighPerformanceFactor*avgCTR:
		return maxBid
	case ctr <= BidLowPerformanceFactor*avgCTR:
		return maxBid * BidMinMultiplier
	default:
		minBid := maxBid * BidMinMultiplier
		ratio := (ctr - BidLowPerformanceFactor*avgCTR) / ((BidHighPerformanceFactor - BidLowPerformanceFactor) * avgCTR)
		return minBid + ratio*(maxBid-minBid)
	}
}

func calculateRateWithFallbacks(fallbacks []model.EventCounts, threshold int, extract func(model.EventCounts) int) float64 {
	for _, data := range fallbacks {
		if data.Impressions >= threshold {
			return float64(extract(data)) / float64(max(1, data.Impressions))
		}
	}
	for _, data := range fallbacks {
		if data.Impressions > 0 {
			return float64(extract(data)) / float64(max(1, data.Impressions))
		}
	}
	return 0
}
