package db

type DB interface {
  Import(string) 
  Export()
  Query()
}

type Column interface {
  CreateTableQuery() string

}
