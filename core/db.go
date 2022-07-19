package core

import (
	"database/sql"
	"sync"
)

func Db(driver, arguments string) (fail error) {
	db, fail = sql.Open(driver, arguments)
	return
}

func dbQueryClose(query dbQuery) {
	query.wg.Wait()
	if query.rows != nil {
		query.rows.Close()
	}
}

func newDbQuery(text string, args ...interface{}) *dbQuery {
	query := dbQuery{}
	query.wg.Add(1)
	go func() {
		query.rows, query.fail = db.Query(text, args...)
		query.wg.Done()
	}()
	return &query
}

func (query *dbQuery) execute(callback func(...interface{}), variables ...interface{}) error {
	query.wg.Wait()
	if query.fail != nil {
		return query.fail
	}
	for query.rows.Next() {
		query.fail = query.rows.Scan(variables...)
		if query.fail != nil {
			return query.fail
		}
		callback(variables...)
	}
	return nil
}

type dbQuery struct {
	fail error
	rows *sql.Rows
	wg   sync.WaitGroup
}

var db *sql.DB
