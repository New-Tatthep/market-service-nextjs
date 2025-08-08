package datastore

import (
	"database/sql"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/New-Tatthep/microservice"
)

const (
	PostgresContextName = "pgdb"
	MySQLContexName     = "mysqldb"

	EmptyJson      = "{}"
	EmptyJsonArray = "[]"
)

var (
	ErrorDBStoreNotFound = fmt.Errorf("db context name %s not found", PostgresContextName)
	ErrorRowNotFound     = fmt.Errorf("row not found")
	ErrorRowNotAffected  = fmt.Errorf("sql: no rows affected")
)

type IAction interface {
	Do(fn func(action IAction) error) error
	dbAction
	ProductAction
}

type action struct {
	dbStore microservice.IDBStore
	txn     *sql.Tx
}

type Option func(act *action) error

func WithDBStore(dbStore microservice.IDBStore) Option {
	return func(act *action) error {
		act.dbStore = dbStore
		return nil
	}
}

func Action(ctx microservice.IContext, options ...Option) IAction {
	act := &action{}

	for _, option := range options {
		if err := option(act); err != nil {
			panic(err)
		}
	}

	if act.dbStore == nil {
		dbStore, err := getDbStore(ctx)
		if err != nil {
			panic(err)
		}
		act.dbStore = dbStore
	}

	return act
}

func (act *action) Do(fn func(action IAction) error) error {
	txn, err := act.dbStore.Conn().Begin()
	if err != nil {
		return err
	}

	act.txn = txn

	if err := fn(act); err != nil {
		if err := act.txn.Rollback(); err != nil {
			act.txn = nil
			return err
		}
		act.txn = nil
		return err
	}

	act.txn = nil

	return txn.Commit()
}

func (act *action) prepare(query string) (*sql.Stmt, error) {
	if act.txn != nil {
		return act.txn.Prepare(query)
	}
	return act.dbStore.Conn().Prepare(query)
}

func (act *action) exec(query string, args ...interface{}) (sql.Result, error) {
	if act.txn != nil {
		return act.txn.Exec(query, args...)
	}
	return act.dbStore.Conn().Exec(query, args...)
}

func (act *action) query(query string, args ...interface{}) (*sql.Rows, error) {
	if act.txn != nil {
		return act.txn.Query(query, args...)
	}
	return act.dbStore.Conn().Query(query, args...)
}

func (act *action) connSquirrelQueryable() squirrel.BaseRunner {
	if act.txn != nil {
		return act.txn
	}
	return act.dbStore.Conn()
}

func getDbStore(ctx microservice.IContext) (microservice.IDBStore, error) {
	dbStore, found := ctx.DB(PostgresContextName)
	if !found {
		return nil, ErrorDBStoreNotFound
	}

	return dbStore, nil
}

func checkErrorNoRow(err error) error {
	if err == sql.ErrNoRows {
		return ErrorRowNotFound
	}

	return err
}
