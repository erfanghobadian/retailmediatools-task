package service

import (
	"sweng-task/internal/model"
	"sweng-task/internal/repository"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// LineItemService provides operations for line items
type LineItemService struct {
	repo repository.LineItemRepository
	log  *zap.SugaredLogger
}

// NewLineItemService creates a new LineItemService
func NewLineItemService(repo repository.LineItemRepository, log *zap.SugaredLogger) *LineItemService {
	return &LineItemService{
		repo: repo,
		log:  log,
	}
}

// Create creates a new line item
func (s *LineItemService) Create(input model.LineItemCreate) (*model.LineItem, error) {
	now := time.Now()

	// Map to entity and populate defaults
	lineItem := model.ToLineItemEntityFromCreate(input)
	lineItem.ID = "li_" + uuid.New().String()
	lineItem.CreatedAt = now
	lineItem.UpdatedAt = now

	// Save to DB
	if err := s.repo.Create(&lineItem); err != nil {
		return nil, err
	}

	s.log.Infow("Line item created",
		"id", lineItem.ID,
		"name", lineItem.Name,
		"advertiser_id", lineItem.AdvertiserID,
		"placement", lineItem.Placement,
	)

	// Convert to API response DTO
	dto := model.ToDTOLineItem(lineItem)
	return &dto, nil
}

// GetByID retrieves a line item by ID
func (s *LineItemService) GetByID(id string) (*model.LineItem, error) {
	item, err := s.repo.GetByID(id)
	if err != nil {
		return nil, ErrLineItemNotFound

	}
	dto := model.ToDTOLineItem(*item)
	return &dto, nil
}

// GetAll retrieves all line items, optionally filtered by advertiser ID and placement
func (s *LineItemService) GetAll(advertiserID, placement string) ([]*model.LineItem, error) {
	entityItems, err := s.repo.GetAll(advertiserID, placement)
	if err != nil {
		return nil, err
	}

	var dtoItems []*model.LineItem
	for _, entityItem := range entityItems {
		dto := model.ToDTOLineItem(*entityItem)
		dtoItems = append(dtoItems, &dto)
	}

	return dtoItems, nil
}

// FindMatchingLineItems finds line items matching the given placement and filters
// This method will be used by the AdService when implementing the ad selection logic
func (s *LineItemService) FindMatchingLineItems(placement string, category, keyword string) ([]*model.LineItemEntity, error) {
	entityItems, err := s.repo.FindMatchingLineItems(placement, category, keyword)
	if err != nil {
		return nil, err
	}

	return entityItems, nil
}

func (s *LineItemService) ResetDailySpending() error {
	err := s.repo.ResetDailySpending()
	if err != nil {
		s.log.Errorw("ResetDailySpending failed", "error", err)
		return err
	}

	s.log.Info("✅ Daily budgets reset successfully")
	return nil
}

func (s *LineItemService) IncreaseDailySpending(lineItemID string, costPerEvent float64) error {
	if err := s.repo.IncreaseDailySpending(lineItemID, costPerEvent); err != nil {
		s.log.Errorw("Failed to increase daily spending", "line_item_id", lineItemID, "error", err)
		return err
	}
	return nil
}
