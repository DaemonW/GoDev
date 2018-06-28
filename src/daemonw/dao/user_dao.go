package dao

import (
	"daemonw/model"
	"database/sql"
)

type UserDao struct {
	*baseDao
}

func NewUserDao() *UserDao {
	return &UserDao{newBaseDao()}
}

func (dao *UserDao) Get(id uint64) (*model.User, error) {
	user := &model.User{}
	err := dao.conn.Get(user, `select * from users where id=$1`, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return user, err
}

func (dao *UserDao) GetByName(username string) (*model.User, error) {
	user := &model.User{}
	err := dao.conn.Get(user, `select * from users where username=$1`, username)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return user, err
}

func (dao *UserDao) GetAll() ([]model.User, error) {
	users := []model.User{}
	err := dao.conn.Select(&users, `select * from users`)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return users, err
}

func (dao *UserDao) CreateUser(user *model.User) error {
	schema := `insert into users(username,password,salt,login_ip,create_at,update_at) 
						values (:username,:password,:salt,:login_ip,:create_at,:update_at)`
	_, err := dao.conn.NamedExec(schema, user)
	return err
}

func (dao *UserDao) DeleteUser(id int64) error {
	_, err := dao.conn.Exec(`delete from users where id=$1`, id)
	return err
}
