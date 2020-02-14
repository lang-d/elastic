package elastic

import (
	"errors"
	"fmt"
	"log"
)

const (
	DEFAULT_SCROLL      = "5m"
	DEFAULT_SCROLL_SIZE = "1000"
)

type ScrollResp struct {
	hits     chan *Hit
	done     bool
	scrollId string
}

func (this *ScrollResp) Pull() *Hit {
	hit := <-this.hits
	if hit == nil {
		this.done = true
	}
	return hit

}

func (this *ScrollResp) push(hit *Hit) {
	this.hits <- hit
}

func (this *ScrollResp) closeChan() {
	close(this.hits)
}

func (this *ScrollResp) setScrollId(scrollId string) {
	this.scrollId = scrollId
}

func newScollResp() *ScrollResp {
	return &ScrollResp{hits: make(chan *Hit), done: false}
}

type ClearScrollResp struct {
	Error     *Error `json:"error,omitempty"`
	Succeeded *bool  `json:"succeeded,omitempty"`
	NumFreed  *int   `json:"num_freed,omitempty"`
}

// https://www.elastic.co/guide/en/elasticsearch/reference/6.3/search-request-scroll.html
func Scan(client *Client, query *QueryBody, index string, docType string, params map[string]string) (scrollResp *ScrollResp, err error) {
	if params == nil {
		params = map[string]string{}
	}

	if _, ok := params["scroll"]; !ok {
		params["scroll"] = DEFAULT_SCROLL
	}

	if _, ok := params["size"]; !ok {
		params["size"] = DEFAULT_SCROLL_SIZE
	}

	if preserveOrder, ok := params["preserve_order"]; ok {
		if preserveOrder != "true" {
			query.SetSort([]map[string]string{{"_doc": "asc"}})
		}

		delete(params, "preserve_order")
	} else {
		if preserveOrder != "true" {
			query.SetSort([]map[string]string{{"_doc": "asc"}})
		}
	}

	paramsList := make([]string, len(params))

	count := 0
	for k, v := range params {
		paramsList[count] = fmt.Sprintf("%s=%s", k, v)
		count += 1
	}

	resp, err := client.Search(index, docType, query, paramsList...)
	if err != nil {
		return nil, err
	}

	if resp.Error != nil {
		return nil, errors.New(resp.Error.Reason)
	}

	scrollResp = newScollResp()

	if "" == resp.ScrollId || len(resp.Hits.Hits) == 0 {
		scrollResp.closeChan()
		return scrollResp, nil
	}

	if fault := recover(); fault != nil {
		scrollResp.closeChan()
		return scrollResp, fault.(error)
	}

	go func(client *Client, params map[string]string, resp *SearchResult, scrollResp *ScrollResp) {
		scrollParams := map[string]string{}
		scrollParams["scroll"] = params["scroll"]
		for {
			if "" == resp.ScrollId || len(resp.Hits.Hits) == 0 {
				clearScroll(scrollResp.scrollId, client)
				scrollResp.closeChan()
				break
			}
			for _, hit := range resp.Hits.Hits {
				scrollResp.push(hit)
			}
			scrollParams["scroll_id"] = resp.ScrollId
			scrollResp.scrollId = resp.ScrollId

			resp, err = client.Scroll(scrollParams)
			if err != nil {
				panic(err)
			}

			if resp.Error != nil {
				panic(errors.New(resp.Error.Reason))
			}
		}

	}(client, params, resp, scrollResp)

	return scrollResp, nil
}

func clearScroll(scrollId string, client *Client) {
	if "" != scrollId {
		_, err := client.ClearScroll(scrollId)
		if err != nil {
			log.Printf("clear scroll fail,reason:%s\n", err)
		}

	}
}
