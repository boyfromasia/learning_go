package repository

import (
	"avito_internship/pkg/models"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type UserPostgres struct {
	db *sqlx.DB
}

func NewUserPostgres(db *sqlx.DB) *UserPostgres {
	return &UserPostgres{db: db}
}

func (r *UserPostgres) GetBalanceUser(user models.UserGetBalanceRequest) (models.UserGetBalanceResponse, error) {
	var response models.UserGetBalanceResponse
	var balance float64
	var flagExist bool

	query := fmt.Sprintf("SELECT EXISTS(SELECT * from %s WHERE userid=$1)", usersTable)
	row := r.db.QueryRow(query, user.UserId)

	if err := row.Scan(&flagExist); err != nil {
		response.Balance = 0
		return response, err
	}

	if !flagExist {
		response.Balance = 0
		return response, errors.New("Error, wrong id of user")
	} else {
		r.db.QueryRow(fmt.Sprintf("SELECT balance from %s WHERE userid=$1", usersTable), user.UserId).Scan(&balance)
		response.Balance = balance
		return response, nil
	}

}

func (r *UserPostgres) AddBalanceUser(user models.UserAddBalanceRequest) (models.StatusResponse, error) {
	var response models.StatusResponse
	var flagExist bool
	var id int

	query := fmt.Sprintf("SELECT EXISTS(SELECT * from %s WHERE userid=$1)", usersTable)
	row := r.db.QueryRow(query, user.UserId)

	if err := row.Scan(&flagExist); err != nil {
		response.Status = "Error"
		return response, err
	}

	if flagExist {
		query := fmt.Sprintf("UPDATE %s SET balance=$1+balance WHERE userid=$2 RETURNING userid", usersTable)
		row := r.db.QueryRow(query, user.Balance, user.UserId)

		if err := row.Scan(&id); err != nil {
			response.Status = "Error"
			return response, err
		}
	} else {
		query := fmt.Sprintf("INSERT INTO %s (userid, balance, reserve) VALUES ($1, $2, 0) RETURNING userid", usersTable)
		row := r.db.QueryRow(query, user.UserId, user.Balance)

		if err := row.Scan(&id); err != nil {
			response.Status = "Error"
			return response, err
		}
	}

	response.Status = "OK"
	return response, nil
}

func (r *UserPostgres) ReserveMoneyUser(user models.UserReserveMoneyRequest) (models.StatusResponse, error) {
	var response models.StatusResponse
	var flagExist bool
	var balance, reserve float64

	rowExist := fmt.Sprintf("SELECT EXISTS(SELECT * from %s WHERE userid=$1)", usersTable)
	queryExist := r.db.QueryRow(rowExist, user.UserId)

	if err := queryExist.Scan(&flagExist); err != nil {
		response.Status = "Error"
		return response, err
	}

	query := fmt.Sprintf("SELECT balance, reserve from %s WHERE userid=$1", usersTable)
	row := r.db.QueryRow(query, user.UserId)

	if err := row.Scan(&balance, &reserve); err != nil {
		response.Status = "Error"
		return response, err
	}

	if balance-reserve-user.Price < 0 {
		response.Status = "Error"
		return response, errors.New("Error, not enough money")
	}

	r.db.QueryRow(fmt.Sprintf("UPDATE %s SET reserve=reserve + $1 WHERE userid=$2 RETURNING userid", usersTable), user.Price, user.UserId)
	response.Status = "OK"
	return response, nil
}

func (r *UserPostgres) ApproveReserveUser(user models.UserDecisionRequest) (models.StatusResponse, error) {
	var response models.StatusResponse
	var id int

	query := fmt.Sprintf("UPDATE %s SET reserve=reserve - $1, balance=balance - $1 WHERE userid=$2 RETURNING userid", usersTable)
	row := r.db.QueryRow(query, user.Cost, user.UserId)

	if err := row.Scan(&id); err != nil {
		response.Status = "Error"
		return response, err
	}

	response.Status = "OK"
	return response, nil
}

func (r *UserPostgres) RejectReserveUser(user models.UserDecisionRequest) (models.StatusResponse, error) {
	var response models.StatusResponse
	var id int

	query := fmt.Sprintf("UPDATE %s SET reserve=reserve - $1 WHERE userid=$2 RETURNING userid", usersTable)
	row := r.db.QueryRow(query, user.Cost, user.UserId)

	if err := row.Scan(&id); err != nil {
		response.Status = "Error"
		return response, err
	}

	response.Status = "OK"
	return response, nil
}

func (r *UserPostgres) TransferMoneyUsers(user models.UserTransferRequest) (models.StatusResponse, error) {
	var response models.StatusResponse

	tx, err := r.db.Begin()
	if err != nil {
		response.Status = "Error"
		return response, errors.New("Error with DB")
	}

	_, err = tx.Exec(fmt.Sprintf("UPDATE %s SET balance = balance - $1 WHERE userid = $2", usersTable), user.Amount, user.FromId)
	_, err = tx.Exec(fmt.Sprintf("UPDATE %s SET balance = balance + $1 WHERE userid = $2", usersTable), user.Amount, user.ToId)
	if err != nil {
		response.Status = "Error"
		return response, err
	}

	err = tx.Commit()
	if err != nil {
		response.Status = "Error"
		return response, err
	}

	response.Status = "OK"

	return response, nil
}
