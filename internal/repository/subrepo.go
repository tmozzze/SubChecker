package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tmozzze/SubChecker/internal/model"
)

type SubRepository interface {
	Create(ctx context.Context, s *model.Sub) error
	GetById(ctx context.Context, id int) (*model.Sub, error)
	Update(ctx context.Context, s *model.Sub) error
	Delete(ctx context.Context, id int) error
	List(ctx context.Context, limit, offset int) ([]model.Sub, error)
	FindForUserAndService(ctx context.Context, userID, serviceName string) ([]model.Sub, error)
}

type subRepository struct {
	pool *pgxpool.Pool
}

func NewSubRepository(pool *pgxpool.Pool) SubRepository {
	return &subRepository{pool: pool}
}

func (r *subRepository) Create(ctx context.Context, s *model.Sub) error {
	var endDate any
	if s.EndDate != nil {
		endDate = *s.EndDate
	}

	query := `
		INSERT INTO subs (service_name, price, user_id, start_date, end_date)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING sub_id
	`
	return r.pool.QueryRow(ctx, query, s.ServiceName, s.Price, s.UserId, s.StartDate, endDate).Scan()
}

func (r *subRepository) GetById(ctx context.Context, id int) (*model.Sub, error) {
	var s model.Sub
	query := `
		SELECT sub_id, service_name, price, user_id, start_date, end_date
		FROM subs WHERE sub_id = $1
	`
	row := r.pool.QueryRow(ctx, query, id)
	err := row.Scan(&s.SubId, &s.ServiceName, &s.Price, &s.UserId, &s.StartDate, &s.EndDate)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *subRepository) Update(ctx context.Context, s *model.Sub) error {
	query := `
		UPDATE subs
		SET service_name=$1, price=$2, user_id=$3, start_date=$4, end_date=$5
		WHERE sub_id=$6
	`

	_, err := r.pool.Exec(ctx, query, s.ServiceName, s.Price, s.UserId, s.StartDate, s.EndDate, s.SubId)
	return err
}

func (r *subRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM subs WHERE sub_id=$1`
	_, err := r.pool.Exec(ctx, query, id)

	return err
}

func (r *subRepository) List(ctx context.Context, limit, offset int) ([]model.Sub, error) {
	rows, err := r.pool.Query(ctx, `
        SELECT sub_id, service_name, price, user_id, start_date, end_date
        FROM subs ORDER BY sub_id LIMIT $1 OFFSET $2
    `, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []model.Sub
	for rows.Next() {
		var s model.Sub
		err := rows.Scan(&s.SubId, &s.ServiceName, &s.Price, &s.UserId, &s.StartDate, &s.EndDate)
		if err != nil {
			return nil, err
		}
		result = append(result, s)
	}
	return result, nil
}

func (r *subRepository) FindForUserAndService(ctx context.Context, userId, serviceName string) ([]model.Sub, error) {
	base := `SELECT sub_id, service_name, price, user_id, start_date, end_date FROM subs WHERE 1=1`
	args := []interface{}{}
	i := 1
	if userId != "" {
		base += fmt.Sprintf(" AND user_id = $%d", i)
		args = append(args, userId)
		i++
	}
	if serviceName != "" {
		base += fmt.Sprintf(" AND service_name = $%d", i)
		args = append(args, serviceName)
		i++
	}

	rows, err := r.pool.Query(ctx, base, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []model.Sub
	for rows.Next() {
		var s model.Sub
		err := rows.Scan(&s.SubId, &s.ServiceName, &s.Price, &s.UserId, &s.StartDate, &s.EndDate)
		if err != nil {
			return nil, err
		}
		result = append(result, s)
	}
	return result, nil
}
