package main

import (
	"bytes"
	"encoding/json"
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
