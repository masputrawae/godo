package repo

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"godo/internal/model"

	"github.com/DATA-DOG/go-sqlmock"
)

func openDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	t.Helper()

	db, mock, err := sqlmock.New(
		sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual),
	)
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}

	t.Cleanup(func() {
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet SQL expectations: %v", err)
		}

		db.Close()
	})

	return db, mock
}

func userRows(now time.Time) *sqlmock.Rows {
	return sqlmock.
		NewRows([]string{"id", "email", "username", "password", "created_at", "updated_at"}).
		AddRow(1, "abc@def.ghi", "abc-user", "abc-pass", now, now)
}

func assertUser(t *testing.T, got, want *model.User) {
	t.Helper()
	if got.ID != want.ID {
		t.Errorf("ID: want %d, got %d", want.ID, got.ID)
	}
	if got.Email != want.Email {
		t.Errorf("Email: want %s, got %s", want.Email, got.Email)
	}
	if got.Username != want.Username {
		t.Errorf("Username: want %s, got %s", want.Username, got.Username)
	}
	if got.Password != want.Password {
		t.Errorf("Password: want %s, got %s", want.Password, got.Password)
	}
	if !got.CreatedAt.Equal(want.CreatedAt) {
		t.Errorf("CreatedAt: want %v, got %v", want.CreatedAt, got.CreatedAt)
	}
	if !got.UpdatedAt.Equal(want.UpdatedAt) {
		t.Errorf("UpdatedAt: want %v, got %v", want.UpdatedAt, got.UpdatedAt)
	}
}

func TestRepoUserCreate(t *testing.T) {
	ctx := context.Background()
	db, mock := openDB(t)
	repo := NewUser(db)
	req := model.UserRequestRegister{
		Email:    "abc@def.ghi",
		Username: "abc-user",
		Password: "abc-pass",
	}
	now := time.Now()
	query := `
		INSERT INTO users (email,username,password)
		VALUES (?,?,?)
		RETURNING id,email,username,password,created_at,updated_at
	`
	mock.ExpectQuery(query).
		WithArgs(
			"abc@def.ghi",
			"abc-user",
			"abc-pass",
		).
		WillReturnRows(userRows(now))

	got, err := repo.Create(ctx, req)
	if err != nil {
		t.Fatalf("want no error, got %v", err)
	}

	want := &model.User{
		ID:        1,
		Email:     "abc@def.ghi",
		Username:  "abc-user",
		Password:  "abc-pass",
		CreatedAt: now,
		UpdatedAt: now,
	}

	assertUser(t, got, want)
}

func TestRepoUserCreateEmailAlreadyUsed(t *testing.T) {
	ctx := context.Background()
	db, mock := openDB(t)
	repo := NewUser(db)

	req := model.UserRequestRegister{
		Email:    "abc@def.ghi",
		Username: "abc-user",
		Password: "abc-pass",
	}

	query := `
		INSERT INTO users (email,username,password)
		VALUES (?,?,?)
		RETURNING id,email,username,password,created_at,updated_at
	`

	mock.
		ExpectQuery(query).
		WithArgs(
			"abc@def.ghi",
			"abc-user",
			"abc-pass",
		).
		WillReturnError(errors.New(
			"UNIQUE constraint failed: users.email",
		))

	_, err := repo.Create(ctx, req)
	if !errors.Is(err, ErrEmailAlreadyUsed) {
		t.Fatalf("want ErrEmailAlreadyUsed, got %v", err)
	}
}

func TestRepoUserFindByID(t *testing.T) {
	ctx := context.Background()
	db, mock := openDB(t)

	repo := NewUser(db)

	now := time.Now()

	query := `
		SELECT id, email, username, password, created_at, updated_at
		FROM users
		WHERE id = ?
	`

	mock.ExpectQuery(query).
		WithArgs(1).
		WillReturnRows(userRows(now))

	got, err := repo.FindByID(ctx, 1)
	if err != nil {
		t.Fatalf("want no error, got %v", err)
	}

	want := &model.User{
		ID:        1,
		Email:     "abc@def.ghi",
		Username:  "abc-user",
		Password:  "abc-pass",
		CreatedAt: now,
		UpdatedAt: now,
	}
	assertUser(t, got, want)
}

func TestRepoUserFindByUsername(t *testing.T) {
	ctx := context.Background()
	db, mock := openDB(t)
	repo := NewUser(db)
	now := time.Now()
	query := `
		SELECT id, email, username, password, created_at, updated_at
		FROM users
		WHERE username = ?
	`
	mock.ExpectQuery(query).
		WithArgs("abc-user").
		WillReturnRows(userRows(now))

	got, err := repo.FindByUsername(ctx, "abc-user")
	if err != nil {
		t.Fatalf("want no error, got %v", err)
	}

	want := &model.User{
		ID:        1,
		Email:     "abc@def.ghi",
		Username:  "abc-user",
		Password:  "abc-pass",
		CreatedAt: now,
		UpdatedAt: now,
	}

	assertUser(t, got, want)
}

func TestRepoUserDeleteByID(t *testing.T) {
	ctx := context.Background()
	db, mock := openDB(t)
	repo := NewUser(db)

	query := `DELETE FROM users WHERE id = ?`

	mock.
		ExpectExec(query).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.DeleteByID(ctx, 1)
	if err != nil {
		t.Fatalf("want no error, got %v", err)
	}
}

func TestRepoUserUpdatePartials(t *testing.T) {
	ctx := context.Background()
	db, mock := openDB(t)
	repo := NewUser(db)
	now := time.Now()

	payload := model.UserRequestUpdate{
		Username: new("jono"),
	}

	query := `
		UPDATE users 
		SET 
			updated_at = CURRENT_TIMESTAMP,
			username = ?
		WHERE id = ?
		RETURNING id,email,username,password,created_at,updated_at`

	mock.
		ExpectQuery(query).
		WithArgs("jono", 1).
		WillReturnRows(sqlmock.
			NewRows([]string{"id", "email", "username", "password", "created_at", "updated_at"}).
			AddRow(1, "abc@def.ghi", "jono", "abc-pass", now, now),
		)

	got, err := repo.Update(ctx, 1, payload)
	if err != nil {
		t.Fatalf("want no error got %s", err)
	}

	want := &model.User{
		ID:        1,
		Email:     "abc@def.ghi",
		Username:  "jono",
		Password:  "abc-pass",
		CreatedAt: now,
		UpdatedAt: now,
	}

	assertUser(t, got, want)
}

func TestRepoUserUpdateAll(t *testing.T) {
	ctx := context.Background()
	db, mock := openDB(t)
	repo := NewUser(db)
	now := time.Now()
	payload := model.UserRequestUpdate{
		Username: new("jono"),
		Password: new("password"),
		Email:    new("email@update.com"),
	}
	query := `
		UPDATE users 
		SET 
			updated_at = CURRENT_TIMESTAMP,
			password = ?,
			username = ?,
			email = ?
		WHERE id = ?
		RETURNING id,email,username,password,created_at,updated_at`

	mock.
		ExpectQuery(query).
		WithArgs("password", "jono", "email@update.com", 1).
		WillReturnRows(sqlmock.
			NewRows([]string{"id", "email", "username", "password", "created_at", "updated_at"}).
			AddRow(1, "email@update.com", "jono", "password", now, now),
		)

	got, err := repo.Update(ctx, 1, payload)
	if err != nil {
		t.Fatalf("want no error got %s", err)
	}

	want := &model.User{
		ID:        1,
		Email:     "email@update.com",
		Username:  "jono",
		Password:  "password",
		CreatedAt: now,
		UpdatedAt: now,
	}

	assertUser(t, got, want)
}
