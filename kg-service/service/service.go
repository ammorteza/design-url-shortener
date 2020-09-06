package service

import (
	job_queue "github.com/ammorteza/clean_architecture/job-queue"
	"github.com/ammorteza/clean_architecture/repository"
	"log"
)

type service struct {
	repo 		repository.DbRepository
	cache 		repository.Cache
	jobQueue 	job_queue.JobQueue
}

type AppService interface {
	BeginTx() (repository.Tx, error)
	RollbackTx() error
	CommitTx() error
	WithTx(tx repository.Tx) AppService
	keyService
}

func New(_repo repository.DbRepository, _cache repository.Cache, jobQueue job_queue.JobQueue) AppService{
	service := &service{
		repo : _repo,
		cache: _cache,
		jobQueue: jobQueue,
	}

	if err := service.jobQueue.CreateQueue(job_queue.KEY_GENERATION_QUEUE); err != nil{
		log.Fatal(err)
	}

	err := service.jobQueue.Consume(job_queue.KEY_GENERATION_QUEUE, func(msg string) {
		switch msg {
		case job_queue.KG_PM:
			if err := service.PrepareKeys(); err != nil{
				log.Fatal(err)
			}
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