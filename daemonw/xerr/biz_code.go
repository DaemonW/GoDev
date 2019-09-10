package xerr

//business code
const (
	CodeInternal = 1000 + iota
)
const (
	CodeBiz = 2000 + iota
	CodeRateLimit
	CodeAuth
	CodeVerify
	CodeLogin
	CodeCreateUser
	CodeQueryUser
	CodeDelUser
	CodeInsertUser
	CodeUpdateUser

	CodeCreateFile
	CodeQueryFile
	CodeDeleteFile
	CodeRenameFile
	CodeMoveFile
)
