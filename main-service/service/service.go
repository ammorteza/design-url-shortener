package service

import (
	"encoding/json"
	"github.com/ammorteza/clean_architecture/entity"
	job_queue "github.com/ammorteza/clean_architecture/job-queue"
	"github.com/ammorteza/clean_architecture/repository"
	"log"
)

type service struct {
	cache 	repository.Cache
	repo  	repository.DbRepository
	queue 	job_queue.JobQueue
}

type AppService interface {
	BeginTx() (repository.Tx, error)
	RollbackTx() error
	CommitTx() error
	WithTx(tx repository.Tx) AppService
	UrlService
}

func New(_repo repository.DbRepository, queue job_queue.JobQueue, cache repository.Cache) AppService{
	service := &service{
		repo : _repo,
		queue: queue,
		cache: cache,
	}

	if err := service.queue.CreateQueue(job_queue.URL_SHORTENER_QUEUE); err != nil{
		log.Fatal(err)
	}

	err := service.queue.Consume(job_queue.URL_SHORTENER_QUEUE, func(id string, body []byte) {
		switch id {
		case job_queue.INCREASE_VISIT_COUNT_PM:
			var url entity.Url
			if err := json.Unmarshal(body, &url); err != nil{
				log.Println(err)
			}
			service.IncreaseVisitCount(url)
		}
	})

	if err != nil{
		log.Fatal(err)
	}

	return service
}

func (s service)WithTx(tx repository.Tx) AppService{
	temp := s
	temp.repo = temp.repo.WithTx(tx)
	return temp
}

func (s service)BeginTx() (repository.Tx, error){
	return s.repo.Begin()
}

func (s service)RollbackTx() error{
	return s.repo.Rollback()
}

func (s service)CommitTx() error{
	return s.repo.Commit()
}