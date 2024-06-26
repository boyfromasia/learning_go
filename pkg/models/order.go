package models

import (
	"errors"
	"time"
)

type AddRecordRequest struct {
	OrderId     int       `json:"order_id" binding:"required"`
	UserId      int       `json:"user_id" binding:"required"`
	PurchaseId  int       `json:"purchase_id" binding:"required"`
	Price       float64   `json:"price" binding:"required"`
	TimeCreated time.Time `json:"time_created"`
	Status      string    `json:"status"`
}

type AddRecordResponse struct {
	Status string `json:"status"`
}

type DecisionReserveRequest struct {
	OrderId    int     `json:"order_id" binding:"required"`
	UserId     int     `json:"user_id" binding:"required"`
	PurchaseId int     `json:"purchase_id" binding:"required"`
	Price      float64 `json:"price" binding:"required"`
	Decision   string  `json:"decision" binding:"required"`
}

type DecisionReserveResponse struct {
	Status string `json:"status"`
}

func (i AddRecordRequest) Validate() error {
	if i.Price < 0 {
		return errors.New("Wrong input - 'price'")
	}

	return nil
}

func (i DecisionReserveRequest) Validate() error {
	if !(i.Decision == "approved" || i.Decision == "cancelled") {
		return errors.New("Wrong input - 'decision'")
	}
	return nil
}
