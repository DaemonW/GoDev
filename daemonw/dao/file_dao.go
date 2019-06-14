package dao

import "daemonw/model"

type fileDao struct {
	*baseDao
}

func newFileDao() *fileDao {
	return &fileDao{baseDao: newBaseDao()}
}

func (dao *fileDao) AddFile(f model.File) error {
	schema := `INSERT INTO files(id,name,size,parent_id,type,create_at, meta_data) 
						VALUES (:id,:name,:size,:parent_id,:type,:create_at,:meta_data)`
	_, err := dao.CreateObj(schema, f)
	return err
}
