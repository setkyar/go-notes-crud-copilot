package main

import (
	"database/sql"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

type Note struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("sqlite3", "./notes.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db.Exec("CREATE TABLE IF NOT EXISTS notes (id INTEGER PRIMARY KEY, title TEXT, content TEXT);")

	r := gin.Default()

	r.GET("/notes", func(c *gin.Context) {
		var notes []Note
		rows, _ := db.Query("SELECT * FROM notes")
		for rows.Next() {
			var note Note
			rows.Scan(&note.ID, &note.Title, &note.Content)
			notes = append(notes, note)
		}
		c.JSON(200, notes)
	})

	r.GET("/notes/:id", func(c *gin.Context) {
		var note Note
		id := c.Param("id")
		row := db.QueryRow("SELECT * FROM notes WHERE id = ?", id)
		row.Scan(&note.ID, &note.Title, &note.Content)
		c.JSON(200, note)
	})

	r.POST("/notes", func(c *gin.Context) {
		var note Note
		c.BindJSON(&note)
		result, _ := db.Exec("INSERT INTO notes (title, content) VALUES (?, ?)", note.Title, note.Content)
		id, _ := result.LastInsertId()
		note.ID = int(id)
		c.JSON(200, note)
	})

	r.PUT("/notes/:id", func(c *gin.Context) {
		var note Note
		c.BindJSON(&note)
		id := c.Param("id")
		db.Exec("UPDATE notes SET title = ?, content = ? WHERE id = ?", note.Title, note.Content, id)
		row := db.QueryRow("SELECT * FROM notes WHERE id = ?", id)
		row.Scan(&note.ID, &note.Title, &note.Content)
		c.JSON(200, note)
	})

	r.DELETE("/notes/:id", func(c *gin.Context) {
		id := c.Param("id")
		db.Exec("DELETE FROM notes WHERE id = ?", id)
		c.JSON(200, gin.H{"message": "Note deleted"})
	})

	r.Run(":8080")
}
