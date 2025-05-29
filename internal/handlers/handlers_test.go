package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
	"net/http"
	"net/http/httptest"
	"quotes_service/internal/models"
	"quotes_service/internal/models/errs"
	"testing"
)

type mockStorage struct {
	InMemory
	addQuoteFunc       func(quote models.Quote) models.Quote
	getAllQuotesFunc   func(author string) ([]models.Quote, error)
	getRandomQuoteFunc func() (models.Quote, error)
	deleteQuoteFunc    func(id int) error
}

func (m *mockStorage) AddQuote(quote models.Quote) models.Quote {
	return m.addQuoteFunc(quote)
}

func (m *mockStorage) GetAllQuotes(author string) ([]models.Quote, error) {
	return m.getAllQuotesFunc(author)
}

func (m *mockStorage) GetRandomQuote() (models.Quote, error) {
	return m.getRandomQuoteFunc()
}

func (m *mockStorage) DeleteQuote(id int) error {
	return m.deleteQuoteFunc(id)
}

func setupHandler(storage InMemory) (*Handler, *observer.ObservedLogs) {
	loggerCore, observed := observer.New(zap.InfoLevel)
	logger := zap.New(loggerCore)
	handler := NewHandler(storage, logger)
	return handler, observed
}

func TestCreateQuote_Success(t *testing.T) {
	mock := &mockStorage{
		addQuoteFunc: func(quote models.Quote) models.Quote {
			quote.ID = 1
			return quote
		},
	}

	handler, _ := setupHandler(mock)

	reqBody := `{"author":"Leo","quote":"Test quote"}`
	req := httptest.NewRequest(http.MethodPost, "/quotes", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.CreateQuote(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var res models.Quote
	_ = json.NewDecoder(w.Body).Decode(&res)
	assert.Equal(t, "Leo", res.Author)
	assert.Equal(t, "Test quote", res.Quote)
	assert.NotZero(t, res.CreatedAt)
}

func TestCreateQuote_MissingFields(t *testing.T) {
	handler, _ := setupHandler(&mockStorage{})

	reqBody := `{"author":""}`
	req := httptest.NewRequest(http.MethodPost, "/quotes", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.CreateQuote(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetAllQuotesHandler_Success(t *testing.T) {
	mock := &mockStorage{
		getAllQuotesFunc: func(author string) ([]models.Quote, error) {
			return []models.Quote{
				{ID: 1, Author: "A", Quote: "Q1"},
				{ID: 2, Author: "B", Quote: "Q2"},
			}, nil
		},
	}

	handler, _ := setupHandler(mock)

	req := httptest.NewRequest(http.MethodGet, "/quotes", nil)
	w := httptest.NewRecorder()

	handler.GetAllQuotesHandler(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var res []models.Quote
	_ = json.NewDecoder(w.Body).Decode(&res)
	assert.Len(t, res, 2)
}

func TestGetAllQuotesHandler_EmptyList(t *testing.T) {
	mock := &mockStorage{
		getAllQuotesFunc: func(author string) ([]models.Quote, error) {
			return nil, errs.ErrQuotesIsEmpty
		},
	}

	handler, _ := setupHandler(mock)

	req := httptest.NewRequest(http.MethodGet, "/quotes", nil)
	w := httptest.NewRecorder()

	handler.GetAllQuotesHandler(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var res []models.Quote
	_ = json.NewDecoder(w.Body).Decode(&res)
	assert.Empty(t, res)
}

func TestGetRandomQuoteHandler_Success(t *testing.T) {
	mock := &mockStorage{
		getRandomQuoteFunc: func() (models.Quote, error) {
			return models.Quote{ID: 1, Author: "A", Quote: "Q"}, nil
		},
	}

	handler, _ := setupHandler(mock)

	req := httptest.NewRequest(http.MethodGet, "/quotes/random", nil)
	w := httptest.NewRecorder()

	handler.GetRandomQuoteHandler(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var res models.Quote
	_ = json.NewDecoder(w.Body).Decode(&res)
	assert.Equal(t, "A", res.Author)
	assert.Equal(t, "Q", res.Quote)
}

func TestGetRandomQuoteHandler_NoQuotes(t *testing.T) {
	mock := &mockStorage{
		getRandomQuoteFunc: func() (models.Quote, error) {
			return models.Quote{}, errs.ErrQuotesIsEmpty
		},
	}

	handler, _ := setupHandler(mock)

	req := httptest.NewRequest(http.MethodGet, "/quotes/random", nil)
	w := httptest.NewRecorder()

	handler.GetRandomQuoteHandler(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	var body map[string]string
	_ = json.NewDecoder(w.Body).Decode(&body)
	assert.Equal(t, "No quotes available", body["error"])
}

func TestDeleteQuoteHandler_Success(t *testing.T) {
	mock := &mockStorage{
		deleteQuoteFunc: func(id int) error {
			return nil
		},
	}

	handler, _ := setupHandler(mock)

	req := httptest.NewRequest(http.MethodDelete, "/quotes/1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	w := httptest.NewRecorder()

	handler.DeleteQuoteHandler(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeleteQuoteHandler_NotFound(t *testing.T) {
	mock := &mockStorage{
		deleteQuoteFunc: func(id int) error {
			return errs.ErrQuoteNotFound
		},
	}

	handler, _ := setupHandler(mock)

	req := httptest.NewRequest(http.MethodDelete, "/quotes/999", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "999"})
	w := httptest.NewRecorder()

	handler.DeleteQuoteHandler(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}
