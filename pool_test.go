package elastic

import (
	"encoding/json"
	"fmt"
	"github.com/lang-d/elastic/pool"
	"io/ioutil"
	"log"
	"sync"
	"testing"
)

func FormatHit(hits Hits) []map[string]interface{} {
	data := make([]map[string]interface{}, len(hits.Hits))
	for i, hit := range hits.Hits {
		source := hit.Source
		data[i] = source
	}
	return data
}

var MyPool pool.Pool

func iniPool() {
	// 	config.json={
	// 	"host": "",
	// 	"port": ,
	// 	"user": "",
	// 	"passwd": "",
	// 	"count": ,
	// 	"max_count": ,
	// 	"min_count": ,
	// 	"timeout":
	//   }
	data, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Fatalf("%s", err)
	}
	var config EsConfig
	err = json.Unmarshal(data, &config)
	if err != nil {
		log.Fatalf("%s", err)
	}
	MyPool = InitClientPool(config)
}

func init() {
	iniPool()
}

func TestMyPool(t *testing.T) {
	var wg sync.WaitGroup

	for i := 0; i < 10000; i++ {
		wg.Add(1)
		go func() {
			// defer wg.Done()
			defer wg.Add(-1)
			get_mypoolsearch()
		}()
	}

	wg.Wait()
}

func BenchmarkMyPool(b *testing.B) {
	b.SetParallelism(10)
	b.RunParallel(func(pb *testing.PB) {
		// 每个 goroutine 有属于自己的 bytes.Buffer.
		for pb.Next() {
			// 所有 goroutine 一起，循环一共执行 b.N 次
			get_mypoolsearch()
		}
	})
}

func get_mypoolsearch() {

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
	body := NewBoolQuery()
	body.Filter(NewTermQuery("aweme_id", "6787574817835994368"))
	query.Query(body).Source("includes", "uid")

	searchResult, err := client.Search("search_douyin_aweme", "_doc", query)
	if err != nil {
		log.Fatalf(err.Error())
	}
	hits := searchResult.Hits

	fmt.Printf("%s\n", FormatHit(*hits))
}
