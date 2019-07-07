package dao

import (
	. "daemonw/model"
	"database/sql"
)

type userInfoDao struct {
	*baseDao
}

func newUserInfoDao() *userInfoDao {
	return &userInfoDao{newBaseDao()}
}

func (dao *userInfoDao) Get(id uint64) (*User, error) {
	user := &User{}
	err := dao.baseDao.Get(user, `SELECT * FROM user_infos WHERE user_id=?`, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return user, err
}

func (dao *userInfoDao) GetByEmail(email string) (*User, error) {
	user := &User{}
	err := dao.baseDao.Get(user, `SELECT * FROM user_infos WHERE email=?`, email)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return user, err
}

func (dao *userInfoDao) GetByPhone(phone string) (*User, error) {
	user := &User{}
	err := dao.baseDao.Get(user, `SELECT * FROM user_infos WHERE phone=?`, phone)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return user, err
}

func (dao *userInfoDao) GetAll() ([]User, error) {
	users := []User{}
	err := dao.baseDao.Select(&users, `SELECT * FROM user_infos`)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return users, err
}

func (dao *userInfoDao) AddUserInfo(info *UserInfo) error {
	schema := `INSERT INTO user_infos(user_id,nickname,sex,age,email,phone,address,ip,meta) 
						VALUES (:user_id,:nickname,:sex,:age,:email,:phone,:address,:ip,:meta)`
	_, err := dao.baseDao.NamedExec(schema, info)
	return err
}

func (dao *userInfoDao) DeleteUserInfo(id int64) error {
	_, err := dao.baseDao.Exec(`DELETE FROM user_infos WHERE user_id=?`, id)
	return err
}

func (dao *userInfoDao) UpdateUserInfo(id int64, info UserInfo) error {
	_, err := dao.baseDao.NamedExec(`UPDATE user_infos SET nickname=:nickname, sex=:sex, age=:age, email=:email,
						phone=:phone, address=:address, ip=:ip, meta=:meta WHERE id=?`, id)
	return err
}
