package model

import "time"

// UserAccount maps to the `users` table (migration 012).
// PasswordHash is never exposed in JSON (json:"-").
type UserAccount struct {
	Id              string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Phone           string    `json:"phone" gorm:"type:numeric(10);not null;unique"`
	Email           string    `json:"email,omitempty" gorm:"type:varchar(150)"`
	Name            string    `json:"name" gorm:"type:varchar(100);not null"`
	AvatarUrl       string    `json:"-" gorm:"type:text"`
	PasswordHash    string    `json:"-" gorm:"type:text"`
	CounterNo       string    `json:"counterNo,omitempty" gorm:"type:varchar(50)"`
	IsPhoneVerified bool      `json:"isPhoneVerified" gorm:"default:false"`
	IsEmailVerified bool      `json:"isEmailVerified" gorm:"default:false"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

func (UserAccount) TableName() string {
	return "users"
}
