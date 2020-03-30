package ctx

import (
	"github.com/Ghamster0/os-rq-fsender/conf"
	"github.com/Ghamster0/os-rq-fsender/sender"
	"github.com/Ghamster0/os-rq-fsender/task"
	"github.com/go-redis/redis/v7"
)

// ApplicationContext TODO
type ApplicationContext struct {
	SC      chan *task.Feed
	Tasks   *task.Container
	RClient *redis.Client
}

// StartApplication TODO
func StartApplication() (ctx *ApplicationContext, err error) {
	threads := conf.SendThreadPool
	sc := make(chan *task.Feed, 2*threads)
	// start sender runner
	runner := &sender.Runner{SC: &sc}
	runner.Start(threads)
	// init Redis
	client := redis.NewClient(&redis.Options{
		Addr:     conf.RedisHost,
		Password: conf.RedisPass,
		DB:       conf.RedisDB,
	})
	_, err = client.Ping().Result()
	if err != nil {
		panic(err.Error())
	}
	ctx = &ApplicationContext{
		SC:      sc,
		Tasks:   task.NewContainer(),
		RClient: client,
	}
	return
}
