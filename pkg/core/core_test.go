package test_core

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"testing"
)

func createDBinMemory(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("can't open db: %v", err)
	}
	return db
}

func TestCheckClientExists(t *testing.T) {
	db := createDBinMemory(t)
	defer db.Close()

	id := CheckClientExists("login", db)
	if id != 0 {
		t.Error("want 0, got: ", id)
	}

	_, err := db.Exec(clientsDDL)
	if err != nil {
		t.Fatal(err)
	}

	id = CheckClientExists("login", db)
	if id != 0 {
		t.Error("want 0, got: ", id)
	}

	client := Client{
		Login:    "Aziz",
		Password: "",
		Name:     "",
		Surname:  "",
		Phone:    "",
	}
	_, err = db.Exec(insertClientSQL,
		sql.Named("login", client.Login),
		sql.Named("password", client.Password),
		sql.Named("name", client.Name),
		sql.Named("surname", client.Surname),
		sql.Named("phone", client.Phone),
	)
	if err != nil {
		t.Fatal(err)
	}

	id = CheckClientExists("Aziz", db)
	if id != 1 {
		t.Error("want 1, got: ", id)
	}
}

func TestClient_AddClient(t *testing.T) {
	db := createDBinMemory(t)
	defer db.Close()

	client := Client{
		ID:       1,
		Login:    "Aziz",
		Password: "abc",
		Name:     "a",
		Surname:  "b",
		Phone:    "c",
	}
	err := client.AddClient(db)
	if err == nil {
		t.Error("want not nil, got", err)
	}
	_, err = db.Exec(clientsDDL)
	if err != nil {
		t.Fatal(err)
	}
	err = client.AddClient(db)
	if err != nil {
		t.Error("want nil, got", err)
	}

	clientGot := Client{}
	err = db.QueryRow(`SELECT * FROM clients`,
	).Scan(
		&clientGot.ID,
		&clientGot.Name,
		&clientGot.Surname,
		&clientGot.Login,
		&clientGot.Password,
		&clientGot.Phone,
	)
	if err != nil {
		t.Fatal(err)
	}

	if client != clientGot {
		t.Errorf("want: \n%v\ngot: \n%v\n",
			client, clientGot)
	}
}

func TestAddServices(t *testing.T) {
	db := createDBinMemory(t)
	defer db.Close()
	err := AddServices("abc", 0, db)
	if err == nil{
		t.Errorf("want err, got:",)
	}


	_, err = db.Exec(servicesDDL)
	if err != nil {
		t.Fatal(err)
	}
	err = AddServices("water", 0,db)
	if err != nil {
		t.Error("want not nil, got", err)
	}


	name := "name"
	balance := "balance"

	_, err = db.Exec(insertServicesSQL,
		sql.Named("name", name ),
		sql.Named("balance", balance),
	)

}

func TestAddAtm(t *testing.T) {
	db := createDBinMemory(t)
	defer db.Close()
	err := AddAtm("a","b", db)
		if  err == nil{
			t.Error("want not nil, got:", err)
		}
	_, err = db.Exec(ATMsDDL)
	if err != nil {
		t.Fatal(err)
	}
	err = AddAtm("a", "b", db)
	if err != nil {
		t.Error("want not nil, got", err)
	}


	name := "name"
	address:= "address"

	_, err = db.Exec(insertServicesSQL,
		sql.Named("name", name ),
		sql.Named("address", address),
	)


}


func TestCheckAccountExists(t *testing.T){
db := createDBinMemory(t)
defer db.Close()

id := CheckAccountExists("1234", db)
if id != 0 {
t.Error("want 0, got: ", id)
}
	_, err := db.Exec(accountsDDL)
	if err != nil {
		t.Fatal(err)
	}

	id = CheckAccountExists("1234", db)
	if id != 0 {
		t.Error("want 0, got: ", id)
	}

	var accountNumber string

	_, err = db.Exec(CheckAccountsExists,
		sql.Named("accountNumber", accountNumber),
		)
	if err == nil {
		t.Errorf("This accountNumber not exist")
	}

}



func TestCheckLogin(t *testing.T){
	db := createDBinMemory(t)
	defer db.Close()
	logins := "Aziz"
	passwords := "1234"

	_, err2 := Login(logins, passwords, db)
	if err2 == nil {
		t.Error("this client is not exist: ", logins)
	}



	_, err := db.Exec(accountsDDL)
	if err != nil {
		t.Fatal(err)
	}
	//
	//id = CheckAccountExists("1234", db)
	//if id != 0 {
	//	t.Error("want 0, got: ", id)
	//}
	//
	//var accountNumber string
	//
	//_, err = db.Exec(CheckAccountsExists,
	//	sql.Named("accountNumber", accountNumber),
	//)
	//if err == nil {
	//	t.Errorf("This accountNumber not exist")
	//}

}
