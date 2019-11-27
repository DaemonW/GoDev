package dao

import (
	. "daemonw/entity"
	"database/sql"
)

type appDao struct {
	*baseDao
}

func NewAppDao() *appDao {
	return &appDao{baseDao: NewDao()}
}

func (dao *appDao) CreateApp(app *App) error {
	smt := `INSERT INTO apps(app_id,version,version_code,name,size,hash,encrypted,url) 
			VALUES (:app_id,:version,:version_code,:name,:size,:hash,:encrypted,:url)
			RETURNING id`
	_, err := dao.CreateObj(smt, app)
	return err
}

func (dao *appDao) GetApp(appId string, version string) (*App, error) {
	app := &App{}
	smt := `SELECT * FROM apps WHERE app_id=? AND version=?`
	err := dao.Get(app, smt, appId, version)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return app, err
}

func (dao *appDao) GetAllApp(appId string) ([]App, error) {
	apps := make([]App, 0)
	smt := `SELECT * FROM apps WHERE app_id=?`
	err := dao.Select(&apps, smt, appId)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return apps, err
}

func (dao *appDao) GetAppById(id uint64) (*App, error) {
	app := &App{}
	smt := `SELECT * FROM apps WHERE id=?`
	err := dao.Get(app, smt, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return app, err
}

func (dao *appDao) GetAllApps() ([]App, error) {
	apps := make([]App, 0)
	smt := `SELECT * FROM apps`
	err := dao.Select(&apps, smt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return apps, err
}

func (dao *appDao) DeleteApp(appId string, version string) error {
	smt := `DELETE FROM apps WHERE app_id=? AND version=?`
	_, err := dao.Delete(smt, appId, version)
	return err
}

func (dao *appDao) DeleteAppById(id uint64) error {
	smt := `DELETE FROM apps WHERE id=?`
	_, err := dao.Delete(smt, id)
	return err
}
