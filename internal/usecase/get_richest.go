package usecase

import (
	"context"
	"sort"

	"github.com/luique16/quitocoin/internal/domain/utxo"
)

type GetRichestUseCase struct {
	utxoService utxo.Service
}

func NewGetRichestUseCase(utxoService utxo.Service) *GetRichestUseCase {
	return &GetRichestUseCase{utxoService: utxoService}
}

type RichestEntry struct {
	PublicID string  `json:"public_id"`
	Balance  float32 `json:"balance"`
}

type GetRichestOutput struct {
	Richest []RichestEntry `json:"richest"`
}

func (uc *GetRichestUseCase) Execute(ctx context.Context) (*GetRichestOutput, error) {
	entries, err := uc.utxoService.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Amount != entries[j].Amount {
			return entries[i].Amount > entries[j].Amount
		}
		return entries[i].UserId < entries[j].UserId
	})

	limit := min(10, len(entries))

	result := make([]RichestEntry, limit)
	for i := range result {
		result[i] = RichestEntry{
			PublicID: entries[i].UserId,
			Balance:  entries[i].Amount,
		}
	}

	return &GetRichestOutput{Richest: result}, nil
}
