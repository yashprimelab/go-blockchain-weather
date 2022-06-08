package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"golang-blockchain/blockchain"

	"github.com/gorilla/mux"
)

type BlockSuccess struct {
	Message string `json:"message"`
}
type ServerSetup struct {
	Status string
}

func getBlocks(w http.ResponseWriter, r *http.Request) {

	var tmpRecords []blockchain.Block
	iter := blockchain.InitBlockChain().Iterator()

	for {
		block := iter.Next()

		pow := blockchain.NewProof(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		tmpRecords = append(tmpRecords, *block)
		if len(block.PrevHash) == 0 {
			break
		}
	}
	setupHeader(w)
	json.NewEncoder(w).Encode(tmpRecords)
	iter.Database.Close()
}

func checkServer(w http.ResponseWriter, r *http.Request) {

	var newEmployee = ServerSetup{Status: "Server is in running state"}
	setupHeader(w)
	json.NewEncoder(w).Encode(newEmployee)
}

func createBlock(w http.ResponseWriter, r *http.Request) {
	setupHeader(w)
	var inst = blockchain.InitBlockChain()
	city := r.FormValue("city")
	response, err := http.Get("http://api.weatherapi.com/v1/current.json?key=b66dc8096d394253871181414220806&q=" + city + "&aqi=yes")

	if err != nil {
		fmt.Print(err.Error())
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Weather API's Response: \n")
	fmt.Println(string(responseData))
	fmt.Printf("\n")

	inst.AddBlock(string(responseData))
	json.NewEncoder(w).Encode(BlockSuccess{Message: "Block has been added successfully."})
	inst.Database.Close()
}

func setupHeader(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")
}

func main() {

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", checkServer)
	router.HandleFunc("/add-weather-block", createBlock).Methods("POST")
	router.HandleFunc("/get-blocks", getBlocks).Methods("GET")
	log.Fatal(http.ListenAndServe(":8080", router))
}
