package xerr

//business code
const (
	CodeInternal = 1000 + iota

	CodeBiz = 2000 + iota
	CodeAuth
	CodeLogin
	CodeCreateUser
	CodeQueryUser
	CodeDelUser
	CodeInsertUser
	CodeUpdateUser
)
