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
	smt := `INSERT INTO apps(app_id,version,version_code,name,size,hash,encrypted,url,create_at) 
			VALUES (:app_id,:version,:version_code,:name,:size,:hash,:encrypted,:url,:created_at)
			RETURNING id`
	_, err := dao.CreateObj(smt, app)
	return err
}

func (dao *appDao) CreateAppInfo(app *AppInfo) error {
	smt := `INSERT INTO app_infos(id,name,package,version,description,change_log,image_detail,language,country,vendor,category) 
			VALUES (:id,:name,:package,:version,:description,:change_log,:image_detail,:language,:country,:vendor,:category)
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

func (dao *appDao) GetAppInfoById(id uint64) (*AppInfo, error) {
	info := &AppInfo{}
	smt := `SELECT * FROM app_infos WHERE id=?`
	err := dao.Get(info, smt, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return info, err
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

func (dao *appDao) GetLatestApps() ([]App, error) {
	apps := make([]App, 0)
	//smt := `SELECT * FROM apps INNER JOIN updates ON apps.app_id=updates.app_id, app.version_code=updates.latest`
	smt := `SELECT apps.id,apps.app_id,apps.version,apps.version_code,apps.name,apps.size,apps.hash,apps.encrypted,apps.url
			FROM apps INNER JOIN updates ON apps.app_id=updates.app_id AND apps.version_code=updates.latest`
	err := dao.Select(&apps, smt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return apps, err
}

func (dao *appDao) GetLatestApp(appId string) ([]App, error) {
	apps := make([]App, 0)
	smt := `SELECT apps.id,apps.app_id,apps.version,apps.version_code,apps.name,apps.size,apps.hash,apps.encrypted,apps.url 
			FROM apps INNER JOIN updates ON apps.app_id=updates.app_id AND apps.version_code=updates.latest WHERE apps.app_id=?`
	err := dao.Select(&apps, smt, appId)
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
