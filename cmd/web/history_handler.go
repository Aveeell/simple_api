package main

import (
	"avito/pkg/models"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

func (app *application) getHistory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "Forbidden method", http.StatusBadRequest)
		return
	}

	type Request struct {
		ID        int    `json:id`
		Offset    int    `json:offset`
		Limit     int    `json:limit`
		Sort_type string `json:sort`
		Sort_col  string `json:sort_col`
	}
	var req Request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "can't parse .json", http.StatusInternalServerError)
		return
	}
	if req.ID < 1 {
		http.Error(w, "wrong params", http.StatusBadRequest)
		return
	}

	data, err := app.db.GetHistory(req.ID, req.Limit, req.Offset, req.Sort_type, req.Sort_col)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		} else {
			fmt.Fprintf(w, "%s", err)
		}
		return
	}

	responce, _ := json.Marshal(data)
	w.Header().Set("Content-type", "application/json")
	fmt.Fprintf(w, "%s", string(responce))
}
