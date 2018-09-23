package main

import (
	"fmt"
	"github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/dialers/mysql"
	"github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/proxy"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func main() {
	project := os.Getenv("CLOUD_PROJECT")
	region := os.Getenv("CLOUD_REGION")
	root := os.Getenv("CLOUD_ROOT")
	dbName := os.Getenv("SQL_DBNAME")
	user := os.Getenv("SQL_USERNAME")
	password := os.Getenv("SQL_PASSWORD")
	authFile := os.Getenv("AUTH_FILE_PATH")
	instance := project + ":" + region + ":" + root
	SQLScope := "https://www.googleapis.com/auth/sqlservice.admin"
	ctx := context.Background()
	creds, err := ioutil.ReadFile(authFile)
	if err != nil {
		log.Fatal(err)
	}
	jwt, err := google.JWTConfigFromJSON(creds, SQLScope)
	if err != nil {
		log.Fatal(err)
	}
	client := jwt.Client(ctx)
	proxy.Init(client, nil, nil)
	cfg := mysql.Cfg(instance, user, password)
	cfg.DBName = dbName
	db, err := mysql.DialCfg(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(fmt.Sprintf("USE %s", dbName))
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("DROP TABLE IF EXISTS customers")
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec("CREATE TABLE customers(id SERIAL, first_name VARCHAR(255), last_name VARCHAR(255))")
	if err != nil {
		log.Fatal(err)
	}
	fullNames := []string{
		"John Woo",
		"Jeff Dean",
		"Josh Bloch",
		"Josh Long"}
	stmt, err := db.Prepare("INSERT INTO customers(first_name, last_name) VALUES (?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	for _, fullName := range fullNames {
		log.Printf("Name: %s", fullName)
		name := strings.Split(fullName, " ")
		_, err := stmt.Query(name[0], name[1])
		if err != nil {
			log.Fatal(err)
		}
	}
	rows, err := db.Query("SELECT id, first_name, last_name FROM customers")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var (
		id        string
		firstName string
		lastName  string
	)
	for rows.Next() {
		err := rows.Scan(&id, &firstName, &lastName)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("%v: %s %s", id, firstName, lastName)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec("DROP TABLE customers")
}
