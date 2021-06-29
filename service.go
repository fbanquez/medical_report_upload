package main

import (
	"crypto/tls"
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

/*
SendMR receives a JSON object with the data of the medical report and sends it
as a body to the call of the microservice hosted into viewmed platform.
*/
func sendMR(jsonObj string) (body []byte, err error) {

	// Creating the microservice's URI
	urlString := config.Service.Uri + ":" + config.Service.Port
	if strings.HasPrefix(config.Service.Endpoint, "/") {
		urlString += config.Service.Endpoint
	} else {
		urlString += "/" + config.Service.Endpoint
	}

	// Defining the HTTP method to use
	method := "POST"

	// Creating the proxy's URI
	urlProxy, err := url.Parse(config.Proxy.Host + ":" + config.Proxy.Port)
	if err != nil {
		Error.Println("Problems parsing URL proxy. ", err)
	}

	// Defining proxy credentials
	auth := config.Proxy.User + ":" + config.Proxy.Passwd
	basicAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))

	payload := strings.NewReader(jsonObj)

	// Defining the http client
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(urlProxy),
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	// Building a new request
	req, err := http.NewRequest(method, urlString, payload)
	if err != nil {
		Error.Println("Problems creating HTTP request. ", err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", config.Service.Auth)
	req.Header.Add("User-Agent", config.Service.Agent)
	req.Header.Add("Proxy-Authorization", basicAuth)

	// Performing the request
	Info.Println("Calling the viewmed's microservice.")
	res, err := client.Do(req)
	if err != nil {
		Error.Println("Problem sending HTTP request. ", err)
	}

	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		Error.Println("Problem manipulating response. ", err)
	}

	defer res.Body.Close()

	return
}
