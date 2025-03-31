package service

import (
	"sort"

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

	lineItems, err := s.fetchEligibleLineItems(placement, category, keyword)
	if err != nil {
		return nil, err
	}

	scoredItems := s.applyDynamicBidding(lineItems, placement)
	selected := s.scoreAndSelectTopAds(scoredItems, limit)

	return s.mapToAds(selected), nil
}

func (s *AdService) fetchEligibleLineItems(placement, category, keyword string) ([]*model.LineItemEntity, error) {
	return s.lineItemService.FindMatchingLineItems(placement, category, keyword)
}

func (s *AdService) applyDynamicBidding(items []*model.LineItemEntity, placement string) []*model.LineItemEntity {
	globalEventCounts, _ := s.trackingService.GetEventCounts("", "")
	placementEventCounts, _ := s.trackingService.GetEventCounts("", placement)

	strategy := utils.AvgConversionRateStrategy{}

	for _, item := range items {
		itemEventCounts, _ := s.trackingService.GetEventCounts(item.ID, "")
		itemPlacementEventCounts, _ := s.trackingService.GetEventCounts(item.ID, placement)

		dynamicBid := strategy.Calculate(
			item.Bid,
			globalEventCounts,
			placementEventCounts,
			itemEventCounts,
			itemPlacementEventCounts,
		)

		item.Bid = dynamicBid
	}

	return items
}

func (s *AdService) scoreAndSelectTopAds(items []*model.LineItemEntity, limit int) []*model.LineItemEntity {
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
