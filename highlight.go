package elastic

import "errors"

type SearchHighlight struct {
	numberOfFragments *int
	fragmentSize      *int
	preTags           []string
	postTags          []string
	fields            []SearchHighlightField
	params            []map[string]interface{}
}

func NewSearchHighlight() *SearchHighlight {
	return &SearchHighlight{}
}

// add one or many preTag
func (this *SearchHighlight) PreTags(preTags ...string) *SearchHighlight {
	this.preTags = append(this.preTags, preTags...)
	return this
}

// add one or many postTag
func (this *SearchHighlight) PostTags(postTags ...string) *SearchHighlight {
	this.postTags = append(this.postTags, postTags...)
	return this
}

func (this *SearchHighlight) NumberOfFragments(numberOfFragments int) *SearchHighlight {
	this.numberOfFragments = &numberOfFragments
	return this
}

func (this *SearchHighlight) FragmentSize(fragmentSize int) *SearchHighlight {
	this.fragmentSize = &fragmentSize
	return this
}

func (this *SearchHighlight) Fields(fields ...SearchHighlightField) *SearchHighlight {
	this.fields = append(this.fields, fields...)
	return this
}

// other extra condition add by this func
func (this *SearchHighlight) Params(key string, value interface{}) *SearchHighlight {
	param := make(map[string]interface{})
	param[key] = value
	this.params = append(this.params, param)
	return this
}

func (this *SearchHighlight) BuildBody() (map[string]interface{}, error) {
	highlight := make(map[string]interface{})
	if this.numberOfFragments != nil {
		highlight["number_of_fragments"] = this.numberOfFragments
	}
	if this.fragmentSize != nil {
		highlight["fragment_size"] = this.fragmentSize
	}
	if len(this.fields) > 0 {
		fields := make([]map[string]interface{}, len(this.fields))
		for i, field := range this.fields {
			f, err := field.BuildBody()
			if err != nil {
				return nil, err
			}
			fields[i] = f
		}
		highlight["fields"] = fields
	}
	if len(this.preTags) > 0 {
		highlight["pre_tags"] = this.preTags
	}
	if len(this.postTags) > 0 {
		highlight["post_tags"] = this.postTags
	}
	// other conditions
	for _, param := range this.params {
		for k, v := range param {
			highlight[k] = v
		}
	}

	return highlight, nil
}

type SearchHighlightField struct {
	field             string
	numberOfFragments *int
	fragmentSize      *int
	preTags           []string
	postTags          []string
	highlightQuery    Query
	params            []map[string]interface{}
}

func NewSearchHighlightField(field string) *SearchHighlightField {
	return &SearchHighlightField{field: field}
}

// add one or many preTag
func (this *SearchHighlightField) PreTags(preTags ...string) *SearchHighlightField {
	this.preTags = append(this.preTags, preTags...)
	return this
}

// add one or many postTag
func (this *SearchHighlightField) PostTags(postTags ...string) *SearchHighlightField {
	this.postTags = append(this.postTags, postTags...)
	return this
}

func (this *SearchHighlightField) NumberOfFragments(numberOfFragments int) *SearchHighlightField {
	this.numberOfFragments = &numberOfFragments
	return this
}

func (this *SearchHighlightField) FragmentSize(fragmentSize int) *SearchHighlightField {
	this.fragmentSize = &fragmentSize
	return this
}

// other extra condition add by this func
func (this *SearchHighlightField) Params(key string, value interface{}) *SearchHighlightField {
	param := make(map[string]interface{})
	param[key] = value
	this.params = append(this.params, param)
	return this
}

//{"field":{"pre_tags":[],"post_tags":[]}}
func (this *SearchHighlightField) BuildBody() (map[string]interface{}, error) {
	highlight := make(map[string]interface{})
	highlightField := make(map[string]interface{})
	if "" == this.field {
		return nil, errors.New("highlight field cannot be null!")
	}
	if len(this.preTags) > 0 {
		highlightField["pre_tags"] = this.preTags
	}
	if len(this.postTags) > 0 {
		highlightField["post_tags"] = this.postTags
	}
	if this.fragmentSize != nil {
		highlightField["fragment_size"] = this.fragmentSize
	}
	if this.numberOfFragments != nil {
		highlightField["number_of_fragments"] = this.numberOfFragments
	}
	if this.highlightQuery != nil {
		highlightQuery, err := this.highlightQuery.BuildBody()
		if err != nil {
			return nil, err
		}
		highlightField["highlight_query"] = highlightQuery
	}

	// other conditions
	for _, param := range this.params {
		for k, v := range param {
			highlightField[k] = v
		}
	}
	highlight[this.field] = highlightField
	return highlight, nil
}
