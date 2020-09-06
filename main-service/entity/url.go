package entity

import "time"

type Url struct {
	ID 					uint				`sql:"primary_key"`
	User 				string 				`gorm:"column:user;type:varchar(50);NOT NULL"json:"user"`
	UniqueKey 			string				`gorm:"column:unique_key;unique;type:varchar(6);NOT NULL"json:"unique_key"`
	OriginalUrl			string				`gorm:"column:original_url;type:text;NOT NULL"json:"original_url"`
	VisitCount			uint64				`gorm:"column:visit_count;NOT NULL;default:0"json:"visit_count"`

	CreatedAt 			time.Time
	ExpiredAt			time.Time
}