package handlers

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"path"
)

const (
	NumOfRetry = 10
)

func proxyRequest(r *http.Request) (*http.Response, error) {
	body, _ := ioutil.ReadAll(r.Body)
	for i := 0; i <= NumOfRetry; i++ {
		u := url.URL{
			Scheme:   "http",
			Host:     currentLeader.String(),
			Path:     path.Join("v2", "keys", "_etcd", "registry", r.URL.Path),
			RawQuery: r.URL.RawQuery,
		}

		buf := bytes.NewBuffer(body)
		outreq, err := http.NewRequest(r.Method, u.String(), buf)
		if err != nil {
			return nil, err
		}

		copyHeader(outreq.Header, r.Header)

		client := http.Client{}
		resp, err := client.Do(outreq)
		if err != nil {
			return nil, err
		}
		log.Println("[proxyRequest] url is: ", u)

		// Try again on the next host
		if resp.StatusCode == 307 &&
			(r.Method == "PUT" || r.Method == "DELETE") {
			u, err := resp.Location()
			if err != nil {
				return nil, err
			}
			currentLeader.Set(u.Host)
			continue
		}

		return resp, nil
	}

	return nil, errors.New("All attempts at proxying to etcd failed")
}

// copyHeader copies all of the headers from dst to src.
func copyHeader(dst, src http.Header) {
	for k, v := range src {
		for _, q := range v {
			dst.Add(k, q)
		}
	}
}

func (h *Handler) TokenHandler(w http.ResponseWriter, r *http.Request) {
	resp, err := proxyRequest(r)
	if err != nil {
		log.Printf("Error making request: %v", err)
		http.Error(w, "", 500)
	}

	copyHeader(w.Header(), resp.Header)
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}
