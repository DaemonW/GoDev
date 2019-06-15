package dao

import . "daemonw/model"

type fileDao struct {
	*baseDao
}

func newFileDao() *fileDao {
	return &fileDao{baseDao: newBaseDao()}
}

func (dao *fileDao) AddFile(f File) error {
	smt := `INSERT INTO files(id,name,size,parent_id,type,create_at, meta_data) 
						VALUES (:id,:name,:size,:parent_id,:type,:create_at,:meta_data) return id`
	_, err := dao.CreateObj(smt, f)
	return err
}


func (dao *fileDao) DeleteFile(pid uint64, name string) error {
	smt := `DELETE FROM files where name=':name' AND parent_id=':parent_id'`
	_, err := dao.Delete(smt)
	return err
}
