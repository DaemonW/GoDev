package entity

import "time"

type App struct {
	Id          uint64    `db:"id"`
	AppId       string    `db:"app_id"`
	Version     string    `db:"version"`
	VersionCode int32     `db:"version_code"`
	Name        string    `db:"name"`
	Icon        string    `db:"icon"`
	Size        int64     `db:"size"`
	Hash        string    `db:"hash"`
	Encrypted   bool      `db:"encrypted"`
	Url         string    `db:"url"`
	CreateAt    time.Time `db:"create_at"`
}

type Update struct {
	AppId  string `db:"app_id"`
	Latest int32  `db:"latest"`
}
