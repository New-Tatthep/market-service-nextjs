package datastore

import "database/sql"

type LogAction interface {
	InsertSignOnLog(tx *sql.Tx, input InsertSignOnLogInput) error
}

func (st *store) Log() LogAction {
	return st
}

func (st *store) InsertSignOnLog(tx *sql.Tx, input InsertSignOnLogInput) error {
	sqlCmd := `
	INSERT INTO log_sign_on (
		log_code,
		username,
		ip_address,
		machine_id,
		source_action,
		action,
		login_type,
		status,
		last_update,
		message,
		user_agent
	) VALUES (
		$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
	)
	`

	var stmt *sql.Stmt
	var err error
	if tx != nil {
		stmt, err = tx.Prepare(sqlCmd)
	} else {
		stmt, err = st.conn.Prepare(sqlCmd)
	}

	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		input.LogCode,
		input.Username,
		input.IpAddress,
		input.MachineId,
		input.SourceAction,
		input.Action,
		input.LoginType,
		input.Status,
		input.LastUpdate,
		input.Message,
		input.UserAgent,
	)
	if err != nil {
		return err
	}

	return nil
}
