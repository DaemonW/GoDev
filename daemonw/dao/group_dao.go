package dao

import (
	"daemonw/entity"
)

type groupDao struct {
	*baseDao
}

func NewGroupDao() *groupDao {
	return &groupDao{baseDao: NewDao()}
}

func (dao *groupDao) CreateGroup(group *entity.AppGroup) error {
	smt := `INSERT INTO groups(name, priority) VALUES (:name, :priority) RETURNING id`
	row, err := dao.NamedQuery(smt, group)
	if row != nil && row.Next() {
		row.Scan(&group.Id)
	}
	return err
}

func (dao *groupDao) DeleteGroup(id uint64) error {
	smt := `DELETE FROM groups WHERE id = ?`

}

func (dao *groupDao) GetAllGroups() ([]entity.AppGroup, error) {
	smt := `SELECT * from groups`
	groups := make([]entity.AppGroup, 0)
	err := dao.Select(&groups, smt)
	return groups, err
}
