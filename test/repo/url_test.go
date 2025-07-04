package repo

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
	repo := repo.NewUrlRepo(db)

	// create test user
	user := &model.User{ID: "1234", Email: "alice@example.com", Password: "12345678"}
	err := userRepo.Create(user)
	require.NoError(t, err)

	t.Run("success", func(t *testing.T) {
		url := &model.Url{ID: "alias", Link: "https://google.com", UserID: user.ID}
		// Create
		err := repo.Create(url)
		assert.NoError(t, err)

		// ByID
		found, err := repo.ByID("alias")
		assert.NoError(t, err)
		assert.Equal(t, url.Link, found.Link)
		assert.Equal(t, url.UserID, found.UserID)

		// Delete
		err = repo.Delete("alias", "1234")
		assert.NoError(t, err)
		_, err = repo.ByID("alias")
		assert.ErrorIs(t, gorm.ErrRecordNotFound, err)

		// create urls for test
		testUrls := []model.Url{}
		for i := range 10 {
			url := model.Url{ID: "alias" + strconv.Itoa(i), Link: "https://google.com", UserID: user.ID}
			testUrls = append(testUrls, url)
			err := repo.Create(&url)
			assert.NoError(t, err)
		}
		// ByUserID
		urls, err := repo.ByUserID(user.ID, 5, 0)
		assert.NoError(t, err)
		assert.Equal(t, testUrls[:5], urls)

		urls, err = repo.ByUserID(user.ID, 5, 5)
		assert.NoError(t, err)
		assert.Equal(t, testUrls[5:], urls)
	})

	t.Run("error", func(t *testing.T) {
		// Create
		err := repo.Create(&model.Url{ID: "alias", Link: "https://google.com", UserID: user.ID})
		assert.NoError(t, err)

		// Create duplicate
		err = repo.Create(&model.Url{ID: "alias", Link: "https://google.com", UserID: user.ID})
		assert.Equal(t, "23505", pg.ParsePGError(err).Code)

		// ByID
		_, err = repo.ByID("notfound")
		assert.ErrorIs(t, gorm.ErrRecordNotFound, err)
	})
}
