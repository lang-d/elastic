package elastic

import (
	"encoding/json"
	"log"
	"strings"
)

// use to packet elastic search result
type SearchResult struct {
	Took         int                  `json:"took,omitempty"`
	TimeOut      bool                 `json:"time_out,omitempty"`
	Shards       *Shards              `json:"_shards,omitempty"`
	Clusters     *Clusters            `json:"_clusters,omitempty"`
	Hits         *Hits                `json:"hits,omitempty"`
	Aggregations Aggregations         `json:"aggregations,omitempty"`
	Error        *Error               `json:"error,omitempty"`
	Status       int                  `json:"status,omitempty"`
	Suggest      map[string][]Suggest `json:"suggest,omitempty"`
	ScrollId     string               `json:"_scroll_id,omitempty"`
}

type Suggest struct {
	Text    string    `json:"text,omitempty"`
	Offset  int       `json:"offset,omitempty"`
	Length  int       `json:"length,omitempty"`
	Options []*Option `json:"options,omitempty"`
}

type Option struct {
	Text   string  `json:"text,omitempty"`
	Type   string  `json:"_type,omitempty"`
	Index  string  `json:"_index,omitempty"`
	Id     string  `json:"_id,omitempty"`
	Score  float64 `json:"_score,omitempty"`
	Source Source  `json:"_source,omitempty"`
}

type Shards struct {
	Total      int `json:"total,omitempty"`
	Successful int `json:"successful,omitempty"`
	Skipped    int `json:"skipped,omitempty"`
	Failed     int `json:"failed,omitempty"`
}
type Clusters struct {
	Total      int `json:"total,omitempty"`
	Successful int `json:"successful,omitempty"`
	Skipped    int `json:"skipped,omitempty"`
}

type Hits struct {
	Total    int     `json:"total"`
	MaxScore float64 `json:"max_score"`
	Hits     []*Hit  `json:"hits"`
}
type Hit struct {
	Index     string    `json:"_index,omitempty"`
	Type      string    `json:"_type,omitempty"`
	Id        string    `json:"_id,omitempty"`
	Score     float64   `json:"_score"`
	Source    Source    `json:"_source,omitempty"`
	Highlight Highlight `json:"highlight,omitempty"`
}

type Source map[string]interface{}

type Highlight map[string][]string

type Aggregations map[string]*json.RawMessage

// from a point get aggregations data
// if the aggs not just one ,the path must like this "emoji.top_sticker.buckets","emoji.sum_sticker.buckets"
func (this *Aggregations) GetData(paths ...string) (map[string]interface{}, error) {
	aggses := make(map[string]interface{})
	for _, path := range paths {
		aggs := *this
		a := make(Aggregations)
		ps := strings.Split(path, ".")
		for _, p := range ps {
			if v, ok := aggs[p]; ok {
				if v != nil {
					if err := json.Unmarshal(*v, &a); err != nil {
						aggses[ps[0]] = *aggs[p]
					}
					aggs = a
				} else {
					aggses[ps[0]] = nil
				}

			} else {
				log.Printf("[warn] path:%s not in aggregations", p)
				aggses[ps[0]] = nil
				break
			}

		}

	}
	return aggses, nil
}

type RootCause []map[string]interface{}

type Error struct {
	RootCause RootCause `json:"root_cause,omitempty"`
	Type      string    `json:"type,omitempty"`
	Reason    string    `json:"reason,omitempty"`
	Line      int       `json:"line,omitempty"`
	Col       int       `json:"col,omitempty"`
	CausedBy  CausedBy  `json:"caused_by,omitempty"`
}

type CausedBy struct {
	Type   string `json:"type,omitempty"`
	Reason string `json:"reason,omitempty"`
}

type BulkResult struct {
	Took   int         `json:"took,omitempty"`
	Errors bool        `json:"errors,omitempty"`
	Items  []*BulkItem `json:"items,omitempty"`
	Error  *Error      `json:"error,omitempty"`
}

type BulkItem struct {
	Index  *SubBulkItem `json:"index,omitempty"`
	Delete *SubBulkItem `json:"delete,omitempty"`
	Create *SubBulkItem `json:"create,omitempty"`
	Update *SubBulkItem `json:"update,omitempty"`
}

type SubBulkItem struct {
	Index       string `json:"_index,omitempty"`
	DocType     string `json:"_type,omitempty"`
	Id          string `json:"_id,omitempty"`
	Version     int64  `json:"_version,omitempty"`
	Result      string `json:"result,omitempty"`
	Shards      Shards `json:"_shards,omitempty"`
	Status      int    `json:"status,omitempty"`
	SeqNo       int    `json:"_seq_no,omitempty"`
	PrimaryTerm int    `json:"_primary_term,omitempty"`
	Error       Error  `json:"error,omitempty"`
}
