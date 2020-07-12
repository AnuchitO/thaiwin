package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCheckInHandler(t *testing.T) {
	payload := new(bytes.Buffer)
	json.NewEncoder(payload).Encode(Check{ID: 1234, PlaceID: 4321})
	req := httptest.NewRequest("POST", "/checkin", payload)
	w := httptest.NewRecorder()

	var fn InFunc = func(id, placeID int64) error {
		return nil
	}

	handler := CheckIn(fn)

	handler(w, req)

	if w.Code != http.StatusOK {
		t.Error("not ok")
	}
}

func TestSealMiddleware(t *testing.T) {
	payload := bytes.NewBufferString("ewogICAgImlkIjogMTIzNCwKICAgICJwbGFjZV9pZCI6IDQzMjEKfQ==")

	req := httptest.NewRequest("POST", "/checkin", payload)
	w := httptest.NewRecorder()

	handler := SealMiddleware()

	handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Error(err)
			return
		}
		w.Write(b)
	})).ServeHTTP(w, req)

	body, _ := ioutil.ReadAll(w.Result().Body)
	fmt.Println(string(body))
}
