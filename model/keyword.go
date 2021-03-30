package model

// Keyword is model for keywords
type Keyword struct {
	ID      uint   `json:"-" gorm:"primaryKey"`
	Keyword string `json:"keyword" gorm:"size:25;unique;not null"`
}
