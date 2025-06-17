package models

type Company struct {
	ID   string `json:"id" gorm:"primaryKey"`
	Name string `json:"name"`
}
