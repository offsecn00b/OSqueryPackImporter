package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	//"net/http/httputil"
)

var (
	token    = flag.String("token", "", "Kolide authentication token")
	hostName = flag.String("hostname", "https://localhost:8080", "Kolide server hostname")
	packDir  = flag.String("pack_dir", "", "Directory of packs")
)

var httpClient = &http.Client{
	Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	},
}

type serverError struct {
	Message string `json:"message"`
	Errors  []struct {
		Name   string `json:"name"`
		Reason string `json:"reason"`
	} `json:"errors"`
}

func getQuery(name string) (uint, error) {
	//Define JSON struct for the Queries returned
	type Message struct {
		CreatedAt   string        `json:"created_at"`
		UpdatedAt   string        `json:"updated_at"`
		DeletedAt   interface{}   `json:"deleted_at"`
		Deleted     bool          `json:"deleted"`
		ID          int           `json:"id"`
		Name        string        `json:"name"`
		Description string        `json:"description"`
		Query       string        `json:"query"`
		Saved       bool          `json:"saved"`
		AuthorID    int           `json:"author_id"`
		AuthorName  string        `json:"author_name"`
		Packs       []interface{} `json:"packs"`
	}
	//Define array of query structs returned
	type Queries struct {
		Message []Message `json:"queries"`
	}
	//Build new request to get all queries into an array
	request, err := http.NewRequest(
		"GET",
		*hostName+"/api/v1/kolide/queries", nil)

	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", *token))

	if err != nil {
		return 0, errors.New("Error creating request for query list lookup: " + err.Error())
	}

	//debug(httputil.DumpRequestOut(request, true))
	response, err := httpClient.Do(request)
	if err != nil {
		return 0, errors.New("Error making request:" + err.Error())
	}
	defer response.Body.Close()
	//debug(httputil.DumpResponse(response, true))

	//read the response json data and format/decode as array of Message structs
	jsonDataFromHttp, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}
	var msg Queries
	err = json.Unmarshal([]byte(jsonDataFromHttp), &msg) // here!

	if err != nil {
		panic(err)
	}

	//loop through all queries and see if they match the query name passed to the function
	//If so return that query ID
	for i := range msg.Message {
		if msg.Message[i].Name == name {
			fmt.Println("[+] We found the query ID: %s", msg.Message[i].ID)
			return uint(msg.Message[i].ID), nil
		}
	}

	return 100, nil

}

//dump full requests responses for debugging
func debug(data []byte, err error) {
	if err == nil {
		fmt.Printf("%s\n\n", data)
	} else {
		log.Fatalf("%s\n\n", err)
	}
}

func createPack(name, description string) (uint, error) {
	type createPackRequest struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	type createPackResponse struct {
		Pack struct {
			ID uint `json:"id"`
		} `json:"pack"`
		serverError
	}

	body := createPackRequest{
		Name:        name,
		Description: description,
	}
	b, err := json.Marshal(body)
	if err != nil {
		return 0, err
	}
	request, err := http.NewRequest(
		"POST",
		*hostName+"/api/v1/kolide/packs",
		bytes.NewBuffer(b),
	)
	if err != nil {
		return 0, errors.New("Error creating request object: " + err.Error())
	}

	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", *token))
	response, err := httpClient.Do(request)
	if err != nil {
		return 0, errors.New("Error making request:" + err.Error())
	}
	defer response.Body.Close()

	var responseBody createPackResponse
	err = json.NewDecoder(response.Body).Decode(&responseBody)
	if err != nil {
		return 0, errors.New("Error decoding HTTP response body")
	}

	if len(responseBody.Errors) != 0 {
		errs := []string{}
		for _, e := range responseBody.Errors {
			errs = append(errs, e.Reason)
		}
		return 0, errors.New(strings.Join(errs, ";"))
	}

	return responseBody.Pack.ID, nil
}

func createQuery(name, query, description string) (uint, error) {
	type createQueryRequest struct {
		Name        string `json:"name"`
		Query       string `json:"query"`
		Description string `json:"description"`
	}
	type createQueryResponse struct {
		Query struct {
			ID uint `json:"id"`
		} `json:"query"`
		serverError
	}

	body := createQueryRequest{
		Name:        name,
		Query:       query,
		Description: description,
	}

	b, err := json.Marshal(body)
	if err != nil {
		return 0, err
	}
	request, err := http.NewRequest(
		"POST",
		*hostName+"/api/v1/kolide/queries",
		bytes.NewBuffer(b),
	)
	if err != nil {
		return 0, errors.New("Error creating request object: " + err.Error())
	}

	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", *token))

	//debug(httputil.DumpRequestOut(request, true))

	response, err := httpClient.Do(request)
	if err != nil {
		return 0, errors.New("Error making request:" + err.Error())
	}
	defer response.Body.Close()
	//debug(httputil.DumpResponse(response, true))

	thiscode := strconv.Itoa(response.StatusCode)
	var responseBody createQueryResponse
	err = json.NewDecoder(response.Body).Decode(&responseBody)
	// HTTP.StatusCode 409 is returned when the create query api call doesn't execute because there is already a record under that query name
	// So we must call a new function "getQuery" by passing the query name, to do a lookup by name and return the query ID that already exists and return that ID so the query can be added to the pack
	if thiscode == "409" {
		fmt.Println("[+] We can't add this query name because it already exists....")
		fmt.Println("[+] Attempting to enumerate the Query ID by Name.....")
		testID, err := getQuery(name)
		if err != nil {
			return 0, errors.New("Error retrieving Existing query ID:" + err.Error())
		}
		return testID, nil
	}

	if err != nil {
		return 0, errors.New("Error decoding HTTP response body")
	}

	if len(responseBody.Errors) != 0 {
		errs := []string{}
		for _, e := range responseBody.Errors {
			errs = append(errs, e.Reason)
		}
		return 0, errors.New(strings.Join(errs, ";"))
	}

	return responseBody.Query.ID, nil
}

func addQueryToPack(packID, queryID uint, interval uint64, snapshot, removed bool) (uint, error) {
	type addQueryToPackRequest struct {
		PackID   uint   `json:"pack_id"`
		QueryID  uint   `json:"query_id"`
		Interval uint64 `json:"interval"`
		Snapshot bool   `json:"snapshot"`
		Removed  bool   `json:"removed"`
	}
	type addQueryToPackResponse struct {
		Scheduled struct {
			ID uint `json:"id"`
		} `json:"scheduled"`
		serverError
	}

	body := addQueryToPackRequest{
		PackID:   packID,
		QueryID:  queryID,
		Interval: interval,
		Snapshot: snapshot,
		Removed:  removed,
	}
	b, err := json.Marshal(body)
	if err != nil {
		return 0, err
	}
	request, err := http.NewRequest(
		"POST",
		*hostName+"/api/v1/kolide/schedule",
		bytes.NewBuffer(b),
	)
	if err != nil {
		return 0, errors.New("Error creating request object: " + err.Error())
	}

	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", *token))
	response, err := httpClient.Do(request)
	if err != nil {
		return 0, errors.New("Error making request:" + err.Error())
	}
	defer response.Body.Close()

	var responseBody addQueryToPackResponse
	err = json.NewDecoder(response.Body).Decode(&responseBody)
	if err != nil {
		return 0, errors.New("Error decoding HTTP response body")
	}

	if len(responseBody.Errors) != 0 {
		errs := []string{}
		for _, e := range responseBody.Errors {
			errs = append(errs, e.Reason)
		}
		return 0, errors.New(strings.Join(errs, ";"))
	}

	return responseBody.Scheduled.ID, nil
}

type queryStanza struct {
	Query       string      `json:"query"`
	Interval    interface{} `json:"interval"`
	Description string      `json:"description"`
	Version     string      `json:"version"`
	Value       string      `json:"value"`
	Snapshot    *bool       `json:"snapshot"`
	Removed     *bool       `json:"removed"`
}

type packFile struct {
	Queries map[string]queryStanza `json:"queries"`
}

func init() {
	flag.Parse()

	if packDir == nil || *packDir == "" {
		*packDir = "."
	}
}

func main() {
	files, err := ioutil.ReadDir(*packDir)
	if err != nil {
		log.Fatalln("Could not list files:", err)
	}
	for _, file := range files {
		content, err := os.Open(file.Name())
		if err != nil {
			log.Fatalf("Could not open file at path %s: %s", file, err)
		}
		defer content.Close()

		var pack packFile
		jsonParser := json.NewDecoder(content)
		if err = jsonParser.Decode(&pack); err != nil {
			log.Printf("%s is not a query pack... skipping.", file.Name())
			continue
		}

		packID, err := createPack(file.Name(), "")
		if err != nil {
			log.Fatalf("Error creating pack %s: %s", file.Name(), err)
		}
		log.Printf("Created pack %s (%d)", file.Name(), packID)
		for name, query := range pack.Queries {
			interval, err := convertToUint64(query.Interval)
			if err != nil {
				log.Fatalln(err)
			}
			queryID, err := createQuery(name, query.Query, query.Description)
			if err != nil {
				log.Fatalf("Error creating query %s: %s", name, err)
			}
			log.Printf("Created query %s (%d)", name, queryID)

			removed := true
			if query.Removed != nil {
				removed = *query.Removed
			}

			snapshot := false
			if query.Snapshot != nil {
				snapshot = *query.Snapshot
			}

			if queryID != 0 {
				scheduledID, err := addQueryToPack(packID, queryID, interval, snapshot, removed)
				if err != nil {
					log.Fatalf("Error scheduling query %s: %s", name, err)
				}

				log.Printf("Added query %s to pack %s (%d)", name, file.Name(), scheduledID)
			}

		}
	}
}

func convertToUint64(input interface{}) (uint64, error) {
	switch i := input.(type) {
	case int, float64, uint, uint64:
		return uint64(i.(float64)), nil
	case string:
		value, err := strconv.ParseUint(i, 10, 64)
		if err != nil {
			return 0, errors.New("Error converting string to uint: " + err.Error())
		}
		return value, nil
	default:
		return 0, errors.New("Got an unacceptable type for interval")
	}

}
