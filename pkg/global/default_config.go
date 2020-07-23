package global

// EnvPrefix for env
const EnvPrefix = "FSENDER"

// DefaultConfig TODO
var DefaultConfig = map[string]interface{}{
	"debug":             false,
	"log.level":         "info",
	"http.addr":         ":6789",
	"http.log.enable":   true,
	"http.api.path":     "",
	"db.mysql.host":     "tcp(127.0.0.1:3306)",
	"db.mysql.user":     "root",
	"db.mysql.password": "root",
	"db.mysql.dbname":   "fsender",
	"upload.path":       "./upload_files/",
}
