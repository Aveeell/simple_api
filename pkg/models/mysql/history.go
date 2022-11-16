package mysql

import (
	"avito/pkg/models"
	"database/sql"
	"errors"
	"strconv"
)

func (m *AvitoModel) AddToHistory(user_id, transaction int, comment string) (int, error) {
	_, err := m.DB.Exec(`INSERT INTO history (user_id, transaction, date, comment) VALUES (?, ?, UTC_TIMESTAMP(), ?)`, user_id, transaction, comment)
	if err != nil {
		return 0, err
	}
	return 0, nil
}

func (m *AvitoModel) GetHistory(id, limit, offset int, sort_type, sort_col string) ([]*models.History, error) {
	command := `SELECT id, transaction, date, comment FROM history WHERE user_id = ?`
	if len(sort_col) != 0 {
		command += ` ORDER BY ` + sort_col
		if len(sort_type) != 0 {
			command += ` ` + sort_type
		}
	}
	if limit != 0 {
		command += ` LIMIT ` + strconv.Itoa(limit)
	}
	if offset != 0 {
		command += ` OFFSET ` + strconv.Itoa(offset)
	}

	rows, _ := m.DB.Query(command, id)
	defer rows.Close()
	var history []*models.History

	for rows.Next() {
		str := &models.History{}
		err1 := rows.Scan(&str.ID, &str.Transaction, &str.Date, &str.Comment)
		if err1 != nil {
			if errors.Is(err1, sql.ErrNoRows) {
				return nil, models.ErrNoRecord
			} else {
				return nil, err1
			}
		}
		history = append(history, str)
	}

	return history, nil
}
