package main

import (
	"avito/pkg/models"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type Request struct {
	Order_id   int `json: order_id`
	User_id    int `json: user_id`
	Product_id int `json: product_id`
	Sum        int `json: sum`
}

func (app *application) approveOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "Forbidden method", http.StatusBadRequest)
		return
	}

	var req Request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "wrong json configuration", http.StatusBadRequest)
		return
	}
	if req.Order_id < 1 || req.Product_id < 1 || req.User_id < 1 || req.Sum < 1 {
		http.Error(w, "wrong params", http.StatusBadRequest)
		return
	}

	user, err := app.db.GetBalance(req.User_id)
	product, _ := app.db.GetProduct(req.Product_id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) && req.Sum > 0 {
			http.Error(w, "no users with this id", http.StatusNotFound)
		} else {
			fmt.Fprintf(w, "%s", err)
		}
		return
	} else {
		if user.Reserved-req.Sum < 0 {
			http.Error(w, "reserved money can't be less than 0", http.StatusBadRequest)
			return
		}
		app.db.ApproveOrder(req.User_id, user.Reserved-req.Sum, req.Sum, product.Title)
		app.db.AddToHistory(req.User_id, req.Sum, "approved purchase "+product.Title)
	}
}

func (app *application) createOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "Forbidden method", http.StatusBadRequest)
		return
	}

	var req Request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "can't parse .json", http.StatusInternalServerError)
		return
	}
	if req.Order_id < 1 || req.Product_id < 1 || req.User_id < 1 || req.Sum < 1 {
		http.Error(w, "wrong params", http.StatusBadRequest)
		return
	}

	user, err := app.db.GetBalance(req.User_id)
	product, _ := app.db.GetProduct(req.Product_id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) && req.Sum > 0 {
			http.Error(w, "no users with this id", http.StatusNotFound)
		} else {
			fmt.Fprintf(w, "%s", err)
		}
		return
	} else {
		if user.Balance-req.Sum < 0 {
			http.Error(w, "Not enough money on balance", http.StatusBadRequest)
			return
		} else {
			app.db.CreateOrder(user.Balance-req.Sum, user.Reserved+req.Sum, user.ID)
			app.db.AddToHistory(req.User_id, req.Sum, "trying to buy "+product.Title)
		}

	}
}
