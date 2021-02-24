package entity


type AppGroup struct{
	Id uint64 `db:"id"`
	Name string `db:"name"`
	Priority uint8 `db:"priority"`
}
