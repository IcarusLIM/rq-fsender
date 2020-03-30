package conf

// const
const (
	FileStore      string = "/search/odin/data/fsender/"
	SendThreadPool int    = 2

	// Redis Config
	RedisHost string = "localhost:6379"
	RedisPass string = "root"
	RedisDB   int    = 0

	FilesKey string = "os:rq:fsender:files"
)
