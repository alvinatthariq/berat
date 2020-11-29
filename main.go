package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

var (
	Weights  map[string]Weight
	myRouter *mux.Router
)

type Weight struct {
	Date string `json:"date"`
	Min  int    `json:"min"`
	Max  int    `json:"max"`
	Diff int    `json:"diff"`
}

type WeightResp struct {
	Weight Weight `json:"weight"`
}

type WeightsResp struct {
	Weights []Weight `json:"weights"`
}

func getWeightList(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: returnWeightList")
	weightsData := make([]Weight, 0)
	for _, weight := range Weights {
		weightsData = append(weightsData, weight)
	}
	json.NewEncoder(w).Encode(WeightsResp{Weights: weightsData})
}

func getWeightByDate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	date := vars["date"]

	// validate date
	_, err := time.Parse("2006-01-02", date)
	if err != nil {
		fmt.Printf("Invalid Date : %v", date)
		w.WriteHeader(400)
		return
	}

	if _, ok := Weights[date]; !ok {
		w.WriteHeader(404)
		return
	}

	json.NewEncoder(w).Encode(WeightResp{Weight: Weights[date]})
}

func createNewWeight(w http.ResponseWriter, r *http.Request) {
	// read request body
	reqBody, _ := ioutil.ReadAll(r.Body)
	var weight Weight
	json.Unmarshal(reqBody, &weight)

	// validate date
	_, err := time.Parse("2006-01-02", weight.Date)
	if err != nil {
		fmt.Printf("Invalid Date %v", weight.Date)
		w.WriteHeader(400)
		return
	}

	// count diff
	weight.Diff = weight.Max - weight.Min
	Weights[weight.Date] = weight

	json.NewEncoder(w).Encode(WeightResp{Weight: weight})
}

func deleteWeightByDate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	date := vars["date"]

	// validate date
	_, err := time.Parse("2006-01-02", date)
	if err != nil {
		fmt.Printf("Invalid Date %v", date)
		w.WriteHeader(400)
		return
	}

	delete(Weights, date)
}

func handleRequests() {
	myRouter = mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/weights", getWeightList).Methods("GET")
	myRouter.HandleFunc("/weight", createNewWeight).Methods("POST")
	myRouter.HandleFunc("/weight/{date}", deleteWeightByDate).Methods("DELETE")
	myRouter.HandleFunc("/weight/{date}", getWeightByDate).Methods("GET")
	log.Fatal(http.ListenAndServe(":10000", myRouter))
}

func main() {
	Weights = make(map[string]Weight)
	handleRequests()
}
