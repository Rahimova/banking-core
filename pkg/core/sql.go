package test_core


const clientsDDL  = `create table if not exists clients (
    id integer primary key AUTOINCREMENT,
    name text not null,
    surname text not null,
	login text not null unique,
    password text not null ,
	phone text not null
);`
const accountsDDL = `create table if not exists accounts (
    id integer primary key AUTOINCREMENT,
    user_id integer not null references clients,
    name text not null,
    accountNumber text not null,
	accountBalance integer not null
);`

const ATMsDDL  = `create table if not exists ATMs(
    id integer primary key AUTOINCREMENT,
    name text not null,
	address text not null
);`

const servicesDDL =  `create table if not exists services (
    id integer primary key AUTOINCREMENT,
    name text not null,
	balance integer not null
);`



const checkClientExists = `SELECT id FROM clients WHERE login = ?`
const getBalance = `SELECT accountBalance FROM accounts WHERE user_id = ? and accountNumber = ?`
const insertClientSQL = `INSERT INTO clients(name, surname, login, password, phone) VALUES (:name, :surname, :login, :password, :phone);`
const insertAccountSQL  = `INSERT INTO accounts(user_id, name, accountNumber, accountBalance) VALUES (:user_id, :name, :accountNumber, :accountBalance);`
const insertServicesSQL  = `INSERT INTO services(name, balance) VALUES (:name, :balance);`
const insertAtmSQL  = `INSERT INTO ATMs(name, address) VALUES (:name, :address);`
const loginSQL = `SELECT id, name, surname, login, password, phone FROM clients WHERE login = ?`
const getAccountsSQL = `SELECT id, user_id, name, accountNumber, accountBalance FROM accounts WHERE user_id = ?`
const CheckAccountsExists  =`SELECT id From accounts WHERE accountNumber =  ?`
const transferSQL  = `update accounts set accountBalance = accountBalance - ? where  accountNumber = ?`
const addTransferSQL = `update accounts set accountBalance = accountBalance + ? where  accountNumber = ?`