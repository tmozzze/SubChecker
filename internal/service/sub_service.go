package service

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
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
	log        *logrus.Logger
}

func NewSubService(r repository.SubRepository, log *logrus.Logger) SubService {
	return &subService{repository: r, log: log}
}

func (s *subService) Create(ctx context.Context, sub *model.Sub) error {
	s.log.WithFields(logrus.Fields{
		"user_id": sub.UserId,
		"service": sub.ServiceName,
	}).Info("Creating new subscription")

	return s.repository.Create(ctx, sub)
}

func (s *subService) GetById(ctx context.Context, id int) (*model.Sub, error) {
	s.log.WithFields(logrus.Fields{
		"sub_id": id,
	}).Info("Getting subscription by id")

	return s.repository.GetById(ctx, id)
}

func (s *subService) Update(ctx context.Context, sub *model.Sub) error {
	s.log.WithFields(logrus.Fields{
		"sub_id": sub.SubId,
	}).Info("Updating subscription")

	return s.repository.Update(ctx, sub)
}

func (s *subService) Delete(ctx context.Context, id int) error {
	s.log.WithFields(logrus.Fields{
		"sub_id": id,
	}).Info("Deleting subscription")

	return s.repository.Delete(ctx, id)
}

func (s *subService) List(ctx context.Context, limit, offset int) ([]model.Sub, error) {
	s.log.Info("Getting list of subscriptions")

	return s.repository.List(ctx, limit, offset)
}

func (s *subService) SumCost(ctx context.Context, userId, serviceName string, startDate, endDate time.Time) (int64, error) {
	s.log.WithFields(logrus.Fields{
		"user_id": userId,
		"service": serviceName,
	}).Info("Calculating total subscription cost")

	subs, err := s.repository.SumCost(ctx, userId, serviceName)
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
			continue
		}
	}

	return total, nil
}
