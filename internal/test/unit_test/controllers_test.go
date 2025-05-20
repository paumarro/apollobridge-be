package controllers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/paumarro/apollo-be/internal/controllers"
	"github.com/paumarro/apollo-be/internal/dto"
	"github.com/paumarro/apollo-be/internal/initializers"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open GORM DB: %v", err)
	}

	initializers.DB = gormDB
	return gormDB, mock
}

func TestArtworkCreate(t *testing.T) {
	_, mock := setupTestDB(t)

	t.Run("Successful Creation", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "artworks"`)).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), "Test Artwork", "Test Artist", sqlmock.AnyArg(), "Test Description", "http://test.com/image.jpg").
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectCommit()

		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		artwork := dto.ArtworkRequest{
			Title:       "Test Artwork",
			Artist:      "Test Artist",
			Description: "Test Description",
			Image:       "http://test.com/image.jpg",
		}

		c.Set("sanitizedArtwork", artwork)
		controllers.ArtworkCreate(c)

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.Contains(t, w.Body.String(), "Test Artwork")
	})

	t.Run("Missing SanitizedArtwork", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		controllers.ArtworkCreate(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "Failed to retrieve sanitized input")
	})

	t.Run("Database Error", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "artworks"`)).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), "Test Artwork", "Test Artist", sqlmock.AnyArg(), "Test Description", "http://test.com/image.jpg").
			WillReturnError(assert.AnError)
		mock.ExpectRollback()

		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		artwork := dto.ArtworkRequest{
			Title:       "Test Artwork",
			Artist:      "Test Artist",
			Description: "Test Description",
			Image:       "http://test.com/image.jpg",
		}

		c.Set("sanitizedArtwork", artwork)
		controllers.ArtworkCreate(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "Failed to create artwork")
	})
}

func TestArtworkIndex(t *testing.T) {
	_, mock := setupTestDB(t)

	// Edge Case 1: Successful Fetch with Artworks
	t.Run("Successful Fetch with Artworks", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "title", "artist", "date", "description", "image"}).
			AddRow(1, time.Now(), time.Now(), nil, "Artwork 1", "Artist 1", time.Now(), "Description 1", "http://image1.jpg").
			AddRow(2, time.Now(), time.Now(), nil, "Artwork 2", "Artist 2", time.Now(), "Description 2", "http://image2.jpg")

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "artworks" WHERE "artworks"."deleted_at" IS NULL`)).WillReturnRows(rows)

		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		controllers.ArtworkIndex(c)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Artwork 1")
		assert.Contains(t, w.Body.String(), "Artwork 2")
	})

	// Edge Case 2: Empty Database
	t.Run("Empty Database", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "title", "artist", "date", "description", "image"})
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "artworks" WHERE "artworks"."deleted_at" IS NULL`)).WillReturnRows(rows)

		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		controllers.ArtworkIndex(c)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), `"artworks":[]`) // Ensure the response contains an empty array
	})

	// Edge Case 3: Database Error
	t.Run("Database Error", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "artworks" WHERE "artworks"."deleted_at" IS NULL`)).
			WillReturnError(assert.AnError)

		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		controllers.ArtworkIndex(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "Failed to fetch artworks")
	})

	// Edge Case 4: Large Dataset
	t.Run("Large Dataset", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "title", "artist", "date", "description", "image"})
		for i := 1; i <= 1000; i++ {
			rows.AddRow(i, time.Now(), time.Now(), nil, fmt.Sprintf("Artwork %d", i), fmt.Sprintf("Artist %d", i), time.Now(), fmt.Sprintf("Description %d", i), fmt.Sprintf("http://image%d.jpg", i))
		}
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "artworks" WHERE "artworks"."deleted_at" IS NULL`)).WillReturnRows(rows)

		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		controllers.ArtworkIndex(c)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Artwork 1")
		assert.Contains(t, w.Body.String(), "Artwork 1000")
	})

}

func TestArtworkFind(t *testing.T) {
	_, mock := setupTestDB(t)

	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "title", "artist", "date", "description", "image"}).
		AddRow(1, time.Now(), time.Now(), nil, "Artwork 1", "Artist 1", time.Now(), "Description 1", "http://image1.jpg")

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "artworks" WHERE "artworks"."id" = $1 AND "artworks"."deleted_at" IS NULL ORDER BY "artworks"."id" LIMIT $2`)).
		WithArgs("1", 1).
		WillReturnRows(rows)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: "1"}}

	controllers.ArtworkFind(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Artwork 1")
}

func TestArtworkUpdate(t *testing.T) {
	_, mock := setupTestDB(t)

	// Edge Case 1: Successful Update
	t.Run("Successful Update", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "artworks" WHERE "artworks"."id" = $1 AND "artworks"."deleted_at" IS NULL ORDER BY "artworks"."id" LIMIT $2`)).
			WithArgs("1", 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "title", "artist", "date", "description", "image"}).
				AddRow(1, time.Now(), time.Now(), nil, "Old Title", "Old Artist", time.Now(), "Old Description", "http://old-image.jpg"))

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "artworks" SET "created_at"=$1,"updated_at"=$2,"deleted_at"=$3,"title"=$4,"artist"=$5,"date"=$6,"description"=$7,"image"=$8 WHERE "artworks"."deleted_at" IS NULL AND "id" = $9`)).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), "Updated Title", "Updated Artist", sqlmock.AnyArg(), "Updated Description", "http://updated-image.jpg", 1).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: "1"}}

		updatedArtwork := dto.ArtworkRequest{
			Title:       "Updated Title",
			Artist:      "Updated Artist",
			Description: "Updated Description",
			Image:       "http://updated-image.jpg",
		}

		c.Set("sanitizedArtwork", updatedArtwork)

		controllers.ArtworkUpdate(c)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Updated Title")
	})

	// Edge Case 2: Artwork Not Found
	t.Run("Artwork Not Found", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "artworks" WHERE "artworks"."id" = $1 AND "artworks"."deleted_at" IS NULL ORDER BY "artworks"."id" LIMIT $2`)).
			WithArgs("999", 1).
			WillReturnError(gorm.ErrRecordNotFound)

		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: "999"}}

		controllers.ArtworkUpdate(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Contains(t, w.Body.String(), "Artwork not found")
	})

	// Edge Case 3: Missing SanitizedArtwork Context
	t.Run("Missing SanitizedArtwork Context", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "artworks" WHERE "artworks"."id" = $1 AND "artworks"."deleted_at" IS NULL ORDER BY "artworks"."id" LIMIT $2`)).
			WithArgs("1", 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "title", "artist", "date", "description", "image"}).
				AddRow(1, time.Now(), time.Now(), nil, "Old Title", "Old Artist", time.Now(), "Old Description", "http://old-image.jpg"))

		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: "1"}}

		controllers.ArtworkUpdate(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "Failed to retrieve sanitized input")
	})

	// Edge Case 4: Invalid ID Format
	t.Run("Invalid ID Format", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: "abc"}} // Invalid ID format

		controllers.ArtworkUpdate(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "Failed to find artwork")
	})

	// Edge Case 5: Database Error During Update
	t.Run("Database Error During Update", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "artworks" WHERE "artworks"."id" = $1 AND "artworks"."deleted_at" IS NULL ORDER BY "artworks"."id" LIMIT $2`)).
			WithArgs("1", 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "title", "artist", "date", "description", "image"}).
				AddRow(1, time.Now(), time.Now(), nil, "Old Title", "Old Artist", time.Now(), "Old Description", "http://old-image.jpg"))

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "artworks" SET "created_at"=$1,"updated_at"=$2,"deleted_at"=$3,"title"=$4,"artist"=$5,"date"=$6,"description"=$7,"image"=$8 WHERE "artworks"."deleted_at" IS NULL AND "id" = $9`)).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), "Updated Title", "Updated Artist", sqlmock.AnyArg(), "Updated Description", "http://updated-image.jpg", 1).
			WillReturnError(assert.AnError)
		mock.ExpectRollback()

		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: "1"}}

		updatedArtwork := dto.ArtworkRequest{
			Title:       "Updated Title",
			Artist:      "Updated Artist",
			Description: "Updated Description",
			Image:       "http://updated-image.jpg",
		}

		c.Set("sanitizedArtwork", updatedArtwork)

		controllers.ArtworkUpdate(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "Failed to update artwork")
	})
}

func TestArtworkDelete(t *testing.T) {
	_, mock := setupTestDB(t)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "artworks" WHERE "artworks"."id" = $1 AND "artworks"."deleted_at" IS NULL ORDER BY "artworks"."id" LIMIT $2`)).
		WithArgs("1", 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "title", "artist", "date", "description", "image"}).
			AddRow(1, time.Now(), time.Now(), nil, "Artwork 1", "Artist 1", time.Now(), "Description 1", "http://image1.jpg"))

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "artworks" SET "deleted_at"=$1 WHERE "artworks"."id" = $2 AND "artworks"."deleted_at" IS NULL`)).
		WithArgs(sqlmock.AnyArg(), 1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: "1"}}

	controllers.ArtworkDelete(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Artwork successfully deleted")
}
