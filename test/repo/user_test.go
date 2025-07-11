package repo_test

import (
	"strconv"
	"testing"
	"url-shortener/internal/database/repo"
	"url-shortener/internal/lib/pg"
	"url-shortener/internal/model"
	"url-shortener/internal/testutils/testdb"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestUserRepo(t *testing.T) {
	db := testdb.New(t)

	testdb.TruncateTables(t, "users", "urls")

	userRepo := repo.NewUserRepo(db)
	urlRepo := repo.NewUrlRepo(db)

	t.Run("success", func(t *testing.T) {
		user := &model.User{ID: "abc123", Email: "alice@example.com", Password: "12345"}

		// Create
		err := userRepo.Create(user)
		assert.NoError(t, err)
		assert.NotZero(t, user.ID)

		// urls for test
		urls := make([]model.Url, 20)
		for i := range 20 {
			urls[i] = model.Url{ID: "a" + strconv.Itoa(i), Link: "https://google.com", UserID: user.ID}
			err := urlRepo.Create(&urls[i])
			require.NoError(t, err)
		}

		// ByEmail
		found, err := userRepo.ByEmail("alice@example.com")
		assert.NoError(t, err)
		assert.Equal(t, "12345", found.Password)
		assert.Equal(t, user.ID, found.ID)

		// Update
		user.Password = "1234567"
		err = userRepo.Update(user)
		assert.NoError(t, err)

		// ById
		found, err = userRepo.ById("abc123")
		assert.NoError(t, err)
		assert.Equal(t, "1234567", found.Password)
		assert.Equal(t, user.Email, found.Email)

		// Update with zero values
		err = userRepo.Update(&model.User{ID: "abc123", Email: "", Password: "newpass"})
		assert.NoError(t, err)
		found, err = userRepo.ById("abc123")
		assert.NoError(t, err)
		assert.Equal(t, "newpass", found.Password)
		assert.Equal(t, user.Email, found.Email)

		// ContextById
		found, err = userRepo.ContextById("abc123")
		assert.NoError(t, err)
		assert.Equal(t, "newpass", found.Password)
		assert.Equal(t, user.Email, found.Email)
		assert.Equal(t, urls[0:16], found.Urls)

		// ContextByEmail
		found, err = userRepo.ContextByEmail("alice@example.com")
		assert.NoError(t, err)
		assert.Equal(t, "newpass", found.Password)
		assert.Equal(t, user.Email, found.Email)
		assert.Equal(t, urls[0:16], found.Urls)

		// Delete
		err = userRepo.Delete("abc123")
		assert.NoError(t, err)
		_, err = userRepo.ById("abc123")
		assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
	})

	t.Run("error", func(t *testing.T) {
		// Create
		err := userRepo.Create(&model.User{ID: "abc123", Email: "alice@example.com", Password: "12345"})
		assert.NoError(t, err)

		err = userRepo.Create(&model.User{ID: "abc123", Email: "another@example.com", Password: "12345"})
		assert.Equal(t, pg.ParsePGError(err).Code, "23505")

		err = userRepo.Create(&model.User{ID: "another", Email: "alice@example.com", Password: "12345"})
		assert.Equal(t, pg.ParsePGError(err).Code, "23505")

		// ByEmail
		_, err = userRepo.ByEmail("notfound@example.com")
		assert.ErrorIs(t, err, gorm.ErrRecordNotFound)

		// ById
		_, err = userRepo.ById("notfound")
		assert.ErrorIs(t, err, gorm.ErrRecordNotFound)

		// Update
		err = userRepo.Update(&model.User{ID: "notfound", Email: "test@example.com", Password: "1234567"})
		assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
	})
}
