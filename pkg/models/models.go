package models

import (
	"errors"
	"time"
)

var ErrNoRecord = errors.New("models: подходящей записи не найдено")

type Products struct {
	ID    int
	Title string
	Price int
}

type User struct {
	ID       int
	Balance  int
	Reserved int
}

type History struct {
	ID          int
	Transaction int
	Date        time.Time
	Comment     string
}

type Bookkeeping struct {
	Product string
	Total   int
}
