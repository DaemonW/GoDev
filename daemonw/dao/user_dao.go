package dao

import (
	. "daemonw/entity"
	"database/sql"
)

type userDao struct {
	*baseDao
}

func newUserDao() *userDao {
	return &userDao{newBaseDao()}
}

func (dao *userDao) Get(id uint64) (*User, error) {
	user := &User{}
	err := dao.baseDao.Get(user, `SELECT * FROM users WHERE id=?`, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return user, err
}

func (dao *userDao) GetByName(username string) (*User, error) {
	user := &User{}
	err := dao.baseDao.Get(user, `SELECT * FROM users WHERE username=?`, username)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return user, err
}

func (dao *userDao) GetLikeName(username string) ([]User, error) {
	users := []User{}
	err := dao.baseDao.Select(&users, `SELECT * FROM users WHERE username LIKE ?`, username+"%")
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return users, err
}

func (dao *userDao) GetAll() ([]User, error) {
	users := []User{}
	err := dao.baseDao.Select(&users, `SELECT * FROM users`)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return users, err
}

func (dao *userDao) CreateUser(user *User) error {
	schema := `INSERT INTO users(username,password,salt,status,role,create_at,update_at) 
						VALUES (:username,:password,:salt,:status,:role,:create_at,:update_at)`
	_, err := dao.baseDao.NamedExec(schema, user)
	return err
}

func (dao *userDao) DeleteUser(id int64) error {
	_, err := dao.baseDao.Exec(`DELETE FROM users WHERE id=?`, id)
	return err
}

func (dao *userDao) UpdateStatus(id uint64, status uint8) error {
	_, err := dao.baseDao.Exec(`UPDATE users SET status=? WHERE id=?`, status, id)
	return err
}

func (dao *userDao) UpdatePassword(user *User) error {
	schema := `UPDATE users set password=:passowrd`
	_, err := dao.baseDao.NamedExec(schema, user)
	return err
}

func (dao *userDao) UpdateRole(user *User) error {
	schema := `UPDATE users set role=:role`
	_, err := dao.baseDao.NamedExec(schema, user)
	return err
}
