package entity

type UniqueKey struct {
	ID 				uint			`sql:"primary_key"`
	Key 			string 			`gorm:"column:key;unique;type:varchar(8);NOT NULL"json:"key"`
	State 			bool 			`gorm:"column:state;default:0;NOT NULL"`
}