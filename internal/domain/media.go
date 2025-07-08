package domain

type Media struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Link string `gorm:"size:500,unique" json:"link"`
	Sent bool   `gorm:"default:false"`
}

type MediaRepository interface {
	Create(media Media) error
	Delete(media Media) error
	GetByLink(link string) (Media, error)
	Random() (Media, error)
	ClearStatus() error
}
