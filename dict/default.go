package dict

import (
	"sync"
	"time"

	"github.com/shupkg/trier/json"
)

const keyType = "__KeyTypes"

func New() Dict {
	return &stat{}
}

type stat struct {
	data map[string]interface{}
	once sync.Once
}

func (s *stat) init() {
	s.data = map[string]interface{}{keyType: map[string]Kind{}}
}

func (s *stat) getKeyType(key string) Kind {
	return s.Types()[key]
}

func (s *stat) setKeyType(key string, kind Kind) {
	types := s.Types()
	types[key] = kind
	s.data[keyType] = types
}

func (s *stat) Set(key string, valueKind Kind, value interface{}) {
	if key == keyType {
		return
	}
	s.once.Do(s.init)
	keyType := s.getKeyType(key)
	if keyType == Invalid {
		s.setKeyType(key, valueKind)
		keyType = valueKind
	}
	if keyType != valueKind {
		return
	}
	s.data[key] = value
}

func (s *stat) Get(key string, valueKind Kind) (value interface{}) {
	s.once.Do(s.init)
	if keyType := s.getKeyType(key); keyType != valueKind {
		return
	}
	return s.data[key]
}

func (s *stat) Del(key string) {
	if key == keyType {
		return
	}
	delete(s.Types(), key)
	delete(s.data, key)
}

func (s *stat) Types() map[string]Kind {
	s.once.Do(s.init)
	return s.data[keyType].(map[string]Kind)
}

func (s *stat) Keys() []string {
	var keys []string
	for key := range s.Types() {
		keys = append(keys, key)
	}
	return keys
}

func (s *stat) SetString(key string, value string) {
	if value != "" {
		s.Set(key, String, value)
	} else {
		s.Del(key)
	}
}

func (s *stat) SetInt(key string, value int64) {
	if value != 0 {
		s.Set(key, Int64, value)
	} else {
		s.Del(key)
	}
}

func (s *stat) SetFloat(key string, value float64) {
	if value != 0 {
		s.Set(key, Float64, value)
	} else {
		s.Del(key)
	}
}

func (s *stat) SetBool(key string, value bool) {
	if value {
		s.Set(key, Bool, value)
	} else {
		s.Del(key)
	}
}

func (s *stat) SetTime(key string, value time.Time) {
	if !value.IsZero() {
		s.Set(key, Time, value.Format(time.RFC3339))
	} else {
		s.Del(key)
	}
}

func (s *stat) SetDuration(key string, value time.Duration) {
	if value > 0 {
		s.Set(key, Duration, value.String())
	} else {
		s.Del(key)
	}
}

func (s *stat) GetString(key string) (value string) {
	if v := s.Get(key, String); v != nil {
		return v.(string)
	}
	return ""
}

func (s *stat) GetInt(key string) (value int64) {
	if v := s.Get(key, Int64); v != nil {
		return v.(int64)
	}
	return 0
}

func (s *stat) GetFloat(key string) (value float64) {
	if v := s.Get(key, Float64); v != nil {
		return v.(float64)
	}
	return 0
}

func (s *stat) GetBool(key string) (value bool) {
	if v := s.Get(key, Bool); v != nil {
		return v.(bool)
	}
	return false
}

func (s *stat) GetTime(key string) (value time.Time) {
	if s := s.Get(key, Time); s != nil {
		value, _ = time.Parse(time.RFC3339, s.(string))
	}
	return
}

func (s *stat) GetDuration(key string) (value time.Duration) {
	if s := s.Get(key, Duration); s != "" {
		value, _ = time.ParseDuration(s.(string))
	}
	return
}

func (s *stat) Copy(stat Dict) {
	stat.Each(func(key string, kind Kind, value interface{}) bool {
		s.Set(key, kind, value)
		return true
	})
}

func (s *stat) AddInt(key string, stat Dict) {
	s.SetInt(key, s.GetInt(key)+stat.GetInt(key))
}

func (s *stat) Clone() Dict {
	var ns = &stat{}
	ns.Copy(s)
	return ns
}

func (s *stat) Each(eachFunc func(key string, kind Kind, value interface{}) bool) {
	types := s.Types()
	for key, kind := range types {
		if value, find := s.data[key]; find {
			if !eachFunc(key, kind, value) {
				break
			}
		}
	}
}

func (s *stat) MarshalBinary() (data []byte, err error) {
	return json.MarshalIndent(s.data, "", "  ")
}

func (s *stat) UnmarshalBinary(data []byte) error {
	s.once.Do(s.init)
	return json.Unmarshal(data, &s.data)
}

var _ Dict = (*stat)(nil)
