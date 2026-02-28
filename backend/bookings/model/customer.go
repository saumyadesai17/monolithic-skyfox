package model

type Customer struct {
	Id          int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Oid         string `json:"-" gorm:"type:varchar(50);not null;default:'system'"`
	Name        string `json:"name" binding:"required,min=2,max=15"`
	PhoneNumber string `json:"phoneNumber" binding:"phoneNumber"`
}

func (Customer) TableName() string {
	return "customer"
}
