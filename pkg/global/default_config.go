package global

// EnvPrefix for env
const EnvPrefix = "FSENDER"

// DefaultConfig TODO
var DefaultConfig = map[string]interface{}{
	"log.level":         "info",
	"http.port":         "8080",
	"db.mysql.host":     "tcp(127.0.0.1:3306)",
	"db.mysql.user":     "root",
	"db.mysql.password": "root",
	"db.mysql.dbname":   "fsender",
	"upload.file":       "./upload_files/",
	"task.concurrent":   10,
}
