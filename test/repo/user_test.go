package repo_test

import (
	"testing"
	"url-shortener/internal/database/repo"
	"url-shortener/internal/lib/pg"
	"url-shortener/internal/model"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func cleanUsersTable(db *gorm.DB) {
	db.Exec("DELETE FROM users")
}

func TestUserRepo(t *testing.T) {
	db := setupTestDB(t)

	cleanUsersTable(db)

	repo := repo.NewUserRepo(db)

	t.Run("success", func(t *testing.T) {
		user := &model.User{ID: "abc123", Email: "alice@example.com", Password: "12345"}
		// Create
		err := repo.Create(user)
		assert.NoError(t, err)
		assert.NotZero(t, user.ID)

		// ByEmail
		found, err := repo.ByEmail("alice@example.com")
		assert.NoError(t, err)
		assert.Equal(t, "12345", found.Password)
		assert.Equal(t, user.ID, found.ID)

		// Update
		user.Password = "1234567"
		err = repo.Update(user)
		assert.NoError(t, err)

		// ById
		found, err = repo.ById("abc123")
		assert.NoError(t, err)
		assert.Equal(t, "1234567", found.Password)
		assert.Equal(t, user.Email, found.Email)

		// Update with zero values
		err = repo.Update(&model.User{ID: "abc123", Email: "", Password: "newpass"})
		assert.NoError(t, err)
		found, err = repo.ById("abc123")
		assert.NoError(t, err)
		assert.Equal(t, "newpass", found.Password)
		assert.Equal(t, user.Email, found.Email)

		// Delete
		err = repo.Delete("abc123")
		assert.NoError(t, err)
		_, err = repo.ById("abc123")
		assert.ErrorIs(t, err, gorm.ErrRecordNotFound)

		// TODO: ContextById
		// TODO: ContextByEmail
	})

	t.Run("error", func(t *testing.T) {
		// Create
		err := repo.Create(&model.User{ID: "abc123", Email: "alice@example.com", Password: "12345"})
		assert.NoError(t, err)

		err = repo.Create(&model.User{ID: "abc123", Email: "another@example.com", Password: "12345"})
		assert.Equal(t, pg.ParsePGError(err).Code, "23505")

		err = repo.Create(&model.User{ID: "another", Email: "alice@example.com", Password: "12345"})
		assert.Equal(t, pg.ParsePGError(err).Code, "23505")

		// ByEmail
		_, err = repo.ByEmail("notfound@example.com")
		assert.ErrorIs(t, err, gorm.ErrRecordNotFound)

		// ById
		_, err = repo.ById("notfound")
		assert.ErrorIs(t, err, gorm.ErrRecordNotFound)

		// Update
		err = repo.Update(&model.User{ID: "notfound", Email: "test@example.com", Password: "1234567"})
		assert.ErrorIs(t, err, gorm.ErrRecordNotFound)

		// TODO: ContextById
		// TODO: ContextByEmail
	})
}
