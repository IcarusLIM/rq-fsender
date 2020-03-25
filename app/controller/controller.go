package controller

import (
	"net/http"

	"github.com/Ghamster0/os-rq-fsender/app/ctx"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

// Upload TODO
func Upload(c *gin.Context, app *ctx.ApplicationContext) {
	file, _ := c.FormFile("file")
	fid := "upload_" + uuid.NewV4().String()
	c.SaveUploadedFile(file, "/search/odin/data/fsender/"+fid)
	fileInfo := map[string]interface{}{
		"fid":  fid,
		"name": file.Filename,
		"type": file.Header.Get("Content-Type"),
		"size": file.Size,
	}
	c.JSON(http.StatusOK, gin.H{
		"res":       true,
		"file_info": fileInfo,
	})
}

// SendFromLocal TODO
func SendFromLocal(c *gin.Context, app *ctx.ApplicationContext) {
	c.String(http.StatusOK, "Task success")
}

// SendFromHDFS TODO
func SendFromHDFS(c *gin.Context, app *ctx.ApplicationContext) {

}
