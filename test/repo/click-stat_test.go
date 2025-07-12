package repo_test

import (
	"testing"
	"time"
	"url-shortener/internal/database/repo"
	"url-shortener/internal/lib/pg"
	"url-shortener/internal/model"
	"url-shortener/internal/testutils/testdb"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClickStatRepo(t *testing.T) {
	db := testdb.New(t)

	testdb.TruncateTables(t, "users")

	userRepo := repo.NewUserRepo(db)
	urlRepo := repo.NewUrlRepo(db)
	repo := repo.NewClickStatRepo(db)

	// create test user
	user := &model.User{ID: "1234", Email: "alice@example.com", Password: "12345678"}
	require.NoError(t, userRepo.Create(user))
	// create test url
	url := &model.Url{ID: "alias", Link: "https://google.com", UserID: user.ID}
	require.NoError(t, urlRepo.Create(url))

	t.Run("success", func(t *testing.T) {
		clickStat := &model.ClickStat{UrlID: url.ID}
		// Create
		err := repo.Create(clickStat)
		assert.NoError(t, err)

		// clicks for test
		for range 50 {
			err := repo.Create(&model.ClickStat{UrlID: url.ID})
			require.NoError(t, err)
		}

		stats, err := repo.ByUrlID(url.ID, user.ID)
		assert.NoError(t, err)
		today := time.Now().Truncate(24 * time.Hour).UTC()
		assert.Equal(t, today, stats[len(stats)-1].Day)
		assert.Equal(t, int64(51), stats[len(stats)-1].Count)
	})

	t.Run("error", func(t *testing.T) {
		// Create with url id that doesn't exist
		err := repo.Create(&model.ClickStat{UrlID: "notfound"})
		assert.Equal(t, "23503", pg.ParsePGError(err).Code) // 23503 = foreign_key_violation
		assert.Equal(t, "fk_urls_click_stats", pg.ParsePGError(err).ConstraintName)
	})
}
