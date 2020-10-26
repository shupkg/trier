package dict

import (
	"reflect"

	"github.com/shupkg/trier/json"
)

type Kind uint

const (
	Invalid  = Kind(reflect.Invalid)
	Bool     = Kind(reflect.Bool)
	Int64    = Kind(reflect.Int64)
	Float64  = Kind(reflect.Float64)
	String   = Kind(reflect.String)
	Time     = Kind(reflect.UnsafePointer + 1)
	Duration = Kind(reflect.UnsafePointer + 2)
)

var kindNames = map[string]Kind{
	"invalid":  Invalid,
	"bool":     Bool,
	"int64":    Int64,
	"float64":  Float64,
	"string":   String,
	"time":     Time,
	"duration": Duration,
}

func (kind Kind) MarshalBinary() (data []byte, err error) {
	return kind.MarshalJSON()
}

func (kind *Kind) UnmarshalBinary(data []byte) error {
	return kind.UnmarshalJSON(data)
}

func (kind Kind) MarshalJSON() ([]byte, error) {
	return json.Marshal(kind.String())
}

func (kind *Kind) UnmarshalJSON(bytes []byte) error {
	var s string
	if err := json.Unmarshal(bytes, &s); err != nil {
		return err
	}
	if k, find := kindNames[s]; find {
		*kind = k
	} else {
		*kind = Invalid
	}

	return nil
}

func (kind Kind) String() string {
	if kind < Time {
		return reflect.Kind(kind).String()
	}
	switch kind {
	case Time:
		return "time"
	case Duration:
		return "duration"
	default:
		return reflect.Invalid.String()
	}
}
