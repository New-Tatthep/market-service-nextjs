package datastore

import (
	"database/sql"
	"fmt"
	"new-service/binding"
	"new-service/custom_error"
	"strconv"
	"strings"

	"github.com/Masterminds/squirrel"
)

type dbAction interface {
	InsertEmployee(input EmployeeModel) error
	UpdateEmployee(input EmployeeModel) error
	DeleteEmployee(id string) error
	FilterEmployee(filterData *FilterData) ([]*EmployeeModel, int64, error)
}

func (act *action) InsertEmployee(input EmployeeModel) error {

	builder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	insertBuilder := builder.Insert(
		binding.TblEmployeeTableName,
	).Columns(
		"id",
		"name",
	).Values(
		input.ID,
		input.Name,
	)

	sqlCmd, values, err := insertBuilder.ToSql()
	if err != nil {
		return custom_error.Wrap(err)
	}

	stmt, err := act.prepare(sqlCmd)
	if err != nil {
		return custom_error.Wrap(err)
	}
	defer stmt.Close()

	result, err := stmt.Exec(values...)
	if err != nil {
		return custom_error.Wrap(err)
	}

	_, err = result.RowsAffected()
	if err != nil {
		return custom_error.Wrap(err)
	}

	return nil
}

func (act *action) UpdateEmployee(input EmployeeModel) error {
	sqlCmd := `
	UPDATE tbl_employee SET name = $2  WHERE id = $1
    `

	stmt, err := act.prepare(sqlCmd)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		input.ID,
		input.Name,
	)
	if err != nil {
		return err
	}

	return nil
}

func (act *action) DeleteEmployee(id string) error {
	sqlCmd := `
	DELETE FROM tbl_employee  WHERE id = $1
    `

	stmt, err := act.prepare(sqlCmd)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		id,
	)
	if err != nil {
		return err
	}

	return nil
}

func (act *action) FilterEmployee(filterData *FilterData) ([]*EmployeeModel, int64, error) {
	searchBuilder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	totalBuilder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	searchQuery := searchBuilder.Select(
		"id",
		"name",
	).From(`tbl_employee`)

	totalQuery := totalBuilder.Select("COUNT(id) as total").From(`tbl_employee`)

	// Define filter variables for date and time ranges
	var limit, offset int64

	for key, val := range filterData.Filters {

		if val == "" {
			continue
		}

		switch key {
		case "name":
			name := val.(string)
			searchQuery = searchQuery.Where(
				squirrel.And{
					squirrel.Like{"LOWER(name)": fmt.Sprintf("%%%s%%", strings.ToLower(name))},
				})
			totalQuery = totalQuery.Where(
				squirrel.And{
					squirrel.Like{"LOWER(name)": fmt.Sprintf("%%%s%%", strings.ToLower(name))},
				})
		case "limit":
			valStr := val.(float64)
			conv, err := strconv.ParseInt(fmt.Sprintf("%.0f", valStr), 10, 64)
			if err != nil {
				return nil, 0, err
			}
			limit = conv
		case "offset":
			valStr := val.(float64)
			conv, err := strconv.ParseInt(fmt.Sprintf("%.0f", valStr), 10, 64)
			if err != nil {
				return nil, 0, err
			}
			offset = conv
		}
	}

	if limit > 0 {
		searchQuery = searchQuery.Limit(uint64(limit))
	}
	if offset > 0 {
		searchQuery = searchQuery.Offset(uint64(offset))
	}

	sqlCmd, values, err := searchQuery.ToSql()
	if err != nil {
		return nil, -1, err
	}

	stmt, err := act.prepare(sqlCmd)
	if err != nil {
		return nil, -1, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(values...)
	if err != nil {
		return nil, -1, err
	}
	defer rows.Close()
	results := make([]*EmployeeModel, 0)
	for rows.Next() {
		item := &EmployeeModel{}
		if err := rows.Scan(
			&item.ID,
			&item.Name,
		); err != nil {
			return nil, -1, err
		}

		results = append(results, item)
	}

	totalResult, err := totalQuery.RunWith(act.dbStore.Conn()).Query()
	if err != nil {
		return nil, -1, err
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(totalResult)

	var total int64
	for totalResult.Next() {
		if err := totalResult.Scan(&total); err != nil {
			return nil, -1, err
		}
	}

	return results, total, nil
}
