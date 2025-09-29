package repositories

import (
	"context"

	"github.com/m8ypie/mtg-sale-manager/internal/models"
	"gorm.io/gorm"
)

type GlobalRepository struct {
	DB          *gorm.DB
	EbayListing *EbayListing
	// Add new repository here
}

type EbayListing struct {
	DB *gorm.DB
}

func CreateNewEbayListing(DB *gorm.DB) *EbayListing {
	return &EbayListing{
		DB: DB,
	}
}

func (r *EbayListing) Upsert(ctx context.Context, e *models.EbayListing) int64 {
	tx := r.DB.WithContext(ctx).Where("ebay_listing_id = ? AND ebay_offer_id = ?", e.EbayListingId, e.EbayOfferId).FirstOrCreate(e).Updates(e)
	return tx.RowsAffected
}

func (r *EbayListing) GetAll(ctx context.Context) []*models.EbayListing {
	listings := []*models.EbayListing{}
	r.DB.WithContext(ctx).Find(&listings)
	return listings
}

func NewGlobalRepository(DB *gorm.DB) *GlobalRepository {
	gr := &GlobalRepository{
		DB:          DB,
		EbayListing: CreateNewEbayListing(DB),
		// Add new repository here
	}
	return gr
}
