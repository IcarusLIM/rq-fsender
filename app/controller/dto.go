package controller

// LocalTaskMeta TODO
type LocalTaskMeta struct {
	ID       string   `json:"id" binding:"required"`
	Receiver string   `json:"receiver" binding:"required"`
	Fids     []string `json:"fids" binding:"required"`
}

// HDFSTaskMeta TODO
type HDFSTaskMeta struct {
	Receiver string   `json:"receiver" binding:"required"`
	NameNode string   `json:"name_node" binding:"required"`
	Dirs     []string `json:"dirs,omitempty"`
	Paths    []string `json:"paths,omitempty"`
}
