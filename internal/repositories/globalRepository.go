package repositories

import (
	"github.com/m8ypie/mtg-sale-manager/internal/models"
	gorm_generics "github.com/ompluscator/gorm-generics"
	"gorm.io/gorm"
)

type GlobalRepository struct {
	EbayListing *gorm_generics.GormRepository[models.EbayListingGorm, models.EbayListing]

	// Add new repository here
}

func NewGlobalRepository(DB *gorm.DB) *GlobalRepository {
	ebaylisting := gorm_generics.NewRepository[models.EbayListingGorm, models.EbayListing](DB)
	gr := &GlobalRepository{
		EbayListing: ebaylisting,

		// Add new repository here
	}
	return gr
}
