package accrual

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/dglazkoff/go-musthave-diploma-tpl/internal/config"
	"github.com/dglazkoff/go-musthave-diploma-tpl/internal/logger"
	"github.com/dglazkoff/go-musthave-diploma-tpl/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

type storageMock struct {
	mock.Mock
}

func (s *storageMock) UpdateOrderTx(ctx context.Context, tx *sql.Tx, order models.Order) (models.Order, error) {
	args := s.Called(ctx, tx, order)
	return args.Get(0).(models.Order), args.Error(1)
}
func (s *storageMock) BeginTx(ctx context.Context, options *sql.TxOptions) (*sql.Tx, error) {
	args := s.Called(ctx, options)
	return args.Get(0).(*sql.Tx), args.Error(1)
}
func (s *storageMock) GetNotAccrualOrders(ctx context.Context) ([]models.Order, error) {
	return []models.Order{}, nil
}

type externalServiceMock struct{}

func (s *externalServiceMock) UpdateBalanceTx(ctx context.Context, tx *sql.Tx, accrual float64, userLogin string) error {
	return nil
}

func TestService_GetAccrual(t *testing.T) {
	err := logger.Initialize()
	require.NoError(t, err)
	order := models.Order{ID: "123", UserID: "123", UploadedAt: "2021"}
	calculatedOrder := models.Order{ID: "123", UserID: "123", UploadedAt: "2021", Status: models.Processed, Accrual: 100}
	// processingOrder := models.Order{ID: "123", UserID: "123", UploadedAt: "2021", Status: models.Processing, Accrual: 0}

	var handler http.HandlerFunc

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler.ServeHTTP(w, r)
	}))
	defer ts.Close()

	err = os.Setenv("ACCRUAL_SYSTEM_ADDRESS", ts.URL)
	require.NoError(t, err)

	cfg := config.ParseConfig()

	setHandler := func(h http.HandlerFunc) {
		handler = h
	}

	t.Run("success test", func(t *testing.T) {
		setHandler(func(w http.ResponseWriter, r *http.Request) {
			err := json.NewEncoder(w).Encode(AccrualSystemResponse{Order: "123", Status: Processed, Accrual: 100})

			if err != nil {
				logger.Log.Error("Error while encode response: ", err)
				return
			}

			w.WriteHeader(http.StatusOK)
		})

		db, mockSQL, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mockSQL.ExpectBegin()
		mockSQL.ExpectCommit()
		tx, err := db.Begin()
		assert.NoError(t, err)

		storage := storageMock{}
		storage.On("BeginTx", mock.Anything, mock.Anything).Return(tx, nil)
		storage.On("UpdateOrderTx", mock.Anything, tx, calculatedOrder).Return(calculatedOrder, nil)

		externalService := externalServiceMock{}
		accrualService := New(&storage, &externalService, &cfg)

		accrualService.GetAccrual(order, "123")

		storage.AssertCalled(t, "UpdateOrderTx", mock.Anything, tx, calculatedOrder)
		err = mockSQL.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	//t.Run("retry-after test", func(t *testing.T) {
	//	callCount := 0
	//	setHandler(func(w http.ResponseWriter, r *http.Request) {
	//		fmt.Println(callCount)
	//		if callCount == 0 {
	//			w.Header().Set("Retry-After", "5")
	//			w.WriteHeader(http.StatusTooManyRequests)
	//		}
	//
	//		if callCount == 1 {
	//			err := json.NewEncoder(w).Encode(AccrualSystemResponse{Order: "123", Status: Processed, Accrual: 100})
	//
	//			if err != nil {
	//				logger.Log.Error("Error while encode response: ", err)
	//				return
	//			}
	//
	//			w.WriteHeader(http.StatusOK)
	//		}
	//
	//		callCount++
	//	})
	//
	//	db, mockSQL, err := sqlmock.New()
	//	assert.NoError(t, err)
	//	defer db.Close()
	//
	//	tx, err := db.Begin()
	//
	//	mockSQL.ExpectBegin()
	//	mockSQL.ExpectCommit()
	//
	//	assert.NoError(t, err)
	//
	//	storage := storageMock{}
	//	storage.On("BeginTx", mock.Anything, mock.Anything).Return(tx, nil)
	//	storage.On("UpdateOrderTx", mock.Anything, tx, calculatedOrder).Return(calculatedOrder, nil)
	//
	//	externalService := externalServiceMock{}
	//	accrualService := New(&storage, &externalService, &cfg)
	//
	//	accrualService.GetAccrual(order, "123")
	//
	//	storage.AssertNotCalled(t, "UpdateOrderTx", mock.Anything, tx, calculatedOrder)
	//
	//	storage.AssertCalled(t, "UpdateOrderTx", mock.Anything, tx, calculatedOrder)
	//	err = mockSQL.ExpectationsWereMet()
	//	assert.NoError(t, err)
	//
	//	assert.Equal(t, 2, callCount)
	//})
}
