package repository

type DbRepository interface {
	Close() error
	Model(table interface{}) DbRepository
	Where(query interface{}, args ...interface{}) DbRepository
	HasTable(table interface{}) bool
	CreateTable(table interface{}) error
	DropTable(table interface{}) error
	AddForeignKey(model interface{}, field, dest, onDelete, onUpdate string) error
	Create(table interface{}) error
	Find(model interface{}, res interface{}) error
	First(res interface{}) error
	Save(model interface{}) error
	Updates(model interface{}) error
	Update(args ...interface{}) error

	WithTx(tx Tx) DbRepository
	Begin() (Tx, error)
	Rollback() error
	Commit() error
}

type Tx struct {
	ID 				string
}

