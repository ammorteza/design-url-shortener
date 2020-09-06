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
	cache repository.Cache = redis.NewRedisPool()
	queue job_queue.JobQueue = rabbit_mq.New(job_queue.URL_SHORTENER_EXCHANGE)
	repo repository.DbRepository = gorm.New()
	appService service.AppService = service.New(repo, queue, cache)
	httpRouter router.Router = gin.New()
	appController controller.AppController = controller.New(appService)
)

func initRoutes(routes router.Router) router.Router{
	routes.POST("/get_url_randomly", appController.GetUrlRandomly) // this api will used for run unit tests
	routes.POST("/create_url", appController.CreateUrl)
	routes.GET("/:key", appController.RedirectUrl)
	return routes
}

func main()  {
	httpRouter = initRoutes(httpRouter)
	httpRouter.SERVE("8080")
}
