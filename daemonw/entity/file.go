package entity

import "time"

type File struct {
	Id       uint64            `db:"id"`
	Name     string            `db:"name"`
	Size     int               `db:"size"`
	ParentId uint64            `db:"parent_id"`
	Type     uint8             `db:"type"`
	CreateAt time.Time         `db:"create_at"`
	MetaData map[string]string `db:"meta_data"`
}

func newFile(name string, parentId uint64) *File {
	return &File{Name: name, ParentId: parentId}
}
