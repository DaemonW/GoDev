package dao

import . "daemonw/entity"

type fileDao struct {
	*baseDao
}

func newFileDao() *fileDao {
	return &fileDao{baseDao: newBaseDao()}
}

func (dao *fileDao) CreateFile(f File) error {
	smt := `INSERT INTO files(owner,name,size,parent_id,type,create_at, meta)
						VALUES (:owner,:name,:size,:parent_id,:type,:create_at,:meta) returning id`
	_, err := dao.CreateObj(smt, f)
	return err
}

func (dao *fileDao) DeleteFile(pid uint64, name string) error {
	smt := `DELETE FROM files where name=':name' AND parent_id=':parent_id'`
	_, err := dao.Delete(smt)
	return err
}
