package datastore

import "database/sql"

type SessionAction interface {
	GetUserProfile(username string) (*EmployeeModel, error)
	UpdateLoginFail(tx *sql.Tx, userCode string, count int64, status string) error
}

func (st *store) Session() SessionAction {
	return st
}

//  code      varchar(50) NOT NULL,
//             user_code VARCHAR NOT NULL,
//             user_name VARCHAR NOT NULL,
//             user_type varchar(10) NOT NULL,
//             password VARCHAR(255) NOT NULL,
//             first_name VARCHAR(100) NOT NULL,
//             last_name VARCHAR(100) NOT NULL,
//             email VARCHAR(150) NOT NULL,
//             birth_date          bigint NOT NULL,
//             profile_image      JSONB NULL DEFAULT '{}'::JSONB,
//             mobile_no VARCHAR(20) NOT NULL,
//             create_code             varchar(50) NOT NULL,
//             create_time             bigint NOT NULL,
//             update_code             varchar(50) NOT NULL,
//             update_time             bigint NOT NULL,
//             status                  varchar(20) NOT NULL,
//             count_login_fail int4 NOT NULL DEFAULT 0,

func (st *store) GetUserProfile(username string) (*EmployeeModel, error) {
	sqlCmd := `
	SELECT emp_profile.user_code,
		emp_profile.user_type,
		LOWER(emp_profile.user_name),
		emp_profile.first_name,
		emp_profile.last_name,
		emp_profile.email,
		emp_profile.profile_image,
		emp_profile.count_login_fail,
		emp_profile.password,
		emp_profile.status
	FROM employee AS emp_profile
		WHERE LOWER(emp_profile.user_name) = $1
	GROUP BY emp_profile.user_code ,
			emp_profile.user_type,
			emp_profile.user_name,
			emp_profile.first_name,
			emp_profile.last_name,
			emp_profile.email ,
			emp_profile.profile_image,
			emp_profile.count_login_fail,
			emp_profile.password,
			emp_profile.status
	`

	print(sqlCmd)
	stmt, err := st.conn.Prepare(sqlCmd)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result EmployeeModel
	for rows.Next() {
		if err := rows.Scan(
			&result.UserCode,
			&result.UserType,
			&result.UserName,
			&result.FirstName,
			&result.LastName,
			&result.Email,
			&result.ProfileImage,
			&result.CountLoginFail,
			&result.Password,
			&result.Status,
		); err != nil {
			return nil, err
		}
	}

	return &result, nil
}

func (st *store) UpdateLoginFail(tx *sql.Tx, userCode string, count int64, status string) error {
	sqlCmd := `
	UPDATE employee
	SET count_login_fail = $1,
		status = $2
	WHERE user_code = $3
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

	_, err = stmt.Exec(count, status, userCode)
	if err != nil {
		return err
	}

	return nil
}
