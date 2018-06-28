package dao

import "github.com/jmoiron/sqlx"

type Dao interface {
	Get(map[string]interface{}) (interface{}, error)
	Find(map[string]interface{}) ([]interface{}, error)
	Delete(map[string]interface{}) error
	Update(map[string]interface{}) error
	Insert(interface{}) error
	Table() string
}

type BaseDao struct {
	db *sqlx.DB
}
