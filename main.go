package main

import (
	"log"
	"time"

	"github.com/panlw/using-db-v3/dbx"
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

func ormQueryPage(db sqlbuilder.Database) {
	res := db.Collection("iam_grp").Find().
		Where("code like 'SFS_%'").OrderBy("code").
		Paginate(3).Page(2)

	var grps []userGrp
	if page, err := dbx.FetchPage(res, &grps); !dbx.HandleErr(err) {
		log.Printf("Page: %v, Rows: %v", page, grps)
	}
}

func ormQueryRow(db db.Database) {
	res := db.Collection("iam_grp").Find().
		Where("code like 'SFS_%'").OrderBy("code").
		Limit(1)

	var grp userGrp
	if !dbx.HandleErr(res.One(&grp)) {
		log.Println(grp)
	}
}

func rawQuery(db sqlbuilder.Database) {
	sql := `
		select * from iam_grp
		order by code limit 3,2
	`
	stmt, err := db.Prepare(sql)
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	var grps []userGrp
	if !dbx.HandleErr(dbx.QueryRows(stmt, &grps)) {
		log.Printf("Rows: %v", grps)
	}
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
		log.Printf("[NEO] Query: %s\n", code)
		if !dbx.HandleErr(dbx.QueryRow(stmt, &grp, &code)) {
			log.Printf("Row: %v", grp)
		}
	}
}

func rawQueryScan(db sqlbuilder.Database) {
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
		log.Printf("[NEO] Query: %s\n", code)
		if !dbx.HandleErr(stmt.QueryRow(code).
			Scan(&grp.ID, &grp.Code, &grp.Name)) {
			log.Printf("Row: %v", grp)
		}
	}
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

	ormQueryPage(db)
}
