package http

// func TestGetUserByID(t *testing.T) {
// 	id := "1sadf3245df3245"

// 	req, _ := http.NewRequest(http.MethodGet, ("/users" + id), nil)
// 	req.Header.Set("Content-Type", "application/json")

// 	response := executeRequest(req)

// 	checkResponseCode(t, http.StatusOK, response.Code)

// 	//check returned obj
// }

// func TestGetUserByIDNotFound(t *testing.T) {
// 	id := "1sadf3245df3245" + "notfound"

// 	req, _ := http.NewRequest(http.MethodGet, ("/users" + id), nil)
// 	req.Header.Set("Content-Type", "application/json")

// 	response := executeRequest(req)

// 	checkResponseCode(t, http.StatusOK, response.Code)

// 	b, err := ioutil.ReadAll(response.Body)
// 	if err != nil {
// 		t.Fatalf("Expected response to be 'null', got error: %v", err)
// 	}
// 	if string(b) != "null" {
// 		t.Fatalf("Expected response to be 'null', got %q", b)
// 	}
// }

// func TestGetUserByIDInvalidID(t *testing.T) {
// 	id := "1sadf3245df3245" + "INVALIDID"

// 	req, _ := http.NewRequest(http.MethodGet, ("/users" + id), nil)
// 	req.Header.Set("Content-Type", "application/json")

// 	response := executeRequest(req)

// 	checkResponseCode(t, http.StatusBadRequest, response.Code)

// 	expectedErr := repositories.NewErrorInvalidID(id, "user")

// 	b, err := ioutil.ReadAll(response.Body)
// 	if err != nil {
// 		t.Fatalf("want response to be %v, got error: %v", expectedErr, err)
// 	}
// 	var data map[string]interface{}

// 	err = json.Unmarshal(b, &data)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	if data["error"] != expectedErr.Error() {
// 		t.Fatalf("want response to be like {\"error\": \"%s\"}, got %s", expectedErr.Error(), string(b))
// 	}
// }
