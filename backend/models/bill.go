package models

import "time"

type Bill struct {
	BillID         int32     `json:"bill_id" gorm:"column:billId;type:int;primaryKey"`
	BillNo         int64     `json:"bill_no" gorm:"column:billNo;type:bigint;not null"`
	CustomerID     int32     `json:"customer_id" gorm:"column:customerId;type:int;not null"`
	Date           time.Time `json:"date" gorm:"column:date;type:datetime;not null"`
	GSTStatus      int8      `json:"gst_status" gorm:"column:gstStatus;type:tinyint;not null"`
	GSTAmount      int32     `json:"gst_amount" gorm:"column:gstAmount;type:int;not null"`
	Amount         int32     `json:"amount" gorm:"column:amount;type:int;not null"`
	AmountReceived int32     `json:"amount_received" gorm:"column:amountReceived;type:int;not null"`
	AmountPending  int32     `json:"amount_pending" gorm:"column:amountPending;type:int;not null"`
	EngineersName  string    `json:"engineers_name" gorm:"column:engineersName;type:text;not null"`
	SiteName       string    `json:"site_name" gorm:"column:siteName;type:text;not null"`
	Location       string    `json:"location" gorm:"column:location;type:text;not null"`
	ItemCount      int32     `json:"item_count" gorm:"column:itemCount;type:int;not null"`
	Status         int8      `json:"status" gorm:"column:status;type:tinyint;not null"`
}

func (Bill) TableName() string {
	return "bills"
}
