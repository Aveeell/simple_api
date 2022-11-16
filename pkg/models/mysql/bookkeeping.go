package mysql

import (
	"avito/pkg/models"
	"database/sql"
	"errors"
)

func (m *AvitoModel) GetDocument(date string) ([]*models.Bookkeeping, error) {
	date += "%"
	command := `SELECT product, sum(price) AS total FROM bookkeeping WHERE date LIKE ? GROUP BY product;`
	rows, _ := m.DB.Query(command, date)
	defer rows.Close()
	var bookkeeping []*models.Bookkeeping

	for rows.Next() {
		str := &models.Bookkeeping{}
		err1 := rows.Scan(&str.Product, &str.Total)
		if err1 != nil {
			if errors.Is(err1, sql.ErrNoRows) {
				return nil, models.ErrNoRecord
			} else {
				return nil, err1
			}
		}
		bookkeeping = append(bookkeeping, str)
	}

	return bookkeeping, nil
}
