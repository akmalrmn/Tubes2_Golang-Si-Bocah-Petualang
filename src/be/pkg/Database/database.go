package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

// Database represents a database connection
type Database struct {
	db *sql.DB
}

// NewDatabase creates a new database connection
func NewDatabase() *Database {
	db, err := sql.Open("sqlite3", "./wiki.db")
	if err != nil {
		log.Fatal(err)
	}
	return &Database{db}
}

// Close closes the database connection
func (d *Database) Close() {
	d.db.Close()
}

// GetPath returns the path between two nodes
func (d *Database) GetPath(start, end string) []string {
	var path []string
	query := fmt.Sprintf("SELECT value FROM paths WHERE start = %s AND end = %s", start, end)
	rows, err := d.db.Query(query)
	if err != nil {
		log.Println(err)
		return path
	}
	defer rows.Close()
	for rows.Next() {
		var value string
		err = rows.Scan(&value)
		if err != nil {
			log.Println(err)
			return path
		}
		path = append(path, value)
	}
	return path
}

func (d *Database) IsRelationsValidById(start, end int) bool {
	query := fmt.Sprintf("SELECT id FROM paths WHERE start = %d AND end = %d", start, end)
	row := d.db.QueryRow(query)
	var id int
	err := row.Scan(&id)
	if err != nil {
		return false
	}
	return true
}

// SavePath saves the path between two nodes
func (d *Database) SavePath(start, end string) {
	query := fmt.Sprintf("INSERT INTO paths (start, end) VALUES (%s, %s)", start, end)
	_, err := d.db.Exec(query)
	if err != nil {
		log.Println(err)
	}
}

// GetNode returns the node with the given value
func (d *Database) GetNode(value string) int {
	var id int
	query := fmt.Sprintf("SELECT id FROM nodes WHERE value = %s", value)
	err := d.db.QueryRow(query).Scan(&id)
	if err != nil {
		log.Println(err)
	}
	return id
}

func (d *Database) GetNodeById(id int) string {
	var value string
	query := fmt.Sprintf("SELECT value FROM nodes WHERE id = %d", id)
	err := d.db.QueryRow(query).Scan(&value)
	if err != nil {
		log.Println(err)
	}
	return value
}

// SaveNode saves the node with the given value
func (d *Database) SaveNode(value string) {
	query := fmt.Sprintf("INSERT INTO nodes (value) VALUES (%s)", value)
	_, err := d.db.Exec(query)
	if err != nil {
		log.Println(err)
	}
}

// isNodeExists checks if the node with the given value exists
func (d *Database) isNodeExists(value string) bool {
	query := fmt.Sprintf("SELECT id FROM nodes WHERE value = %s", value)
	row := d.db.QueryRow(query)
	var id int
	err := row.Scan(&id)
	if err != nil {
		return false
	}
	return true
}

// CreateTable creates the paths & nodes table
func (d *Database) CreateTable() {
	// Create Paths table
	query := `CREATE TABLE IF NOT EXISTS paths (
		start int,
		end int,
		Unique(start, end),
	)`
	_, err := d.db.Exec(query)
	if err != nil {
		log.Println(err)
	}

	// Create node table
	query = `CREATE TABLE IF NOT EXISTS nodes (
    		id INTEGER PRIMARY KEY AUTOINCREMENT,
    		value TEXT,
    		unique(id, value)
        )`
	_, err = d.db.Exec(query)
	if err != nil {
		log.Println(err)
	}
}

// DropTable drops the paths table
func (d *Database) DropTable() {
	query := "DROP TABLE IF EXISTS paths"
	_, err := d.db.Exec(query)
	if err != nil {
		log.Println(err)
	}

	query = "DROP TABLE IF EXISTS nodes"
	_, err = d.db.Exec(query)
	if err != nil {
		log.Println(err)
	}
}

// Migrate migrates the database
func (d *Database) Migrate() {
	d.DropTable()
	d.CreateTable()
}
