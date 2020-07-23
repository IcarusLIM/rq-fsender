package entity

type BatchModel struct {
	Id        string `gorm:"type:varchar(36);primary_key" json:"id"`
	Api       string `json:"api"`
	TotalSize int64  `json:"total_size"`
	Report    string `json:"report"`
}
