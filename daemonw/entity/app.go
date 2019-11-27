package entity

type App struct {
	Id          uint64 `db:"id" json:"-"`
	AppId       string `db:"app_id"`
	Version     string `db:"version"`
	VersionCode int32  `db:"version_code"`
	Name        string `db:"name"`
	Size        int64  `db:"size"`
	Hash        string `db:"hash"`
	Encrypted   bool   `db:"encrypted"`
	Url         string `db:"url"`
}

type Update struct {
	AppId  string `db:"app_id"`
	Latest string `db:"latest"`
}
