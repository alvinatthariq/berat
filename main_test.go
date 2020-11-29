package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

var testData = Weight{
	Date: "2020-10-11",
	Max:  50,
	Min:  40,
}

func TestWeight(t *testing.T) {
	initTest()

	testCreateWeight(t)
	testGetWeightList(t)
	testGetWeightByDate(t)
}

func testGetWeightList(t *testing.T) {
	req, _ := http.NewRequest("GET", "/weights", nil)
	myRouter.HandleFunc("/weights", getWeightList).Methods("GET")
	response := executeRequest(req)

	// assert http status code
	if http.StatusOK != response.Code {
		t.Errorf("Expected response code %d. Got %d\n", http.StatusOK, response.Code)
	}

	// assert resp
	var result WeightsResp
	rawData, _ := ioutil.ReadAll(response.Body)
	if err := json.Unmarshal(rawData, &result); err != nil {
		t.Error("Unmarshal Resp ", err.Error())
	}

	if len(result.Weights) == 0 {
		t.Error("Expected list not empty ")
	}
}

func testCreateWeight(t *testing.T) {

	body, err := json.Marshal(testData)
	if err != nil {
		t.Error("Marshal Body Req", err.Error())
	}

	req, _ := http.NewRequest("POST", "/weight", bytes.NewBuffer(body))
	myRouter.HandleFunc("/weight", createNewWeight).Methods("POST")
	response := executeRequest(req)

	// assert http status code
	if http.StatusOK != response.Code {
		t.Errorf("Expected response code %d. Got %d\n", http.StatusOK, response.Code)
	}

	// assert resp
	var result WeightResp
	rawData, _ := ioutil.ReadAll(response.Body)
	if err := json.Unmarshal(rawData, &result); err != nil {
		t.Error("Unmarshal Resp ", err.Error())
	}

	if result.Weight.Date != testData.Date {
		t.Errorf("Expected date equal to %s. Got %s", testData.Date, result.Weight.Date)
	}

	if result.Weight.Min != testData.Min {
		t.Errorf("Expected Min equal to %d. Got %d", testData.Min, result.Weight.Min)
	}

	if result.Weight.Max != testData.Max {
		t.Errorf("Expected Max equal to %d. Got %d", testData.Max, result.Weight.Max)
	}
}

func testGetWeightByDate(t *testing.T) {
	path := fmt.Sprintf("/weight/%s", "2020-10-11")

	req, _ := http.NewRequest("GET", path, nil)
	myRouter.HandleFunc("/weight/{date}", getWeightByDate).Methods("GET")
	response := executeRequest(req)

	// assert http status code
	if http.StatusOK != response.Code {
		t.Errorf("Expected response code %d. Got %d\n", http.StatusOK, response.Code)
	}

	// assert resp
	var result WeightResp
	rawData, _ := ioutil.ReadAll(response.Body)
	if err := json.Unmarshal(rawData, &result); err != nil {
		t.Error("Unmarshal Resp ", err.Error())
	}

	if result.Weight.Date != testData.Date {
		t.Errorf("Expected date equal to %s. Got %s", testData.Date, result.Weight.Date)
	}

	if result.Weight.Min != testData.Min {
		t.Errorf("Expected Min equal to %d. Got %d", testData.Min, result.Weight.Min)
	}

	if result.Weight.Max != testData.Max {
		t.Errorf("Expected Max equal to %d. Got %d", testData.Max, result.Weight.Max)
	}
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	myRouter.ServeHTTP(rr, req)
	return rr
}

func initTest() {
	Weights = make(map[string]Weight)
	myRouter = mux.NewRouter().StrictSlash(true)
}
