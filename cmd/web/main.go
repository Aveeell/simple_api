package main

import (
	"avito/pkg/models/mysql"
	"database/sql"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

func openDB(auth string) (*sql.DB, error) {
	db, err := sql.Open("mysql", auth)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

type application struct {
	db *mysql.AvitoModel
}

func main() {

	db, err1 := openDB("avito_user:avitoBackend-2022@tcp(host.docker.internal:3306)/avito?parseTime=True")
	// db, err1 := openDB("root:root@tcp( ?????? )/avito?parseTime=True")

	app := &application{&mysql.AvitoModel{DB: db}}

	if err1 != nil {
		log.Fatal(err1)
	}
	defer db.Close()

	mux := http.NewServeMux()

	mux.HandleFunc("/user/balance", app.getBalance)
	mux.HandleFunc("/user/balance/add", app.addMoney)
	mux.HandleFunc("/history", app.getHistory)
	mux.HandleFunc("/order/create", app.createOrder)
	mux.HandleFunc("/order/approve", app.approveOrder)
	mux.HandleFunc("/bookkeeping", app.getDocument)
	mux.HandleFunc("/bookkeeping/download", app.downloadDocument)

	log.Println("run on http://localhost:8080")
	err := http.ListenAndServe(":8080", mux)
	log.Fatal(err)
}
