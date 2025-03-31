package model

func ToEntityLineItem(dto LineItem) LineItemEntity {
	return LineItemEntity{
		ID:           dto.ID,
		Name:         dto.Name,
		AdvertiserID: dto.AdvertiserID,
		Bid:          dto.Bid,
		Budget:       dto.Budget,
		Placement:    dto.Placement,
		Categories:   dto.Categories,
		Keywords:     dto.Keywords,
		Status:       dto.Status,
		CreatedAt:    dto.CreatedAt,
		UpdatedAt:    dto.UpdatedAt,
	}
}

func ToDTOLineItem(e LineItemEntity) LineItem {
	return LineItem{
		ID:           e.ID,
		Name:         e.Name,
		AdvertiserID: e.AdvertiserID,
		Bid:          e.Bid,
		Budget:       e.Budget,
		Placement:    e.Placement,
		Categories:   e.Categories,
		Keywords:     e.Keywords,
		Status:       e.Status,
		CreatedAt:    e.CreatedAt,
		UpdatedAt:    e.UpdatedAt,
	}
}

func ToLineItemEntityFromCreate(dto LineItemCreate) LineItemEntity {
	return LineItemEntity{
		Name:         dto.Name,
		AdvertiserID: dto.AdvertiserID,
		Bid:          dto.Bid,
		Budget:       dto.Budget,
		Placement:    dto.Placement,
		Categories:   dto.Categories,
		Keywords:     dto.Keywords,
		Status:       LineItemStatusActive,
	}
}

func ToDTOLineItemList(entities []*LineItemEntity) []*LineItem {
	var result []*LineItem
	for _, e := range entities {
		item := ToDTOLineItem(*e)
		result = append(result, &item)
	}
	return result
}

func ToAd(e LineItemEntity) Ad {
	return Ad{
		ID:           e.ID,
		Name:         e.Name,
		AdvertiserID: e.AdvertiserID,
		Bid:          e.Bid,
		Placement:    e.Placement,
		ServeURL:     "a",
	}
}

func ToEntityTrackingEvent(dto TrackingEvent) TrackingEventEntity {
	return TrackingEventEntity{
		EventType:  dto.EventType,
		LineItemID: dto.LineItemID,
		Timestamp:  dto.Timestamp,
		Placement:  dto.Placement,
		UserID:     dto.UserID,
		Metadata:   dto.Metadata,
	}
}

func ToDTOTrackingEvent(e TrackingEventEntity) TrackingEvent {
	return TrackingEvent{
		EventType:  e.EventType,
		LineItemID: e.LineItemID,
		Timestamp:  e.Timestamp,
		Placement:  e.Placement,
		UserID:     e.UserID,
		Metadata:   e.Metadata,
	}
}
