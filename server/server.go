package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Log struct {
	Id string `json:"id"`
	Date string `json:"date"`
	Notes string `json:"notes"`
	Goal string `json:"goal"`
	GoalAccomplished bool `json:"goalAccomplished,boolean"`
}

var LOGS = make([]Log, 0) // use make with size 0 to get an empty array as the JSON output

func genUUID() string {
	randomStr := ""
	for len(randomStr) <= 6 {
		randomStr += strconv.FormatInt(int64(rand.Intn(36)), 36)
	}
	return strconv.FormatInt(time.Now().Unix(), 16) + randomStr
}

func validateDateString(str string) bool {
	_, err := time.Parse("2006-01-02", str)

	return err == nil
}

func validateLog(log Log) bool {
	return validateDateString(log.Date)
}

func respondJSON( writer http.ResponseWriter, body interface{}) {
	writer.Header().Set("Content-Type", "application/json")

	json.NewEncoder(writer).Encode(body)
}

func getLogs(writer http.ResponseWriter, _ *http.Request) {
	respondJSON(writer, LOGS)
}

func getLog(writer http.ResponseWriter, id string) {
	for _, log := range LOGS {
		if log.Id == id {
			respondJSON(writer, log)
			return
		}
	}

	http.Error(writer, fmt.Sprintf("Could not find log with id: %s", id), http.StatusNotFound)
}

func updateLog(writer http.ResponseWriter, request *http.Request, id string) {
	var newLog Log

	err := json.NewDecoder(request.Body).Decode(&newLog)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	if !validateLog(newLog) {
		http.Error(writer, "invalid log", http.StatusBadRequest)
		return
	}

	for i, log := range LOGS {
		if log.Id == id {
			newLog.Id = id
			LOGS[i] = newLog

			respondJSON(writer, log)
		}
	}
}

func postLog(writer http.ResponseWriter, request *http.Request) {
	var newLog Log

	err := json.NewDecoder(request.Body).Decode(&newLog)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	if !validateLog(newLog) {
		http.Error(writer, "invalid log", http.StatusBadRequest)
		return
	}

	newLog.Id = genUUID()
	LOGS = append(LOGS, newLog)

	respondJSON(writer, newLog)
}

func idHandler(writer http.ResponseWriter, request *http.Request, id string) {
	switch request.Method {
	case http.MethodGet:
		getLog(writer, id)
	case http.MethodPut:
		updateLog(writer, request, id)
	default:
		http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func baseHandler(writer http.ResponseWriter, request *http.Request) {
	pathParts := strings.Split(request.URL.Path, "/")
	if len(pathParts) > 2 {
		http.Error(writer, "Invalid URL", http.StatusNotFound)
	} else if len(pathParts) == 2 && request.URL.Path != "/" {
		idHandler(writer, request, pathParts[1])
  	} else {
		switch request.Method {
		case http.MethodGet:
			getLogs(writer, request)
		case http.MethodPost:
			postLog(writer, request)
		default:
			http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}

}

func main() {
	http.HandleFunc("/", baseHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}