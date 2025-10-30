package service

import (
	"context"
	"fmt"
	"time"

	"github.com/tmozzze/SubChecker/internal/model"
	"github.com/tmozzze/SubChecker/internal/repository"
	"github.com/tmozzze/SubChecker/internal/utils"
)

type SubService interface {
	Create(ctx context.Context, s *model.Sub) error
	GetById(ctx context.Context, id int) (*model.Sub, error)
	Update(ctx context.Context, s *model.Sub) error
	Delete(ctx context.Context, id int) error
	List(ctx context.Context, limit, offset int) ([]model.Sub, error)
	SumCost(ctx context.Context, userID, serviceName string, periodStart, periodEnd time.Time) (int64, error)
}

type subService struct {
	repository repository.SubRepository
}

func NewSubService(r repository.SubRepository) SubService {
	return &subService{repository: r}
}

func (s *subService) Create(ctx context.Context, sub *model.Sub) error {
	return s.repository.Create(ctx, sub)
}

func (s *subService) GetById(ctx context.Context, id int) (*model.Sub, error) {
	return s.repository.GetById(ctx, id)
}

func (s *subService) Update(ctx context.Context, sub *model.Sub) error {
	return s.repository.Update(ctx, sub)
}

func (s *subService) Delete(ctx context.Context, id int) error { return s.repository.Delete(ctx, id) }

func (s *subService) List(ctx context.Context, limit, offset int) ([]model.Sub, error) {
	return s.repository.List(ctx, limit, offset)
}

func (s *subService) SumCost(ctx context.Context, userId, serviceName string, startDate, endDate time.Time) (int64, error) {
	subs, err := s.repository.FindForUserAndService(ctx, userId, serviceName)
	if err != nil {
		return 0, err
	}

	var total int64 = 0
	for _, sub := range subs {
		sStart := utils.TruncateToMonth(sub.StartDate)
		var sEnd time.Time
		if sub.EndDate == nil {
			sEnd = time.Now().UTC()
		} else {
			sEnd = *sub.EndDate
		}
		sEnd = utils.TruncateToMonth(sEnd)

		pStart := utils.TruncateToMonth(startDate)
		pEnd := utils.TruncateToMonth(endDate)

		months := utils.MonthsOverlap(sStart, sEnd, pStart, pEnd)
		if months > 0 {
			total += int64(months) * int64(sub.Price)
		} else {
			return 0, fmt.Errorf("failed to count months")
		}
	}
	return total, nil
}
