package controller

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Ghamster0/os-rq-fsender/app/ctx"
	"github.com/Ghamster0/os-rq-fsender/conf"
	"github.com/Ghamster0/os-rq-fsender/task"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

// Upload TODO
func Upload(c *gin.Context, app *ctx.ApplicationContext) {
	file, _ := c.FormFile("file")
	fid := "upload_" + uuid.NewV4().String()
	c.SaveUploadedFile(file, conf.FileStore+fid)
	fileInfo := map[string]interface{}{
		"fid":  fid,
		"name": file.Filename,
		"ftype": file.Header.Get("Content-Type"),
		"size": file.Size,
		"time": time.Now(),
	}
	fileJSON, _ := json.Marshal(fileInfo)
	err := app.RClient.HSet(conf.FilesKey, fid, fileJSON).Err()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"res": false, "err": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"res": true, "file_info": fileInfo})
	}
}

// ListTask TODO
func ListTask(c *gin.Context, app *ctx.ApplicationContext) {
	tasks, _ := app.Tasks.List()
	c.JSON(http.StatusBadRequest, gin.H{"res": true, "list": tasks})
	return
}

// TaskStatus TODO
func TaskStatus(c *gin.Context, app *ctx.ApplicationContext) {
	task, ok := app.Tasks.Get(c.Param("taskID"))
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"res": false, "err": "Task not exist!"})
	} else {
		c.JSON(http.StatusOK, gin.H{"res": true, "status": task.Status()})
	}
}

// CancelTask TODO
func CancelTask(c *gin.Context, app *ctx.ApplicationContext) {
	status, err := app.Tasks.Cancel(c.Param("taskID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"res": false, "err": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"res": true, "status": status})
	}
}

// SendFromLocal TODO
func SendFromLocal(c *gin.Context, app *ctx.ApplicationContext) {
	meta := &LocalTaskMeta{}
	if err := c.ShouldBindJSON(meta); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"res": false})
		return
	}
	files := task.LocalFileLoader(&(meta.Fids), app.RClient)
	err := app.Tasks.Add(&app.SC, meta.ID, meta.Receiver, files)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"res": false, "err": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"res": true})
	}
}

// SendFromHDFS TODO
func SendFromHDFS(c *gin.Context, app *ctx.ApplicationContext) {
	c.String(http.StatusOK, "Task success")
}
