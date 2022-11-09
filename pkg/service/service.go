package service

import (
	"avito_internship/pkg/models"
	"avito_internship/pkg/repository"
)

type User interface {
	GetBalanceUser(user models.UserGetBalanceRequest) (models.UserGetBalanceResponse, error)
	AddBalanceUser(user models.UserAddBalanceRequest) (models.UserAddBalanceResponse, error)
	ReserveMoneyUser(user models.UserReserveMoneyRequest) (models.UserReserveMoneyResponse, error)
}

type Purchase interface {
	GetPurchase(purchase models.GetPurchaseRequest) (models.GetPurchaseResponse, error)
}

type Order interface {
	AddRecord(record models.AddRecordRequest) (models.AddRecordResponse, error)
}

type HistoryUser interface {
	AddRecordHistory(record models.AddHistoryRequest) (models.AddHistoryResponse, error)
}

type Service struct {
	User
	Purchase
	Order
	HistoryUser
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		User:        NewUserService(repos.User),
		Order:       NewOrderService(repos.Order),
		HistoryUser: NewHistoryUserService(repos.HistoryUser),
		Purchase:    NewPurchaseService(repos.Purchase),
	}
}
