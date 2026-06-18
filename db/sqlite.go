package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	_ "github.com/lib/pq"
	_ "modernc.org/sqlite"
)

type User struct {
	ID         int64
	Name       string
	Phone      string
	Token      string
	Platform   string
	LastTime   string
	LastMsg    string
	LastResult string
}

type DB struct {
	*sql.DB
	dialect string
}

func NewDB(dataSourceName string) (*DB, error) {
	if dsn := firstNonEmptyEnv("DATABASE_URL", "SQL_DSN", "AXONHUB_DB_DSN"); dsn != "" {
		dataSourceName = dsn
	}
	if requireDatabaseURL() && !isPostgresDSN(dataSourceName) {
		return nil, fmt.Errorf("生产环境必须设置 DATABASE_URL 或 SQL_DSN 为 Aiven/PostgreSQL DSN")
	}

	dialect := "sqlite"
	if isPostgresDSN(dataSourceName) {
		dialect = "postgres"
	} else {
		dir := filepath.Dir(dataSourceName)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("创建数据库目录失败: %v", err)
		}
	}

	database, err := sql.Open(dialect, dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("打开数据库失败: %v", err)
	}

	if err := database.Ping(); err != nil {
		return nil, fmt.Errorf("连接数据库失败: %v", err)
	}

	return &DB{DB: database, dialect: dialect}, nil
}

func (db *DB) InitTables() error {
	idColumn := "id INTEGER PRIMARY KEY AUTOINCREMENT"
	if db.dialect == "postgres" {
		idColumn = "id SERIAL PRIMARY KEY"
	}

	query := fmt.Sprintf(`
	CREATE TABLE IF NOT EXISTS users (
		%s,
		name TEXT NOT NULL,
		phone TEXT NOT NULL DEFAULT '',
		token TEXT NOT NULL,
		platform TEXT NOT NULL DEFAULT 'smzdm',
		last_time TEXT NOT NULL DEFAULT '',
		last_msg TEXT NOT NULL DEFAULT '',
		last_result TEXT NOT NULL DEFAULT ''
	);`, idColumn)

	if _, err := db.Exec(query); err != nil {
		return err
	}

	if err := db.addColumnIfMissing("phone TEXT NOT NULL DEFAULT ''"); err != nil {
		return err
	}
	if err := db.addColumnIfMissing("platform TEXT NOT NULL DEFAULT 'smzdm'"); err != nil {
		return err
	}
	if err := db.addColumnIfMissing("last_time TEXT NOT NULL DEFAULT ''"); err != nil {
		return err
	}
	if err := db.addColumnIfMissing("last_msg TEXT NOT NULL DEFAULT ''"); err != nil {
		return err
	}
	if err := db.addColumnIfMissing("last_result TEXT NOT NULL DEFAULT ''"); err != nil {
		return err
	}
	return nil
}

func (db *DB) AddUser(user *User) error {
	if strings.TrimSpace(user.Platform) == "" {
		user.Platform = "smzdm"
	}

	if db.dialect == "postgres" {
		query := `
		INSERT INTO users (name, phone, token, platform, last_time, last_msg, last_result)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id`
		return db.QueryRow(query, user.Name, user.Phone, user.Token, user.Platform, user.LastTime, user.LastMsg, user.LastResult).Scan(&user.ID)
	}

	query := `
	INSERT INTO users (name, phone, token, platform, last_time, last_msg, last_result)
	VALUES (?, ?, ?, ?, ?, ?, ?)`

	result, err := db.Exec(query, user.Name, user.Phone, user.Token, user.Platform, user.LastTime, user.LastMsg, user.LastResult)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	user.ID = id
	return nil
}

func (db *DB) GetAllUsers() ([]User, error) {
	query := `
	SELECT id, name, phone, token, platform, last_time, last_msg, last_result
	FROM users
	ORDER BY id`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Name, &user.Phone, &user.Token, &user.Platform, &user.LastTime, &user.LastMsg, &user.LastResult)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, rows.Err()
}

func (db *DB) UpdateUserCheckResult(id int64, resultCode string, resultMsg string, lastTime string) error {
	if id == 0 {
		return nil
	}
	query := "UPDATE users SET last_time = ?, last_msg = ?, last_result = ? WHERE id = ?"
	args := []interface{}{lastTime, resultMsg, resultCode, id}
	if db.dialect == "postgres" {
		query = "UPDATE users SET last_time = $1, last_msg = $2, last_result = $3 WHERE id = $4"
	}
	_, err := db.Exec(query, args...)
	return err
}

func (db *DB) addColumnIfMissing(columnSQL string) error {
	_, err := db.Exec("ALTER TABLE users ADD COLUMN " + columnSQL)
	if err == nil {
		return nil
	}
	errText := strings.ToLower(err.Error())
	if strings.Contains(errText, "duplicate") || strings.Contains(errText, "exists") || strings.Contains(errText, "duplicate column") {
		return nil
	}
	return err
}

func isPostgresDSN(dsn string) bool {
	dsn = strings.ToLower(strings.TrimSpace(dsn))
	return strings.HasPrefix(dsn, "postgres://") || strings.HasPrefix(dsn, "postgresql://")
}

func firstNonEmptyEnv(keys ...string) string {
	for _, key := range keys {
		if value := strings.TrimSpace(os.Getenv(key)); value != "" {
			return value
		}
	}
	return ""
}

func requireDatabaseURL() bool {
	value := strings.TrimSpace(os.Getenv("REQUIRE_DATABASE_URL"))
	if value == "" {
		return false
	}
	enabled, err := strconv.ParseBool(value)
	return err == nil && enabled
}
