package router

import (
	"go-api/middleware"
	"net/http"

	"github.com/gorilla/mux"
)

// Router is exported and used in main.go
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func readinessHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
func Router() *mux.Router {

	router := mux.NewRouter()

	router.HandleFunc("/api/getTayaraUser", middleware.GetTayaraUser).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/getAllTayaraUsers", middleware.GetAllTayaraUsers).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/changeUserPass", middleware.ChangeUserPass).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/users", middleware.GetAllUser).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/addUser", middleware.CreateUser).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/setUserPass", middleware.SetUserPass).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/deleteUser/{id}", middleware.DeleteUser).Methods("DELETE", "OPTIONS")
	router.HandleFunc("/api/deleteAllUsers", middleware.DeleteAllUser).Methods("DELETE", "OPTIONS")
	router.HandleFunc("/health", healthHandler)
	router.HandleFunc("/readiness", readinessHandler)
	return router
}
