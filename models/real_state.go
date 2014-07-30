package models

type RealState struct {
	Id   int64
	Name string
	Uri  string `sql:"not null;unique"`
}
