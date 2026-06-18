package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"

	"ggball.com/smzdm/file"
)

const productConfigKey = "product_config"

func (db *DB) InitSettingsTable() error {
	query := fmt.Sprintf(`
	CREATE TABLE IF NOT EXISTS %s (
		key TEXT PRIMARY KEY,
		value TEXT NOT NULL,
		updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
	);`, db.settingsTableName)
	_, err := db.Exec(query)
	return err
}

func (db *DB) GetProductConfig(fallback file.Config) (file.Config, error) {
	query := "SELECT value FROM " + db.settingsTableName + " WHERE key = ?"
	args := []interface{}{productConfigKey}
	if db.dialect == "postgres" {
		query = "SELECT value FROM " + db.settingsTableName + " WHERE key = $1"
	}

	var raw string
	err := db.QueryRow(query, args...).Scan(&raw)
	if err == sql.ErrNoRows {
		return fallback, nil
	}
	if err != nil {
		return fallback, err
	}

	conf := fallback
	if err := json.Unmarshal([]byte(raw), &conf); err != nil {
		return fallback, err
	}
	file.ApplyEnvOverrides(&conf)
	return conf, nil
}

func (db *DB) SaveProductConfig(conf file.Config) error {
	raw, err := json.Marshal(conf)
	if err != nil {
		return err
	}

	if db.dialect == "postgres" {
		query := "INSERT INTO " + db.settingsTableName + " (key, value, updated_at) VALUES ($1, $2, CURRENT_TIMESTAMP) ON CONFLICT (key) DO UPDATE SET value = EXCLUDED.value, updated_at = CURRENT_TIMESTAMP"
		_, err = db.Exec(query, productConfigKey, string(raw))
		return err
	}

	query := "INSERT INTO " + db.settingsTableName + " (key, value, updated_at) VALUES (?, ?, CURRENT_TIMESTAMP) ON CONFLICT(key) DO UPDATE SET value = excluded.value, updated_at = CURRENT_TIMESTAMP"
	_, err = db.Exec(query, productConfigKey, string(raw))
	return err
}

func settingsTableName() (string, error) {
	tableName := strings.TrimSpace(os.Getenv("SMZDM_SETTINGS_TABLE"))
	if tableName == "" {
		tableName = "smzdm_app_settings"
	}
	if !regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`).MatchString(tableName) {
		return "", fmt.Errorf("非法配置表名: %s", tableName)
	}
	return tableName, nil
}
