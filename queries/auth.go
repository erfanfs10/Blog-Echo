package queries

const (
	Register = `
	INSERT INTO users(username,email,password) values(?,?,?)
	`

	Login = `
	SELECT id, is_active,username,password FROM users WHERE username=?
	`

	UpdateLastLogin = `
	UPDATE users SET last_login=? WHERE id=?
	`

	GetEmail = `
	SELECT email FROM users WHERE email=?
	`

	UpdateVerificationCode = `
	UPDATE users SET verification_code=? WHERE email=?
	`

	GetVerificationCode = `
	SELECT verification_code FROM users WHERE email=?
	`

	UpdatePassword = `
	UPDATE users SET password=?, verification_code='' WHERE email=?
	`

	GetIsActiveStatus = `
	SELECT is_active FROM users WHERE id=?
	`
)
