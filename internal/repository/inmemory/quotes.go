package inmemory

import (
	"go.uber.org/zap"
	"math/rand"
	"quotes_service/internal/models"
	"quotes_service/internal/models/errs"
	"sync"
)

type InMemory struct {
	Data   map[int]models.Quote
	logger *zap.Logger
	mutex  sync.RWMutex
	nextID int
}

func NewInMemory(logger *zap.Logger) *InMemory {
	return &InMemory{
		Data:   make(map[int]models.Quote),
		logger: logger,
		mutex:  sync.RWMutex{},
		nextID: 1,
	}
}

func (i *InMemory) AddQuote(quote models.Quote) models.Quote {
	const op = "inmemory.AddQuote"
	i.mutex.Lock()
	defer i.mutex.Unlock()
	defer i.logger.Debug("Unlocked mutex", zap.String("op", op))
	i.logger.Debug(
		"adding quote",
		zap.String("op", op),
		zap.Int("current id", i.nextID),
		zap.String("author", quote.Author),
		zap.String("text", quote.Quote),
		zap.String("mutex", "locked"),
	)
	quote.ID = i.nextID
	i.Data[i.nextID] = quote
	i.logger.Info("Quote Added",
		zap.String("op", op),
		zap.Int("quoteid", i.nextID))
	i.nextID++

	return quote
}

func (i *InMemory) GetAllQuotes(author string) ([]models.Quote, error) {
	const op = "inmemory.GetAllQuote"
	i.mutex.RLock()
	defer i.mutex.RUnlock()
	defer i.logger.Debug("Unlocked mutex", zap.String("op", op))
	i.logger.Debug(
		"getting all quote",
		zap.String("op", op),
		zap.String("filter by author value", author),
		zap.String("mutex", "locked"),
	)

	var quotes []models.Quote
	if len(i.Data) == 0 {
		i.logger.Warn("main quotes list is empty",
			zap.String("op", op),
			zap.Error(errs.ErrQuotesIsEmpty))
		return nil, errs.ErrQuotesIsEmpty
	}
	for _, q := range i.Data {
		if q.Author == author || author == "" {
			quotes = append(quotes, q)
		}
	}
	if len(quotes) == 0 {
		i.logger.Warn("filtred quotes list is empty",
			zap.String("op", op),
			zap.String("author filter value", author),
			zap.Error(errs.ErrQuotesIsEmpty))
		return nil, errs.ErrQuotesIsEmpty
	}
	i.logger.Info("getting quotes", zap.String("op", op), zap.Int("quotes count", len(quotes)))
	return quotes, nil
}

func (i *InMemory) GetRandomQuote() (models.Quote, error) {
	const op = "inmemory.GetRandomQuote"
	i.logger.Debug("getting random quote", zap.String("op", op))
	all, err := i.GetAllQuotes("")
	if err != nil {
		return models.Quote{}, err
	}
	i.logger.Info("getting random quote", zap.String("op", op))
	return all[rand.Intn(len(all))], nil
}

func (i *InMemory) DeleteQuote(id int) error {
	const op = "inmemory.DeleteQuote"
	i.mutex.Lock()
	defer i.mutex.Unlock()
	defer i.logger.Debug("Unlocked mutex", zap.String("op", op))
	i.logger.Debug(
		"deleting quote",
		zap.String("op", op),
		zap.Int("quote id", id),
		zap.String("mutex", "locked"),
	)
	if _, ok := i.Data[id]; !ok {
		i.logger.Error("quote not found",
			zap.String("op", op),
			zap.Int("quote id", id),
			zap.Error(errs.ErrQuoteNotFound))
		return errs.ErrQuoteNotFound
	}
	delete(i.Data, id)
	i.logger.Info("quote deleted", zap.String("op", op), zap.Int("quote id", id))
	return nil
}
