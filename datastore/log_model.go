package datastore

import "database/sql"

type InsertSignOnLogInput struct {
	LogCode      string
	Username     string
	IpAddress    sql.NullString
	MachineId    sql.NullString
	UserAgent    sql.NullString
	SourceAction string
	Action       string
	LoginType    sql.NullString
	Status       string
	LastUpdate   int64
	Message      sql.NullString
}
