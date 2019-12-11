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

type AppInfo struct {
	Id          uint64 `db:"id"`
	Package     string `db:"package"`
	Version     string `db:"version"`
	Description string `db:"description"`
	ChangeLog   string `db:"change_log"`
	ImageDetail string `db:"image_detail"`
	Language    string `db:"language"`
	Country     string `db:"country"`
}
