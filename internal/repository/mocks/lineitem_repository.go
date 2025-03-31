package mocks

import (
	"errors"
	"strings"
	"sync"

	"sweng-task/internal/model"
)

type LineItemRepository struct {
	mu    sync.RWMutex
	store map[string]*model.LineItemEntity
}

func NewInMemoryLineItemRepository() *LineItemRepository {
	return &LineItemRepository{
		store: make(map[string]*model.LineItemEntity),
	}
}

func (r *LineItemRepository) Create(item *model.LineItemEntity) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.store[item.ID]; exists {
		return errors.New("line item already exists")
	}
	r.store[item.ID] = item
	return nil
}

func (r *LineItemRepository) GetByID(id string) (*model.LineItemEntity, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	item, exists := r.store[id]
	if !exists {
		return nil, errors.New("line item not found")
	}
	return item, nil
}

func (r *LineItemRepository) GetAll(advertiserID, placement string) ([]*model.LineItemEntity, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*model.LineItemEntity
	for _, item := range r.store {
		if (advertiserID == "" || item.AdvertiserID == advertiserID) &&
			(placement == "" || item.Placement == placement) {
			result = append(result, item)
		}
	}
	return result, nil
}

func (r *LineItemRepository) FindMatchingLineItems(placement, category, keyword string) ([]*model.LineItemEntity, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*model.LineItemEntity
	for _, item := range r.store {
		if item.Placement != placement || item.Status != model.LineItemStatusActive {
			continue
		}
		if category != "" && !contains(item.Categories, category) {
			continue
		}
		if keyword != "" && !contains(item.Keywords, keyword) {
			continue
		}
		result = append(result, item)
	}
	return result, nil
}

func (r *LineItemRepository) ResetDailySpending() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, item := range r.store {
		item.DailySpending = 0
	}
	return nil
}

func (r *LineItemRepository) IncreaseDailySpending(lineItemID string, amount float64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	item, exists := r.store[lineItemID]
	if !exists {
		return errors.New("line item not found")
	}
	item.DailySpending += amount
	return nil
}

func contains(slice []string, target string) bool {
	for _, v := range slice {
		if strings.EqualFold(v, target) {
			return true
		}
	}
	return false
}
