package models

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

const dbuser = "root"
const dbpassword = "sersjezakon" // Add your MySQL password
const dbname = "api"

// Customer struct that has all the types from the MySQL table
type Customer struct {
	IDC      int     `json:"idC"`
	Name     string  `json:"name"`
	Surname  string  `json:"surname"`
	Title    string  `json:"title"`
	Address  string  `json:"address"`
	City     string  `json:"city"`
	Zip      int     `json:"zip"`
	Phone    string  `json:"phone"`
	Username string  `json:"username"`
	Password string  `json:"password"`
	Balance  float64 `json:"balance"`
}

// function to create the connection string to the MySQL server
func GetConnectionString() (string, error) {
	dsn := dbuser + ":" + dbpassword + "@tcp(127.0.0.1:3306)/" + dbname
	fmt.Println("DSN:", dsn) // Print connection details for debugging

	db, err := sql.Open("mysql", dsn)

	// if it gets an error while connecting it prints the error message
	if err != nil {
		fmt.Println("Error connecting to database:", err.Error())
		return "", err
	}
	defer db.Close()

	return dsn, nil
}

func GetInfoCustomers() []Customer {
	dsn, err := GetConnectionString()
	if err != nil {
		return nil
	}

	// makes a variable for the open connection to SQL
	db, err := sql.Open("mysql", dsn)

	// string od the SQL query
	query := "SELECT idC, title, name, surname, address, city, zip FROM customer;"
	fmt.Println("Executing query:", query)

	// string of the result of the query and the potential error message from the query
	results, err := db.Query(query)
	// if there is an error message it writes it out
	if err != nil {
		fmt.Println("Error executing query:", err.Error())
		return nil
	}
	defer results.Close()

	// array of Customer structures
	var customers []Customer
	// for each result in the query, makes a Customer and add the queried items to the Customer array
	for results.Next() {
		var customer Customer
		err = results.Scan(&customer.IDC, &customer.Title, &customer.Name, &customer.Surname, &customer.Address, &customer.City, &customer.Zip)
		// if it gets an error in a row, it prints it out
		if err != nil {
			fmt.Println("Error scanning row:", err.Error())
			continue
		}
		// else it adds the Customer to the array of Customer
		customers = append(customers, customer)
	}

	// prints out the array of Customer structures
	fmt.Println("Retrieved customers:", customers)

	return customers
}

func CustomerBalance(id int) (Customer, error) {
	dsn, err := GetConnectionString()
	if err != nil {
		return Customer{}, nil
	}

	// makes a variable for the open connection to SQL
	db, err := sql.Open("mysql", dsn)

	// string od the SQL query
	balanceQuery := "SELECT balance FROM customer WHERE idC = ?;"
	fmt.Println("Executing query:", balanceQuery)

	result, err := db.Query(balanceQuery, id)

	// if there is an error message it writes it out
	if err != nil {
		fmt.Println("Error executing query:", err.Error())
		return Customer{}, nil
	}
	defer result.Close()

	if result.Next() {
		var customer Customer
		err = result.Scan(&customer.Balance)

		if err != nil {
			return Customer{}, nil
		}
		return customer, nil
	} else {
		return Customer{}, nil
	}

}

func CreateCustomer(str Customer) (Customer, error) {
	dsn, err := GetConnectionString()
	if err != nil {
		return Customer{}, err
	}
	// makes a variable for the open connection to SQL
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return Customer{}, err // Return an empty Customer and the error if db.Open fails
	}
	defer db.Close()

	// Prepare the insert query (assuming IDC is auto-incremented)
	insertQuery := "INSERT INTO customer (name, surname, title, address, city, zip, phone, username, password, balance) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	fmt.Println("Executing query:", insertQuery)

	// Execute the query, using placeholders to prevent SQL injection
	result, err := db.Exec(insertQuery, str.Name, str.Surname, str.Title, str.Address, str.City, str.Zip, str.Phone, str.Username, str.Password, str.Balance)
	if err != nil {
		fmt.Println("Error executing insert:", err.Error())
		return Customer{}, err
	}

	// Get the generated ID (assuming IDC is auto-incremented)
	id, err := result.LastInsertId()
	if err != nil {
		fmt.Println("Error getting last insert ID:", err.Error())
		return Customer{}, err
	}

	// Set the generated IDC in the Customer struct
	str.IDC = int(id)
	fmt.Println("Customer added successfully with ID:", str.IDC)

	return str, nil
}

func GetCustomerbyID(id int, user string, pass string) (Customer, error) {
	dsn, err := GetConnectionString()
	if err != nil {
		return Customer{}, err
	}
	// makes a variable for the open connection to SQL
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return Customer{}, err // Return an empty Customer and the error if db.Open fails
	}
	defer db.Close()

	selectQuery := "SELECT idC, name, surname, title, address, city, zip, phone, username, password, balance FROM customer where idC=? AND username=? AND password=?"

	result, err := db.Query(selectQuery, id, user, pass)

	// if there is an error in the query
	if err != nil {
		fmt.Println("Error executing query:", err.Error())
		return Customer{}, err
	}

	defer result.Close()

	// if it finds a customer
	if result.Next() {
		var customer Customer
		err = result.Scan(&customer.IDC, &customer.Name, &customer.Surname, &customer.Title, &customer.Address, &customer.City, &customer.Zip, &customer.Phone, &customer.Username, &customer.Password, &customer.Balance)
		if err != nil {
			return Customer{}, err
		}
		return customer, nil
	} else {
		// if it doesn't find a customer
		return Customer{}, fmt.Errorf("no customer found")
	}
}

func DepositMoney(user, pass string, id int, amount float64) (Customer, error) {
	dsn, err := GetConnectionString()
	if err != nil {
		return Customer{}, err
	}

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return Customer{}, err
	}
	defer db.Close()

	// Query to get the balance
	query := "SELECT balance FROM customer WHERE username = ? AND password = ? AND idC = ?"
	result, err := db.Query(query, user, pass, id)
	if err != nil {
		return Customer{}, fmt.Errorf("failed to query database: %v", err)
	}
	defer result.Close()

	if result.Next() {
		var customer Customer
		err = result.Scan(&customer.Balance)
		if err != nil {
			return Customer{}, fmt.Errorf("failed to scan balance: %v", err)
		}

		// Calculate the new balance
		newBalance := customer.Balance + amount

		// Update the balance
		updateQuery := "UPDATE customer SET balance = ? WHERE username = ? AND password = ? AND idC = ?"
		_, err = db.Exec(updateQuery, newBalance, user, pass, id)
		if err != nil {
			return Customer{}, fmt.Errorf("failed to update balance: %v", err)
		}

		customer.Balance = newBalance
		return customer, nil
	} else {
		return Customer{}, fmt.Errorf("user not found or incorrect credentials")
	}
}
