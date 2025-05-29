package inmemory

import (
	"go.uber.org/zap"
	"quotes_service/internal/models"
	"testing"
)

var (
	logger = zap.NewNop()
)

func BenchmarkAddQuote(b *testing.B) {
	store := NewInMemory(logger)

	quote := models.Quote{
		Author: "Benchmark Author",
		Quote:  "This is a quote for benchmarking.",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		store.AddQuote(quote)
	}
}

func BenchmarkGetAllQuotes_NoFilter(b *testing.B) {
	store := NewInMemory(logger)

	for i := 0; i < 1000; i++ {
		store.Data[i+1] = models.Quote{
			ID:     i + 1,
			Author: "Author",
			Quote:  "Quote text",
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = store.GetAllQuotes("")
	}
}

func BenchmarkGetAllQuotes_FilterByAuthor(b *testing.B) {
	store := NewInMemory(logger)

	for i := 0; i < 500; i++ {
		store.Data[i+1] = models.Quote{
			ID:     i + 1,
			Author: "Matched",
			Quote:  "Quote text",
		}
	}
	for i := 0; i < 500; i++ {
		store.Data[i+501] = models.Quote{
			ID:     i + 501,
			Author: "NotMatched",
			Quote:  "Quote text",
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = store.GetAllQuotes("Matched")
	}
}

func BenchmarkGetRandomQuote(b *testing.B) {
	store := NewInMemory(logger)

	for i := 0; i < 1000; i++ {
		store.Data[i+1] = models.Quote{
			ID:     i + 1,
			Author: "Author",
			Quote:  "Quote text",
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = store.GetRandomQuote()
	}
}

func BenchmarkDeleteQuote(b *testing.B) {
	store := NewInMemory(logger)

	store.Data[1] = models.Quote{
		ID:     1,
		Author: "ToDelete",
		Quote:  "To delete this quote.",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = store.DeleteQuote(1)
	}
}
