package dao

import "database/sql"
import . "daemonw/entity"

type updateDao struct {
	*baseDao
}

func NewUpdateDao() *updateDao {
	return &updateDao{baseDao: NewDao()}
}

func (dao *updateDao) GetApp(appId string) (*App, error) {
	app := &App{}
	smt := `SELECT * FROM apps INNER JOIN updates ON apps.app_id=updates.app_id WHERE apps.app_id=?`
	err := dao.Get(app, smt, appId)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return app, err
}

func (dao *updateDao) GetApps() ([]App, error) {
	apps := make([]App, 0)
	smt := `SELECT * FROM apps INNER JOIN updates ON apps.app_id=updates.app_id`
	err := dao.Select(&apps, smt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return apps, err
}
