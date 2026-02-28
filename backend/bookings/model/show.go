package model

type Show struct {
	Id      int     `json:"id" gorm:"primaryKey" form:"id"`
	MovieId string  `json:"movieId"`
	Date    string  `json:"date"`
	Slot    Slot    `json:"slot"`
	SlotId  int     `gorm:"foreignKey:id"`
	Cost    float64 `json:"cost"`
}

func (Show) TableName() string {
	return "show"
}
