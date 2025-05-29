package handlers

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"net/http"
	"quotes_service/internal/models"
	"quotes_service/internal/models/errs"
	"strconv"
	"time"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

type InMemory interface {
	AddQuote(quote models.Quote) models.Quote
	GetAllQuotes(author string) ([]models.Quote, error)
	GetRandomQuote() (models.Quote, error)
	DeleteQuote(id int) error
}
type Handler struct {
	Storage InMemory
	logger  *zap.Logger
}

func NewHandler(storage InMemory, logger *zap.Logger) *Handler {
	return &Handler{
		Storage: storage,
		logger:  logger,
	}
}

func (h *Handler) CreateQuote(w http.ResponseWriter, r *http.Request) {
	var q models.Quote
	err := json.NewDecoder(r.Body).Decode(&q)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if q.Author == "" || q.Quote == "" {
		http.Error(w, "Missing 'author' or 'quote'", http.StatusBadRequest)
		return
	}

	q.CreatedAt = time.Now()

	saved := h.Storage.AddQuote(q)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(saved)
}

func (h *Handler) GetAllQuotesHandler(w http.ResponseWriter, r *http.Request) {
	author := r.URL.Query().Get("author")
	result, err := h.Storage.GetAllQuotes(author)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		if errors.Is(err, errs.ErrQuotesIsEmpty) {
			result = []models.Quote{}
			json.NewEncoder(w).Encode(result)
			return
		}
	}

	json.NewEncoder(w).Encode(result)
}

func (h *Handler) GetRandomQuoteHandler(w http.ResponseWriter, r *http.Request) {
	q, err := h.Storage.GetRandomQuote()
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		if errors.Is(err, errs.ErrQuotesIsEmpty) {
			w.WriteHeader(http.StatusNotFound)

			json.NewEncoder(w).Encode(ErrorResponse{"No quotes available"})
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(q)
}

func (h *Handler) DeleteQuoteHandler(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.DeleteQuoteHandler"
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.logger.Error("Deleting quote", zap.String("op", op), zap.Error(err))
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if err := h.Storage.DeleteQuote(id); err != nil {
		if errors.Is(err, errs.ErrQuoteNotFound) {
			http.Error(w, "Quote not found", http.StatusNotFound)
			return
		}
		h.logger.Error("Deleting quote", zap.String("op", op), zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
}
