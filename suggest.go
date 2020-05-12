package elastic

type SearchSuggest interface {
	Query
	Name() string
}

// see https://www.elastic.co/guide/en/elasticsearch/reference/6.3/search-suggesters.html
type GlobSuggest struct {
	text      *string
	suggestes []SearchSuggest
}

func NewGlobSuggest() *GlobSuggest {
	return &GlobSuggest{}
}

func (this *GlobSuggest) Text(text string) *GlobSuggest {
	this.text = &text
	return this
}

func (this *GlobSuggest) Suggestes(suggestes ...SearchSuggest) *GlobSuggest {
	this.suggestes = append(this.suggestes, suggestes...)
	return this
}

//{
//  "suggest": {
//    "my-suggest-1" : {
//      "text" : "tring out Elasticsearch",
//      "term" : {
//        "field" : "message"
//      }
//    },
//    "my-suggest-2" : {
//      "text" : "kmichy",
//      "term" : {
//        "field" : "user"
//      }
//    }
//  }
//}
func (this *GlobSuggest) BuildBody() (map[string]interface{}, error) {
	suggest := make(map[string]interface{})
	if this.text != nil {
		suggest["text"] = *this.text
	}
	if this.suggestes != nil {
		for _, s := range this.suggestes {
			body, err := s.BuildBody()
			if err != nil {
				return nil, err
			}
			suggest[s.Name()] = body[s.Name()]
		}
	}
	return suggest, nil
}

// see https://www.elastic.co/guide/en/elasticsearch/reference/6.3/search-suggesters-completion.html
type CompletionSuggest struct {
	prefix         *string
	text           *string
	field          string
	size           *int
	skipDuplicates bool
	fuzzy          bool
	fuzziness      *int
	transpositions *bool
	minLength      *int
	prefixLength   *int
	unicodeAware   bool
	params         map[string]interface{}
	name           string
}

func NewCompletionSuggest(name string, field string) *CompletionSuggest {
	return &CompletionSuggest{name: name, field: field}
}

func (this *CompletionSuggest) Name() string {
	return this.name
}

func (this *CompletionSuggest) Size(size int) *CompletionSuggest {
	this.size = &size
	return this
}

func (this *CompletionSuggest) Prefix(prefix string) *CompletionSuggest {
	this.prefix = &prefix
	return this
}

func (this *CompletionSuggest) Text(text string) *CompletionSuggest {
	this.text = &text
	return this
}

func (this *CompletionSuggest) SkipDuplicates(skipDuplicates bool) *CompletionSuggest {
	this.skipDuplicates = skipDuplicates
	return this
}

func (this *CompletionSuggest) Fuzzy(fuzzy bool) *CompletionSuggest {
	this.fuzzy = fuzzy
	return this
}

func (this *CompletionSuggest) Transpositions(transpositions bool) *CompletionSuggest {
	this.transpositions = &transpositions
	return this
}

func (this *CompletionSuggest) MinLength(minLength int) *CompletionSuggest {
	this.minLength = &minLength
	return this
}

func (this *CompletionSuggest) PrefixLength(prefixLength int) *CompletionSuggest {
	this.prefixLength = &prefixLength
	return this
}

func (this *CompletionSuggest) UnicodeAware(unicodeAware bool) *CompletionSuggest {
	this.unicodeAware = unicodeAware
	return this
}

func (this *CompletionSuggest) Fuzziness(fuzziness int) *CompletionSuggest {
	this.fuzziness = &fuzziness
	return this
}

func (this *CompletionSuggest) Params(key string, value interface{}) *CompletionSuggest {
	if this.params == nil {
		this.params = make(map[string]interface{})
	}

	this.params[key] = value
	return this
}

//    {
//        "song-suggest" : {
//            "prefix" : "nir",
//            "completion" : {
//                "field" : "suggest"
//            }
//        }
//    }
func (this *CompletionSuggest) BuildBody() (map[string]interface{}, error) {
	completion := make(map[string]interface{})
	completion["field"] = this.field
	if this.size != nil {
		completion["size"] = this.size
	}

	if this.skipDuplicates {
		completion["skip_duplicates"] = this.skipDuplicates
	}

	suggestBody := make(map[string]interface{})
	if this.fuzzy {
		completion["fuzzy"] = this.fuzzy
	} else {
		useFuzzy := false
		fuzzy := make(map[string]interface{})
		if this.fuzziness != nil {
			useFuzzy = true
			fuzzy["fuzziness"] = *this.fuzziness
		}
		if this.transpositions != nil {
			useFuzzy = true
			fuzzy["transpositions"] = *this.transpositions
		}
		if this.minLength != nil {
			useFuzzy = true
			fuzzy["min_length"] = *this.minLength
		}
		if this.prefixLength != nil {
			useFuzzy = true
			fuzzy["prefix_length"] = *this.prefixLength
		}
		if this.unicodeAware {
			useFuzzy = true
			fuzzy["unicode_aware"] = this.unicodeAware
		}
		if useFuzzy {
			completion["fuzzy"] = fuzzy
		}
	}

	suggest := make(map[string]interface{})
	if this.prefix != nil {
		suggest["prefix"] = *this.prefix
	}
	if this.text != nil {
		suggest["text"] = *this.text
	}
	suggest["completion"] = completion
	if this.params != nil {
		for k, v := range this.params {
			suggest[k] = v
		}
	}

	suggestBody["suggest"] = map[string]interface{}{this.name: suggest}

	return suggestBody, nil
}
