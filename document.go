package elastic

import (
	"encoding/json"
	"errors"
)

type Action struct {
	OpType          string
	Index           string
	DocType         string
	Id              string
	Data            map[string]interface{}
	Version         string
	Routing         string
	Refresh         *int
	RetryOnConflict *int
	DocAsUpsert     bool
}

func (this Action) Format() ([]byte, error) {
	if this.Index == "" {
		return nil, errors.New("no index set")
	}
	if this.DocType == "" {
		return nil, errors.New("no doc_type set")
	}
	if this.OpType == "" {
		this.OpType = "index"
	}
	op := make(map[string]map[string]interface{})
	op[this.OpType] = map[string]interface{}{"_index": this.Index, "_type": this.DocType}
	if this.Id != "" {
		op[this.OpType]["_id"] = this.Id
	}
	if this.Version != "" {
		op[this.OpType]["version"] = this.Version
	}
	if this.Routing != "" {
		op[this.OpType]["routing"] = this.Routing
	}
	if this.Refresh != nil {
		op[this.OpType]["refresh"] = *this.Refresh
	}
	if this.RetryOnConflict != nil {
		op[this.OpType]["retry_on_conflict"] = *this.Refresh
	}
	if this.DocAsUpsert {
		op[this.OpType]["doc_as_upsert"] = this.DocAsUpsert
	}

	opByte, err := json.Marshal(op)
	if err != nil {
		return nil, err
	}
	if this.OpType == "delete" {
		return opByte, nil
	}

	if this.Data == nil {
		return nil, errors.New("no data set")
	}
	if this.OpType == "update" {
		dataByte, err := json.Marshal(map[string]map[string]interface{}{"doc": this.Data})
		if err != nil {
			return nil, err
		}
		opByte = append(opByte, []byte("\n")...)
		return append(opByte, dataByte...), nil
	}

	dataByte, err := json.Marshal(this.Data)
	if err != nil {
		return nil, err
	}
	opByte = append(opByte, []byte("\n")...)
	return append(opByte, dataByte...), nil

}
