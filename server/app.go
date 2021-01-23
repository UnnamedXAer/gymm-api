package server

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/unnamedxaer/gymm-api/repository"
	"github.com/unnamedxaer/gymm-api/server/controllers"
	"github.com/unnamedxaer/gymm-api/services"
)

type App struct {
	Router *mux.Router
}

func (app *App) InitializeApp() {
	var repo repository.IRepository
	repo = &repository.MongoRepository{}
	err := repo.Initialize(os.Getenv("MONGO_URI"))
	if err != nil {
		log.Panic(err)
	}

	services.UService.SetRepo(repo)

	app.Router = mux.NewRouter()
	app.addHandlers()
}

func (app *App) Run(addr string) {
	log.Println("server is up and running at " + addr)
	log.Fatalln(http.ListenAndServe(addr, app.Router))
}

func (app *App) addHandlers() {
	app.Router.HandleFunc("/users/{id:[0-9a-zA-Z]+}", controllers.GetUser).Methods("GET")
	app.Router.HandleFunc("/users", controllers.CreateUser).Methods("POST")
	app.Router.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		log.Println("[" + r.Method + "/] -> URL: " + r.RequestURI)
		rw.WriteHeader(http.StatusMethodNotAllowed)
	})
}
