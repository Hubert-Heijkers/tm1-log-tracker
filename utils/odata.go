package odata

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

var Verbose = true

type Client struct {
	http.Client
}

func (client *Client) ExecuteGETRequest(urlStr string) *http.Response {
	// Create new, GET, request
	req, _ := http.NewRequest("GET", urlStr, nil)
	// Add the OData-Version header
	req.Header.Add("OData-Version", "4.0")
	// We'll be expecting a JSON formatted response, set Accept header accordingly
	req.Header.Add("Accept", "application/json")
	if Verbose == true {
		fmt.Println(req.Method, req.URL)
	}
	// Execute the request
	resp, err := client.Do(req)
	// If no errors then return the response
	if err != nil {
		log.Fatal(err)
	}
	return resp
}

func (client *Client) ExecuteGETRequestEx(urlStr string, preReq func(*http.Request)) *http.Response {
	// Create new, GET, request
	req, _ := http.NewRequest("GET", urlStr, nil)
	// Add the OData-Version header
	req.Header.Add("OData-Version", "4.0")
	// We'll be expecting a JSON formatted response, set Accept header accordingly
	req.Header.Add("Accept", "application/json")
	// Allow additional processing of the request before actually executing
	preReq(req)
	if Verbose == true {
		fmt.Println(req.Method, req.URL)
	}
	// Execute the request
	resp, err := client.Do(req)
	// If no errors then return the response
	if err != nil {
		log.Fatal(err)
	}
	return resp
}

func (client *Client) ExecutePOSTRequest(urlStr, contentType, body string) *http.Response {
	// Create new, POST, request
	req, _ := http.NewRequest("POST", urlStr, strings.NewReader(body))
	req.Header.Add("Content-Type", contentType)
	// Add the OData-Version header
	req.Header.Add("OData-Version", "4.0")
	// We'll be expecting a JSON formatted response, set Accept header accordingly
	req.Header.Add("Accept", "application/json")
	if Verbose == true {
		fmt.Println(req.Method, req.URL)
		fmt.Println(body)
	}
	// Execute the request
	resp, err := client.Do(req)
	// If no errors then return the response
	if err != nil {
		log.Fatal(err)
	}
	return resp
}

func (client *Client) ExecutePOSTRequestEx(urlStr, contentType, body string, preReq func(*http.Request)) *http.Response {
	// Create new, POST, request
	req, _ := http.NewRequest("POST", urlStr, strings.NewReader(body))
	req.Header.Add("Content-Type", contentType)
	// Add the OData-Version header
	req.Header.Add("OData-Version", "4.0")
	// We'll be expecting a JSON formatted response, set Accept header accordingly
	req.Header.Add("Accept", "application/json")
	// Allow additional processing of the request before actually executing
	preReq(req)
	if Verbose == true {
		fmt.Println(req.Method, req.URL)
		fmt.Println(body)
	}
	// Execute the request
	resp, err := client.Do(req)
	// If no errors then return the response
	if err != nil {
		log.Fatal(err)
	}
	return resp
}

func (client *Client) IterateCollection(datasourceServiceRootURL string, urlStr string, processResponse func([]byte) (int, string)) {
	// Set up the request to retrieve the collection given the passed url
	// Note: While we are requesting the collection completely in one request, the service might
	// opt to apply server driven paging and give us a partial response with a nextLink which
	// subsequently can be used to retrieve the next chunk or remainder of the collection.
	for nextLink := urlStr; nextLink != ""; {
		resp := client.ExecuteGETRequest(datasourceServiceRootURL + nextLink)
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		if Verbose == true {
			fmt.Println(string(body))
		}

		// Process the response
		_, nextLink = processResponse(body)
	}
}

func (client *Client) TrackCollection(serviceRootURL string, urlStr string, interval time.Duration, processResponse func([]byte) (string, string)) {
	// Set up the request to retrieve the collection given the passed url
	// Note: While we are requesting the collection completely in one request, the service might
	// opt to apply server driven paging and give us a partial response with a nextLink which
	// subsequently can be used to retrieve the next chunk or remainder of the collection.
	for urlStr := urlStr; urlStr != ""; {
		resp := client.ExecuteGETRequestEx(serviceRootURL+urlStr, func(req *http.Request) { req.Header.Add("Prefer", "odata.track-changes") })
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		if Verbose == true {
			fmt.Println(string(body))
		}

		// Process the response
		nextLink, deltaLink := processResponse(body)

		// TM1 doesn't but other services could return a nextLink when applying server side windowing
		// while returning the collection. Note that, following OData conventions, only the last
		// window, which does not have a nextLink, contains a deltaLink.
		if nextLink != "" {
			// Continue processing the collection being returned
			urlStr = nextLink
		} else if deltaLink != "" {
			// Wait a second before querying for the next deltaLink
			time.Sleep(interval)

			// Continue with the deltaLink
			urlStr = deltaLink
		} else {
			// Seems the server is no longer willing to give us deltas.
			break
		}
	}
}

func ValidateStatusCode(resp *http.Response, statusCode int, logFmt func() string) {
	if resp.StatusCode != statusCode {
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		log.Fatal(logFmt() + "\r\nServer responded with: " + resp.Status + "\r\n" + string(body))
	}
}
