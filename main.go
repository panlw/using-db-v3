package main

import (
	"database/sql"
	"fmt"
	"log"
	"reflect"
	"time"
	"unsafe"

	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
	"upper.io/db.v3/mysql"
)

var settings = mysql.ConnectionURL{
	Database: `miot_v2`,
	Host:     `localhost:3306`,
	User:     `dev`,
	Password: `Dev.1234`,
}

type userGrp struct {
	ID        int64     `db:"id"`
	TID       int64     `db:"tid"`
	Code      string    `db:"code"`
	Name      string    `db:"name"`
	CreatedAt time.Time `db:"created_at"`
}

func (grp *userGrp) toString() string {
	return fmt.Sprintf("%s %s (ID: %d) %v",
		grp.Code, grp.Name, grp.ID, grp.CreatedAt)
}

// https://golang.org/pkg/database/sql/
func logQueryErr(err error) bool {
	if err == sql.ErrNoRows || err == db.ErrNoMoreRows {
		log.Println(`[Neo] No data`)
		return true
	}
	if err != nil {
		log.Fatal(err)
		return true
	}
	return false
}

func logQueryRows(rows []userGrp, totalRows uint64, totalPages uint) {
	log.Printf("Rows: %d, Total Rows: %d, Total Pages: %d",
		len(rows), totalRows, totalPages)
	for _, row := range rows {
		log.Println(row.toString())
	}
}

func ormQuery(db sqlbuilder.Database) {
	res := db.Collection("iam_grp").Find().
		Where("code like 'SFS_%'").OrderBy("code").
		Paginate(3).Page(2)

	var grps []userGrp
	logQueryErr(res.All(&grps))

	totalRows, err := res.TotalEntries()
	logQueryErr(err)

	totalPages, err := res.TotalPages()
	logQueryErr(err)

	logQueryRows(grps, totalRows, totalPages)
}

func ormQueryRow(db db.Database) {
	res := db.Collection("iam_grp").Find().
		Where("code like 'SFS_%'").OrderBy("code").
		Limit(1)

	var grp userGrp
	logQueryErr(res.One(&grp))
	log.Println(grp)
}

func rawQuery(db sqlbuilder.Database) {
	sql := `
		select * from iam_grp
		order by code limit 3,2
	`
	rows, err := db.Query(sql)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	iter := sqlbuilder.NewIterator(rows)
	var grps []userGrp
	logQueryErr(iter.All(&grps))
	logQueryRows(grps, 0, 0)
}

func rawQueryRow(db sqlbuilder.Database) {
	sql := `
		select id, code, name from iam_grp
		where code = ?
		order by code limit 1
	`

	stmt, err := db.Prepare(sql)
	if err != nil {
		panic(err)
	}

	codes := []string{`X`, `SFS_LY_UG1`, `SFS_LY_UG2`}
	var grp userGrp
	for _, code := range codes {
		log.Printf("[Neo] Query: %s\n", code)
		if !logQueryErr(stmt.QueryRow(code).
			Scan(&grp.ID, &grp.Code, &grp.Name)) {
			log.Println(grp)
		}
	}
}

func rawQueryRowUnsafe(db sqlbuilder.Database) {
	sql := `
		select id, code, name from iam_grp
		where code = ?
		order by code limit 1
	`
	stmt, err := db.Prepare(sql)
	if err != nil {
		panic(err)
	}

	code := `SFS_LY_UG1`
	log.Printf("[Neo] Query: %s\n", code)
	row := stmt.QueryRow(code)

	rows, err := unwrapRow(row)
	logQueryErr(err)

	iter := sqlbuilder.NewIterator(rows)
	var grp userGrp
	logQueryErr(iter.One(&grp))
	log.Println(grp)
}

// unwrap `sql.Row`
func unwrapRow(r *sql.Row) (*sql.Rows, error) {
	rval := reflect.ValueOf(*r)
	fval := rval.FieldByName(`err`)
	if fval.IsNil() {
		fval = rval.FieldByName(`rows`)
		rows := (*sql.Rows)(unsafe.Pointer(fval.Pointer()))
		return rows, nil
	}
	return nil, fval.Interface().(error)
}

func main() {
	db, err := mysql.Open(settings)
	if err != nil {
		log.Printf("Error: %v\n", err)
	}
	defer db.Close()

	// Set this to true to enable the query logger which will print all SQL
	// statements to stdout.
	db.SetLogging(true)

	rawQueryRowUnsafe(db)
}
