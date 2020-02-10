package elastic

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	client          *http.Client
	url             string
	basicAuthUser   string
	basicAuthPasswd string
	timeOut         time.Duration
}

type ClientOptionFunc func(*Client)

// create a Client
func NewClient(options ...ClientOptionFunc) (*Client, error) {
	this := &Client{}
	for _, op := range options {
		op(this)
	}

	this.client = &http.Client{Timeout: this.timeOut}

	return this, nil
}

func SetBasicAuth(user string, passwd string) ClientOptionFunc {
	return func(this *Client) {
		this.basicAuthUser = user
		this.basicAuthPasswd = passwd
	}
}

func SetUrl(url string) ClientOptionFunc {
	return func(this *Client) {
		this.url = url
	}
}

func SetTimeOut(timeOut int) ClientOptionFunc {
	return func(this *Client) {
		this.timeOut = time.Duration(timeOut) * time.Second
	}
}

func (this *Client) buildUrl(index string, docType string, params ...string) string {
	if len(params) > 0 {
		return fmt.Sprintf("%s/%s/%s/_search?%s", this.url, index, docType, strings.Join(params, "&"))
	} else {
		return fmt.Sprintf("%s/%s/%s/_search", this.url, index, docType)
	}

}
func (this *Client) buildRequest(url string, query Query) (*http.Request, error) {
	b, err := query.BuildBody()
	if err != nil {
		return nil, err
	}
	body, err := json.Marshal(b)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(this.basicAuthUser, this.basicAuthPasswd)
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	return req, nil
}

func (this *Client) Search(index string, docType string, query Query, params ...string) (*SearchResult, error) {
	url := this.buildUrl(index, docType, params...)
	req, err := this.buildRequest(url, query)
	if err != nil {
		return nil, err
	}
	resp, err := this.client.Do(req)
	if err != nil {
		return nil, err
	}
	result := new(SearchResult)
	return this.buildResult(resp, result)
}

func (this *Client) buildResult(response *http.Response, result *SearchResult) (*SearchResult, error) {
	defer response.Body.Close()

	err := json.NewDecoder(response.Body).Decode(result)

	if err != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		reason := fmt.Sprintf("search fail,reason:%s", result.Error.Reason)
		return result, errors.New(reason)
	}
	return result, nil
}

// check if cant connect to elastic
// if can't return false
func (this *Client) Ping() (bool, error) {
	req, err := http.NewRequest("GET", this.url, nil)
	if err != nil {
		return false, err
	}
	req.SetBasicAuth(this.basicAuthUser, this.basicAuthPasswd)

	resp, err := this.client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return false, err
	}
	if resp.StatusCode != 200 {
		errorStr := fmt.Sprintf("can't connect to elastic,bad code:%d", resp.StatusCode)
		return false, errors.New(errorStr)
	}
	return true, nil
}

// todo need dubug and fix
func (this *Client) Bulk(actions []Action) (*BulkResult, error) {
	body := []byte{}
	for _, a := range actions {
		data, err := a.Format()
		if err != nil {
			return nil, err
		}
		body = append(body, []byte("\n")...)
		body = append(body, data...)
	}
	body = append(body, []byte("\n")...)

	req, err := http.NewRequest("POST", this.url+"/_bulk", bufio.NewReader(bytes.NewReader(body)))
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(this.basicAuthUser, this.basicAuthPasswd)
	req.Header.Set("Content-Type", "application/x-ndjson")

	response, err := this.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	bulkResult := new(BulkResult)

	err = json.NewDecoder(response.Body).Decode(bulkResult)

	if err != nil {
		return nil, err
	}

	return bulkResult, nil

}
