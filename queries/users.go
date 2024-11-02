package queries

const (
	UserList = `
	SELECT id, username,email FROM users LIMIT ? OFFSET ?
	`

	UserSearch = `
	SELECT 
		users.id,
		username,
		email,
		avatar,
		COUNT(posts.id) AS post_count
	FROM users 
	LEFT JOIN posts ON users.id = posts.user_id
	WHERE username LIKE CONCAT('%',?,'%')
	GROUP BY users.id
	`

	UserMy = `
	SELECT id,username,is_active,is_admin,email,avatar,last_login,created,updated FROM users WHERE id=?
	`
)
