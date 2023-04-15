package store

type Store interface {
	GetAll(tableName string) (interface{}, error)
	GetByID(entityID int, tableName string) (interface{}, error)
	Save(entity interface{}, tableName string) (interface{}, error)
	Update(entityID int, entity interface{}, tableName string) (interface{}, error)
	Delete(entityID int, tableName string) error
}
