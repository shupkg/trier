package snowflake

import (
	"sync"
	"time"
)

/*
关于节点位长度取值问题(不借秒的情况下)
节点位  支持节点数  每节点最大秒并发
 0     1         2097152
 1     2         1048576
 2     4         524288
 3     8         262144
 4     16        131072
 5     32        65536
 6     64        32768
 7     128       16384
 8     256       8192
 9     512       4096
10     1024      2048
*/

const (
	nodeBits      = 5
	sequenceBits  = 16
	sequenceMask  = -1 ^ (-1 << sequenceBits)
	sequenceShift = nodeBits
	timeShift     = sequenceShift + sequenceBits
	defaultEpoch  = 1515000000 //2018-01-04T01:20:00+08:00 //为何取这个时间，想在2018年初找一个好记的时间戳, 这个点挺好，1515+6个0
)

func New() *Snowflake {
	return &Snowflake{node: 0, sequence: 1, epoch: defaultEpoch}
}

//32位秒+16位自增+5位机器 = 53位， 兼容js
type Snowflake struct {
	node     int64 //节点
	epoch    int64 //纪元时间（序号时间起始）
	time     int64 //时间
	sequence int64 //序列号 2^16
	mu       sync.Mutex
}

func (s *Snowflake) SetEpoch(startAt time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.epoch = startAt.Unix()
}

func (s *Snowflake) SetNode(node int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.node = node
}

func (s *Snowflake) Next() ID {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := time.Now().Unix()
	seq := s.sequence
	if now <= s.time {
		if now < s.time {
			now = s.time
		}
		//达到最大值，借秒
		if seq = (seq + 1) & sequenceMask; seq == 0 {
			now++
		}
	} else {
		seq = 0
	}
	s.time = now
	s.sequence = seq
	return ID{
		node:     s.node,
		epoch:    s.epoch,
		time:     time.Unix(now, 0),
		sequence: seq,
	}
}
