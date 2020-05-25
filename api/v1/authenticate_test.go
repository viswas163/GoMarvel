package v1

import (
	"net/http"
	"testing"
)

func TestRunHello(t *testing.T) {
	request, _ := http.NewRequest("GET", "http://localhost:3001/", nil)
	response := executeRequest(request)
	checkResponseCode(http.StatusOK, response)
	_, err := checkResponseBody(response)
	if err != nil {
		t.Fatal("Error checking response : ", err)
	}
}

func TestRunAuth(t *testing.T) {
	request, _ := GetAuthRequest("characters")
	response := executeRequest(request)
	checkResponseCode(http.StatusOK, response)
	_, err := checkResponseBody(response)
	if err != nil {
		t.Fatal("Error checking response : ", err)
	}
}

// func testCheckResponseCode(t *testing.T, expected int, response *httptest.ResponseRecorder) {
// 	actual := response.Result().StatusCode
// 	if expected != actual {
// 		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
// 	}
// }

// func testCheckResponseBody(t *testing.T, response *httptest.ResponseRecorder) {
// 	body, err := ioutil.ReadAll(response.Body)
// 	if err != nil {
// 		t.Error("Error reading response body : ", err)
// 	}
// 	t.Log(string(body))
// }
