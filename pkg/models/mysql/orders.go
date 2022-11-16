package mysql

func (m *AvitoModel) CreateOrder(balance, sum, id int) (int, error) {
	_, err := m.DB.Exec(`UPDATE users SET balance = ?, reserved = ? WHERE id = ?`, balance, sum, id)
	if err != nil {
		return 0, err
	}
	return 0, nil
}

func (m *AvitoModel) ApproveOrder(id, reserved, sum int, product string) (int, error) {
	_, err := m.DB.Exec(`UPDATE users SET reserved = ? WHERE id = ?`, reserved, id)
	_, err1 := m.DB.Exec(`INSERT INTO bookkeeping (user_id, product, price, date) VALUES (?, ?, ?, UTC_TIMESTAMP())`, id, product, sum)
	if err != nil || err1 != nil {
		return 0, err
	}
	return 0, nil
}
