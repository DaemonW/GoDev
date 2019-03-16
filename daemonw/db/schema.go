package db

const (
	//create table users
	SCHEMA_CREATE_USER_TABLE = `create table users (
      id SERIAL primary key,
      username varchar(64) unique ,
      password varchar(64),
      salt bytea,
      create_at timestamp ,
      update_at timestamp
	)`
)
