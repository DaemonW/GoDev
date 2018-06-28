package dao

import (
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

const (
	DB_CONN_WITHOUT_PASSWORD = "postgres://%s@%s:%d/%s?sslmode=%s"
	DB_CONN_WITH_PASSWORD    = "postgres://%s:%s@%s:%d/%s?sslmode=%s"
)
