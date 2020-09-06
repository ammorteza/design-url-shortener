package service

import (
	"github.com/ammorteza/clean_architecture/entity"
	job_queue "github.com/ammorteza/clean_architecture/job-queue"
	"log"
	"math/rand"
	"time"
)

type keyService interface {
	MigrateKey(key *entity.UniqueKey) error
	ResetKey(key *entity.UniqueKey) error
	SeedKey() error
	FindUniqueKey() (entity.UniqueKey, error)
	PrepareKeys() error
}

func (s service)FindUniqueKey() (entity.UniqueKey, error){
	var key entity.UniqueKey

	//if err := s.repo.Where("state = ?" , 0).First(&key); err != nil{
	//	return entity.UniqueKey{}, err
	//}

	//if err := s.cache.Ping(); err != nil{
	//	return entity.UniqueKey{}, err
	//}
	//if err := s.cache.Push("unique_key", "morteza"); err != nil{
	//	return entity.UniqueKey{}, err
	//}
	//
	val, err := s.cache.Pop("unique_key")
	if err != nil{
		return entity.UniqueKey{}, err
	}

	count, err := s.cache.ListLen("unique_key")
	log.Println("generated keys in redis queue: ", count)
	if count <= 50{
		s.PrepareKeys()
	}
	key.Key = val
	//if err := s.cache.Set("key1", "sdfsdfsdf"); err != nil{
	//	return entity.UniqueKey{}, err
	//}
	//if err := s.cache.Remove("key1"); err != nil{
	//	return entity.UniqueKey{}, err
	//}
	//val, err := s.cache.Get("key1")
	//if err != nil{
	//	return entity.UniqueKey{}, err
	//}
	//log.Println("key1 in redis is = ", val)
	return key, nil
}

func generateKey() string{
	var chars = []byte("qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM0123456789")
	rand.Seed(time.Now().UnixNano())
	res := ""
	for i := 0 ; i < 6; i++{
		res += string(chars[rand.Intn(62)])
	}
	return res
}

func (s service)PrepareKeys() error{
	if err := s.cache.Ping(); err != nil{
		return err
	}

	for i := 0; i < 100; i++{
		item := entity.UniqueKey{}
		item.Key = generateKey()
		item.State = true
		if err := s.repo.Create(&item); err != nil{
			return err
		}else if err := s.cache.Push("unique_key", item.Key); err != nil{
			return err
		}
	}

	return nil
}

func (s service)SeedKey() error{
	if err := s.jobQueue.Publish(job_queue.KG_PM); err != nil{
		log.Fatal(err)
	}
	return nil
}

func (s service)MigrateKey(key *entity.UniqueKey) error{
	if !s.repo.HasTable(key){
		return s.repo.CreateTable(key)
	}
	return nil
}

func (s service)ResetKey(key *entity.UniqueKey) error{
	if s.repo.HasTable(key){
		return s.repo.DropTable(key)
	}

	return nil
}