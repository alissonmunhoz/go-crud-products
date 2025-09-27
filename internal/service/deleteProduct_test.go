package service

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
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

const useSoftDelete = true

func init() { logger = config.GetLogger("test") }

func setupGinDelete() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.DELETE("/v1/product", DeleteProductService)
	return r
}

func newMockGormDelete(t *testing.T) (*gorm.DB, sqlmock.Sqlmock, *sql.DB) {
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

func TestDeleteProductHandler(t *testing.T) {
	r := setupGinDelete()

	t.Run("retorna 400 quando id não é informado", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/v1/product", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)
		require.Equal(t, http.StatusBadRequest, w.Code)
		require.Contains(t, w.Body.String(), "id")
	})

	t.Run("retorna 404 quando produto não existe", func(t *testing.T) {
		gdb, mock, sqlDB := newMockGormDelete(t)
		defer sqlDB.Close()
		orig := db
		db = gdb
		defer func() { db = orig }()

		id := "123"

		// SELECT ultra-tolerante (case-insensitive + dotall)
		selectRegex := `(?is)SELECT.*FROM.*products.*WHERE.*id`
		// força "não encontrado"
		mock.ExpectQuery(selectRegex).WillReturnError(sql.ErrNoRows)

		u := url.URL{Path: "/v1/product"}
		q := u.Query()
		q.Set("id", id)
		u.RawQuery = q.Encode()

		req := httptest.NewRequest(http.MethodDelete, u.String(), nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)
		require.Equal(t, http.StatusNotFound, w.Code)
		require.Contains(t, w.Body.String(), "not found")
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("retorna 500 quando Delete falha no DB", func(t *testing.T) {
		gdb, mock, sqlDB := newMockGormDelete(t)
		defer sqlDB.Close()
		orig := db
		db = gdb
		defer func() { db = orig }()

		id := "42"
		selectRegex := `(?is)SELECT.*FROM.*products.*WHERE.*id`

		cols := []string{"id", "name", "price", "quantity", "description", "created_at", "updated_at", "deleted_at"}
		now := time.Now()
		row := sqlmock.NewRows(cols).AddRow(42, "Mouse", 199, 3, "sem fio", now, now, nil)

		mock.ExpectQuery(selectRegex).WillReturnRows(row)

		mock.ExpectBegin()
		if useSoftDelete {

			updateRegex := `(?is)UPDATE.*products.*SET.*deleted_at.*WHERE.*id`
			mock.ExpectExec(updateRegex).
				WithArgs(sqlmock.AnyArg(), 42).
				WillReturnError(errors.New("delete failed"))
		} else {

			deleteRegex := `(?is)DELETE.*FROM.*products.*WHERE.*id`
			mock.ExpectExec(deleteRegex).
				WithArgs(42).
				WillReturnError(errors.New("delete failed"))
		}
		mock.ExpectRollback()

		u := url.URL{Path: "/v1/product"}
		q := u.Query()
		q.Set("id", id)
		u.RawQuery = q.Encode()

		req := httptest.NewRequest(http.MethodDelete, u.String(), nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)
		require.Equal(t, http.StatusInternalServerError, w.Code)
		require.Contains(t, w.Body.String(), "error deleting product")
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("retorna 200 quando deleta com sucesso", func(t *testing.T) {
		gdb, mock, sqlDB := newMockGormDelete(t)
		defer sqlDB.Close()
		orig := db
		db = gdb
		defer func() { db = orig }()

		id := "7"
		selectRegex := `(?is)SELECT.*FROM.*products.*WHERE.*id`

		cols := []string{"id", "name", "price", "quantity", "description", "created_at", "updated_at", "deleted_at"}
		now := time.Now()
		row := sqlmock.NewRows(cols).AddRow(7, "Teclado", 299, 5, "ABNT2", now, now, nil)

		mock.ExpectQuery(selectRegex).WillReturnRows(row)

		mock.ExpectBegin()
		if useSoftDelete {
			updateRegex := `(?is)UPDATE.*products.*SET.*deleted_at.*WHERE.*id`
			mock.ExpectExec(updateRegex).
				WithArgs(sqlmock.AnyArg(), 7).
				WillReturnResult(sqlmock.NewResult(0, 1))
		} else {
			deleteRegex := `(?is)DELETE.*FROM.*products.*WHERE.*id`
			mock.ExpectExec(deleteRegex).
				WithArgs(7).
				WillReturnResult(sqlmock.NewResult(0, 1))
		}
		mock.ExpectCommit()

		u := url.URL{Path: "/v1/product"}
		q := u.Query()
		q.Set("id", id)
		u.RawQuery = q.Encode()

		req := httptest.NewRequest(http.MethodDelete, u.String(), nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)
		require.Equal(t, http.StatusOK, w.Code)
		require.Contains(t, w.Body.String(), "delete-product")

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
		require.Equal(t, "Teclado", body.Data.Name)
		require.Equal(t, float64(299), body.Data.Price)
		require.Equal(t, int64(5), body.Data.Quantity)
		require.Equal(t, "ABNT2", body.Data.Description)

		require.NoError(t, mock.ExpectationsWereMet())
	})
}
