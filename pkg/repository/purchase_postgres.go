package repository

import (
	"avito_internship/pkg/models"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type PurchasePostgres struct {
	db *sqlx.DB
}

func NewPurchasePostgres(db *sqlx.DB) *PurchasePostgres {
	return &PurchasePostgres{db: db}
}

func (r *PurchasePostgres) GetPurchase(purchase models.GetPurchaseRequest) (models.GetPurchaseResponse, error) {
	var response models.GetPurchaseResponse
	var flagExist bool

	query := fmt.Sprintf("SELECT EXISTS(SELECT * from %s WHERE purchaseid=$1)", purchaseTable)
	row := r.db.QueryRow(query, purchase.PurchaseId)

	if err := row.Scan(&flagExist); err != nil {
		response.Status = "Error"
		return response, err
	}

	if !flagExist {
		response.Status = "Error"
		return response, errors.New("Error, wrong id of purchase")
	}

	response.Status = "OK"

	return response, nil
}
