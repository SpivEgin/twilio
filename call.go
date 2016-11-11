package twilio

import (
	"errors"
	"net/http"
	"strings"
)

type CallService struct {
	client *Client
}

type Call struct {
	Sid            string    `json:"sid"`
	ParentCallSid  string    `json:"parent_call_sid"`
	DateCreated    Timestamp `json:"date_created,omitempty"`
	DateUpdated    Timestamp `json:"date_updated,omitempty"`
	AccountSid     string    `json:"account_sid"`
	To             string    `json:"to"`
	From           string    `json:"from"`
	PhoneNumberSid string    `json:"phone_number_sid"`
	Status         string    `json:"status"`
	StartTime      Timestamp `json:"start_time,omitempty"`
	EndTime        Timestamp `json:"end_time,omitempty"`
	Duration       int       `json:"duration"`
	Price          Price     `json:"price"`
	PriceUnit      string    `json:"price_unit"`
	Direction      string    `json:"direction"`
	AnsweredBy     string    `json:"answered_by,omitempty"`
	ForwardedFrom  string    `json:"forwarded_from"`
	ToFormatted    string    `json:"to_formatted"`
	FromFormatted  string    `json:"from_formatted"`
	CallerName     string    `json:"caller_name"`
	Uri            string    `json:"uri"`
}

type CallParams struct {
	From                          string
	To                            string
	Url                           string
	ApplicationSid                string
	Method                        string
	FallbackUrl                   string
	FallbackMethod                string
	StatusCallback                string
	StatusCallbackMethod          string
	StatusCallbackEvent           []string
	SendDigits                    string
	IfMachine                     string
	Timeout                       int
	Record                        bool
	RecordingChannels             string
	RecordingStatusCallback       string
	RecordingStatusCallbackMethod string
}

var ErrUrlXorApplicationSid = errors.New(`Exactly one of "Url" or "ApplicationSid" is required"`)

func (p CallParams) Validates() error {
	if p.Url != "" && p.ApplicationSid != "" {
		return ErrUrlXorApplicationSid
	}
	return nil
}

func (s *CallService) Call(params CallParams) (*Call, *Response, error) {
	if err := params.Validates(); err != nil {
		return nil, nil, err
	}

	u := s.client.EndPoint("Calls")

	paramsStr := structToUrlValues(&params).Encode()
	req, err := s.client.NewRequest(http.MethodPost,
		u.String(), strings.NewReader(paramsStr))
	if err != nil {
		return nil, nil, err
	}

	c := new(Call)
	resp, err := s.client.Do(req, c)
	if err != nil {
		return nil, resp, err
	}

	return c, resp, nil
}

func (s *CallService) Get(sid string) (*Call, *Response, error) {
	u := s.client.EndPoint("Calls", sid)

	req, err := s.client.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, nil, err
	}

	c := new(Call)
	resp, err := s.client.Do(req, c)
	if err != nil {
		return nil, resp, err
	}

	return c, resp, nil
}
