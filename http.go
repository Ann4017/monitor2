package main

import (
	"net/http"
	"time"
)

type C_http struct {
	s_url       string
	s_status    string
	s_time      string
	pc_response *http.Response
}

func (c *C_http) Get_http_status(_s_url string) error {
	client := &http.Client{
		Timeout: time.Duration(time.Second * 15),
	}

	resp, err := client.Get(_s_url)
	if err != nil {
		return err
	}

	resp_time, err := http.ParseTime(resp.Header.Get("Date"))
	if err != nil {
		return err
	}

	c.pc_response = resp
	c.s_url = _s_url
	c.s_status = resp.Status
	c.s_time = resp_time.In(time.FixedZone("KST", 9*60*60)).Format(time.RFC3339)

	return nil
}

func (c *C_http) Close_resp_body() error {
	if c.pc_response != nil {
		return c.pc_response.Body.Close()
	}

	return nil
}
