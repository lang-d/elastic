package elastic

import "errors"

type ScriptInter interface {
	Query
}

type Script struct {
	lang   string
	source string
	id     string
	params []map[string]interface{}
}

func NewScript() *Script {
	return &Script{}
}

func (this *Script) Lang(lang string) *Script {
	this.lang = lang
	return this
}

func (this *Script) Source(source string) *Script {
	this.source = source
	return this
}

func (this *Script) Id(id string) *Script {
	this.id = id
	return this
}

// other extra condition add by this func
func (this *Script) Params(key string, value interface{}) *Script {
	param := make(map[string]interface{})
	param[key] = value
	this.params = append(this.params, param)
	return this
}

// return {"lang":"painless","id":"xx","source":"doc['time'].value"}
func (this *Script) BuildBody() (map[string]interface{}, error) {
	if "" == this.source && "" == this.id {
		return nil, errors.New("script must have id or source")
	}

	query := make(map[string]interface{})
	params := make(map[string]interface{})
	if this.id != "" {
		query["id"] = this.id
	}
	if this.source != "" {
		query["source"] = this.source
	}
	if this.lang != "" {
		query["lang"] = this.lang
	}

	if len(this.params) > 0 {
		for _, param := range this.params {
			for k, v := range param {
				params[k] = v
			}
		}
	}

	return query, nil
}
