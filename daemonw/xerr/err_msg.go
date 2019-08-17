package xerr

//context message
const (
	MsgInternal           = "internal server error"
	MsgUserExist          = "user already exists"
	MsgUserNotExist       = "user does not exist"
	MsgUserNotActive      = "user is not active"
	MsgFreezeUser         = "user is frozen"
	MsgIncorrectAuth      = "incorrect username or password"
	MsgCreateUserFail     = "failed to create user"
	MsgActiveUserFail     = "failed to active user"
	MsgBadParam           = "bad param format"
	MsgPermissionDenied   = "permission denied"
	MsgIllegalRequestCode = "illegal request code"
	MsgIllegalVerifyScope = "illegal verify scope"
	MsgAccessFrequency    = "access too frequency"
)
