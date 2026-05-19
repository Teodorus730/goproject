package repository

const (
	createUserQuery = `INSERT INTO users(name,email,phone,password_hash) VALUES($1,$2,$3,$4) RETURNING id, name, email, phone, password_hash, created_at, updated_at;`
	getUserByIDQuery = `SELECT id, name, email, phone, password_hash, created_at, updated_at FROM users WHERE id=$1`
	getUserByEmailQuery = `SELECT id, name, email, phone, password_hash, created_at, updated_at FROM users WHERE email=$1`

	createOrderQuery = `INSERT INTO orders(user_id,tariff_id,amount,address,total_price,status) VALUES($1,$2,$3,$4,$5,$6)`
	getOrdersByUserIDQuery = `SELECT * FROM orders WHERE user_id=$1 ORDER BY created_at DESC`


	createSessionQuery = `INSERT INTO sessions (session_id, user_id, created_at, expires_at) VALUES ($1, $2, NOW(), NOW() + INTERVAL '24 hours')`
	getSessionByUserIDQuery = `SELECT session_id, user_id, created_at, expires_at FROM sessions WHERE user_id = $1`
	updateSessionExpiryQuery = `UPDATE sessions SET expires_at = NOW() + INTERVAL '24 hours' WHERE session_id = $1`

	GetTariffByIDQuery = `SELECT * FROM tariffes WHERE id=$1`
)