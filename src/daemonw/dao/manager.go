package dao

var UserDao *userDao

func InitDaoManager(){
	UserDao = newUserDao()
}