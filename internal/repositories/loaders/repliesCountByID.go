package loaders

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/vikstrous/dataloadgen"
)

// Это методы, котоыре должен реализовывать передаваемый репозиторий
// В моем случае метод GetRepliesCounts должен быть в реализациях pg и inmemory
type RepliesCountRepo interface {
	GetRepliesCounts(ctx context.Context, ids []uuid.UUID) (map[string]int, error)
}

type RepliesCountReader struct {
	Repo RepliesCountRepo
}

type RepliesCountLoader struct {
	L *dataloadgen.Loader[uuid.UUID, int]
}

func NewRepliesCountLoader(r RepliesCountRepo, maxBatch int) *RepliesCountLoader {
	rcr := &RepliesCountReader{Repo: r}
	return &RepliesCountLoader{
		L: dataloadgen.NewLoader(rcr.BatchFn, dataloadgen.WithWait(time.Millisecond*2), dataloadgen.WithBatchCapacity(maxBatch)),
	}
}

func (rcr *RepliesCountReader) BatchFn(ctx context.Context, ids []uuid.UUID) ([]int, []error) {
	// Будет использоваться на слое сервисов
	// Полдождет пока накопятся id (uuid) в services.RepliesCount
	// Выполнит общий запрос для этих айдишников
	// Метод для этого реализован в обоих репозиториях (pg и inmemory)
	counts, err := rcr.Repo.GetRepliesCounts(ctx, ids)
	if err != nil {
		errs := make([]error, len(ids))
		errs[0] = fmt.Errorf("batch error: %w", err)
		return nil, errs
	}

	results := make([]int, len(ids))
	errs := make([]error, len(ids))
	for i, key := range ids {
		if cnt, ok := counts[key.String()]; ok {
			results[i] = cnt
		} else {
			results[i] = 0
		}
		errs[i] = nil
	}
	fmt.Println(errs)
	return results, errs
}
