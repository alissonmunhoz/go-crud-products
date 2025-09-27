// handler/create_product_test.go
package service

import (
	"bytes"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"

	"github.com/alissonmunhoz/go-crud-products/internal/config"
)

func init() { logger = config.GetLogger("test") }

func setupGin() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/v1/product", CreateProductService)
	return r
}

func newMockGorm(t *testing.T) (*gorm.DB, sqlmock.Sqlmock, *sql.DB) {
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

func TestCreateProductsHandler(t *testing.T) {
	r := setupGin()

	t.Run("retorna 400 quando JSON é inválido (bind error)", func(t *testing.T) {
		body := bytes.NewBufferString(`{"name": "Mouse", "price": 199,`)
		req := httptest.NewRequest(http.MethodPost, "/v1/product", body)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		require.Equal(t, http.StatusBadRequest, w.Code)
		require.Contains(t, w.Body.String(), "invalid request body")
	})

	t.Run("retorna 400 quando faltam campos obrigatórios (validação do bind)", func(t *testing.T) {

		body := bytes.NewBufferString(`{"description":"qualquer"}`)
		req := httptest.NewRequest(http.MethodPost, "/v1/product", body)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		require.Equal(t, http.StatusBadRequest, w.Code)
		require.Contains(t, w.Body.String(), "invalid request body")
	})

	t.Run("retorna 500 quando o insert no DB falha", func(t *testing.T) {
		gdb, mock, sqlDB := newMockGorm(t)
		defer sqlDB.Close()

		originalDB := db
		db = gdb
		defer func() { db = originalDB }()

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `products`")).
			WillReturnError(assertiveDBErr("insert failed"))
		mock.ExpectRollback()

		body := bytes.NewBufferString(`{"name":"Teclado","price":299,"quantity":5,"description":"Mecânico ABNT2"}`)
		req := httptest.NewRequest(http.MethodPost, "/v1/product", body)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		require.Equal(t, http.StatusInternalServerError, w.Code)
		require.Contains(t, w.Body.String(), "error creating product on database")
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("retorna 200 quando cria com sucesso", func(t *testing.T) {
		gdb, mock, sqlDB := newMockGorm(t)
		defer sqlDB.Close()

		originalDB := db
		db = gdb
		defer func() { db = originalDB }()

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `products`")).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		body := bytes.NewBufferString(`{"name":"Monitor","price":1299,"quantity":10,"description":"27\" 144Hz"}`)
		req := httptest.NewRequest(http.MethodPost, "/v1/product", body)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)
		require.Contains(t, w.Body.String(), "operation from handler: create-product successful")
		require.Contains(t, w.Body.String(), `"name":"Monitor"`)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}

type assertiveDBErr string

func (e assertiveDBErr) Error() string { return string(e) }
