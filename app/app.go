package app

import (
	h "../handlers"
	mw "../middlewares"
	m "../models"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"log"
	"net/http"
)

// App has router and db instances
type App struct {
	Router *mux.Router
	DB     *gorm.DB
}

// App initialize with predefined configuration
func (a *App) Initialize() {

	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		log.Fatal("Could not connect database")
	}

	db.AutoMigrate(&m.User{})
	a.DB = db
	a.Router = mux.NewRouter()
	a.setRouters()
}

func (a *App) Get(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("GET")
}

// Wrap the router for POST method
func (a *App) Post(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("POST")
}

// Wrap the router for PUT method
func (a *App) Put(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("PUT")
}

// Wrap the router for DELETE method
func (a *App) Delete(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("DELETE")
}

func (a *App) Run(host string) {
	a.Router.Use(mw.LoggingMiddleware)
	a.Router.Use(mw.AuthMiddleware)
	log.Fatal(http.ListenAndServe(host, a.Router))
}

func (a *App) setRouters() {
	a.Get("/api/clear", a.Clear)
	a.Post("/api/users/login", a.Login)
	a.Post("/api/users/register", a.Register)
	a.Post("/api/account", a.CreateAccount)
}

func (a *App) Login(w http.ResponseWriter, r *http.Request) {
	h.UserLogin(a.DB, w, r)
}
func (a *App) Register(w http.ResponseWriter, r *http.Request) {
	h.UserRegister(a.DB, w, r)
}

func (a *App) CreateAccount(w http.ResponseWriter, r *http.Request) {
	h.CreateAccount(a.DB, w, r)
}

func (a *App) Clear(w http.ResponseWriter, r *http.Request) {
	h.Clear(a.DB, w, r)
}