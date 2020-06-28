package elastic

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"
	"time"
)

func init() {
	iniPool()
}

func TestTermsScript(t *testing.T) {
	//从连接池中取得一个连接
	v, err := MyPool.GetClient()
	if err != nil {
		fmt.Println(err.Error())
	}
	defer MyPool.PutClient(v)
	if v == nil {
		log.Fatalf("\n\n\n\n*******\n\n\n")
	}

	client := v.(*Client)
	query := NewQueryBody()
	boolQuery := NewBoolQuery()

	userId := "1591411397"
	boolQuery.Filter(NewTermQuery("userId", userId))

	timeType := ""

	var lastTime string
	switch timeType {
	case "yesterday":
		lastTime = time.Now().AddDate(0, 0, -1).Local().Format("2006-01-02 00:00:00")
	case "seven":
		lastTime = time.Now().AddDate(0, 0, -7).Local().Format("2006-01-02 00:00:00")
	case "thirty":
		lastTime = time.Now().AddDate(0, 0, -30).Local().Format("2006-01-02 00:00:00")
	}
	if "" != lastTime {
		boolQuery.Filter(NewRangeQuery("time").Gte(lastTime))
	}
	today := time.Now().Local().Format("2006-01-02 00:00:00")
	boolQuery.Filter(NewRangeQuery("time").Lt(today))
	durationsAggs := NewRangeAggs("durations", "duration").Ranges([]map[string]interface{}{
		{
			"to":  15000,
			"key": "<15s",
		},
		{
			"from": 15000,
			"to":   30000,
			"key":  "15-30s",
		},
		{
			"from": 30000,
			"to":   60000,
			"key":  "30-60s",
		},
		{
			"from": 60000,
			"key":  ">60s",
		},
	})
	dateAggs := NewTermsAggs("date", "").Script(NewScript().Source("doc['time'].value.hourOfDay<10?'0'+doc['time'].value.hourOfDay:doc['time'].value.hourOfDay").
		Lang("painless"))

	query.Aggs(NewAvgAggs("avgView", "viewCount"), NewAvgAggs("avgLike", "likeCount"), NewAvgAggs("avgComment", "commentCount"),
		NewAvgAggs("avgShare", "shareCount"), durationsAggs, dateAggs)

	query.Query(boolQuery).Size(0)

	searchResult, err := client.Search("search_kuaishou_photo", "_doc", query)
	if err != nil {
		log.Fatalf(err.Error())
	}
	result, err := Format(searchResult, "avgView", "avgLike", "avgComment", "avgShare", "durations.buckets", "date.buckets")
	if err != nil {
		log.Fatalf(err.Error())
	}
	resultStr, err := json.Marshal(result)
	if err != nil {
		log.Fatalf(err.Error())
	}
	fmt.Printf("%s\n", resultStr)
}

func Format(searchResult *SearchResult, paths ...string) (map[string]interface{}, error) {
	hits := searchResult.Hits
	total := hits.Total
	took := searchResult.Took
	aggs := searchResult.Aggregations
	result := make(map[string]interface{})
	data := FormatHit(*hits)
	result["datas"] = data
	result["total"] = total
	result["took"] = took
	if aggs != nil {
		aggsData, err := aggs.GetData(paths...)
		if err != nil {
			return nil, err
		}
		result["aggs"] = aggsData
	} else {
		result["aggs"] = nil
	}

	return result, nil
}
