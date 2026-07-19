package usecase

import (
	"context"
	"fmt"

	"github.com/luique16/quitocoin/internal/domain/block"
	"github.com/luique16/quitocoin/internal/domain/userblock"
)

type GetUserBlocksUseCase struct {
	blockService     block.Service
	userBlockService userblock.Service
}

func NewGetUserBlocksUseCase(blockService block.Service, userBlockService userblock.Service) *GetUserBlocksUseCase {
	return &GetUserBlocksUseCase{
		blockService:     blockService,
		userBlockService: userBlockService,
	}
}

type UserBlockRecord struct {
	Type       string  `json:"type"`
	Value      float64 `json:"value"`
	Date       string  `json:"date"`
	OtherParty string  `json:"other_party"`
	BlockIndex int     `json:"block_index"`
}

type GetUserBlocksOutput struct {
	Blocks []UserBlockRecord `json:"blocks"`
}

type GetUserBlocksInput struct {
	PublicID string
	Role     string
	Limit    int
}

func (uc *GetUserBlocksUseCase) Execute(ctx context.Context, input GetUserBlocksInput) (*GetUserBlocksOutput, error) {
	refs, err := uc.userBlockService.GetBlocks(ctx, input.PublicID, input.Role, input.Limit)
	if err != nil {
		return nil, fmt.Errorf("get blocks: %w", err)
	}

	records := make([]UserBlockRecord, 0, len(refs))
	for _, ref := range refs {
		b, err := uc.blockService.GetBlockByIndex(ctx, ref.Index)
		if err != nil {
			return nil, fmt.Errorf("get block %d: %w", ref.Index, err)
		}

		rec := UserBlockRecord{
			Type:       ref.Role,
			BlockIndex: ref.Index,
			Date:       b.CreatedAt.Format("2006-01-02 15:04:05"),
		}

		switch ref.Role {
		case "miner":
			rec.Value = b.Reward
			rec.OtherParty = ""
		case "sender":
			for _, tx := range b.Transactions {
				if tx.From == input.PublicID {
					rec.Value = float64(tx.Amount + 1)
					rec.OtherParty = tx.To
					break
				}
			}
		case "receiver":
			for _, tx := range b.Transactions {
				if tx.To == input.PublicID {
					rec.Value = float64(tx.Amount)
					rec.OtherParty = tx.From
					break
				}
			}
		}

		records = append(records, rec)
	}

	return &GetUserBlocksOutput{Blocks: records}, nil
}
