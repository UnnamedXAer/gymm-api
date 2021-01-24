package controllers

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
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

func TestGetUserById(t *testing.T) {
	log.Println(" [TestGetUserById] start")
	req, err := http.NewRequest("GET", "/users/600dab27ec9c7bb884438d78", nil)
	if err != nil {
		t.Fatal(err)
	}
	q := req.URL.Query()
	q.Add("id", "600dab27ec9c7bb884438d78")
	req.URL.RawQuery = q.Encode()
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetUserById)
	http.Handle("/users/{id}", handler)
	handler.ServeHTTP(rr, req)
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
