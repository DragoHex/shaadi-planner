package db

type DB interface {
	Import()
	Export()
	Query()
}
