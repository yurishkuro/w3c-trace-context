package xhttp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// PostJSON converts in to JSON and posts to url, the parses the response into out.
func PostJSON(ctx context.Context, url string, decorate func(*http.Request), in, out interface{}) error {
	jsonBytes, err := json.Marshal(in)
	if err != nil {
		return err
	}

	r, err := http.NewRequest("POST", url, bytes.NewReader(jsonBytes))
	if err != nil {
		return err
	}

	if decorate != nil {
		decorate(r)
	}

	log.Printf("PUT %s %s", url, string(jsonBytes))
	r = r.WithContext(ctx)
	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return err
	}

	if c := resp.StatusCode; c != 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("expecting 200, received %d: %s", c, string(body))
	}

	return json.NewDecoder(resp.Body).Decode(out)
}

type handlerFunc func(r *http.Request, in interface{}) (interface{}, error)

// HandleJSON performs JSON in and out marshalling and deelegates to the handler function h
func HandleJSON(w http.ResponseWriter, r *http.Request, in interface{}, h handlerFunc) {
	d := json.NewDecoder(r.Body)
	err := d.Decode(in)
	if err != nil {
		http.Error(w, "cannot parse request: "+err.Error(), http.StatusBadRequest)
		return
	}
	out, err := h(r, in)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	jsonBytes, err := json.Marshal(out)
	if err != nil {
		http.Error(w, "cannot marshal response: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(jsonBytes)
}
