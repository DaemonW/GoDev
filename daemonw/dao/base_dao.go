package dao

import (
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
	"strings"
)

type baseDao struct {
	db       *sqlx.DB
	enableTx bool
	tx       *sqlx.Tx
}

func newBaseDao() *baseDao {
	return &baseDao{db: dbConn, enableTx: false}
}

func (dao *baseDao) BeginTx() error {
	var err error
	dao.enableTx = true
	dao.tx, err = dao.db.Beginx()
	return err
}

func (dao *baseDao) RollBack() error {
	if dao.enableTx {
		dao.enableTx = false
		return dao.tx.Rollback()
	}
	return nil
}

func (dao *baseDao) Commit() error {
	if dao.enableTx {
		dao.enableTx = false
		return dao.tx.Commit()
	}
	return nil
}

func (dao *baseDao) Get(obj interface{}, sqlStatement string, args ...interface{}) error {
	if dao.enableTx {
		return dao.tx.Get(obj, dao.tx.Rebind(sqlStatement), args...)
	}
	dao.db.Rebind(sqlStatement)
	return dao.db.Get(obj, dao.db.Rebind(sqlStatement), args...)
}

func (dao *baseDao) Select(obj interface{}, sqlStatement string, args ...interface{}) error {
	if dao.enableTx {
		dao.tx.Rebind(sqlStatement)
		return dao.tx.Select(obj, dao.tx.Rebind(sqlStatement), args...)
	}
	return dao.db.Select(obj, dao.db.Rebind(sqlStatement), args...)
}

func (dao *baseDao) Delete(sqlStatement string, args ...interface{}) (sql.Result, error) {
	if strings.ToUpper(subString(sqlStatement, 0, 6)) != "DELETE" {
		return nil, errors.New("illegal statement")
	}
	return dao.Exec(sqlStatement, args...)
}

func (dao *baseDao) DeleteObj(sqlStatement string, arg interface{}) (sql.Result, error) {
	if strings.ToUpper(subString(sqlStatement, 0, 6)) != "DELETE" {
		return nil, errors.New("illegal statement")
	}
	return dao.NamedExec(sqlStatement, arg)
}

func (dao *baseDao) Create(sqlStatement string, args ...interface{}) (sql.Result, error) {
	if strings.ToUpper(subString(sqlStatement, 0, 6)) != "INSERT" {
		return nil, errors.New("illegal statement")
	}
	return dao.Exec(sqlStatement, args...)
}

func (dao *baseDao) CreateObj(sqlStatement string, arg interface{}) (sql.Result, error) {
	if strings.ToUpper(subString(sqlStatement, 0, 6)) != "INSERT" {
		return nil, errors.New("illegal statement")
	}
	return dao.NamedExec(sqlStatement, arg)
}

func (dao *baseDao) Update(sqlStatement string, args ...interface{}) (sql.Result, error) {
	if strings.ToUpper(subString(sqlStatement, 0, 6)) != "UPDATE" {
		return nil, errors.New("illegal statement")
	}
	return dao.Exec(sqlStatement, args...)
}

func (dao *baseDao) UpdateObj(sqlStatement string, args ...interface{}) (sql.Result, error) {
	if strings.ToUpper(subString(sqlStatement, 0, 6)) != "UPDATE" {
		return nil, errors.New("illegal statement")
	}
	return dao.NamedExec(sqlStatement, args)
}

func (dao *baseDao) Exec(sqlStatement string, args ...interface{}) (sql.Result, error) {
	if dao.enableTx {
		return dao.tx.Exec(dao.tx.Rebind(sqlStatement), args...)
	}
	return dao.db.Exec(dao.db.Rebind(sqlStatement), args...)
}

func (dao *baseDao) NamedExec(sqlStatement string, args interface{}) (sql.Result, error) {
	if dao.enableTx {
		return dao.tx.NamedExec(dao.tx.Rebind(sqlStatement), args)
	}
	return dao.db.NamedExec(dao.db.Rebind(sqlStatement), args)
}

func subString(Str string, s, t int) string {
	sub := []byte(Str)[s:t]
	return string(sub)
}
