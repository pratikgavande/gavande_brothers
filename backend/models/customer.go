package models

type Customer struct {
	CustomerID int32  `json:"customer_id" gorm:"column:customerId;type:int;primaryKey"`
	Name       string `json:"name" gorm:"column:name;type:varchar(250);not null"`
	Mobile     string `json:"mobile" gorm:"column:mobile;type:varchar(250);not null"`
	GSTNo      string `json:"gst_no" gorm:"column:gst_no;type:varchar(250);not null"`
	Email      string `json:"email" gorm:"column:email;type:varchar(250);not null"`
	Address    string `json:"address" gorm:"column:address;type:text;not null"`
	Status     int8   `json:"status" gorm:"column:status;type:tinyint;not null"`
}

func (Customer) TableName() string {
	return "customer"
}
