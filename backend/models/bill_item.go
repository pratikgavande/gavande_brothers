package models

type BillItem struct {
	ItemID int32   `json:"item_id" gorm:"column:itemId;type:int;primaryKey"`
	BillID int32   `json:"bill_id" gorm:"column:billId;type:int;not null"`
	Name   string  `json:"name" gorm:"column:name;type:varchar(250);not null"`
	Rate   string  `json:"rate" gorm:"column:rate;type:text;not null"`
	Unit   string  `json:"unit" gorm:"column:unit;type:text;not null"`
	Amount int32   `json:"amount" gorm:"column:amount;type:int;not null"`
	Date   *string `json:"date" gorm:"column:date;type:text"`
	Status int8    `json:"status" gorm:"column:status;type:tinyint;not null"`
}

func (BillItem) TableName() string {
	return "bill_items"
}
