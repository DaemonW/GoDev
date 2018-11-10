package errors

//business code
const (
	Internal = 1000 + iota
	Postgres
	MySql
	Orcle
	Redis
	Mongo
	Config

	Biz = 2000 + iota
	Login
	CreateUser
	QueryUser
	DelUser
	InsertUser
	RenameUser
)
