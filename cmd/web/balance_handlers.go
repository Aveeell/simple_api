package main

import (
	"avito/pkg/models"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

func (app *application) getBalance(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "Forbidden method", http.StatusBadRequest)
		return
	}

	type Request struct {
		ID int `json:id`
	}
	var req Request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "can't parse .json", http.StatusInternalServerError)
		return
	}

	data, err := app.db.GetBalance(req.ID)
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

func (app *application) addMoney(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "Forbidden method", http.StatusBadRequest)
		return
	}

	type Request struct {
		ID  int `json:id`
		Sum int `json:sum`
	}

	var req Request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "can't parse .json", http.StatusInternalServerError)
		return
	}

	user, err := app.db.GetBalance(req.ID)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) && req.Sum > 0 {
			app.db.CreateAndAdd(req.ID, req.Sum)
		} else {
			fmt.Fprintf(w, "%s", err)
		}
		return
	} else {
		if user.Balance+req.Sum < 0 {
			http.Error(w, "\nBalance < 0", 400)
			return
		}
		app.db.UpdateBalance(req.ID, user.Balance+req.Sum)
		app.db.AddToHistory(user.ID, req.Sum, "edit balance")
	}
}
