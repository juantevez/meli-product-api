package json

import (
	"context"
	"encoding/json"
	"meli-product-api/internal/domain/model"
	"os"
	"sync"
	"time"
)

type QuestionRepository struct {
	mu        sync.RWMutex
	questions []model.Question
	filePath  string
}

func NewQuestionRepository(filePath string) (*QuestionRepository, error) {
	repo := &QuestionRepository{
		filePath:  filePath,
		questions: make([]model.Question, 0),
	}

	if err := repo.load(); err != nil {
		return nil, err
	}

	return repo, nil
}

func (r *QuestionRepository) load() error {
	data, err := os.ReadFile(r.filePath)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, &r.questions)
}

func (r *QuestionRepository) GetByProductID(ctx context.Context, productID string, limit int) ([]model.Question, error) {
	// Simulate network latency
	time.Sleep(18 * time.Millisecond)

	r.mu.RLock()
	defer r.mu.RUnlock()

	var results []model.Question
	count := 0

	for _, question := range r.questions {
		if question.ProductID == productID {
			results = append(results, question)
			count++
			if count >= limit {
				break
			}
		}
	}

	return results, nil
}
