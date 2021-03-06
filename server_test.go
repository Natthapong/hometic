package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockPairDevice struct{}

func (mockPairDevice) Pair(p Pair) error {
	return nil
}

func TestCreatePairDevice(t *testing.T) {
	payload := new(bytes.Buffer)
	json.NewEncoder(payload).Encode(Pair{DeviceID: 1234, UserID: 4433})
	req := httptest.NewRequest(http.MethodPost, "/pair-device", payload)
	rec := httptest.NewRecorder()

	//handler := PairDeviceHandler(mockPairDevice{})

	// handler := PairDeviceHandler(CreatePairDeviceFunc(func(p Pair) error {
	// 	return nil
	// }))

	create := func(p Pair) error {
		return nil
	}
	handler := CustomHandlerFunc(PairDeviceHandler(CreatePairDeviceFunc(create)))

	handler.ServeHTTP(rec, req)

	if http.StatusOK != rec.Code {
		t.Error("expect 200 OK but got", rec.Code)
	}

	expected := fmt.Sprintf("%s\n", `{"status":"active"}`)
	if rec.Body.String() != expected {
		t.Errorf("expected %q but got %q\n", expected, rec.Body.String())
	}
}
