package twilio

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestCallService_Call(t *testing.T) {
	setup()
	defer teardown()

	u := client.EndPoint("Calls")

	output := `{
		"sid": "abcdef",
		"price": "0.74"
	}`

	mux.HandleFunc(u.String(), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		fmt.Fprint(w, output)
	})

	params := CallParams{
		From: "+14158141829",
		To:   "+15558675309",
		Url:  "http://www.example.com/call.xml",
	}

	m, _, err := client.Calls.Call(params)

	if err != nil {
		t.Errorf("Calls.Call returned error: %q", err)
	}

	want := &Call{
		Sid:   "abcdef",
		Price: 0.74,
	}

	if !reflect.DeepEqual(m, want) {
		t.Errorf("Calls.Call returned %+v, want %+v", m, want)
	}
}

func TestCallService_Call_httpError(t *testing.T) {
	setup()
	defer teardown()

	u := client.EndPoint("Calls")

	mux.HandleFunc(u.String(), func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Bad Request", http.StatusBadRequest)
	})

	params := CallParams{
		From: "+14158141829",
		To:   "+15558675309",
		Url:  "http://www.example.com/call.xml",
	}

	_, _, err := client.Calls.Call(params)

	if err == nil {
		t.Fatal("Expected HTTP 400 error.")
	}
}

func TestCallService_Call_invalidParams(t *testing.T) {
	setup()
	defer teardown()

	u := client.EndPoint("Calls")

	mux.HandleFunc(u.String(), func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Bad Request", http.StatusBadRequest)
	})

	params := CallParams{
		From:           "+14158141829",
		To:             "+15558675309",
		Url:            "http://www.example.com/call.xml",
		ApplicationSid: "1234",
	}

	_, _, err := client.Calls.Call(params)

	if err != ErrUrlXorApplicationSid {
		t.Fatalf("Expected error: %s", ErrUrlXorApplicationSid.Error())
	}
}

func TestCallService_Get(t *testing.T) {
	setup()
	defer teardown()

	sid := "MM90c6fc909d8504d45ecdb3a3d5b3556e"
	u := client.EndPoint("Calls", sid)

	output := `{
		"account_sid": "AC5ef8732a3c49700934481addd5ce1659",
		"date_created": "Wed, 18 Aug 2010 20:01:40 +0000",
		"date_sent": null,
		"date_updated": "Wed, 18 Aug 2010 20:01:40 +0000",
		"direction": "outbound-api",
		"from": "+14158141829",
		"price": null,
		"sid": "MM90c6fc909d8504d45ecdb3a3d5b3556e",
		"status": "queued",
		"to": "+15558675309",
		"uri": "/2010-04-01/Accounts/AC5ef87/Calls/MM90c6.json"
	}`

	mux.HandleFunc(u.String(), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, output)
	})

	m, _, err := client.Calls.Get(sid)

	if err != nil {
		t.Errorf("Call.SendSMS returned error: %v", err)
	}

	tm := parseTimestamp("Wed, 18 Aug 2010 20:01:40 +0000")
	want := &Call{
		AccountSid:  "AC5ef8732a3c49700934481addd5ce1659",
		DateCreated: tm,
		DateUpdated: tm,
		Direction:   "outbound-api",
		From:        "+14158141829",
		Price:       0,
		Sid:         "MM90c6fc909d8504d45ecdb3a3d5b3556e",
		Status:      "queued",
		To:          "+15558675309",
		Uri:         "/2010-04-01/Accounts/AC5ef87/Calls/MM90c6.json",
	}

	if !reflect.DeepEqual(m, want) {
		t.Errorf("Call.Get returned %+v, want %+v", m, want)
	}
}

func TestCallService_Get_httpError(t *testing.T) {
	setup()
	defer teardown()

	u := client.EndPoint("Calls", "abc")

	mux.HandleFunc(u.String(), func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Bad Request", http.StatusBadRequest)
	})

	_, _, err := client.Calls.Get("abc")

	if err == nil {
		t.Error("Expected HTTP 400 errror.")
	}
}
