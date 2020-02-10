package elastic

import (
	"errors"
)

type SearchAggregations interface {
	Query
	Name() string
}

// terms aggregations
type TermsAggs struct {
	name  string
	field string
	size  *int
	order []map[string]string
	aggs  []SearchAggregations
}

func NewTermsAggs(name string, field string) *TermsAggs {
	return &TermsAggs{name: name, field: field}
}
func (this *TermsAggs) Aggs(aggs ...SearchAggregations) *TermsAggs {
	this.aggs = append(this.aggs, aggs...)
	return this
}

// return this aggs's name
func (this *TermsAggs) Name() string {
	return this.name
}

func (this *TermsAggs) Field(field string) *TermsAggs {
	this.field = field
	return this
}

func (this *TermsAggs) Size(size int) *TermsAggs {
	this.size = &size
	return this
}
func (this *TermsAggs) Order(field string, order string) *TermsAggs {
	item := make(map[string]string)
	item[field] = order
	this.order = append(this.order, item)
	return this
}

// {"name":{"terms":{"field":"field","size":10,"order":[{"_key":"asc"}]}}}
func (this *TermsAggs) BuildBody() (map[string]interface{}, error) {
	if "" == this.name {
		return nil, errors.New("terms aggs name can't be ''")
	}
	if "" == this.field {
		return nil, errors.New("field can't be null")
	}
	query := make(map[string]interface{})
	terms := make(map[string]interface{})
	subTerms := make(map[string]interface{})
	subTerms["field"] = this.field
	if this.size != nil {
		subTerms["size"] = this.size
	}
	if len(this.order) > 0 {
		subTerms["order"] = this.order
	}
	terms["terms"] = subTerms
	if this.aggs != nil {
		aggses := make(map[string]interface{})
		for _, a := range this.aggs {
			subAggs, err := a.BuildBody()
			if err != nil {
				return nil, err
			}
			aggses[a.Name()] = subAggs[a.Name()]
		}
		terms["aggs"] = aggses
	}
	query[this.name] = terms

	return query, nil
}

// filter aggregations
type FilterAggs struct {
	name   string
	filter Query
	aggs   []SearchAggregations
}

func NewFilterAggs(name string) *FilterAggs {
	return &FilterAggs{name: name}
}
func (this *FilterAggs) Filter(query Query) *FilterAggs {
	this.filter = query
	return this
}
func (this *FilterAggs) Aggs(aggs ...SearchAggregations) *FilterAggs {
	this.aggs = append(this.aggs, aggs...)
	return this
}

// return this aggs's name
func (this *FilterAggs) Name() string {
	return this.name
}

// {"name":{"filter":[],"aggs":{}}}
func (this *FilterAggs) BuildBody() (map[string]interface{}, error) {
	if "" == this.name {
		return nil, errors.New("filter aggs name can't be ''")
	}
	query := make(map[string]interface{})
	filterAggs := make(map[string]interface{})
	if this.filter != nil {
		filter, err := this.filter.BuildBody()
		if err != nil {
			return nil, err
		}
		filterAggs["filter"] = filter
	}
	if this.aggs != nil {
		aggses := make(map[string]interface{})
		for _, a := range this.aggs {
			subAggs, err := a.BuildBody()
			if err != nil {
				return nil, err
			}
			aggses[a.Name()] = subAggs[a.Name()]
		}
		filterAggs["aggs"] = aggses
	}
	query[this.name] = filterAggs
	return query, nil
}

// date_histogram
type DateHistogramAggs struct {
	name     string
	field    string
	interval string
	params   []map[string]interface{}
	aggs     []SearchAggregations
}

func NewDateHistogramAggs(name string, field string, interval string) *DateHistogramAggs {
	return &DateHistogramAggs{name: name, field: field, interval: interval}
}

// return this aggs's name
func (this *DateHistogramAggs) Name() string {
	return this.name
}
func (this *DateHistogramAggs) Field(field string) *DateHistogramAggs {
	this.field = field
	return this
}
func (this *DateHistogramAggs) Interval(interval string) *DateHistogramAggs {
	this.interval = interval
	return this
}

func (this *DateHistogramAggs) Aggs(aggs ...SearchAggregations) *DateHistogramAggs {
	this.aggs = append(this.aggs, aggs...)
	return this
}

// other extra condition add by this func
func (this *DateHistogramAggs) Params(key string, value interface{}) *DateHistogramAggs {
	param := make(map[string]interface{})
	param[key] = value
	this.params = append(this.params, param)
	return this
}

// {"name":{"date_histogram":{}}}
func (this *DateHistogramAggs) BuildBody() (map[string]interface{}, error) {
	if "" == this.name {
		return nil, errors.New("date_histogram aggs name can't be ''")
	}
	if "" == this.field {
		return nil, errors.New("date histogram must give field")
	}
	if "" == this.interval {
		return nil, errors.New("date histogram must give interval")
	}
	query := make(map[string]interface{})
	dateHistogramAggs := make(map[string]interface{})
	subDateHistogramAggs := make(map[string]interface{})
	subDateHistogramAggs["field"] = this.field
	subDateHistogramAggs["interval"] = this.interval
	if len(this.params) > 0 {
		for _, param := range this.params {
			for k, v := range param {
				subDateHistogramAggs[k] = v
			}
		}
	}
	if this.aggs != nil {
		aggses := make(map[string]interface{})
		for _, a := range this.aggs {
			subAggs, err := a.BuildBody()
			if err != nil {
				return nil, err
			}
			aggses[a.Name()] = subAggs[a.Name()]
		}
		dateHistogramAggs["aggs"] = aggses
	}
	dateHistogramAggs["date_histogram"] = subDateHistogramAggs
	query[this.name] = dateHistogramAggs

	return query, nil
}

// histogram
type HistogramAggs struct {
	name     string
	field    string
	interval string
	order    []map[string]string
	params   []map[string]interface{}
	aggs     []SearchAggregations
}

func NewHistogramAggs(name string, field string, interval string) *HistogramAggs {
	return &HistogramAggs{name: name, field: field, interval: interval}
}

// return this aggs's name
func (this *HistogramAggs) Name() string {
	return this.name
}
func (this *HistogramAggs) Field(field string) *HistogramAggs {
	this.field = field
	return this
}
func (this *HistogramAggs) Interval(interval string) *HistogramAggs {
	this.interval = interval
	return this
}

func (this *HistogramAggs) Aggs(aggs ...SearchAggregations) *HistogramAggs {
	this.aggs = append(this.aggs, aggs...)
	return this
}

func (this *HistogramAggs) Order(field, order string) *HistogramAggs {
	item := make(map[string]string)
	item[field] = order
	this.order = append(this.order, item)
	return this
}

// other extra condition add by this func
func (this *HistogramAggs) Params(key string, value interface{}) *HistogramAggs {
	param := make(map[string]interface{})
	param[key] = value
	this.params = append(this.params, param)
	return this
}

// {"name":{"histogram":{}}}
func (this *HistogramAggs) BuildBody() (map[string]interface{}, error) {
	if "" == this.name {
		return nil, errors.New("histogram aggs name can't be ''")
	}
	if "" == this.field {
		return nil, errors.New("histogram must give field")
	}
	if "" == this.interval {
		return nil, errors.New("histogram must give interval")
	}
	query := make(map[string]interface{})
	HistogramAggs := make(map[string]interface{})
	subHistogramAggs := make(map[string]interface{})

	subHistogramAggs["field"] = this.field
	subHistogramAggs["interval"] = this.interval

	if len(this.order) > 0 {
		subHistogramAggs["order"] = this.order
	}

	if len(this.params) > 0 {
		for _, param := range this.params {
			for k, v := range param {
				subHistogramAggs[k] = v
			}
		}
	}
	if this.aggs != nil {
		aggses := make(map[string]interface{})
		for _, a := range this.aggs {
			subAggs, err := a.BuildBody()
			if err != nil {
				return nil, err
			}
			aggses[a.Name()] = subAggs[a.Name()]
		}
		HistogramAggs["aggs"] = aggses
	}
	HistogramAggs["histogram"] = subHistogramAggs
	query[this.name] = HistogramAggs

	return query, nil
}

//nested aggregations
type NestedAggs struct {
	field string
	aggs  []SearchAggregations
}

func NewNestedAggs(field string) *NestedAggs {
	return &NestedAggs{field: field}
}

func (this *NestedAggs) Aggs(aggs ...SearchAggregations) *NestedAggs {
	this.aggs = append(this.aggs, aggs...)
	return this
}

// return this aggs's name
func (this *NestedAggs) Name() string {
	return this.field
}

// {"field":{"nested":{"path":"field"},"aggs":{}}}
func (this *NestedAggs) BuildBody() (map[string]interface{}, error) {
	if "" == this.field {
		return nil, errors.New("must set field to nested aggs")
	}
	query := make(map[string]interface{})
	nested := make(map[string]interface{})
	nested["nested"] = map[string]string{"path": this.field}
	if this.aggs != nil {
		aggses := make(map[string]interface{})
		for _, a := range this.aggs {
			subAggs, err := a.BuildBody()
			if err != nil {
				return nil, err
			}
			aggses[a.Name()] = subAggs[a.Name()]
		}
		nested["aggs"] = aggses
	}
	query[this.field] = nested
	return query, nil
}

// max aggregations
type MaxAggs struct {
	name   string
	field  string
	params []map[string]interface{}
	aggs   []SearchAggregations
}

func NewMaxAggs(name string, field string) *MaxAggs {
	return &MaxAggs{name: name, field: field}
}

// return this aggs's name
func (this *MaxAggs) Name() string {
	return this.name
}

func (this *MaxAggs) Field(field string) *MaxAggs {
	this.field = field
	return this
}
func (this *MaxAggs) Aggs(aggs ...SearchAggregations) *MaxAggs {
	this.aggs = append(this.aggs, aggs...)
	return this
}

// other extra condition add by this func
func (this *MaxAggs) Params(key string, value interface{}) *MaxAggs {
	param := make(map[string]interface{})
	param[key] = value
	this.params = append(this.params, param)
	return this
}

// {"name":{"max":{"field":"field"},"aggs":{}}}
func (this *MaxAggs) BuildBody() (map[string]interface{}, error) {
	if "" == this.name {
		return nil, errors.New("max aggs name can't be ''")
	}
	if "" == this.field {
		return nil, errors.New("max aggregations must give field")
	}
	query := make(map[string]interface{})
	maxAggs := make(map[string]interface{})
	subMaxAggs := make(map[string]interface{})
	subMaxAggs["field"] = this.field

	if len(this.params) > 0 {
		for _, param := range this.params {
			for k, v := range param {
				subMaxAggs[k] = v
			}
		}
	}
	maxAggs["max"] = subMaxAggs
	if this.aggs != nil {
		aggses := make(map[string]interface{})
		for _, a := range this.aggs {
			subAggs, err := a.BuildBody()
			if err != nil {
				return nil, err
			}
			aggses[a.Name()] = subAggs[a.Name()]
		}
		maxAggs["aggs"] = aggses
	}

	query[this.name] = maxAggs
	return query, nil
}

// min aggregations
type MinAggs struct {
	name   string
	field  string
	params []map[string]interface{}
	aggs   []SearchAggregations
}

func NewMinAggs(name string, field string) *MinAggs {
	return &MinAggs{name: name, field: field}
}

// return this aggs's name
func (this *MinAggs) Name() string {
	return this.name
}

func (this *MinAggs) Field(field string) *MinAggs {
	this.field = field
	return this
}
func (this *MinAggs) Aggs(aggs ...SearchAggregations) *MinAggs {
	this.aggs = append(this.aggs, aggs...)
	return this
}

// other extra condition add by this func
func (this *MinAggs) Params(key string, value interface{}) *MinAggs {
	param := make(map[string]interface{})
	param[key] = value
	this.params = append(this.params, param)
	return this
}

// {"name":{"min":{"field":"field"},"aggs":{}}}
func (this *MinAggs) BuildBody() (map[string]interface{}, error) {
	if "" == this.name {
		return nil, errors.New("min aggs name can't be ''")
	}
	if "" == this.field {
		return nil, errors.New("min aggregations must give field")
	}
	query := make(map[string]interface{})
	minAggs := make(map[string]interface{})
	subMinAggs := make(map[string]interface{})
	subMinAggs["field"] = this.field

	if len(this.params) > 0 {
		for _, param := range this.params {
			for k, v := range param {
				subMinAggs[k] = v
			}
		}
	}
	minAggs["min"] = subMinAggs
	if this.aggs != nil {
		aggses := make(map[string]interface{})
		for _, a := range this.aggs {
			subAggs, err := a.BuildBody()
			if err != nil {
				return nil, err
			}
			aggses[a.Name()] = subAggs[a.Name()]
		}
		minAggs["aggs"] = aggses
	}

	query[this.name] = minAggs
	return query, nil
}

// sum aggregations
type SumAggs struct {
	name   string
	field  string
	params []map[string]interface{}
	aggs   []SearchAggregations
}

func NewSumAggs(name string, field string) *SumAggs {
	return &SumAggs{name: name, field: field}
}

// return this aggs's name
func (this *SumAggs) Name() string {
	return this.name
}

func (this *SumAggs) Field(field string) *SumAggs {
	this.field = field
	return this
}
func (this *SumAggs) Aggs(aggs ...SearchAggregations) *SumAggs {
	this.aggs = append(this.aggs, aggs...)
	return this
}

// other extra condition add by this func
func (this *SumAggs) Params(key string, value interface{}) *SumAggs {
	param := make(map[string]interface{})
	param[key] = value
	this.params = append(this.params, param)
	return this
}

// {"name":{"sum":{"field":"field"},"aggs":{}}}
func (this *SumAggs) BuildBody() (map[string]interface{}, error) {
	if "" == this.name {
		return nil, errors.New("sum aggs name can't be ''")
	}
	if "" == this.field {
		return nil, errors.New("sum aggregations must give field")
	}
	query := make(map[string]interface{})
	sumAggs := make(map[string]interface{})
	subSumAggs := make(map[string]interface{})
	subSumAggs["field"] = this.field

	if len(this.params) > 0 {
		for _, param := range this.params {
			for k, v := range param {
				subSumAggs[k] = v
			}
		}
	}
	sumAggs["sum"] = subSumAggs
	if this.aggs != nil {
		aggses := make(map[string]interface{})
		for _, a := range this.aggs {
			subAggs, err := a.BuildBody()
			if err != nil {
				return nil, err
			}
			aggses[a.Name()] = subAggs[a.Name()]
		}
		sumAggs["aggs"] = aggses
	}

	query[this.name] = sumAggs
	return query, nil
}

// avg aggregations
type AvgAggs struct {
	name   string
	field  string
	params []map[string]interface{}
	aggs   []SearchAggregations
}

func NewAvgAggs(name string, field string) *AvgAggs {
	return &AvgAggs{name: name, field: field}
}

// return this aggs's name
func (this *AvgAggs) Name() string {
	return this.name
}

func (this *AvgAggs) Field(field string) *AvgAggs {
	this.field = field
	return this
}
func (this *AvgAggs) Aggs(aggs ...SearchAggregations) *AvgAggs {
	this.aggs = append(this.aggs, aggs...)
	return this
}

// other extra condition add by this func
func (this *AvgAggs) Params(key string, value interface{}) *AvgAggs {
	param := make(map[string]interface{})
	param[key] = value
	this.params = append(this.params, param)
	return this
}

// {"name":{"avg":{"field":"field"},"aggs":{}}}
func (this *AvgAggs) BuildBody() (map[string]interface{}, error) {
	if "" == this.name {
		return nil, errors.New("avg aggs name can't be ''")
	}
	if "" == this.field {
		return nil, errors.New("avg aggregations must give field")
	}
	query := make(map[string]interface{})
	avgAggs := make(map[string]interface{})
	subAvgAggs := make(map[string]interface{})
	subAvgAggs["field"] = this.field

	if len(this.params) > 0 {
		for _, param := range this.params {
			for k, v := range param {
				subAvgAggs[k] = v
			}
		}
	}
	avgAggs["avg"] = subAvgAggs
	if this.aggs != nil {
		aggses := make(map[string]interface{})
		for _, a := range this.aggs {
			subAggs, err := a.BuildBody()
			if err != nil {
				return nil, err
			}
			aggses[a.Name()] = subAggs[a.Name()]
		}
		avgAggs["aggs"] = aggses
	}

	query[this.name] = avgAggs
	return query, nil
}

// cardinality aggregations
type CardinalityAggs struct {
	name               string
	field              string
	precisionThreshold *int
	params             []map[string]interface{}
	aggs               []SearchAggregations
}

func NewCardinalityAggs(name string, field string) *CardinalityAggs {
	return &CardinalityAggs{name: name, field: field}
}

// return this aggs's name
func (this *CardinalityAggs) Name() string {
	return this.name
}

func (this *CardinalityAggs) Field(field string) *CardinalityAggs {
	this.field = field
	return this
}
func (this *CardinalityAggs) Aggs(aggs ...SearchAggregations) *CardinalityAggs {
	this.aggs = append(this.aggs, aggs...)
	return this
}

func (this *CardinalityAggs) PrecisionThreshold(precisionThreshold int) *CardinalityAggs {
	this.precisionThreshold = &precisionThreshold
	return this
}

// other extra condition add by this func
func (this *CardinalityAggs) Params(key string, value interface{}) *CardinalityAggs {
	param := make(map[string]interface{})
	param[key] = value
	this.params = append(this.params, param)
	return this
}

// {"name":{"cardinality":{"field":"field","precision_threshold":100},"aggs":{}}}
func (this *CardinalityAggs) BuildBody() (map[string]interface{}, error) {
	if "" == this.name {
		return nil, errors.New("cardinality aggs name can't be ''")
	}
	if "" == this.field {
		return nil, errors.New("cardinality aggregations must give field")
	}
	query := make(map[string]interface{})
	cardinalityAggs := make(map[string]interface{})
	subcardinalityAggs := make(map[string]interface{})
	subcardinalityAggs["field"] = this.field

	if len(this.params) > 0 {
		for _, param := range this.params {
			for k, v := range param {
				subcardinalityAggs[k] = v
			}
		}
	}

	if this.precisionThreshold != nil {
		subcardinalityAggs["precision_threshold"] = this.precisionThreshold
	}

	cardinalityAggs["cardinality"] = subcardinalityAggs
	if this.aggs != nil {
		aggses := make(map[string]interface{})
		for _, a := range this.aggs {
			subAggs, err := a.BuildBody()
			if err != nil {
				return nil, err
			}
			aggses[a.Name()] = subAggs[a.Name()]
		}
		cardinalityAggs["aggs"] = aggses
	}

	query[this.name] = cardinalityAggs
	return query, nil
}

// top_hit aggs
// see https://www.elastic.co/guide/en/elasticsearch/reference/6.3/search-aggregations-metrics-top-hits-aggregation.html
type TopHitsAggs struct {
	name   string
	size   *int
	source map[string][]string
	sort   []map[string]map[string]string
	aggs   []SearchAggregations
}

func NewTopHitsAggs(name string) *TopHitsAggs {
	return &TopHitsAggs{name: name}
}

// return this aggs's name
func (this *TopHitsAggs) Name() string {
	return this.name
}

func (this *TopHitsAggs) Size(size int) *TopHitsAggs {
	this.size = &size
	return this
}

func (this *TopHitsAggs) Sort(field string, order string) *TopHitsAggs {
	var sort = map[string]map[string]string{field: {"order": order}}
	this.sort = append(this.sort, sort)
	return this
}

func (this *TopHitsAggs) Aggs(aggs ...SearchAggregations) *TopHitsAggs {
	this.aggs = append(this.aggs, aggs...)
	return this
}

// set return fields,key must includes or excludes
// if key not in includes or excludes,it will be set includes
func (this *TopHitsAggs) Source(key string, fields ...string) *TopHitsAggs {
	if key != "includes" && key != "excludes" {
		key = "includes"
	}
	for _, field := range fields {
		if this.source == nil {
			this.source = make(map[string][]string)
		}
		this.source[key] = append(this.source[key], field)
	}
	return this
}

// "top_sales_hits": {
//                    "top_hits": {
//                        "sort": [
//                            {
//                                "date": {
//                                    "order": "desc"
//                                }
//                            }
//                        ],
//                        "_source": {
//                            "includes": [ "date", "price" ]
//                        },
//                        "size" : 1
//                    }
//                }

func (this *TopHitsAggs) BuildBody() (map[string]interface{}, error) {
	if "" == this.name {
		return nil, errors.New("top_hit aggs name can't be ''")
	}
	aggs := make(map[string]interface{})
	topHitsAggs := make(map[string]interface{})
	subTopHitsAggs := make(map[string]interface{})
	if this.sort != nil {
		subTopHitsAggs["sort"] = this.sort
	}
	if this.source != nil {
		subTopHitsAggs["_source"] = this.source
	}
	if this.size != nil {
		subTopHitsAggs["size"] = this.size
	}
	topHitsAggs["top_hits"] = subTopHitsAggs
	if this.aggs != nil {
		aggses := make(map[string]interface{})
		for _, a := range this.aggs {
			subAggs, err := a.BuildBody()
			if err != nil {
				return nil, err
			}
			aggses[a.Name()] = subAggs[a.Name()]
		}
		topHitsAggs["aggs"] = aggses
	}
	aggs[this.name] = topHitsAggs
	return aggs, nil
}
