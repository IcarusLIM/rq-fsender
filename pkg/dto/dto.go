package dto

type BatchReq struct {
	Id      string   `json:"id"`
	FileIds []string `json:"file_ids" binding:"required"`
	Api     string   `json:"api" binding:"required"`
}
