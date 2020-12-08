package elastic

import (
	"encoding/json"
	"errors"
)

// basic query interface,all kind of query need inplement it
type Query interface {
	BuildBody() (map[string]interface{}, error)
}

//	the query json's last part
type QueryBody struct {
	query       Query
	aggs        []SearchAggregations
	from        int
	size        *int
	sort        []map[string]interface{}
	source      map[string][]string
	highlight   Query
	suggest     *GlobSuggest
	searchAfter []interface{}
}

func NewQueryBody() *QueryBody {
	return &QueryBody{}
}
func (this *QueryBody) Query(query Query) *QueryBody {
	this.query = query
	return this
}

func (this *QueryBody) From(from int) *QueryBody {
	this.from = from
	return this
}

func (this *QueryBody) Size(size int) *QueryBody {
	this.size = &size
	return this
}
func (this *QueryBody) Sort(sort map[string]string) *QueryBody {
	tempSort := make(map[string]interface{})
	for k, v := range sort {
		tempSort[k] = v
	}
	this.sort = append(this.sort, tempSort)
	return this
}

func (this *QueryBody) SortInterface(sort map[string]interface{}) *QueryBody {
	this.sort = append(this.sort, sort)
	return this
}

func (this *QueryBody) Suggest(suggest *GlobSuggest) *QueryBody {
	this.suggest = suggest
	return this
}

func (this *QueryBody) Aggs(aggs ...SearchAggregations) *QueryBody {
	this.aggs = append(this.aggs, aggs...)
	return this
}

func (this *QueryBody) Highlight(query Query) *QueryBody {
	this.highlight = query
	return this
}

func (this *QueryBody) SearchAfter(searchAfter []interface{}) *QueryBody {
	this.searchAfter = searchAfter
	return this
}

// set return fields,key must includes or excludes
// if key not in includes or excludes,it will be set includes
func (this *QueryBody) Source(key string, fields ...string) *QueryBody {
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

// sometime maybe need reset sort or set not just one rule
// it is for this situation
func (this *QueryBody) SetSort(sort []map[string]string) *QueryBody {
	tempSorts := make([]map[string]interface{}, 0)
	for _, s := range sort {
		t := make(map[string]interface{})
		for k, v := range s {
			t[k] = v
		}
		tempSorts = append(tempSorts, t)
	}

	this.sort = tempSorts
	return this
}

// format dsl map to json string
func (this *QueryBody) String() string {
	body, err := this.BuildBody()
	if err != nil {
		panic(err)
	}
	jsonStr, err := json.MarshalIndent(&body, "", "\t")
	if err != nil {
		panic(err)
	}
	return string(jsonStr)
}

func (this *QueryBody) SimpleString() string {
	body, err := this.BuildBody()
	if err != nil {
		panic(err)
	}
	jsonStr, err := json.Marshal(&body)
	if err != nil {
		panic(err)
	}
	return string(jsonStr)
}

func (this *QueryBody) BuildBody() (map[string]interface{}, error) {
	var queryBody = make(map[string]interface{})
	if this.query != nil {
		query, err := this.query.BuildBody()
		if err != nil {
			return nil, err
		}
		queryBody["query"] = query
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
		queryBody["aggs"] = aggses
	}

	queryBody["from"] = this.from
	if this.size != nil {
		queryBody["size"] = this.size
	} else {
		// set default size
		queryBody["size"] = 10
	}
	if len(this.sort) != 0 {
		queryBody["sort"] = this.sort
	}

	if this.source != nil {
		queryBody["_source"] = this.source
	}

	if this.highlight != nil {
		highlight, err := this.highlight.BuildBody()
		if err != nil {
			return nil, err
		}
		queryBody["highlight"] = highlight
	}

	if this.suggest != nil {
		suggest, err := this.suggest.BuildBody()
		if err != nil {
			return nil, err
		}
		queryBody["suggest"] = suggest
	}

	if this.searchAfter != nil {
		queryBody["search_after"] = this.searchAfter
	}

	return queryBody, nil
}

// it's bool query
type BoolQuery struct {
	filter             []Query
	must               []Query
	mustNot            []Query
	should             []Query
	minimumShouldMatch *int
	boost              *float64
}

func NewBoolQuery() *BoolQuery {
	return &BoolQuery{
		filter:  make([]Query, 0),
		must:    make([]Query, 0),
		mustNot: make([]Query, 0),
		should:  make([]Query, 0),
	}
}

func (this *BoolQuery) Filter(query Query) *BoolQuery {
	this.filter = append(this.filter, query)
	return this
}

func (this *BoolQuery) Must(query Query) *BoolQuery {
	this.must = append(this.must, query)
	return this
}

func (this *BoolQuery) MustNot(query Query) *BoolQuery {
	this.mustNot = append(this.mustNot, query)
	return this
}

func (this *BoolQuery) Should(query Query) *BoolQuery {
	this.should = append(this.should, query)
	return this
}

func (this *BoolQuery) MinimumShouldMatch(minimumShouldMatch int) *BoolQuery {
	this.minimumShouldMatch = &minimumShouldMatch
	return this
}

func (this *BoolQuery) Boost(boost float64) *BoolQuery {
	this.boost = &boost
	return this
}

func (this *BoolQuery) BuildBody() (map[string]interface{}, error) {
	boolQuery := make(map[string]interface{})
	query := make(map[string]interface{})
	// filter
	if len(this.filter) > 0 {
		var conditions []map[string]interface{}
		for _, item := range this.filter {
			q, err := item.BuildBody()
			if err != nil {
				return nil, err
			}
			conditions = append(conditions, q)
		}
		boolQuery["filter"] = conditions
	}
	// must
	if len(this.must) > 0 {
		var conditions []map[string]interface{}
		for _, item := range this.must {
			q, err := item.BuildBody()
			if err != nil {
				return nil, err
			}
			conditions = append(conditions, q)
		}
		boolQuery["must"] = conditions
	}
	// must_not
	if len(this.mustNot) > 0 {
		var conditions []map[string]interface{}
		for _, item := range this.mustNot {
			q, err := item.BuildBody()
			if err != nil {
				return nil, err
			}
			conditions = append(conditions, q)
		}
		boolQuery["must_not"] = conditions
	}
	// should
	if len(this.should) > 0 {
		var conditions []map[string]interface{}
		for _, item := range this.should {
			q, err := item.BuildBody()
			if err != nil {
				return nil, err
			}
			conditions = append(conditions, q)
		}
		boolQuery["should"] = conditions
		if this.minimumShouldMatch != nil {
			boolQuery["minimum_should_match"] = this.minimumShouldMatch
		} else {
			boolQuery["minimum_should_match"] = 1
		}
	}

	if this.boost != nil {
		boolQuery["boost"] = this.boost
	}
	query["bool"] = boolQuery

	return query, nil
}

// term query
type TermQuery struct {
	field string
	value interface{}
	boost *float64
}

func NewTermQuery(field string, value interface{}) *TermQuery {
	return &TermQuery{field: field, value: value}
}

func (this *TermQuery) Boost(boost float64) *TermQuery {
	this.boost = &boost
	return this
}

//{"term": {field: {"value": 1,"boost":1}}}
func (this *TermQuery) BuildBody() (map[string]interface{}, error) {
	if this.field == "" || this.value == nil {
		return nil, errors.New("a term query must have field and value")
	}
	query := make(map[string]interface{})
	termQuery := make(map[string]interface{})
	subQuery := make(map[string]interface{})

	subQuery["value"] = this.value
	if this.boost != nil {
		subQuery["boost"] = this.boost
	}

	termQuery[this.field] = subQuery
	query["term"] = termQuery

	return query, nil
}

// terms query
type TermsQuery struct {
	field string
	value interface{}
	boost *float64
}

func NewTermsQuery(field string, value interface{}) *TermsQuery {
	return &TermsQuery{field: field, value: value}
}

func (this *TermsQuery) Boost(boost float64) *TermsQuery {
	this.boost = &boost
	return this
}

//{"terms": {field: {"value": 1,"boost":1}}}
func (this *TermsQuery) BuildBody() (map[string]interface{}, error) {
	if this.field == "" || this.value == nil {
		return nil, errors.New("a terms query must have field and value")
	}
	query := make(map[string]interface{})
	termsQuery := make(map[string]interface{})
	termsQuery[this.field] = this.value
	if this.boost != nil {
		termsQuery["boost"] = this.boost
	}
	query["terms"] = termsQuery

	return query, nil
}

type RangeQuery struct {
	field string
	gt    interface{}
	lt    interface{}
	gte   interface{}
	lte   interface{}
	boost *float64
}

func NewRangeQuery(field string) *RangeQuery {
	return &RangeQuery{field: field}
}

// set range,gte or gt,lt or lte
// gte and gt,lt and lte, can't exists at the same time
func (this *RangeQuery) Gt(gt interface{}) *RangeQuery {
	if gt != nil && gt != "" {
		this.gt = gt
		this.gte = nil
	}

	return this
}

func (this *RangeQuery) Gte(gte interface{}) *RangeQuery {
	if gte != nil && gte != "" {
		this.gte = gte
		this.gt = nil
	}
	return this
}

func (this *RangeQuery) Lt(lt interface{}) *RangeQuery {
	if lt != nil && lt != "" {
		this.lt = lt
		this.lte = nil
	}

	return this
}

func (this *RangeQuery) Lte(lte interface{}) *RangeQuery {
	if lte != nil && lte != "" {
		this.lte = lte
		this.lt = nil
	}
	return this
}

// set boost
func (this *RangeQuery) Boost(boost float64) *RangeQuery {
	this.boost = &boost
	return this
}

func (this *RangeQuery) BuildBody() (map[string]interface{}, error) {
	if this.field == "" {
		return nil, errors.New("must set Rangequery's field")
	}
	if this.lt == nil && this.lte == nil && this.gte == nil && this.gt == nil {
		return nil, errors.New("must set Rangequery's range")
	}

	rangeMap := make(map[string]interface{})
	rangeField := make(map[string]interface{})
	rangeItem := make(map[string]interface{})

	if this.gte != nil {
		rangeItem["gte"] = this.gte
	}
	if this.gt != nil {
		rangeItem["gt"] = this.gt
	}
	if this.lte != nil {
		rangeItem["lte"] = this.lte
	}
	if this.lt != nil {
		rangeItem["lt"] = this.lt
	}

	if this.boost != nil {
		rangeItem["boost"] = this.boost
	}

	rangeField[this.field] = rangeItem
	rangeMap["range"] = rangeField

	return rangeMap, nil
}

// nested query
type NestedQuery struct {
	path  string
	query Query
	boost *float64
}

func NewNestedQuery(path string) *NestedQuery {
	return &NestedQuery{path: path}
}
func (this *NestedQuery) Query(query Query) *NestedQuery {
	this.query = query
	return this
}
func (this *NestedQuery) Boost(boost float64) *NestedQuery {
	this.boost = &boost
	return this
}
func (this *NestedQuery) BuildBody() (map[string]interface{}, error) {
	nestedQuery := make(map[string]interface{})
	query := make(map[string]interface{})
	if this.path == "" {
		return nil, errors.New("must set NestedQuery's path")
	}
	if this.query == nil {
		return nil, errors.New("must set NestedQuery's path")
	}
	nestedQuery["path"] = this.path
	nestedQuery["query"], _ = this.query.BuildBody()

	if this.boost != nil {
		nestedQuery["boost"] = this.boost
	}

	query["nested"] = nestedQuery

	return query, nil

}

// exists query
type ExistsQuery struct {
	field string
	boost *float64
}

func NewExistsQuery(field string) *ExistsQuery {
	return &ExistsQuery{field: field}
}
func (this *ExistsQuery) Boost(boost float64) *ExistsQuery {
	this.boost = &boost
	return this
}

func (this *ExistsQuery) BuildBody() (map[string]interface{}, error) {
	existsQuery := make(map[string]interface{})
	query := make(map[string]interface{})
	if this.field == "" {
		return nil, errors.New("must set ExistsQuery's field")
	}
	existsQuery["field"] = this.field

	if this.boost != nil {
		existsQuery["boost"] = this.boost
	}
	query["exists"] = existsQuery

	return query, nil
}

// match query
type MatchQuery struct {
	field    string
	keyword  string
	boost    *float64
	analyzer string
	params   []map[string]interface{}
}

func NewMatchQuery(field string, keyword string) *MatchQuery {
	return &MatchQuery{field: field, keyword: keyword}
}
func (this *MatchQuery) Analyzer(analyzer string) *MatchQuery {
	this.analyzer = analyzer
	return this
}
func (this *MatchQuery) Boost(boost float64) *MatchQuery {
	this.boost = &boost
	return this
}

// this func use to add more search conditions
// like operator,zero_terms_query,cutoff_frequency etc.
func (this *MatchQuery) Params(key string, value interface{}) *MatchQuery {
	param := make(map[string]interface{})
	param[key] = value
	this.params = append(this.params, param)
	return this
}

/**
{"match": {field: {"query": keyword, "analyzer": analyzer, "boost": boost}}}
*/
func (this *MatchQuery) BuildBody() (map[string]interface{}, error) {
	matchQuery := make(map[string]interface{})
	subMatchQuery := make(map[string]interface{})
	query := make(map[string]interface{})
	if this.field == "" {
		return nil, errors.New("must set MatchQuery's field")
	}
	subMatchQuery["query"] = this.keyword
	if this.analyzer != "" {
		subMatchQuery["analyzer"] = this.analyzer
	}
	if this.boost != nil {
		subMatchQuery["boost"] = this.boost
	}
	//add anther conditions
	for _, param := range this.params {
		for k, v := range param {
			subMatchQuery[k] = v
		}
	}

	matchQuery[this.field] = subMatchQuery
	query["match"] = matchQuery

	return query, nil
}

//match_phrase
type MatchPhraseQuery struct {
	field    string
	keyword  string
	boost    *float64
	analyzer string
	slop     *int
}

func NewMatchPhraseQuery(field string, keyword string) *MatchPhraseQuery {
	return &MatchPhraseQuery{field: field, keyword: keyword}
}

func (this *MatchPhraseQuery) Boost(boost float64) *MatchPhraseQuery {
	this.boost = &boost
	return this
}

func (this *MatchPhraseQuery) Analyzer(analyzer string) *MatchPhraseQuery {
	this.analyzer = analyzer
	return this
}

func (this *MatchPhraseQuery) Slop(slop int) *MatchPhraseQuery {
	this.slop = &slop
	return this
}

/**
{"match_phrase": {field: {"query": keyword, "slop": 5, "boost": 2}}}
*/
func (this *MatchPhraseQuery) BuildBody() (map[string]interface{}, error) {
	matchPhraseQuery := make(map[string]interface{})
	subMatchPhraseQuery := make(map[string]interface{})
	query := make(map[string]interface{})
	if this.field == "" {
		return nil, errors.New("must set MatchPhraseQuery's field")
	}
	subMatchPhraseQuery["query"] = this.keyword
	if this.analyzer != "" {
		subMatchPhraseQuery["analyzer"] = this.analyzer
	}

	if this.boost != nil {
		subMatchPhraseQuery["boost"] = this.boost
	}

	if this.slop != nil {
		subMatchPhraseQuery["slop"] = this.slop
	}

	matchPhraseQuery[this.field] = subMatchPhraseQuery
	query["match_phrase"] = matchPhraseQuery

	return query, nil
}

// multi_macth query
type MultiMacthQuery struct {
	query      string
	analyzer   string
	searchType string
	fields     []string
	tieBreaker *float64
	params     []map[string]interface{}
}

func NewMultiMacthQuery(keyword string, fields []string, searchType string) *MultiMacthQuery {
	return &MultiMacthQuery{query: keyword, fields: fields, searchType: searchType}
}

func (this *MultiMacthQuery) Analyzer(analyzer string) *MultiMacthQuery {
	this.analyzer = analyzer
	return this
}

func (this *MultiMacthQuery) TieBreaker(tieBreaker float64) *MultiMacthQuery {
	this.tieBreaker = &tieBreaker
	return this
}

// this func use to add more search conditions
func (this *MultiMacthQuery) Params(key string, value interface{}) *MultiMacthQuery {
	param := make(map[string]interface{})
	param[key] = value
	this.params = append(this.params, param)
	return this
}

/**
	{
      "query":      "Will Smith",
      "type":       "cross_fields",
      "fields":     [ "first_name", "last_name" ],
      "operator":   "and"
    }
*/
func (this *MultiMacthQuery) BuildBody() (map[string]interface{}, error) {
	multiMatchQuery := make(map[string]interface{})
	query := make(map[string]interface{})
	if this.fields == nil {
		return nil, errors.New("a mulit_macth query must give fields")
	}
	multiMatchQuery["query"] = this.query
	multiMatchQuery["fields"] = this.fields

	if this.searchType != "" {
		multiMatchQuery["type"] = this.searchType
	}

	if this.tieBreaker != nil {
		multiMatchQuery["tie_breaker"] = this.tieBreaker
	}

	if this.params != nil {
		for _, param := range this.params {
			for k, v := range param {
				multiMatchQuery[k] = v
			}
		}
	}
	query["multi_match"] = multiMatchQuery

	return query, nil
}

// constant score
// see https://www.elastic.co/guide/en/elasticsearch/reference/6.3/query-dsl-constant-score-query.html
type ConStantScoreQuery struct {
	filter Query
	boost  *float64
}

func NewConstantScoreQuery() *ConStantScoreQuery {
	return &ConStantScoreQuery{}
}

func (this *ConStantScoreQuery) Filter(query Query) *ConStantScoreQuery {
	this.filter = query
	return this
}

func (this *ConStantScoreQuery) Boost(boost float64) *ConStantScoreQuery {
	this.boost = &boost
	return this
}

// {"constant_score":{"filter":{},"boost":1}}
func (this *ConStantScoreQuery) BuildBody() (map[string]interface{}, error) {
	query := make(map[string]interface{})
	constantScoreQuery := make(map[string]interface{})
	if this.filter == nil {
		return nil, errors.New("a constant_score query must set filter!")
	}
	filter, err := this.filter.BuildBody()
	if err != nil {
		return nil, err
	}
	constantScoreQuery["filter"] = filter

	if this.boost != nil {
		constantScoreQuery["boost"] = this.boost
	}
	query["constant_score"] = constantScoreQuery
	return query, nil
}

// function_score
// see https://www.elastic.co/guide/en/elasticsearch/reference/6.3/query-dsl-function-score-query.html
type FunctionScoreQuery struct {
	query            Query
	boost            *float64
	randomScore      *RandomScoreQuery
	boostMode        *string
	functions        []*FunctionQuery
	maxBoost         *float64
	minScore         *float64
	scoreMode        *string
	fieldValueFactor *FieldValueFactorQuery
}

func NewFunctionScoreQuery() *FunctionScoreQuery {
	return &FunctionScoreQuery{}
}

func (this *FunctionScoreQuery) Query(query Query) *FunctionScoreQuery {
	this.query = query
	return this
}

func (this *FunctionScoreQuery) Boost(boost float64) *FunctionScoreQuery {
	this.boost = &boost
	return this
}

func (this *FunctionScoreQuery) RandomScore(randomScore *RandomScoreQuery) *FunctionScoreQuery {
	this.randomScore = randomScore
	return this
}

func (this *FunctionScoreQuery) BoostMode(boostMode string) *FunctionScoreQuery {
	this.boostMode = &boostMode
	return this
}

func (this *FunctionScoreQuery) Functions(function *FunctionQuery) *FunctionScoreQuery {
	this.functions = append(this.functions, function)
	return this
}

func (this *FunctionScoreQuery) MaxBoost(maxBoost float64) *FunctionScoreQuery {
	this.maxBoost = &maxBoost
	return this
}

func (this *FunctionScoreQuery) MinScore(minScore float64) *FunctionScoreQuery {
	this.minScore = &minScore
	return this
}

func (this *FunctionScoreQuery) ScoreMode(scoreMode string) *FunctionScoreQuery {
	this.scoreMode = &scoreMode
	return this
}

func (this *FunctionScoreQuery) FieldValueFactor(fieldValueFactorQuery *FieldValueFactorQuery) *FunctionScoreQuery {
	this.fieldValueFactor = fieldValueFactorQuery
	return this
}

// {
//        "function_score": {
//          "query": { "match_all": {} },
//          "boost": "5",
//          "functions": [
//              {
//                  "filter": { "match": { "test": "bar" } },
//                  "random_score": {},
//                  "weight": 23
//              },
//              {
//                  "filter": { "match": { "test": "cat" } },
//                  "weight": 42
//              }
//          ],
//          "max_boost": 42,
//          "score_mode": "max",
//          "boost_mode": "multiply",
//          "min_score" : 42
//        }
//    }
func (this *FunctionScoreQuery) BuildBody() (map[string]interface{}, error) {
	query := make(map[string]interface{})
	functionScore := make(map[string]interface{})

	if this.query != nil {
		q, err := this.query.BuildBody()
		if err != nil {
			return nil, err
		}
		functionScore["query"] = q
	}

	if this.boost != nil {
		functionScore["boost"] = this.boost
	}

	if this.functions != nil {
		functions := make([]map[string]interface{}, len(this.functions))

		for i, f := range this.functions {
			body, err := f.BuildBody()
			if err != nil {
				return nil, err
			}
			functions[i] = body
		}

		functionScore["functions"] = functions
	}

	if this.randomScore != nil {
		randomScore, err := this.randomScore.BuildBody()
		if err != nil {
			return nil, err
		}
		functionScore["random_score"] = randomScore
	}

	if this.fieldValueFactor != nil {
		fieldValueFactor, err := this.fieldValueFactor.BuildBody()
		if err != nil {
			return nil, err
		}
		functionScore["field_value_factor"] = fieldValueFactor["field_value_factor"]
	}

	if this.maxBoost != nil {
		functionScore["max_boost"] = this.maxBoost
	}
	if this.scoreMode != nil {
		functionScore["score_mode"] = *this.scoreMode
	}
	if this.boostMode != nil {
		functionScore["boost_mode"] = *this.boostMode
	}
	if this.minScore != nil {
		functionScore["min_score"] = this.minScore
	}

	query["function_score"] = functionScore
	return query, nil
}

// see https://www.elastic.co/guide/en/elasticsearch/reference/6.3/query-dsl-function-score-query.html#function-field-value-factor
type FieldValueFactorQuery struct {
	field    string
	factor   float64
	modifier string
	missing  float64
}

func NewFieldValueFactorQuery(field string, factor float64, modifier string, missing float64) *FieldValueFactorQuery {
	return &FieldValueFactorQuery{field: field, factor: factor, modifier: modifier, missing: missing}
}

// {"field_value_factor":{"field": "likes","factor": 1.2,"modifier": "sqrt","missing": 1}}
func (this *FieldValueFactorQuery) BuildBody() (map[string]interface{}, error) {
	query := make(map[string]interface{})
	fieldValueFactorQuery := make(map[string]interface{})
	if "" == this.field {
		return nil, errors.New("field_value_factor's field can't be ''!")
	}
	fieldValueFactorQuery["field"] = this.field
	fieldValueFactorQuery["factor"] = this.factor
	fieldValueFactorQuery["modifier"] = this.modifier
	fieldValueFactorQuery["missing"] = this.missing

	query["field_value_factor"] = fieldValueFactorQuery
	return query, nil
}

type ScriptScoreQuery struct {
	params map[string]interface{}
	source string
	lang   *string
}

func NewScriptScoreQuery(source string) *ScriptScoreQuery {
	return &ScriptScoreQuery{source: source}
}
func (this *ScriptScoreQuery) Params(key string, value interface{}) *ScriptScoreQuery {
	if this.params == nil {
		this.params = make(map[string]interface{})
	}
	this.params[key] = value
	return this
}
func (this *ScriptScoreQuery) Lang(lang string) *ScriptScoreQuery {
	this.lang = &lang
	return this
}

// {"script_score":{"script":{"params":{...},"source":"..."}}}
func (this *ScriptScoreQuery) BuildBody() (map[string]interface{}, error) {
	query := make(map[string]interface{})
	script := make(map[string]interface{})
	scriptScore := make(map[string]interface{})
	if "" == this.source {
		return nil, errors.New("script_score's source can't be ''!")
	}
	if this.params != nil {
		script["params"] = this.params
	}
	script["source"] = this.source
	if this.lang != nil {
		script["lang"] = this.lang
	}
	scriptScore["script"] = script
	query["script_score"] = scriptScore
	return query, nil
}

// see https://www.elastic.co/guide/en/elasticsearch/reference/6.3/query-dsl-function-score-query.html#function-random
type RandomScoreQuery struct {
	seed  *float64
	field string
}

func NewRandomScoreQuery() *RandomScoreQuery {
	return &RandomScoreQuery{}
}

func (this *RandomScoreQuery) Seed(seed float64) *RandomScoreQuery {
	this.seed = &seed
	return this
}

func (this *RandomScoreQuery) Field(field string) *RandomScoreQuery {
	this.field = field
	return this
}

// {"random_score":{"seed":10,"field":"_seq_no"}}
func (this RandomScoreQuery) BuildBody() (map[string]interface{}, error) {
	query := make(map[string]interface{})
	if this.seed != nil {
		query["seed"] = this.seed
	}
	if "" != this.field {
		query["field"] = this.field
	}
	return query, nil
}

type FunctionQuery struct {
	filter      Query
	randomScore *RandomScoreQuery
	scriptScore *ScriptScoreQuery
	weight      *float64
}

func NewFunctionQuery() *FunctionQuery {
	return &FunctionQuery{}
}
func (this *FunctionQuery) Filter(query Query) *FunctionQuery {
	this.filter = query
	return this
}

func (this *FunctionQuery) RandomScore(randomScore *RandomScoreQuery) *FunctionQuery {
	this.randomScore = randomScore
	return this
}

func (this *FunctionQuery) ScriptScore(scriptScore *ScriptScoreQuery) *FunctionQuery {
	this.scriptScore = scriptScore
	return this
}

func (this *FunctionQuery) Weight(weight float64) *FunctionQuery {
	this.weight = &weight
	return this
}

// {"filter":[...],"weight":1,...}
func (this *FunctionQuery) BuildBody() (map[string]interface{}, error) {
	query := make(map[string]interface{})
	if this.filter != nil {
		if filter, err := this.filter.BuildBody(); err != nil {
			return nil, err
		} else {
			query["filter"] = filter
		}
	}
	if this.randomScore != nil {
		randomScore, err := this.randomScore.BuildBody()
		if err != nil {
			return nil, err
		}
		query["random_score"] = randomScore["random_score"]
	}

	if this.scriptScore != nil {
		if scriptScore, err := this.scriptScore.BuildBody(); err != nil {
			return nil, err
		} else {
			query["script_score"] = scriptScore["script_score"]
		}
	}

	if this.weight != nil {
		query["weight"] = this.weight
	}

	return query, nil
}

// match_all query
type MatchAllQuery struct {
}

// {"match_all":{}}
func (this *MatchAllQuery) BuildBody() (map[string]interface{}, error) {
	query := make(map[string]interface{})
	query["match_all"] = query
	return query, nil
}

//https://www.elastic.co/guide/en/elasticsearch/reference/6.3/query-dsl-wildcard-query.html
type WildcardQuery struct {
	field   string
	keyword string
	boost   *float64
}

func NewWildcardQuery(field string, keyword string) *WildcardQuery {
	return &WildcardQuery{field: field, keyword: keyword}
}

func (this *WildcardQuery) Boost(boost float64) *WildcardQuery {
	this.boost = &boost
	return this
}

/**
{
    "query": {
        "wildcard" : { "user" : { "value" : "ki*y", "boost" : 2.0 } }
    }
}
*/
func (this *WildcardQuery) BuildBody() (map[string]interface{}, error) {
	wildcardQuery := make(map[string]interface{})
	subMatchwildcard := make(map[string]interface{})
	query := make(map[string]interface{})
	if this.field == "" {
		return nil, errors.New("must set wildcard's field")
	}
	subMatchwildcard["value"] = this.keyword

	if this.boost != nil {
		subMatchwildcard["boost"] = this.boost
	}

	wildcardQuery[this.field] = subMatchwildcard
	query["wildcard"] = wildcardQuery

	return query, nil
}
