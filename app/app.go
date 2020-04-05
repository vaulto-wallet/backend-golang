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
	a.Post("/api/seeds", a.CreateSeed)
	a.Get("/api/seeds", a.GetSeeds)
	a.Post("/api/wallets", a.CreateWallet)
	a.Get("/api/wallets", a.GetWallets)
	a.Get("/api/wallets/{asset}", a.GetWalletsForAsset)
	a.Post("/api/assets", a.CreateAsset)
	a.Get("/api/assets", a.GetAssets)
	a.Post("/api/orders", a.CreateOrder)
	a.Get("/api/orders", a.GetOrders)
	a.Put("/api/orders", a.UpdateOrder)
	a.Post("/api/address", a.CreateAddress)
	a.Get("/api/address/{wallet}", a.GetAddressList)
	a.Put("/api/address", a.UpdateAddress)
	a.Post("/api/transactions", a.CreateTransaction)
	a.Get("/api/transactions", a.GetTransactions)
	a.Put("/api/transactions", a.UpdateTransaction)
	a.Get("/api/transaction/{transaction}", a.GetTransaction)

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

//Assets
func (a *App) CreateAsset(w http.ResponseWriter, r *http.Request) {
	h.CreateAsset(a.DB, w, r)
}

func (a *App) GetAssets(w http.ResponseWriter, r *http.Request) {
	h.GetAssets(a.DB, w, r)
}

// Seeds
func (a *App) CreateSeed(w http.ResponseWriter, r *http.Request) {
	h.CreateSeed(a.DB, w, r)
}

func (a *App) GetSeeds(w http.ResponseWriter, r *http.Request) {
	h.GetSeeds(a.DB, w, r)
}

// Wallets
func (a *App) CreateWallet(w http.ResponseWriter, r *http.Request) {
	h.CreateWallet(a.DB, w, r)
}

func (a *App) GetWallets(w http.ResponseWriter, r *http.Request) {
	h.GetWallets(a.DB, w, r)
}

func (a *App) GetWalletsForAsset(w http.ResponseWriter, r *http.Request) {
	h.GetWalletsForAsset(a.DB, w, r)
}

// Orders
func (a *App) CreateOrder(w http.ResponseWriter, r *http.Request) {
	h.CreateOrder(a.DB, w, r)
}

func (a *App) GetOrders(w http.ResponseWriter, r *http.Request) {
	h.GetOrders(a.DB, w, r)
}

func (a *App) UpdateOrder(w http.ResponseWriter, r *http.Request) {
	h.UpdateOrder(a.DB, w, r)
}

// Addresses
func (a *App) CreateAddress(w http.ResponseWriter, r *http.Request) {
	h.CreateAddress(a.DB, w, r)
}

func (a *App) GetAddressList(w http.ResponseWriter, r *http.Request) {
	h.GetAddress(a.DB, w, r)
}

func (a *App) UpdateAddress(w http.ResponseWriter, r *http.Request) {
	h.UpdateAddress(a.DB, w, r)
}

// transactions
func (a *App) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	h.CreateTransaction(a.DB, w, r)
}

func (a *App) GetTransactions(w http.ResponseWriter, r *http.Request) {
	h.GetTransactions(a.DB, w, r)
}

func (a *App) GetTransaction(w http.ResponseWriter, r *http.Request) {
	h.GetTransaction(a.DB, w, r)
}

func (a *App) UpdateTransaction(w http.ResponseWriter, r *http.Request) {
	h.UpdateTransaction(a.DB, w, r)
}
