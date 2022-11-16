package mysql

import (
	"avito/pkg/models"
	"database/sql"
	"errors"
)

func (m *AvitoModel) UpdateBalance(id, sum int) (int, error) {
	_, err := m.DB.Exec(`UPDATE users SET balance = ? WHERE id = ?`, sum, id)
	if err != nil {
		return 0, err
	}
	return 0, nil
}

func (m *AvitoModel) CreateAndAdd(id, sum int) (int, error) {
	_, err := m.DB.Exec(`INSERT INTO users (id, balance, reserved) VALUES (?, ?, 0)`, id, sum)
	m.AddToHistory(id, sum, "user creation")
	if err != nil {
		return 0, err
	}
	return 0, nil
}

func (m *AvitoModel) GetBalance(id int) (*models.User, error) {
	row := m.DB.QueryRow(`SELECT * FROM users WHERE id = ?`, id)
	user := &models.User{}

	err := row.Scan(&user.ID, &user.Balance, &user.Reserved)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}

	return user, nil
}
