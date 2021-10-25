package main

import (
	"net/http"

	//_ "github.com/go-sql-driver/mysql"

	"github.com/meltemseyhan/inventoryservice/database"
	"github.com/meltemseyhan/inventoryservice/product"
)

const apiBasePath = "/api"

func main() {
	database.SetupDatabase()
	product.SetupRoutes(apiBasePath)
	http.ListenAndServe(":5000", nil)
}
