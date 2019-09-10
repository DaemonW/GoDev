package entity

import "time"

const (
	FileTypeFile = iota
	FileTypeFolder
)

type File struct {
	Id       uint64    `db:"id"`
	Owner	 uint64    `db:"owner"`
	Name     string    `db:"name"`
	Size     int       `db:"size"`
	ParentId uint64    `db:"parent_id"`
	Type     uint8     `db:"type"`
	CreateAt time.Time `db:"create_at"`
	Meta     string    `db:"meta"`
}

func NewFile(name string, parentId uint64) *File {
	return &File{Name: name, ParentId: parentId}
}

func NewFolder(name string, parentId uint64) *File {
	return &File{Name: name, ParentId: parentId, Type:FileTypeFolder}
}
