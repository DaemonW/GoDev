package db

const (
	//create table users
	SCHEMA_CREATE_USER_TABLE = `create table users (
      id SERIAL primary key,
      username text unique ,
      password text,
      salt bytea,
      login_ip text,
      create_at timestamp ,
      update_at timestamp
	)`
)
