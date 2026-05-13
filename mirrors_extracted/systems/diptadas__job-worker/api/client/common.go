package client

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"time"
)

// getHTTPClient returns a http.Client for a Client object.
// It returns error when TLS certificates can not be loaded.
func (c Client) getHTTPClient() (*http.Client, error) {
	cert, err := ioutil.ReadFile(c.CaCert)
	if err != nil {
		return nil, fmt.Errorf("could not read CA certificate, reason: %v", err)
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(cert)

	certificate, err := tls.LoadX509KeyPair(c.ClientCert, c.ClientKey)
	if err != nil {
		return nil, fmt.Errorf("could not load certificate, reason: %v", err)
	}

	return &http.Client{
		Timeout: time.Minute * 3,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs:      caCertPool,
				Certificates: []tls.Certificate{certificate},
			},
		},
	}, nil
}

// prepareRequest returns a http.Request for a given path, http method, and request body.
// It adds username to the Authorization header.
func (c Client) prepareRequest(apiPath string, method string, body interface{}) (*http.Request, error) {
	reqUrl, err := getURL(c.Address, apiPath)
	if err != nil {
		return nil, err
	}

	r, w := io.Pipe()

	if body != nil {
		go func() {
			_ = json.NewEncoder(w).Encode(body) // TODO: error handling
			_ = w.Close()                       // TODO: error handling
		}()
	} else {
		_ = w.Close() // TODO: error handling
	}

	req, err := http.NewRequest(method, reqUrl, r)
	if err != nil {
		return nil, fmt.Errorf("could not prepare request, reason: %v", err)
	}

	req.Header.Add("Content-Type", "application/json")

	return req, nil
}

// makeAPICall sends a request to API server and returns the response.
func (c Client) makeAPICall(apiPath string, method string, body interface{}) (string, error) {
	client, err := c.getHTTPClient()
	if err != nil {
		return "", fmt.Errorf("could not prepare client, reason: %v", err)
	}

	req, err := c.prepareRequest(apiPath, method, body)
	if err != nil {
		return "", err
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response, reason: %v", err)
	}
	return string(respBody), nil
}

// getURL merges address with path.
func getURL(address string, apiPath string) (string, error) {
	u, err := url.Parse(address)
	if err != nil {
		return "", fmt.Errorf("could not parse address, reason: %v", err)
	}
	u.Path = path.Join(u.Path, apiPath)
	return u.String(), nil
}
