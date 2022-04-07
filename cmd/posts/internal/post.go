package internal

import (
	"blog-example/internal/platform/database"
	"database/sql"
)

type Post struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
}

type PostsList []Post

// Insert save new post
func (p *Post) Insert(mysql database.DBConnectionInterface) error {
	err := mysql.Execute(func(sql *sql.DB) error {
		_, err := sql.Exec("INSERT INTO posts (title, description) VALUES (?, ?)", p.Title, p.Description)

		return err
	})

	return err
}

func GetAllPosts(mysql database.DBConnectionInterface) (*PostsList, error) {
	var list PostsList

	err := mysql.Execute(func(sql *sql.DB) error {
		rows, err := sql.Query("SELECT title, description FROM posts")

		if err != nil {
			return err
		}

		for rows.Next() {
			var p = Post{}

			if err = rows.Scan(&p.Title, &p.Description); err != nil {
				return err
			}

			list = append(list, p)
		}

		return nil
	})

	return &list, err
}
