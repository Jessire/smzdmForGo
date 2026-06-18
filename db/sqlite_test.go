package db

import (
	"path/filepath"
	"testing"
)

func TestUserPersistenceWithSQLiteFallback(t *testing.T) {
	t.Setenv("DATABASE_URL", "")
	t.Setenv("SQL_DSN", "")
	t.Setenv("AXONHUB_DB_DSN", "")
	t.Setenv("REQUIRE_DATABASE_URL", "")

	database, err := NewDB(filepath.Join(t.TempDir(), "users.db"))
	if err != nil {
		t.Fatalf("NewDB() error = %v", err)
	}
	defer database.Close()

	if err := database.InitTables(); err != nil {
		t.Fatalf("InitTables() error = %v", err)
	}

	user := &User{Name: "Jessire", Token: "cookie-value", Platform: "smzdm"}
	if err := database.AddUser(user); err != nil {
		t.Fatalf("AddUser() error = %v", err)
	}
	if user.ID == 0 {
		t.Fatal("AddUser() did not populate user ID")
	}

	if err := database.UpdateUserCheckResult(user.ID, "success", "ok", "2026-06-18 12:00:00"); err != nil {
		t.Fatalf("UpdateUserCheckResult() error = %v", err)
	}

	users, err := database.GetAllUsers()
	if err != nil {
		t.Fatalf("GetAllUsers() error = %v", err)
	}
	if len(users) != 1 {
		t.Fatalf("GetAllUsers() len = %d, want 1", len(users))
	}
	if users[0].Name != "Jessire" || users[0].Token != "cookie-value" || users[0].LastResult != "success" {
		t.Fatalf("persisted user = %+v", users[0])
	}
}

func TestRequireDatabaseURLRejectsSQLiteFallback(t *testing.T) {
	t.Setenv("DATABASE_URL", "")
	t.Setenv("SQL_DSN", "")
	t.Setenv("AXONHUB_DB_DSN", "")
	t.Setenv("REQUIRE_DATABASE_URL", "true")

	if _, err := NewDB(filepath.Join(t.TempDir(), "users.db")); err == nil {
		t.Fatal("NewDB() error = nil, want error when REQUIRE_DATABASE_URL=true and no Postgres DSN is set")
	}
}
