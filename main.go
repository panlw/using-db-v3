package main

import (
	"log"
	"time"

	"upper.io/db.v3/mysql"
)

type userGrp struct {
	ID        int64     `db:"id"`
	TID       int64     `db:"tid"`
	Code      string    `db:"code"`
	Name      string    `db:"name"`
	CreatedAt time.Time `db:"created_at"`
}

var settings = mysql.ConnectionURL{
	Database: `miot_v2`,
	Host:     `localhost:3306`,
	User:     `dev`,
	Password: `Dev.1234`,
}

func main() {
	sess, err := mysql.Open(settings)
	if err != nil {
		log.Fatal(err)
	}
	defer sess.Close()

	// Set this to true to enable the query logger which will print all SQL
	// statements to stdout.
	sess.SetLogging(false)

	// Define a result set without passing a condition to Find(), this means we
	// want to match all the elements on the books table.
	res := sess.Collection("iam_grp").Find()

	// We can use this res object later in different queries, here we'll use it
	// to fetch all the books on our catalog in descending order.
	var grps []userGrp
	if err := res.OrderBy("code").All(&grps); err != nil {
		log.Fatal(err)
	}

	// The books slice has been populated!
	log.Println("User Groups:")
	for _, g := range grps {
		log.Printf("%s %s (ID: %d) %v\n", g.Code, g.Name, g.ID, g.CreatedAt)
	}
}
