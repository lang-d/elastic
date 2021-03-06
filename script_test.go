package elastic

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
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

func TestFloat64Return(t *testing.T) {
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
	boolQuery.Filter(NewTermQuery("userId", "188888880"))

	query.Query(boolQuery).Aggs(NewAvgAggs("avgKsCoins", "totalKsCoin")).Size(0)

	searchResult, err := client.Search("search_kuaishou_live", "_doc", query)
	if err != nil {
		log.Fatalf(err.Error())
	}
	println(query.String())
	fmt.Printf("%s\n", searchResult.Aggregations["avgKsCoins"])
	result, err := Format(searchResult, "avgKsCoins.value")
	if err != nil {
		log.Fatalf(err.Error())
	}

	fmt.Printf("%s\n", result)
}

func TestSortInterface(t *testing.T) {
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

	boolQuery.Filter(NewTermQuery("commodity.uniqueItemId", "zCGuBjHKlcQ"))
	nestedQuery := NewNestedQuery("commodity").Query(boolQuery)

	query.Query(nestedQuery).SortInterface(map[string]interface{}{
		"commodity.settlementPrice": map[string]interface{}{
			"order": "desc",
			"nested": map[string]interface{}{
				"path": "commodity",
			},
		},
	}).Source("includes", "liveStreamId")

	searchResult, err := client.Search("search_kuaishou_live", "_doc", query)
	if err != nil {
		log.Fatalf(err.Error())
	}
	println(query.String())
	result, err := Format(searchResult)
	if err != nil {
		log.Fatalf(err.Error())
	}

	fmt.Printf("%s\n", result)
}

func TestSort(t *testing.T) {
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

	query.Sort(map[string]string{"anaTime": "desc"}).Source("includes", "liveStreamId")

	searchResult, err := client.Search("search_kuaishou_live", "_doc", query)
	if err != nil {
		log.Fatalf(err.Error())
	}
	println(query.String())
	result, err := Format(searchResult)
	if err != nil {
		log.Fatalf(err.Error())
	}

	fmt.Printf("%s\n", result)
}

func TestRangeAggs(t *testing.T) {
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

	boolQuery.Filter(NewTermQuery("room_id", "6896315940368354055"))
	query.Query(boolQuery).Size(0).Aggs(
		NewTermsAggs("promotion_type_v1", "promotion_type_common"),
		NewTermsAggs("brand_name", "brand_common"),
		NewTermsAggs("goods_source", "goods_source"),
		NewRangeAggs("ana_price", "ana_price").Ranges([]map[string]interface{}{
			{
				"key": "<50",
				"to":  50,
			},
			{
				"key":  "50-100",
				"from": 50,
				"to":   100,
			},
			{
				"key":  "100-300",
				"from": 100,
				"to":   300,
			},
			{
				"key":  "300-500",
				"from": 300,
				"to":   500,
			},
			{
				"key":  "500-1000",
				"from": 500,
				"to":   1000,
			},
			{
				"key":  ">1000",
				"from": 1000,
			},
		}),
	)

	searchResult, err := client.Search("search_douyin_webcast_promotion", "_doc", query)
	if err != nil {
		log.Fatalf(err.Error())
	}
	println(query.String())
	result, err := Format(searchResult, "promotion_type_v1.buckets", "brand_name.buckets", "goods_source.buckets", "ana_price.buckets")
	if err != nil {
		log.Fatalf(err.Error())
	}

	fmt.Printf("%s\n", result)
}

func TestFunctionScoreQuery(t *testing.T) {
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

	boolQuery.Must(NewTermQuery("is_claim", 1))
	boolQuery.Must(NewTermQuery("claim_success", 1))
	boolQuery.Must(NewTermQuery("is_weixin_friend", 1))

	functionScoreQuery := NewFunctionScoreQuery()
	rand.Seed(time.Now().Unix())
	functionScoreQuery.Query(boolQuery).RandomScore(NewRandomScoreQuery().Field("_seq_no").Seed(float64(rand.Intn(10000))))

	query.Source("includes", "uid", "nickname")

	query.Query(functionScoreQuery)
	println(query.String())
	searchResult, err := client.Search("search_douyin_user", "_doc", query)
	if err != nil {
		log.Fatalf(err.Error())
	}
	result, err := Format(searchResult)
	if err != nil {
		log.Fatalf(err.Error())
	}

	fmt.Printf("%s\n", result)
}

func TestExistIndex(t *testing.T) {
	v, err := MyPool.GetClient()
	if err != nil {
		fmt.Println(err.Error())
	}
	defer MyPool.PutClient(v)
	if v == nil {
		log.Fatalf("\n\n\n\n*******\n\n\n")
	}

	client := v.(*Client)
	ok, err := client.ExistIndex("search_test")
	fmt.Printf("%v, %v\n", ok, err)
	if err != nil {
		log.Fatalf("%v", err)
		return
	}
	if !ok {
		file, err := ioutil.ReadFile("test.json")
		if err != nil {
			log.Fatalf("%v", err)
			return
		}
		err = client.CreateIndex("search_test_20210113", file)
		if err != nil {
			log.Fatalf("%v", err)
			return
		}
		err = client.PutAlias("search_test_20210113", "search_test1", "search_test2")
		if err != nil {
			log.Fatalf("%v", err)
			return
		}
	}
}
