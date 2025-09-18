package models

import (
	"time"

	"gorm.io/gorm"
)

type Model struct {
	ID        uint
	createdAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}

type EbayListingGorm struct {
	gorm.Model
	EbayListingID      uint   `gorm:"column:ebay_listing_id"`
	EbayOfferId        string `gorm:"column:ebay_offer_id"`
	Sku                string `gorm:"column:sku"`
	ScryFallId         string `gorm:"column:scryfall_id"`
	EchoMtgInventoryId string `gorm:"column:echo_mtg_inventory_id"`
}

type EbayListing struct {
	Model
	EbayListingID      uint
	EbayOfferId        string
	Sku                string
	ScryFallId         string
	EchoMtgInventoryId string
}

func (m EbayListingGorm) ToEntity() EbayListing {
	return EbayListing{
		Model: Model{
			ID:        m.ID,
			createdAt: m.CreatedAt,
			UpdatedAt: m.UpdatedAt,
			DeletedAt: m.DeletedAt,
		},
		EbayListingID:      m.EbayListingID,
		EbayOfferId:        m.EbayOfferId,
		Sku:                m.Sku,
		ScryFallId:         m.ScryFallId,
		EchoMtgInventoryId: m.EchoMtgInventoryId,
	}
}

func (m EbayListingGorm) FromEntity(entity EbayListing) interface{} {
	return EbayListingGorm{
		Model: gorm.Model{
			ID:        entity.ID,
			CreatedAt: entity.createdAt,
			UpdatedAt: entity.UpdatedAt,
			DeletedAt: entity.DeletedAt,
		},
		EbayListingID:      entity.EbayListingID,
		EbayOfferId:        entity.EbayOfferId,
		Sku:                entity.Sku,
		ScryFallId:         entity.ScryFallId,
		EchoMtgInventoryId: entity.EchoMtgInventoryId,
	}
}
