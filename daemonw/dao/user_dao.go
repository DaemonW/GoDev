package dao

import (
	"daemonw/model"
	"database/sql"
)

type userDao struct {
	*baseDao
}

func newUserDao() *userDao {
	return &userDao{newBaseDao()}
}

func (dao *userDao) Get(id uint64) (*model.User, error) {
	user := &model.User{}
	err := dao.conn.Get(user, `SELECT * FROM users WHERE id=$1`, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return user, err
}

func (dao *userDao) GetByName(username string) (*model.User, error) {
	user := &model.User{}
	err := dao.conn.Get(user, `SELECT * FROM users WHERE username=$1`, username)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return user, err
}

func (dao *userDao) GetLikeName(username string) ([]model.User, error) {
	users := []model.User{}
	err := dao.conn.Select(&users, `SELECT * FROM users WHERE username LIKE $1`, username+"%")
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return users, err
}

func (dao *userDao) GetAll() ([]model.User, error) {
	users := []model.User{}
	err := dao.conn.Select(&users, `SELECT * FROM users`)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return users, err
}

func (dao *userDao) CreateUser(user *model.User) error {
	schema := `INSERT INTO users(username,password,salt,create_at,update_at) 
						VALUES (:username,:password,:salt,:create_at,:update_at)`
	_, err := dao.conn.NamedExec(schema, user)
	return err
}

func (dao *userDao) DeleteUser(id int64) error {
	_, err := dao.conn.Exec(`DELETE FROM users WHERE id=$1`, id)
	return err
}

func (dao *userDao) ActiveUser(id int64) error {
	_, err := dao.conn.Exec(`UPDATE users SET status=0 WHERE id=$1`, id)
	return err
}
