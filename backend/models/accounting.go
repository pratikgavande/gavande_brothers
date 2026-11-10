package models

import "time"

type Accounting struct {
	AccID      int32     `json:"acc_id" gorm:"column:accId;type:int;primaryKey"`
	ShortDesc  string    `json:"short_desc" gorm:"column:shortDesc;type:text;not null"`
	Amount     int64     `json:"amount" gorm:"column:amount;type:bigint;not null"`
	DetailDesc string    `json:"detail_desc" gorm:"column:detailDesc;type:text;not null"`
	Date       time.Time `json:"date" gorm:"column:date;type:date;not null"`
	Status     int8      `json:"status" gorm:"column:status;type:tinyint;not null"`
}

func (Accounting) TableName() string {
	return "accounting"
}
