package errors

//context message
const (
	MsgInternal           = "internal server error"
	MsgUserExist          = "user already exists"
	MsgUserNotExist       = "user does not exist"
	MsgIncorrectAuth      = "incorrect username or password"
	MsgCreateUserFail     = "failed to create user"
	MsgActiveUserFail     = "failed to active user"
	MsgBadParam           = "bad param format"
	MsgIllegalVerifyScope = "illegal verify scope"
)
