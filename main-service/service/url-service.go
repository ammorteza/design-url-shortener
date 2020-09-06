package service

import (
	"encoding/json"
	"github.com/ammorteza/clean_architecture/entity"
	job_queue "github.com/ammorteza/clean_architecture/job-queue"
	"log"
	"math/rand"
	"time"
)

type UrlService interface {
	MigrateUrl(url *entity.Url) error
	ResetUrl(url *entity.Url) error
	CreateUrl(input *entity.Url) error
	FetchUrl(input entity.Url) (entity.Url, error)
	IncreaseVisitCount(input entity.Url)
	GetUrlRandomly() (entity.Url, error)
}

func (s service)GetUrlRandomly() (entity.Url, error){
	var url entity.Url
	var count int64
	if err := s.repo.Model(&url).Count(&count); err != nil{
		return entity.Url{}, err
	}

	if count > 2{
		rand.Seed(time.Now().UnixNano())
		if err := s.repo.Offset(rand.Int63n(count - 1)).First(&url); err != nil{
			return entity.Url{}, err
		}
	}

	return url, nil
}

func (s service)IncreaseVisitCount(input entity.Url){
	var row entity.Url
	if err := s.repo.Where("id = ?", input.ID).First(&row); err != nil{
		log.Println(err)
	}

	row.VisitCount++
	if err := s.repo.Model(&input).Updates(&row); err != nil{
		log.Println(err)
	}
}

func (s service)FetchUrl(input entity.Url) (entity.Url, error){
	var url entity.Url
	tempUrl, err := s.cache.Get(input.UniqueKey)
	if err == nil{
		err = json.Unmarshal([]byte(tempUrl), &url)
		if err == nil{
			temp, err := json.Marshal(entity.Url{
				ID: url.ID,
			})
			if err == nil {
				if err := s.queue.Publish(job_queue.INCREASE_VISIT_COUNT_PM, temp); err != nil {
					log.Println(err)
				}
			}
			return url, nil
		}
	}

	err = s.repo.Where("unique_key = ?", input.UniqueKey).First(&url)
	if err == nil{
		temp, err := json.Marshal(entity.Url{
			ID: url.ID,
		})
		if err == nil{
			if err := s.queue.Publish(job_queue.INCREASE_VISIT_COUNT_PM, temp); err != nil{
				log.Println(err)
			}
		}
		temp, err = json.Marshal(url)
		s.cache.SetWithEp(url.UniqueKey, string(temp), "60")

	}
	return url, err
}

func (s service)CreateUrl(url *entity.Url) error{
	return s.repo.Create(url)
}

func (s service)MigrateUrl(url *entity.Url) error{
	if !s.repo.HasTable(url) {
		if err := s.repo.CreateTable(url); err != nil{
			return err
		}
	}

	return nil
}

func (s service)ResetUrl(url * entity.Url) error{
	if s.repo.HasTable(url){
		if err := s.repo.DropTable(url); err != nil{
			return err
		}
	}

	return nil
}
