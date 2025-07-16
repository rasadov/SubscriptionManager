package models

import (
	"time"
)

type Subscription struct {
	ID          uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	ServiceName string    `json:"service_name" gorm:"type:varchar(255);not null;index"`
	Price       int       `json:"price" gorm:"not null"`
	UserID      string    `json:"user_id" gorm:"type:uuid;not null;index"`
	StartDate   string    `json:"start_date" gorm:"type:varchar(7);not null;index"`
	EndDate     *string   `json:"end_date,omitempty" gorm:"type:varchar(7);default:null"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
