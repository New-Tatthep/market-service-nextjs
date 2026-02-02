package datastore

import (
	"database/sql"
	"errors"
	"fmt"
	"market-service/custom_error"

	"github.com/New-Tatthep/microservice"
)

const (
	PostgresContextName = "pgdb"
	MySQLContexName     = "mysqldb"
)

type StoreAction interface {
	Validate() error
	Session() SessionAction
	// AppConfig() AppConfigAction
	Log() LogAction
	// Counter() CounterAction
	// MobileSession() MobileSessionAction
	ProductAction() ProductDataStoreAction
	EmployeeAction() EmployeeDataStoreAction
}

type store struct {
	conn   *sql.DB
	logger microservice.IContextLogger
}

type Option func(st *store) error

func New(options ...Option) (StoreAction, error) {
	st := new(store)

	for _, opt := range options {
		if err := opt(st); err != nil {
			return nil, custom_error.Wrap(err)
		}
	}

	if err := st.Validate(); err != nil {
		return nil, custom_error.Wrap(err)
	}

	return st, nil
}

func WithDBConnection(conn *sql.DB) Option {
	return func(st *store) error {
		st.conn = conn

		return nil
	}
}

func WithDBContext(ctx microservice.IContext, dbContextName string) Option {
	return func(st *store) error {
		dbStore, found := ctx.DB(dbContextName)
		if !found {
			return fmt.Errorf("db context name %s not found", dbContextName)
		}

		st.conn = dbStore.Conn()

		return nil
	}
}

func WithMicroserviceDB(db microservice.IDBStore) Option {
	return func(st *store) error {
		st.conn = db.Conn()

		return nil
	}
}

func WithMicroserviceLogger(logger microservice.IContextLogger) Option {
	return func(st *store) error {
		st.logger = logger

		return nil
	}
}

func (st *store) Validate() error {
	if st.conn == nil {
		return errors.New("not found db connection")
	}

	return nil
}
