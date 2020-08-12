package main

import (
	"database/sql" //new import
	"flag"
	"log"
	"net/http"
	"os"
	_ "github.com/go-sql-driver/mysql" //new import
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP Network Address")
	//Define a new command-line flag for the mysql DNS string
	dsn := flag.String("dsn", "web:Luisa2012@/snippetbox?parseTime=true", "MySQL data source name")
	flag.Parse()

	//db, err := sql.Open("mysql", "web:Luisa2012@tcp(127.0.0.1:3308)/snippetbox?parseTime=true")
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// to keep the main() function tidy I've put the code for creating a connection
	// pool into the separate openDB() function below, We pass openDB() the DSN from
	// the command line flag
	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}

	// We also defer a call to openDB(), so that the connection pool is closed
	// before the main() function exits
	defer db.Close()

	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
	}

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Starting server on %s", *addr)
	// Because the err variable is now already declared in the code above, we need
	// to use the assignment operator = here, instead of the := 'declare and assign' operator
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}

// The openDB() function wraps sql.Open() and returns a sql.DB connection pool
// for a given DSN
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	// if err = db.Ping(); err != nil {
	// 	return nil, err
	// }
	return db, nil
}