package js

import (
	"encoding/json"
	"syscall/js"
)

type EventEmitter struct {
	document js.Value
	Name     string
}

func NewJSEventEmitter(name string) (*EventEmitter, error) {
	document := js.Global().Get("document")
	res := EventEmitter{
		Name:     name,
		document: document,
	}
	return &res, nil
}

func (j *EventEmitter) Emit(v interface{}) {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}

	val := js.Global().Get("JSON").Call("parse", string(b))
	j.document.Call(j.Name, val)
}
