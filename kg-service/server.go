package main

import (
	"github.com/ammorteza/clean_architecture/controller"
	router "github.com/ammorteza/clean_architecture/http"
	"github.com/ammorteza/clean_architecture/http/gin"
	job_queue "github.com/ammorteza/clean_architecture/job-queue"
	rabbit_mq "github.com/ammorteza/clean_architecture/job-queue/rabbit-mq"
	"github.com/ammorteza/clean_architecture/repository"
	"github.com/ammorteza/clean_architecture/repository/gorm"
	"github.com/ammorteza/clean_architecture/repository/redis"
	"github.com/ammorteza/clean_architecture/service"
)

var (
	jobQueue job_queue.JobQueue = rabbit_mq.New(job_queue.KEY_GENERATION_EXCHANGE)
	cache repository.Cache = redis.NewRedisPool()
	repo repository.DbRepository = gorm.New()
	appService service.AppService = service.New(repo, cache, jobQueue)
	httpRouter router.Router = gin.New()
	appController controller.AppController = controller.New(appService)
)

func main()  {
	httpRouter.POST("/unique_key", appController.GetUniqueKey)
	httpRouter.SERVE("8080")
}
