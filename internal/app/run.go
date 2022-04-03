package app

import (
	"cart_api/internal/controllers"
	"cart_api/internal/repository"
	"cart_api/internal/service"
	"cart_api/internal/storage/mysql"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

var (
	port         = os.Getenv("API_PORT")
	mySqlConnStr = os.Getenv("MYSQL_CONN_STR")
)

func Run() error {
	router := mux.NewRouter()
	mySqlConn, err := mysql.CreateConnection(mySqlConnStr)
	if err != nil {
		return err
	}

	cartRepository := repository.NewCartRepository(mySqlConn)
	cartService := service.NewCartService(cartRepository)
	cartController := controllers.NewCartController(cartService)

	router.HandleFunc("/carts", cartController.GetCarts).Methods("GET")
	router.HandleFunc("/carts/{id}", cartController.GetCart).Methods("GET")

	log.Println("Users api server started on port " + port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		return err
	}

	return nil
}
