package dbx

import (
	"database/sql"
	"log"
	"reflect"
	"unsafe"

	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

// UnwrapRow Unwrap `sql.Row`
func UnwrapRow(r *sql.Row) (*sql.Rows, error) {
	rval := reflect.ValueOf(*r)

	fval := rval.FieldByName(`err`)
	if !fval.IsNil() {
		return nil, fval.Interface().(error)
	}

	fval = rval.FieldByName(`rows`)
	rows := (*sql.Rows)(unsafe.Pointer(fval.Pointer()))
	return rows, nil
}

// HandleErr https://golang.org/pkg/database/sql/
func HandleErr(err error) bool {
	if err == nil {
		return false
	}
	if err == sql.ErrNoRows {
		log.Println(`[DBX] No data`)
		return true
	}
	if err == db.ErrNoMoreRows {
		log.Println(`[DBX] No more data`)
		return true
	}
	log.Fatalf("[DBX] %v", err)
	return true
}

type (
	// Pager 分页条件
	Pager struct {
		size uint
		page uint
	}
	// Page 分页结果
	Page interface {
		Total() uint64
		Pages() uint
	}
	page struct {
		totol uint64
		pages uint
	}
)

func (p *page) Total() uint64 {
	return p.totol
}
func (p *page) Pages() uint {
	return p.pages
}

// QueryPage 分页查询
func QueryPage(res db.Result, dest interface{}) (Page, error) {
	if err := res.All(dest); err != nil {
		return nil, err
	}

	total, err := res.TotalEntries()
	if err != nil {
		return nil, err
	}

	pages, err := res.TotalPages()
	if err != nil {
		return nil, err
	}
	return &page{total, pages}, nil
}

// QueryRow 获取一行数据
func QueryRow(stmt *sql.Stmt, dest interface{}, args ...interface{}) error {
	rows, err := stmt.Query(args...)
	defer rows.Close()
	if err != nil {
		return err
	}
	return sqlbuilder.NewIterator(rows).One(dest)
}

// QueryRows 获取多行数据
func QueryRows(stmt *sql.Stmt, dest interface{}, args ...interface{}) error {
	rows, err := stmt.Query(args...)
	defer rows.Close()
	if err != nil {
		return err
	}
	return sqlbuilder.NewIterator(rows).All(dest)
}
