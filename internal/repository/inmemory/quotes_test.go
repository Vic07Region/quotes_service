package inmemory

import (
	"go.uber.org/zap"
	"quotes_service/internal/models"
	"quotes_service/internal/models/errs"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setup(t *testing.T) *InMemory {
	logger := zap.NewNop()
	store := NewInMemory(logger)
	return store
}

func TestAddQuote_ShouldAssignIDAndSave(t *testing.T) {
	store := setup(t)

	q := models.Quote{
		Author: "Leo Tolstoy",
		Quote:  "All happy families are alike…",
	}

	result := store.AddQuote(q)

	assert.Equal(t, 1, result.ID)
	assert.Equal(t, q.Author, result.Author)
	assert.Equal(t, q.Quote, result.Quote)
	assert.Contains(t, store.Data, 1)
}

func TestGetAllQuotes_NoFilter_ReturnsAll(t *testing.T) {
	store := setup(t)

	store.Data = map[int]models.Quote{
		1: {ID: 1, Author: "A", Quote: "Q1"},
		2: {ID: 2, Author: "B", Quote: "Q2"},
		3: {ID: 3, Author: "A", Quote: "Q3"},
	}

	quotes, err := store.GetAllQuotes("")

	require.NoError(t, err)
	assert.Len(t, quotes, 3)
}

func TestGetAllQuotes_FilterByAuthor_ReturnsMatching(t *testing.T) {
	store := setup(t)

	store.Data = map[int]models.Quote{
		1: {ID: 1, Author: "Alice", Quote: "Q1"},
		2: {ID: 2, Author: "Bob", Quote: "Q2"},
		3: {ID: 3, Author: "Alice", Quote: "Q3"},
	}

	quotes, err := store.GetAllQuotes("Alice")

	require.NoError(t, err)
	assert.Len(t, quotes, 2)
	for _, q := range quotes {
		assert.Equal(t, "Alice", q.Author)
	}
}

func TestGetAllQuotes_EmptyStore_ReturnsError(t *testing.T) {
	store := setup(t)

	quotes, err := store.GetAllQuotes("")

	assert.Nil(t, quotes)
	assert.ErrorIs(t, err, errs.ErrQuotesIsEmpty)
}

func TestGetRandomQuote_SuccessfullyReturnsOne(t *testing.T) {
	store := setup(t)

	store.Data = map[int]models.Quote{
		1: {ID: 1, Author: "A", Quote: "Q1"},
		2: {ID: 2, Author: "B", Quote: "Q2"},
	}

	q, err := store.GetRandomQuote()

	require.NoError(t, err)
	assert.Contains(t, []int{1, 2}, q.ID)
}

func TestGetRandomQuote_EmptyStore_ReturnsError(t *testing.T) {
	store := setup(t)

	q, err := store.GetRandomQuote()

	assert.Empty(t, q)
	assert.ErrorIs(t, err, errs.ErrQuotesIsEmpty)
}

func TestDeleteQuote_ExistingID_RemovesQuote(t *testing.T) {
	store := setup(t)

	store.Data = map[int]models.Quote{
		1: {ID: 1, Author: "A", Quote: "Q1"},
	}

	err := store.DeleteQuote(1)

	assert.NoError(t, err)
	assert.NotContains(t, store.Data, 1)
}

func TestDeleteQuote_NonExistingID_ReturnsError(t *testing.T) {
	store := setup(t)

	err := store.DeleteQuote(999)

	assert.ErrorIs(t, err, errs.ErrQuoteNotFound)
}
