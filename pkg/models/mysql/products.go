package mysql

import (
	"avito/pkg/models"
	"database/sql"
	"errors"
)

func (m *AvitoModel) GetProduct(id int) (*models.Products, error) {
	row := m.DB.QueryRow(`SELECT * FROM products WHERE id = ?`, id)
	s := &models.Products{}

	err := row.Scan(&s.ID, &s.Title, &s.Price)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}

	return s, nil
}
