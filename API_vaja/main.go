package main

import (
	"net/http"
	"strconv"
	"vaja/API/models"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.GET("/GetInfo", getInfoHandler)
	router.POST("/MakeAccount", creatCustomerHandler)
	router.POST("/GetInfoID", customerByIDHandler)
	router.GET("/Balance/:id", checkBalanceHandler)
	router.POST("/ChangeBalance", changeBalanceHandler)
	router.Run("localhost:8085")
}

func getInfoHandler(c *gin.Context) {
	customer := models.GetInfoCustomers()
	if len(customer) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No customers found"})
	} else {
		c.JSON(http.StatusOK, customer)
	}
}

func creatCustomerHandler(c *gin.Context) {
	var str models.Customer

	// Bind the incoming JSON to the Customer struct
	if err := c.BindJSON(&str); err != nil {
		// Return a Bad Request status if there's an error
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// make two variable to get the responce from the function that creates the customer, a Customer struct and an error
	newCustomer, err := models.CreateCustomer(str)
	// if there is an error, write it out in the JSON
	if err != nil {
		// Return an Internal Server Error status if there's an error
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to create customer"})
		return
	}

	// else make the customer and return the IDC (account number) of the created customer
	c.IndentedJSON(http.StatusCreated, gin.H{"account_number": newCustomer.IDC})
}

//-------------------------------------------------------------------------------------

func customerByIDHandler(c *gin.Context) {
	// Define a struct to bind the JSON input
	var request struct {
		ID       int    `json:"id"`
		Username string `json:"username"`
		Password string `json:"password"`
	}

	// Bind JSON to the struct
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	// Call the function to get the customer by ID, username, and password
	customer, err := models.GetCustomerbyID(request.ID, request.Username, request.Password)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Return the customer details in JSON format
	c.JSON(http.StatusOK, customer)
}

//-------------------------------------------------------------------------------------

func checkBalanceHandler(c *gin.Context) {
	// Get the ID as a string from the URL
	idStr := c.Param("id")

	// Convert string to int
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format, must be a number"})
		return
	}

	// Call the function with the converted integer ID
	str, err := models.CustomerBalance(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
		return
	}

	// Return the balance as JSON
	c.IndentedJSON(http.StatusOK, gin.H{"account_balance": str.Balance})
}

func changeBalanceHandler(c *gin.Context) {
	var st struct {
		ID       int     `json:"id"`
		Username string  `json:"username"`
		Password string  `json:"password"`
		Amount   float64 `json:"amount"`
	}

	// Bind the request body to the struct
	if err := c.BindJSON(&st); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	// Call DepositMoney function
	newAmount, err := models.DepositMoney(st.Username, st.Password, st.ID, st.Amount)
	if err != nil {
		// Respond with the error message if any
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Respond with the updated balance
	c.JSON(http.StatusOK, gin.H{
		"message":    "Deposit successful",
		"new amount": newAmount.Balance,
	})
}
