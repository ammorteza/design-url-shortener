package gorm

import (
	"errors"
	"fmt"
	"github.com/ammorteza/clean_architecture/repository"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/nu7hatch/gouuid"
	"time"
)

var (
	txConnections 		map[string]*gorm.DB
)

type repo struct {
	dbConn 				*gorm.DB
	txConn				*gorm.DB
	currentTx			repository.Tx
	txEnable 			bool
}

func init(){
	txConnections = make(map[string]*gorm.DB)
}

func New() repository.DbRepository {
	r := &repo{
		txEnable: false,
	}

	if r.dbConn == nil{
		db := new(gorm.DB)
		var err error
		for {
			fmt.Println("waiting for initializing database ...")
			db, err = gorm.Open("mysql", "root:ush@1234@tcp(172.28.1.21:3306)/ush_db?charset=utf8mb4&parseTime=True&loc=Local")
			if err == nil{
				//log.Fatal(err)
				break
			}
			time.Sleep(time.Second * 10)
		}


		r.dbConn = db
	}
	return r
}

func (r repo)Close() error{
	return r.dbConn.Close()
}

func (r repo)Model(table interface{}) repository.DbRepository{
	temp := r
	temp.dbConn = temp.getConn().Model(table)
	return temp
}

func (r repo)WithTx(tx repository.Tx) repository.DbRepository{
	txConn, ok := txConnections[tx.ID]
	if !ok{
		panic("txConnection id does not exist")
	}
	temp := r
	temp.txEnable = true
	temp.txConn = txConn
	temp.currentTx = tx
	return temp
}

func (r repo)getConn() *gorm.DB{
	if r.txEnable{
		return r.txConn
	}

	return r.dbConn
}

func (r repo)closeTx() error{
	_, ok := txConnections[r.currentTx.ID]
	if !ok{
		return errors.New("tx connection does not exist")
	}

	delete(txConnections, r.currentTx.ID)
	return nil
}

func (r repo)Find(model interface{}, res interface{}) error{
	return r.getConn().Model(&model).Find(res).Error
}

func (r repo)Create(table interface{}) error{
	return r.getConn().Create(table).Error
}

func (r repo)First(res interface{}) error{
	return r.getConn().First(res).Error
}

func (r repo)Offset(off int64) repository.DbRepository{
	temp := r
	temp.dbConn = temp.getConn().Offset(off)
	return temp
}

func (r repo)Count(value interface{}) error{
	return r.getConn().Count(value).Error
}

func (r repo)Save(model interface{}) error{
	return r.getConn().Save(model).Error
}

func (r repo)HasTable(table interface{}) bool{
	return r.getConn().HasTable(table)
}

func (r repo)CreateTable(table interface{}) error{
	return r.getConn().CreateTable(table).Error
}

func (r repo)ResetTable(table interface{}) error{
	return r.getConn().DropTable(table).Error
}

func (r repo)DropTable(table interface{}) error {
	return r.getConn().DropTable(table).Error
}
func (r repo)AddForeignKey(model interface{}, field, dest, onDelete, onUpdate string) error{
	return r.getConn().Model(model).AddForeignKey(field, dest, onDelete, onUpdate).Error
}

func (r repo)Begin() (repository.Tx, error){
	id, err := uuid.NewV4()
	if err != nil{
		return repository.Tx{}, err
	}

	txConn := r.dbConn.Begin()
	if txConn.Error != nil{
		return repository.Tx{}, txConn.Error
	}

	txConnections[id.String()] = txConn
	tx := repository.Tx{
		ID: id.String(),
	}
	return tx, nil
}

func (r repo)Rollback() error{
	if err := r.getConn().Rollback().Error; err != nil{
		return err
	}

	if err := r.closeTx(); err != nil{
		return err
	}

	return nil
}

func (r repo)Commit() error{
	if err := r.getConn().Commit().Error; err != nil{
		return err
	}

	if err := r.closeTx(); err != nil{
		return err
	}

	return nil
}


func (r repo)Where(query interface{}, args ...interface{}) repository.DbRepository{
	temp := r
	temp.dbConn = temp.getConn().Where(query, args...)
	return temp
}

func (r repo)Updates(model interface{}) error{
	return r.getConn().Updates(model).Error
}

func (r repo)Update(args ...interface{}) error{
	return r.getConn().Update(args...).Error
}