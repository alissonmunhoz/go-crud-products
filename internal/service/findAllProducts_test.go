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

func setupGinFindAll() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/v1/products", FindAllProductsService)
	return r
}

func newMockGormFindAll(t *testing.T) (*gorm.DB, sqlmock.Sqlmock, *sql.DB) {
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

func TestFindAllProductsHandler(t *testing.T) {
	r := setupGinFindAll()

	t.Run("retorna 500 quando DB falha", func(t *testing.T) {
		gdb, mock, sqlDB := newMockGormFindAll(t)
		defer sqlDB.Close()
		orig := db
		db = gdb
		defer func() { db = orig }()

		selectRegex := `(?is)SELECT.*FROM.*products.*`
		mock.ExpectQuery(selectRegex).
			WillReturnError(errors.New("db down"))

		req := httptest.NewRequest(http.MethodGet, "/v1/products", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		require.Equal(t, http.StatusInternalServerError, w.Code)
		require.Contains(t, w.Body.String(), "error listing products")
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("retorna 200 com lista de produtos (DTO camelCase)", func(t *testing.T) {
		gdb, mock, sqlDB := newMockGormFindAll(t)
		defer sqlDB.Close()
		orig := db
		db = gdb
		defer func() { db = orig }()

		selectRegex := `(?is)SELECT.*FROM.*products.*`

		cols := []string{"id", "name", "price", "quantity", "description", "created_at", "updated_at", "deleted_at"}
		now := time.Now()
		rows := sqlmock.NewRows(cols).
			AddRow(1, "Mouse", 199, 3, "Sem fio", now, now, nil).
			AddRow(2, "Teclado", 299, 5, "ABNT2", now, now, nil)

		mock.ExpectQuery(selectRegex).WillReturnRows(rows)

		req := httptest.NewRequest(http.MethodGet, "/v1/products", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)
		require.Contains(t, w.Body.String(), "list-products successful")

		var body struct {
			Data []struct {
				ID          int64   `json:"id"`
				Name        string  `json:"name"`
				Price       float64 `json:"price"`
				Quantity    int64   `json:"quantity"`
				Description string  `json:"description"`
			} `json:"data"`
			Message string `json:"message"`
		}
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
		require.Len(t, body.Data, 2)

		require.Equal(t, int64(1), body.Data[0].ID)
		require.Equal(t, "Mouse", body.Data[0].Name)
		require.Equal(t, float64(199), body.Data[0].Price)
		require.Equal(t, int64(3), body.Data[0].Quantity)
		require.Equal(t, "Sem fio", body.Data[0].Description)

		require.Equal(t, int64(2), body.Data[1].ID)
		require.Equal(t, "Teclado", body.Data[1].Name)
		require.Equal(t, float64(299), body.Data[1].Price)
		require.Equal(t, int64(5), body.Data[1].Quantity)
		require.Equal(t, "ABNT2", body.Data[1].Description)

		require.NoError(t, mock.ExpectationsWereMet())
	})
}
