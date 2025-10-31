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
	SubId       int        `json:"id" example:"1"`
	ServiceName string     `json:"service_name" example:"Yandex Plus"`
	Price       int        `json:"price" example:"400"`
	UserId      string     `json:"user_id" example:"60601fee-2bf1-4721-ae6f-7636e79a0cba"`
	StartDate   time.Time  `json:"start_date" example:"2025-07-01T00:00:00Z"`
	EndDate     *time.Time `json:"end_date,omitempty" example:"2025-10-01T00:00:00Z"`
}
