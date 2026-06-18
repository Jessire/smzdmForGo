package db

import (
	"path/filepath"
	"testing"

	"ggball.com/smzdm/file"
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

func TestProductConfigPersistence(t *testing.T) {
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

	lowCommentNum := 8
	maxPrice := 199.9
	conf := file.Config{
		KeyWords:      []string{"显示器", "面包"},
		FilterWords:   []string{"过期"},
		LowCommentNum: 3,
		LowWorthyNum:  6,
		MinPrice:      10,
		MaxPrice:      500,
		SatisfyNum:    5,
		KeywordRules: []file.KeywordRule{
			{
				Words:         []string{"小米"},
				FilterWords:   []string{"二手"},
				LowCommentNum: &lowCommentNum,
				MaxPrice:      &maxPrice,
			},
		},
	}
	if err := database.SaveProductConfig(conf); err != nil {
		t.Fatalf("SaveProductConfig() error = %v", err)
	}

	got, err := database.GetProductConfig(file.Config{})
	if err != nil {
		t.Fatalf("GetProductConfig() error = %v", err)
	}
	if len(got.KeyWords) != 2 || got.KeyWords[0] != "显示器" {
		t.Fatalf("KeyWords = %#v", got.KeyWords)
	}
	if len(got.KeywordRules) != 1 || got.KeywordRules[0].LowCommentNum == nil || *got.KeywordRules[0].LowCommentNum != 8 {
		t.Fatalf("KeywordRules = %#v", got.KeywordRules)
	}
	if got.KeywordRules[0].MaxPrice == nil || *got.KeywordRules[0].MaxPrice != maxPrice {
		t.Fatalf("MaxPrice = %#v", got.KeywordRules[0].MaxPrice)
	}
}
