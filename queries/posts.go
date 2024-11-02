package queries

const (
	PostMy = `
	SELECT 
    	id, 
    	title, 
    	user_id,
		body, 
    	created, 
    	updated
	FROM posts
	WHERE posts.user_id = ?
	ORDER BY created DESC
	`

	PostList = `
	SELECT 
    	posts.id, 
    	posts.title, 
    	posts.body, 
    	posts.created, 
    	posts.updated,
    	users.id AS user_id,
    	users.username AS username,
    	users.email AS email,
		users.avatar AS avatar
	FROM posts
	LEFT JOIN users ON posts.user_id = users.id
	ORDER BY posts.created DESC
	`

	PostCreate = `
	INSERT INTO posts(title,body,user_id) values(?,?,?)
	`

	PostGet = `
	SELECT id,title,body,user_id,created,updated FROM posts WHERE id=?
	`

	PostGetUpdate = `
	SELECT id,user_id,title,body FROM posts WHERE id=?
	`
	PostUpdate = `
	UPDATE posts SET title=?,body=? WHERE id=?
	`

	PostGetDelete = `
	SELECT user_id FROM posts WHERE id=?
	`

	PostDelete = `
	DELETE FROM posts WHERE id=?
	`
)
