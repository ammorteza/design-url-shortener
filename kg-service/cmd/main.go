package main

import (
	"fmt"
	"github.com/ammorteza/clean_architecture/entity"
	job_queue "github.com/ammorteza/clean_architecture/job-queue"
	rabbit_mq "github.com/ammorteza/clean_architecture/job-queue/rabbit-mq"
	"github.com/ammorteza/clean_architecture/repository/gorm"
	"github.com/ammorteza/clean_architecture/repository/redis"
	service2 "github.com/ammorteza/clean_architecture/service"
	"log"
	"os"
)

func main(){
	service := service2.New(gorm.New(), redis.NewRedisPool(), rabbit_mq.New(job_queue.KEY_GENERATION_EXCHANGE))
	args := os.Args[1:]
	if len(args) != 1{
		fmt.Println("this tool is used for manage database")
		fmt.Println("	db:migrate		migrating database")
		fmt.Println("	db:reset 		resetting database")
		return
	}

	switch args[0] {
	case "db:migrate":
		fmt.Print("migrating unique key table ...")
		if err := service.MigrateKey(&entity.UniqueKey{}); err != nil{
			log.Fatal(err)
		}
		fmt.Println(" finished.")
	case "db:reset":
		fmt.Print("resetting unique key table ...")
		if err := service.ResetKey(&entity.UniqueKey{}); err != nil{
			log.Fatal(err)
		}
		fmt.Println(" finished.")
	case "db:seed":
		fmt.Print("seeding unique key table ...")
		if err := service.SeedKey(); err != nil{
			log.Fatal(err)
		}
		fmt.Println(" finished.")
	}
}
