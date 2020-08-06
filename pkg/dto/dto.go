package dto

type BatchReq struct {
	Id      string   `json:"id"`
	FileIds []string `json:"file_ids" binding:"required"`
	Api     string   `json:"api" binding:"required"`
}

type HDFSConfig struct {
	NameNodes []string `json:"name_nodes" binding:"required"`
	User      string   `json:"user" binding:"required"`
	PassWord  string   `json:"password"`
}

type UploadHDFS struct {
	Path     string     `json:"path" binding:"required"`
	HDFSConf HDFSConfig `json:"hdfs" binding:"required"`
}
