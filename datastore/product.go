package datastore

import (
	"fmt"
	"market-service/binding"

	"github.com/Masterminds/squirrel"
	"github.com/New-Tatthep/microservice"
	"github.com/New-Tatthep/microservice/util/stringutil"
)

type ProductDataStoreAction interface {
	FilterProduct(filters []microservice.IQueryFilter, option microservice.IQueryOption) ([]Product, int64, error)
}

func (st *store) ProductAction() ProductDataStoreAction {
	return st
}

func (sv *store) FilterProduct(filters []microservice.IQueryFilter, option microservice.IQueryOption) ([]Product, int64, error) {
	searchBuilder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	totalBuilder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	searchQuery := searchBuilder.Select(
		"code",
		"name",
		"description",
		"price",
		"status",
		"quantity",
		"COALESCE(image, '{}'::jsonb) AS image_json",
	).From(
		binding.TblProductTableName,
	)

	totalQuery := totalBuilder.Select("COUNT(code) as total").From(
		binding.TblProductTableName,
	)

	for _, filter := range filters {
		switch filter.GetField() {
		case "code":
			if stringutil.IsNotEmptyString(fmt.Sprintf("%s", filter.GetValue())) {
				searchQuery = searchQuery.Where(
					squirrel.And{
						squirrel.Eq{fmt.Sprintf("tb1.%s", filter.GetField()): filter.GetValue()},
					})
				totalQuery = totalQuery.Where(
					squirrel.And{
						squirrel.Eq{fmt.Sprintf("tb1.%s", filter.GetField()): filter.GetValue()},
					})
			}
		}
	}

	if option.GetLimit() > 0 {
		searchQuery = searchQuery.Limit(uint64(option.GetLimit()))
	}
	if option.GetOffset() > 0 {
		searchQuery = searchQuery.Offset(uint64(option.GetOffset()))
	}

	searchQuery = queryOptionBuilder(searchQuery, option)
	sqlCmd, values, err := searchQuery.ToSql()
	if err != nil {
		return nil, -1, err
	}

	stmt, err := sv.conn.Prepare(sqlCmd)
	if err != nil {
		return nil, -1, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(values...)
	if err != nil {
		return nil, -1, err
	}
	defer rows.Close()
	results := make([]Product, 0)
	for rows.Next() {
		result := Product{}
		if err := rows.Scan(
			&result.Code,
			&result.Name,
			&result.Description,
			&result.Price,
			&result.Status,
			&result.Quantity,
			&result.Image,
		); err != nil {
			return nil, -1, err
		}

		results = append(results, result)
	}

	totalResult, err := totalQuery.RunWith(sv.conn).Query()
	if err != nil {
		return nil, -1, err
	}
	defer totalResult.Close()
	var total int64
	for totalResult.Next() {
		if err := totalResult.Scan(&total); err != nil {
			return nil, -1, err
		}
	}

	return results, total, nil
}
