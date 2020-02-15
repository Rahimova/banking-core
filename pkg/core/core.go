package test_core

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
)

var (
	ErrInvalidLogin = errors.New("invalid login")
	ErrInvalidPass = errors.New("invalid password")
)

type QueryError struct { // alt + enter
	Query string
	Err   error
}

type DbError struct {
	Err error
}

type DbTxError struct {
	Err         error
	RollbackErr error
}

func (receiver *QueryError) Unwrap() error {
	return receiver.Err
}

func (receiver *QueryError) Error() string {
	return fmt.Sprintf("can't execute query %s: %s", loginSQL, receiver.Err.Error())
}

func queryError(query string, err error) *QueryError {
	return &QueryError{Query: query, Err: err}
}

func (receiver *DbError) Error() string {
	return fmt.Sprintf("can't handle db operation: %v", receiver.Err.Error())
}

func (receiver *DbError) Unwrap() error {
	return receiver.Err
}

func dbError(err error) *DbError {
	return &DbError{Err: err}
}

type Client struct {
	ID int
	Login string
	Password string
	Name string
	Surname string
	Phone string
}

type Account struct {
	ID    int
	ClientID int
	Name  string
	AccountNumber string
	AccountBalance   int
}



func Init(connection *sql.DB) (err error) {
	for _, query := range []string{clientsDDL, accountsDDL, ATMsDDL, servicesDDL} {
		_, err = connection.Exec(query)
		if err != nil {
			return
		}
	}
	return
}

// добавлние пользователей
func (client Client) AddClient(connection *sql.DB) (err error) {
	tx, err := connection.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	_, err = tx.Exec(
		insertClientSQL,
		sql.Named("name", client.Name),
		sql.Named("surname", client.Surname),
		sql.Named("login", client.Login),
		sql.Named("password", client.Password),
		sql.Named("phone", client.Phone),
	)

	return err
}

func CheckClientExists(clientLogin string, connection *sql.DB) (id int) {
	err := connection.QueryRow(checkClientExists, clientLogin).Scan(&id)
	if err != nil {
		return
	}
	return
}

func (account Account) Create(connection *sql.DB) (err error) {
	_, err = connection.Exec(insertAccountSQL,
		sql.Named("user_id", account.ClientID),
		sql.Named("name", account.Name),
		sql.Named("accountNumber", account.AccountNumber),
		sql.Named("accountBalance", account.AccountBalance),

	)
	if err != nil {
		log.Println("core CreateAccountForUser error:", err)
	}
	return err
}

// добавить БАЛАНС
func AddServices(ServiceName string, ServiceBalance int, connection *sql.DB) (err error) {
	_, err = connection.Exec(insertServicesSQL,
		sql.Named("name", ServiceName),
		sql.Named("balance", ServiceBalance),
	)
	if err != nil {
		return
	}
	return err
}

func AddAtm(nameAtm, address string, connection *sql.DB) (err error) {
	_, err = connection.Exec(insertAtmSQL,
		sql.Named("name", nameAtm),
		sql.Named("address", address),

	)
	if err != nil {
		return
	}

	return err
}

func Login(login, password string,  connection *sql.DB) (client Client, err error) {
	/*id, name, surname, login, password, phone*/
	err = connection.QueryRow(
		loginSQL,
		login).Scan(
			&client.ID,
			&client.Name,
			&client.Surname,
			&client.Login,
			&client.Password,
			&client.Phone)

	if err != nil {
		if err != sql.ErrNoRows {
			return client, queryError(loginSQL, err)
		}
		return client, ErrInvalidLogin
	}

	if client.Password != password {
		return client, ErrInvalidPass
	}

	return client, nil
}

//TODO: сделать так чтобы при вызове GetAccount появились  счета 1юзера
func GetAccount(client Client, connection *sql.DB) ([]Account, error){
	var accounts []Account
	rows, err := connection.Query(getAccountsSQL, client.ID)
	if err != nil {
		return nil, queryError(getAccountsSQL, err)
	}

	defer func() {
		if innerErr := rows.Close(); innerErr != nil {
			accounts, err = nil, dbError(innerErr)
		}
	}()

	for rows.Next() {
		account := Account{}
		/*id, user_id, name, accountNumber, accountBalance*/
		err = rows.Scan(
			&account.ID,
			&account.ClientID,
			&account.Name,
			&account.AccountNumber,
			&account.AccountBalance)
		if err != nil {
			return nil, dbError(err)
		}
		accounts = append(accounts, account)
	}

	if rows.Err() != nil {
		return nil, dbError(rows.Err())
	}
	return accounts, nil
}


func TransferByAccount(senderAccountNumber, recipientAccountNumber string, amount int,  connection *sql.DB)( err error)  {
	tx, err := connection.Begin()
	if err != nil{
		return
	}
	res, err := tx.Exec(transferSQL, amount, senderAccountNumber)
	if err != nil{
		tx.Rollback()
		return  err
	}
	if r, _ := res.RowsAffected(); r < 1 {
		tx.Rollback()
		return sql.ErrNoRows
	}

	res, err = tx.Exec(addTransferSQL, amount, recipientAccountNumber)
	if err != nil{
		tx.Rollback()
		return err
	}
	if r, _ := res.RowsAffected(); r <1{
		tx.Rollback()
		return sql.ErrNoRows
	}

	err = tx.Commit()
	return


}


func GetBalance(clientID int, accountNumber string,  connection *sql.DB) (AccountBalance int) {
	err := connection.QueryRow(getBalance, clientID, accountNumber).Scan(&AccountBalance)
	if err != nil {
		return
	}
	return
}

func CheckAccountExists(accountNumber string, connection *sql.DB) (id int) {
	 err := connection.QueryRow(CheckAccountsExists, accountNumber).Scan(&id)
		if err != nil{
			return
		}
			return
	}

