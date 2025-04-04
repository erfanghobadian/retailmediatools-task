package service

import (
	"sort"
	"time"

	"sweng-task/internal/model"
	"sweng-task/internal/utils"

	"go.uber.org/zap"
)

type AdService struct {
	log             *zap.SugaredLogger
	lineItemService *LineItemService
	trackingService *TrackingService
}

func NewAdService(
	lineItemService *LineItemService,
	trackingService *TrackingService,
	log *zap.SugaredLogger,

) *AdService {
	return &AdService{
		lineItemService: lineItemService,
		trackingService: trackingService,
		log:             log,
	}
}

func (s *AdService) GetWinningAds(placement, category, keyword string, limit int) ([]model.Ad, error) {
	s.log.Infow("Selecting winning ads", "placement", placement, "category", category, "keyword", keyword)

	lineItems, err := s.fetchMatchedLineItems(placement, category, keyword)
	if err != nil {
		return nil, err
	}

	scoredItems := s.estimateBid(lineItems, placement)
	selected := s.sortAndSelectAds(scoredItems, limit)

	return s.mapToAds(selected), nil
}

func (s *AdService) fetchMatchedLineItems(placement, category, keyword string) ([]*model.LineItemEntity, error) {
	return s.lineItemService.FindMatchingLineItems(placement, category, keyword)
}

func (s *AdService) estimateBid(items []*model.LineItemEntity, placement string) []*model.LineItemEntity {
	globalEventCounts, _ := s.trackingService.GetEventCounts("", "")
	placementEventCounts, _ := s.trackingService.GetEventCounts("", placement)

	strategy := utils.AvgConversionRateStrategy{}

	for _, item := range items {
		itemEventCounts, _ := s.trackingService.GetEventCounts(item.ID, "")
		itemPlacementEventCounts, _ := s.trackingService.GetEventCounts(item.ID, placement)

		estimatedBid := strategy.Calculate(
			item.Bid,
			globalEventCounts,
			placementEventCounts,
			itemEventCounts,
			itemPlacementEventCounts,
		)

		item.Bid = s.applyPacing(item, estimatedBid)
	}

	return items
}

func (s *AdService) sortAndSelectAds(items []*model.LineItemEntity, limit int) []*model.LineItemEntity {
	sort.SliceStable(items, func(i, j int) bool {
		return items[i].Bid > items[j].Bid
	})

	if len(items) < limit {
		limit = len(items)
	}
	return items[:limit]
}

func (s *AdService) mapToAds(items []*model.LineItemEntity) []model.Ad {
	var ads []model.Ad
	for _, li := range items {
		ads = append(ads, model.Ad{
			ID:           li.ID,
			Name:         li.Name,
			AdvertiserID: li.AdvertiserID,
			Bid:          li.Bid,
			Placement:    li.Placement,
			ServeURL:     "https://ads.cdn/" + li.ID,
		})
	}
	return ads
}

func (s *AdService) applyPacing(item *model.LineItemEntity, bid float64) float64 {
	if item.Budget == 0 {
		return bid
	}

	currentHour := float64(time.Now().Hour())
	pacingRatio := currentHour / 24.0
	expectedSpending := pacingRatio * item.Budget

	if item.DailySpending > expectedSpending {
		reduceFactor := expectedSpending / item.DailySpending
		adjustedBid := bid * reduceFactor

		s.log.Infow("Pacing adjustment applied",
			"line_item_id", item.ID,
			"original_bid", bid,
			"adjusted_bid", adjustedBid,
			"daily_spending", item.DailySpending,
			"expected_spending", expectedSpending,
		)

		return adjustedBid
	}

	return bid
}
