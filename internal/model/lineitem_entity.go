package model

import (
	"time"

	"github.com/lib/pq"
)

type LineItemEntity struct {
	ID            string         `gorm:"primaryKey"`
	Name          string         `gorm:"not null"`
	AdvertiserID  string         `gorm:"not null;index:idx_advertiser_id"`
	Bid           float64        `gorm:"not null;check:bid >= 0"`
	Budget        float64        `gorm:"not null;check:budget >= 0"`
	DailySpending float64        `gorm:"not null;default:0;check:daily_spending >= 0"`
	Placement     string         `gorm:"not null;index:idx_placement"`
	Categories    pq.StringArray `gorm:"type:text[]"`
	Keywords      pq.StringArray `gorm:"type:text[]"`
	Status        LineItemStatus `gorm:"type:text;not null;index:idx_status"`
	CreatedAt     time.Time      `gorm:"index:idx_created_at"`
	UpdatedAt     time.Time
}

func (LineItemEntity) TableName() string {
	return "line_items"
}

type TrackingEventEntity struct {
	ID         uint64            `gorm:"primaryKey"`
	EventType  TrackingEventType `gorm:"type:text;index:idx_event_type"`
	LineItemID string            `gorm:"not null;index:idx_line_item_id"`
	LineItem   LineItemEntity    `gorm:"foreignKey:LineItemID;references:ID;constraint:OnDelete:CASCADE"`
	Timestamp  time.Time         `gorm:"index:idx_timestamp"`
	Placement  string            `gorm:"index:idx_placement"`
	UserID     string
	Metadata   map[string]string `gorm:"type:jsonb"`
}

func (TrackingEventEntity) TableName() string {
	return "tracking_events"
}
