package repo_test

import (
	"testing"
	"url-shortener/internal/database/repo"
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

		// Delete
		err = repo.Delete("abc123")
		_, err = repo.ById("abc123")
		assert.Error(t, err)

		// TODO: ContextById
		// TODO: ContextByEmail
	})

	t.Run("error", func(t *testing.T) {
		// Create
		err := repo.Create(&model.User{ID: "abc123", Email: "alice@example.com", Password: "12345"})
		assert.NoError(t, err)

		err = repo.Create(&model.User{ID: "abc123", Email: "another@example.com", Password: "12345"})
		assert.ErrorIs(t, err, gorm.ErrDuplicatedKey)

		err = repo.Create(&model.User{ID: "another", Email: "alice@example.com", Password: "12345"})
		assert.ErrorIs(t, err, gorm.ErrDuplicatedKey)

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
