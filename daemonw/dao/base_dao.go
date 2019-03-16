package dao

import (
	"github.com/jmoiron/sqlx"
	"database/sql"
	"daemonw/db"
)

type daoConn interface {
	Get(obj interface{}, sqlStatement string, args ...interface{}) error
	Select(obj interface{}, sqlStatement string, args ...interface{}) error
	NamedExec(sqlStatement string, arg interface{}) (sql.Result, error)
	Exec(sqlStatement string, args ...interface{}) (sql.Result, error)
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
