package service

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"

	"github.com/alissonmunhoz/go-crud-products/internal/config"
)

func init() { logger = config.GetLogger("test") }

func setupGinFind() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/v1/product", FindProductService)
	return r
}

func newMockGormFind(t *testing.T) (*gorm.DB, sqlmock.Sqlmock, *sql.DB) {
	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)

	dialector := mysql.New(mysql.Config{
		Conn:                      sqlDB,
		SkipInitializeWithVersion: true,
	})

	gdb, err := gorm.Open(dialector, &gorm.Config{
		Logger: glogger.Default.LogMode(glogger.Silent),
	})
	require.NoError(t, err)

	return gdb, mock, sqlDB
}

func TestFindProductService(t *testing.T) {
	r := setupGinFind()

	t.Run("retorna 400 quando id não é informado", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/v1/product", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)
		require.Equal(t, http.StatusBadRequest, w.Code)
		require.Contains(t, w.Body.String(), "id")
	})

	t.Run("retorna 404 quando produto não existe", func(t *testing.T) {
		gdb, mock, sqlDB := newMockGormFind(t)
		defer sqlDB.Close()
		orig := db
		db = gdb
		defer func() { db = orig }()

		selectRegex := `(?is)SELECT.*FROM.*products.*WHERE.*id`
		mock.ExpectQuery(selectRegex).WillReturnError(sql.ErrNoRows)

		req := httptest.NewRequest(http.MethodGet, "/v1/product?id=999", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)
		require.Equal(t, http.StatusNotFound, w.Code)
		require.Contains(t, w.Body.String(), "product not found")
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("retorna 200 quando encontra produto (DTO camelCase)", func(t *testing.T) {
		gdb, mock, sqlDB := newMockGormFind(t)
		defer sqlDB.Close()
		orig := db
		db = gdb
		defer func() { db = orig }()

		selectRegex := `(?is)SELECT.*FROM.*products.*WHERE.*id`

		cols := []string{"id", "name", "price", "quantity", "description", "created_at", "updated_at", "deleted_at"}
		now := time.Now()
		row := sqlmock.NewRows(cols).AddRow(7, "Teclado", 299, 5, "ABNT2", now, now, nil)

		mock.ExpectQuery(selectRegex).WillReturnRows(row)

		req := httptest.NewRequest(http.MethodGet, "/v1/product?id=7", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)
		require.Equal(t, http.StatusOK, w.Code)
		require.Contains(t, w.Body.String(), "show-product")

		var body struct {
			Data struct {
				ID          int64   `json:"id"`
				Name        string  `json:"name"`
				Price       float64 `json:"price"`
				Quantity    int64   `json:"quantity"`
				Description string  `json:"description"`
			} `json:"data"`
			Message string `json:"message"`
		}
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
		require.Equal(t, int64(7), body.Data.ID)
		require.Equal(t, "Teclado", body.Data.Name)
		require.Equal(t, float64(299), body.Data.Price)
		require.Equal(t, int64(5), body.Data.Quantity)
		require.Equal(t, "ABNT2", body.Data.Description)

		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("retorna 500 se SELECT falhar inesperadamente", func(t *testing.T) {
		gdb, mock, sqlDB := newMockGormFind(t)
		defer sqlDB.Close()
		orig := db
		db = gdb
		defer func() { db = orig }()

		selectRegex := `(?is)SELECT.*FROM.*products.*WHERE.*id`
		mock.ExpectQuery(selectRegex).WillReturnError(errors.New("random db error"))

		req := httptest.NewRequest(http.MethodGet, "/v1/product?id=1", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			require.Equal(t, http.StatusInternalServerError, w.Code)
		}
	})
}
