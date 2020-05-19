package elastic

import (
	"fmt"
	"github.com/lang-d/elastic/pool"
	"log"
)

// var Pool *ClientPool
var Pool pool.Pool

type EsConfig struct {
	Host     string `json:"host,omitempty"`
	Port     int    `json:"port,omitempty"`
	User     string `json:"user,omitempty"`
	Passwd   string `json:"passwd,omitempty"`
	Count    int    `json:"count,omitempty"`
	MaxCount int    `json:"max_count,omitempty"`
	MinCount int    `json:"min_count,omitempty"`
	TimeOut  int    `json:"timeout,omitempty"`
}

// init Client pool for whole project use
// count: evey time create client count
// maxCount: the pool's maximum client count
// minCount: the pool's minimum client count
func InitClientPool(esConfig EsConfig) pool.Pool {
	//factory 创建连接的方法
	factory := func() (interface{}, error) {
		client, err := NewClient(
			SetUrl(fmt.Sprintf("%s:%d", esConfig.Host, esConfig.Port)),
			SetBasicAuth(esConfig.User, esConfig.Passwd),
			SetTimeOut(esConfig.TimeOut),
		)
		return client, err
	}

	//close 关闭连接的方法
	close := func(v interface{}) error {
		return v.(*Client).Close()
	}
	// 连接池中拥有的最小连接数InitialCap
	// 最大并发存活连接数 MaxCap
	// 最大空闲连接MaxIdle
	// 创建一个连接池：
	poolConfig := &pool.Config{
		InitialCap: esConfig.Count,
		MaxIdle:    esConfig.MinCount,
		MaxCap:     esConfig.MaxCount,
		Factory:    factory,
		Close:      close,
		//连接最大空闲时间，超过该时间的连接 将会关闭，可避免空闲时连接EOF，自动失效的问题，elasticPool中设置为-1
		IdleTimeout: -1,
	}
	p, err := pool.NewChannelPool(poolConfig)
	if err != nil {
		log.Fatalf(err.Error())
	}

	Pool = p
	return Pool
}
