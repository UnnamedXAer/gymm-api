package controllers

import (
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/unnamedxaer/gymm-api/repository"
	"github.com/unnamedxaer/gymm-api/services"
)

// func clearTables(repoType string) {
// 	switch repoType {
// 	case "mongo":
// 		clearTablesMongo()
// 	default:
// 		log.Fatalln("unsupported repo type")
// 	}
// }

// func clearTablesMongo() {
// 	repo := repository.MongoRepository{}
// 	err := repo.Initialize(os.Getenv("MONGO_UTI"))

// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer repo.Client.Disconnect(repo.Ctx)

// 	usersCol := repo.DB.Collection("users")

// 	results, err := usersCol.DeleteMany(repo.Ctx, nil)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	log.Println("users collection cleared with results: " + fmt.Sprintf("%v", results))
// }

// func assertResponseCode(t *testing.T, wanted, got int) {
// 	if wanted != got {
// 		t.Errorf("Expected response code to be %d, got %d", wanted, got)
// 	}
// }

// func TestTest(t *testing.T) {
// 	t.Fatal("Test test fatal 😀")
// }

func TestMain(m *testing.M) {
	godotenv.Overload("../../.test.env")
	log.Println(" [TestMain] start " + os.Getenv("ENV"))

	var repo repository.IRepository
	repo = &repository.MongoRepository{}
	mongoURI := os.Getenv("MONGO_URI")
	err := repo.Initialize(mongoURI)
	if err != nil {
		log.Panic(err)
	}

	services.UService.SetRepo(repo)

	code := m.Run()

	log.Println(" [TestMain] end")
	os.Exit(code)
}

func TestGetUserById(t *testing.T) {
	rURI := "/users/600efe8882b4c5ed52e7deb4"
	log.Println(" [TestGetUserById] start")
	router := mux.NewRouter()
	handler := http.HandlerFunc(GetUserById)
	router.Handle("/users/{id:[0-9a-zA-Z]+}", handler)
	req, err := http.NewRequest("GET", rURI, nil)
	if err != nil {
		t.Fatal(err)
	}

	req.RequestURI = rURI
	req.Host = "localhost:" + os.Getenv("PORT")
	req.RemoteAddr = "[::1]:65213"
	req.URL.RawPath = rURI
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := `{"id":1,"first_name":"Krish","last_name":"Bhanushali","email_address":"krishsb2405@gmail.com","phone_number":"0987654321"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
