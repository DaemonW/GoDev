package dao

import (
	"daemonw/entity"
	"database/sql"
)

type userDao struct {
	*baseDao
}

func newUserDao() *userDao {
	return &userDao{newBaseDao()}
}

func (dao *userDao) Get(id uint64) (*entity.User, error) {
	user := &entity.User{}
	err := dao.baseDao.Get(user, `SELECT * FROM users WHERE id=?`, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return user, err
}

func (dao *userDao) GetByName(username string) (*entity.User, error) {
	user := &entity.User{}
	err := dao.baseDao.Get(user, `SELECT * FROM users WHERE username=?`, username)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return user, err
}

func (dao *userDao) GetLikeName(username string) ([]entity.User, error) {
	users := []entity.User{}
	err := dao.baseDao.Select(&users, `SELECT * FROM users WHERE username LIKE ?`, username+"%")
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return users, err
}

func (dao *userDao) GetAll() ([]entity.User, error) {
	users := []entity.User{}
	err := dao.baseDao.Select(&users, `SELECT * FROM users`)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return users, err
}

func (dao *userDao) CreateUser(user *entity.User) error {
	schema := `INSERT INTO users(username,password,salt,status,role,create_at,update_at) 
						VALUES (:username,:password,:salt,:status,:role,:create_at,:update_at)`
	_, err := dao.baseDao.NamedExec(schema, user)
	return err
}

func (dao *userDao) CreateUserIfNotExist(user *entity.User) error {
	schema := `INSERT INTO users(username,password,salt,status,role,create_at,update_at) 
						VALUES (:username,:password,:salt,:status,:role,:create_at,:update_at)
						ON CONFLICT(username)
						DO UPDATE SET update_at = :update_at, role=:role`
	_, err := dao.baseDao.NamedExec(schema, user)
	return err
}

func (dao *userDao) DeleteUser(id uint64) error {
	_, err := dao.baseDao.Exec(`DELETE FROM users WHERE id=?`, id)
	return err
}

func (dao *userDao) UpdateStatus(id uint64, status uint8) error {
	_, err := dao.baseDao.Exec(`UPDATE users SET status=? WHERE id=?`, status, id)
	return err
}

func (dao *userDao) UpdatePassword(id uint64, password string) error {
	p, s := entity.GenPassword(password, nil)
	schema := `UPDATE users SET password=?, salt=? WHERE id=?`
	_, err := dao.baseDao.Exec(schema, p, s, id)
	return err
}

func (dao *userDao) UpdateRole(id uint64, role uint8) error {
	schema := `UPDATE users SET role=? WHERE id=?`
	_, err := dao.baseDao.Exec(schema, role, id)
	return err
}
