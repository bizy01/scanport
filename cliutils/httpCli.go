package cliutils

import (
	"strings"
	"log"
	"net/http"
	"crypto/tls"
	"io/ioutil"
)

func HttpCli(method, path, body string, headers map[string]string) (int, interface{}) {
	if headers == nil {
		headers = map[string]string{}
	}

	method = strings.ToUpper(method)
	if method != "GET" && method != "POST" && method != "PUT" && method != "DELETE" {
		log.Printf("E! Unsupported HTTP Method: %s", method)
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	httpClient := &http.Client{Transport: tr}

	req, err := http.NewRequest(method, path, strings.NewReader(body))
	if err != nil {
		log.Printf("E! http request error: %v", err)
		return -1, err
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		log.Printf("E! http request error: %v", err)
		return -1, err
	}

	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("E! http read body error: %v", err)
		return -1, err
	}

	respStatusCode := resp.StatusCode
	respData := string(respBody)

	return respStatusCode, respData
}
