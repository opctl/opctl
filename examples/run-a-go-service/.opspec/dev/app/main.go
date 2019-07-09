package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

func handler(w http.ResponseWriter, r *http.Request) {
	log.Printf("request from %s", r.RemoteAddr)

	uidQuery := r.URL.Query()["uid"]
	if len(uidQuery) != 1 {
		fmt.Fprintf(w, "server is alive")
	} else {
		uid, err := strconv.ParseInt(uidQuery[0], 10, 64)
		checkErr(err)
		var username string
		// we open a connection to a mysql db
		// feeding credentials through from environment variables is a standard with tooling in the container ecosystem - we do the same
		db, err := sql.Open("mysql", os.Getenv("MYSQL_USER")+":"+os.Getenv("MYSQL_PASSWORD")+"@tcp("+os.Getenv("MYSQL_HOST")+":3306)/"+os.Getenv("MYSQL_DATABASE"))
		checkErr(err)
		err = db.Ping()
		checkErr(err)
		defer db.Close()
		// query
		err = db.QueryRow("SELECT username FROM users WHERE uid=?", uid).Scan(&username)
		switch {
		case err == sql.ErrNoRows:
			fmt.Fprintf(w, "No user with that ID")
		case err != nil:
			log.Fatal(err)

		default:
			fmt.Fprintf(w, "id %d refers to user %s", uid, username)
		}
	}
}

func main() {
	http.HandleFunc("/", handler)
	log.Printf("server started")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
