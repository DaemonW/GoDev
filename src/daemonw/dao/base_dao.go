package dao

import (
	"github.com/jmoiron/sqlx"
	"database/sql"
	"daemonw/db"
)

type daoConn interface {
	Get(dest interface{}, query string, args ...interface{}) error
	Select(dest interface{}, query string, args ...interface{}) error
	NamedExec(query string, arg interface{}) (sql.Result, error)
	Exec(query string, args ...interface{}) (sql.Result, error)
}
type baseDao struct {
	db       *sqlx.DB
	enableTx bool
	tx       *sqlx.Tx
	conn     daoConn
}

func newBaseDao() *baseDao {
	return &baseDao{db: db.GetDB(), enableTx: false, conn: db.GetDB()}
}

func (dao *baseDao) BeginTx() error {
	var err error
	dao.enableTx = true
	dao.tx, err = dao.db.Beginx()
	dao.conn = dao.tx
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
