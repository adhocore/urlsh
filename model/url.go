package model

import "time"

// URL is model for short urls
type URL struct {
	ID        uint      `json:"-" gorm:"primaryKey"`
	ShortCode string    `json:"short_code" gorm:"size:12;uniqueIndex;not null"`
	OriginURL string    `json:"origin_url" gorm:"size:2048;index;not null"`
	Hits      uint      `json:"hits" gorm:"default:0;not null"`
	Deleted   bool      `json:"is_deleted" gorm:"default:false;not null"`
	CreatedAt time.Time `json:"-" gorm:"not null"`
	UpdatedAt time.Time `json:"-" gorm:"not null"`
	ExpiresOn time.Time `json:"expires_on"`
	Keywords  []Keyword `json:"-" gorm:"many2many:url_keywords"`
}

// IsActive checks if the url model is active
// It returns true if url is not marked deleted or expired, false otherwise.
func (urlModel URL) IsActive() bool {
	if urlModel.Deleted {
		return false
	}

	return urlModel.ExpiresOn.In(time.UTC).After(time.Now().In(time.UTC))
}
