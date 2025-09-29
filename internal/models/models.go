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
	EbayListingId      string `gorm:"column:ebay_listing_id;uniqueIndex"`
	EbayOfferId        string `gorm:"column:ebay_offer_id;uniqueIndex"`
	Sku                string `gorm:"column:sku"`
	ScryFallId         string `gorm:"column:scryfall_id"`
	EchoMtgInventoryId string `gorm:"column:echo_mtg_inventory_id"`
}

type EbayListing struct {
	Model
	EbayTitle          string
	EbayListingId      string
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
		EbayListingId:      m.EbayListingId,
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
		EbayListingId: entity.EbayListingId,
		EbayOfferId:   entity.EbayOfferId,
		Sku:           entity.Sku,
		ScryFallId:    entity.ScryFallId,
	}
}

type PriceInfo struct {
	Min     float64
	Max     float64
	Average float64
}
