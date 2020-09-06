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
	service := service2.New(gorm.New(), rabbit_mq.New(job_queue.URL_SHORTENER_EXCHANGE), redis.NewRedisPool())
	args := os.Args[1:]
	if len(args) != 1{
		fmt.Println("this tool is used for manage database")
		fmt.Println("	db:migrate		migrating database")
		fmt.Println("	db:reset 		resetting database")
		return
	}

	switch args[0] {
	case "db:migrate":
		//fmt.Print("migrating user database...")
		//if err := service.MigrateUser(&entity.User{}); err != nil{
		//	log.Fatal(err)
		//}
		//fmt.Println(" finished")
		//fmt.Print("migrating post database...")
		//if err := service.MigratePost(&entity.Post{}); err != nil{
		//	log.Fatal(err)
		//}
		//fmt.Println(" finished")
		fmt.Print("migrating url table ...")
		if err := service.MigrateUrl(&entity.Url{}); err != nil{
			log.Fatal(err)
		}
		fmt.Println(" finished.")
	case "db:reset":
		//fmt.Print("resetting post database...")
		//if err := service.ResetPost(&entity.Post{}); err != nil{
		//	log.Fatal(err)
		//}
		//fmt.Println(" finished")
		//fmt.Print("resetting user database...")
		//if err := service.ResetUser(&entity.User{}); err != nil{
		//	log.Fatal(err)
		//}
		//fmt.Println(" finished")
		fmt.Print("resetting url table ...")
		if err := service.ResetUrl(&entity.Url{}); err != nil{
			log.Fatal(err)
		}
		fmt.Println(" finished.")
	}
}
