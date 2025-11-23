package service

import (
	"context"
	"crypto/rand"
	"math/big"

	"github.com/fpswan/anycraft-backend/internal/model"
	"github.com/fpswan/anycraft-backend/internal/repository"
)

type ComposeService struct {
	repo *repository.ComposeRepository
}

func NewComposeService(repo *repository.ComposeRepository) *ComposeService {
	return &ComposeService{repo: repo}
}

func (s *ComposeService) GetBaseElements(ctx context.Context, gameCode string) ([]model.Element, error) {
	return s.repo.GetBaseElements(ctx, gameCode)
}

func (s *ComposeService) Combine(ctx context.Context, req model.CombineRequest) model.CombineResponse {
	e, err := s.repo.Combine(ctx, req.GameCode, req.ParentAID, req.ParentBID)
	if err != nil {
		code := "NO_RECIPE"
		return model.CombineResponse{OK: false, Error: &code}
	}
	return model.CombineResponse{OK: true, Result: e}
}

func (s *ComposeService) GetChallenges(ctx context.Context, req model.ChallengesRequest) model.ChallengesResponse {
	all, err := s.repo.GetUndiscoveredChallenges(ctx, req.GameCode, req.DiscoveredCharacterIDs)
	if err != nil {
		return model.ChallengesResponse{OK: false}
	}
	return model.ChallengesResponse{OK: true, Items: sampleK(all, 2)}
}

// random sample helper
func sampleK[T any](in []T, k int) []T {
	n := len(in)
	if k >= n { return append([]T{}, in...) }
	out := make([]T, 0, k)
	used := make(map[int]bool)
	for len(out) < k {
		idx, _ := rand.Int(rand.Reader, big.NewInt(int64(n)))
		i := int(idx.Int64())
		if used[i] { continue }
		used[i] = true
		out = append(out, in[i])
	}
	return out
}
