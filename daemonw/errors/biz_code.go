package errors

//business code
const (
	Internal = 1000 + iota
	Postgres
	MySql
	Oracle
	Redis
	Mongo
	Config

	Biz = 2000 + iota
	Auth
	Login
	CreateUser
	QueryUser
	DelUser
	InsertUser
	RenameUser
)
