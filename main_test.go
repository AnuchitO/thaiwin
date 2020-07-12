package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go.uber.org/zap"
)

func TestCheckInHandler(t *testing.T) {
	payload := new(bytes.Buffer)
	json.NewEncoder(payload).Encode(Check{ID: 1234, PlaceID: 4321})
	req := httptest.NewRequest("POST", "/checkin", payload)
	logger, _ := zap.NewDevelopment()
	logger = logger.With(zap.String("hostname", "hostname"))
	req = req.WithContext(context.WithValue(req.Context(), "logger", logger))
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
