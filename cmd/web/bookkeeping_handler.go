package main

import (
	"avito/pkg/models"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

func createCSV(data []*models.Bookkeeping) {
	csvFile, err := os.Create("downloads/data.csv")
	if err != nil {
		fmt.Println(err)
	}
	defer csvFile.Close()
	writer := csv.NewWriter(csvFile)
	for _, item := range data {
		var row []string
		row = append(row, item.Product)
		row = append(row, strconv.Itoa(item.Total))
		writer.Write(row)
	}
	writer.Flush()
}

func (app *application) getDocument(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "Forbidden method", http.StatusBadRequest)
		return
	}

	type Request struct {
		Date string `json:date`
	}
	var req Request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "can't parse .json", http.StatusInternalServerError)
		return
	}

	data, err := app.db.GetDocument(req.Date)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		} else {
			fmt.Fprintf(w, "%s", err)
		}
		return
	}

	createCSV(data)
	type Link struct {
		Link string
	}

	responce, _ := json.Marshal(Link{"http://localhost:8080/bookkeeping/download"})
	w.Header().Set("Content-type", "application/json")
	fmt.Fprintf(w, "%s", string(responce))
}

func (app *application) downloadDocument(w http.ResponseWriter, r *http.Request) {
	const filename = "data.csv"
	file, err := ioutil.ReadFile("downloads/" + filename)
	if err != nil {
		fmt.Fprintf(w, "%s", err)
	}
	w.Header().Set("Accept-ranges", "bytes")
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", "attachment; filename="+filename+"")
	w.Header().Set("Content-Length", strconv.Itoa(len(file)))
	w.WriteHeader(http.StatusOK)
	w.Write(file)
}
