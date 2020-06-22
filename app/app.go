package app

import (
	h "../handlers"
	mw "../middlewares"
	m "../models"
	"context"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"log"
	"net/http"
)

// App has router and db instances
type App struct {
	Router         *mux.Router
	DB             *gorm.DB
	MasterPassword []byte
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
	//a.Router.Use(mw.AccessControlMiddleware)

	a.Router.Use(mw.LoggingMiddleware)
	a.Router.Use(mw.AuthMiddlewareGenerator(a.DB))
	a.Router.Use(mw.TwoFAMiddlewareGenerator(a.DB))

	cors := handlers.CORS(
		handlers.AllowedHeaders([]string{"Origin", "Content-Type", "Authorization"}),
		handlers.AllowedOrigins([]string{"http://localhost:8001"}),
		handlers.AllowedMethods([]string{"POST", "GET", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowCredentials(),
	)(a.Router)

	/*	c := cors.New(cors.Options{
			AllowedOrigins: []string{"http://localhost:8001"},
			AllowCredentials: true,
		})
		handler := c.Handler(a.Router)
	*/

	log.Fatal(http.ListenAndServe(host, cors))
}

func (a *App) setRouters() {
	a.Post("/api/clear", a.Clear)
	a.Post("/api/start", a.Start)
	a.Get("/api/users", a.GetUsers)
	a.Post("/api/users/login", a.Login)
	a.Post("/api/users/register", a.Register)
	a.Get("/api/account", a.GetAccount)
	a.Put("/api/account", a.SetAccount)
	a.Post("/api/seeds", a.CreateSeed)
	a.Get("/api/seeds", a.GetSeeds)
	a.Post("/api/wallets", a.CreateWallet)
	a.Put("/api/wallets/share/{wallet}", a.ShareWallet)
	a.Get("/api/wallets", a.GetWallets)
	a.Get("/api/wallet/{wallet}", a.GetWallet)
	a.Get("/api/wallets/{asset}", a.GetWalletsForAsset)
	a.Get("/api/wallet/orders/{wallet}", a.GetOrders)
	a.Get("/api/wallet/rules/{wallet}", a.GetRules)
	a.Get("/api/wallet/transactions/{wallet}", a.GetTransactions)
	a.Post("/api/assets", a.CreateAsset)
	a.Get("/api/assets", a.GetAssets)
	a.Post("/api/orders", a.CreateOrder)
	a.Get("/api/orders", a.GetOrders)
	a.Get("/api/order/{order}", a.GetOrder)
	a.Post("/api/order/{order}/confirm", a.ConfirmOrder)
	a.Get("/api/order/{order}/txs", a.GetOrderTransactions)
	a.Put("/api/orders", a.UpdateOrder)
	a.Post("/api/address", a.CreateAddress)
	a.Get("/api/address/{wallet}", a.GetAddressList)
	a.Put("/api/address", a.UpdateAddress)
	a.Post("/api/transactions", a.CreateTransaction)
	a.Get("/api/transactions", a.GetTransactions)
	a.Put("/api/transactions", a.UpdateTransaction)
	a.Get("/api/transaction/id/{id}", a.GetTransaction)
	a.Get("/api/transaction/txhash/{txhash}", a.GetTransaction)
	a.Post("/api/firewall", a.CreateRule)
	a.Get("/api/firewall/{rule}", a.GetRule)
}

func (a *App) GetUsers(w http.ResponseWriter, r *http.Request) {
	h.GetUsers(a.DB, w, r)
}

func (a *App) Login(w http.ResponseWriter, r *http.Request) {
	h.UserLogin(a.DB, w, r)
}
func (a *App) Register(w http.ResponseWriter, r *http.Request) {
	h.UserRegister(a.DB, w, r)
}

func (a *App) Start(w http.ResponseWriter, r *http.Request) {
	a.MasterPassword = h.Start(a.DB, w, r)
}

func (a *App) SetAccount(w http.ResponseWriter, r *http.Request) {
	h.SetAccount(a.DB, w, r)
}

func (a *App) GetAccount(w http.ResponseWriter, r *http.Request) {
	h.GetAccount(a.DB, w, r)
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
	ctx := context.WithValue(r.Context(), "masterPassword", a.MasterPassword)
	r = r.WithContext(ctx)
	h.CreateSeed(a.DB, w, r)
}

func (a *App) GetSeeds(w http.ResponseWriter, r *http.Request) {
	h.GetSeeds(a.DB, w, r)
}

// Wallets
func (a *App) CreateWallet(w http.ResponseWriter, r *http.Request) {
	context.WithValue(r.Context(), "masterPassword", a.MasterPassword)
	h.CreateWallet(a.DB, w, r)
}

func (a *App) ShareWallet(w http.ResponseWriter, r *http.Request) {
	context.WithValue(r.Context(), "masterPassword", a.MasterPassword)
	h.ShareWallet(a.DB, w, r)
}

func (a *App) GetWallet(w http.ResponseWriter, r *http.Request) {
	h.GetWallet(a.DB, w, r)
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

func (a *App) GetOrder(w http.ResponseWriter, r *http.Request) {
	h.GetOrder(a.DB, w, r)
}

func (a *App) GetOrderTransactions(w http.ResponseWriter, r *http.Request) {
	h.GetTransactions(a.DB, w, r)
}

func (a *App) UpdateOrder(w http.ResponseWriter, r *http.Request) {
	h.UpdateOrder(a.DB, w, r)
}

func (a *App) ConfirmOrder(w http.ResponseWriter, r *http.Request) {
	h.ConfirmOrder(a.DB, w, r)
}

// Addresses
func (a *App) CreateAddress(w http.ResponseWriter, r *http.Request) {
	ctx := context.WithValue(r.Context(), "masterPassword", a.MasterPassword)
	r = r.WithContext(ctx)
	h.CreateAddress(a.DB, w, r)
}

func (a *App) GetAddressList(w http.ResponseWriter, r *http.Request) {
	h.GetAddresses(a.DB, w, r)
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

func (a *App) GetRules(w http.ResponseWriter, r *http.Request) {
	h.GetRules(a.DB, w, r)
}

func (a *App) GetRule(w http.ResponseWriter, r *http.Request) {
	h.GetRule(a.DB, w, r)
}

func (a *App) CreateRule(w http.ResponseWriter, r *http.Request) {
	h.CreateRule(a.DB, w, r)
}
