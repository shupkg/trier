package uuid

import (
	"sync"
	"time"
)

const (
	nodeBits      = 5
	sequenceBits  = 16
	sequenceMask  = -1 ^ (-1 << sequenceBits)
	sequenceShift = nodeBits
	timeShift     = sequenceShift + sequenceBits
)

var defaultNoe *Snowflake

func init() {
	defaultNoe = &Snowflake{}
	defaultNoe.SetEpoch(time.Date(2018, 1, 1, 0, 0, 0, 0, time.UTC))
}

//32位秒+16位自增+5位机器 = 53位， 兼容js
type Snowflake struct {
	node     int64 //节点
	time     int64 //商家
	sequence int64 //序列号 2^16
	epoch    int64 //纪元时间（序号时间起始）
	sync.Mutex
}

func (s *Snowflake) SetEpoch(epoch time.Time) {
	s.Lock()
	defer s.Unlock()
	s.epoch = epoch.Unix()
}

func (s *Snowflake) SetNode(node int64) {
	s.Lock()
	defer s.Unlock()
	s.node = node
}

func (s *Snowflake) Generate() int64 {
	s.Lock()
	defer s.Unlock()
	now := time.Now().Unix()
	if now <= s.time {
		if s.sequence = (s.sequence + 1) & sequenceMask; s.sequence == 0 {
			now++
		}
	} else {
		s.sequence = 0
	}
	s.time = now
	r := ((now - s.epoch) << timeShift) | (s.sequence << sequenceShift) | s.node
	return r
}

func Apply(apply func(node *Snowflake)) {
	apply(defaultNoe)
}

//32位秒时间戳+16位自增+5位机器
func GetID() int64 {
	return defaultNoe.Generate()
}

func GetTime(id int64) time.Time {
	return time.Unix(id>>timeShift+defaultNoe.epoch, 0)
}
