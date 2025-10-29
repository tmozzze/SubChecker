package model

import "time"

/*
{
“service_name”: “Yandex Plus”,
“price”: 400,
“user_id”: “60601fee-2bf1-4721-ae6f-7636e79a0cba”,
“start_date”: “07-2025”
}
*/

type Sub struct {
	SubId       int        `db:"sub_id" json:"sub_id"`
	ServiceName string     `json:"service_name" db:"service_name"`
	Price       int        `json:"price" db:"price"`
	UserId      string     `json:"user_id" db:"user_id"`
	StartDate   time.Time  `json:"start_date" db:"start_date"`
	EndDate     *time.Time `json:"end_date,omitempty" db:"end_date"`
}
