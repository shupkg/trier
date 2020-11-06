package snowflake

import (
	"fmt"
	"strconv"
	"time"
)

var lCST = time.FixedZone("CST", 28800)

const timeFormat = "20060102150405"

type ID struct {
	node     int64
	epoch    int64
	time     time.Time
	sequence int64
}

//Int int
func (id ID) Int() int64 {
	return ((id.time.Unix() - id.epoch) << timeShift) | (id.sequence << sequenceShift) | id.node
}

//Time 时间
func (id ID) Time() time.Time {
	return id.time
}

//Sequence 序号
func (id ID) Sequence() int64 {
	return id.sequence
}

//Node 节点
func (id ID) Node() int64 {
	return id.node
}

//Format 格式化
func (id ID) Format(radix int) string {
	if radix < 2 {
		radix = 36
	}
	return strconv.FormatInt(id.Int(), radix)
}

func (id ID) String() string {
	return fmt.Sprintf(
		"%s%05d%02d",
		id.time.In(lCST).Format(timeFormat),
		id.sequence,
		id.node,
	)
}

func Parse(id int64, epoch int64) ID {
	if epoch == 0 {
		epoch = defaultNode.epoch
	}
	return ID{
		epoch:    epoch,
		time:     time.Unix(id>>timeShift+epoch, 0),
		sequence: id>>sequenceShift - id>>timeShift<<sequenceBits,
		node:     id - id>>sequenceShift<<sequenceShift,
	}
}

func ParseFormat(s string, radix int, epoch int64) ID {
	if radix < 2 {
		radix = 36
	}
	i, _ := strconv.ParseInt(s, radix, 64)
	return Parse(i, epoch)
}

func ParseString(s string, epoch int64) ID {
	if len(s) == 21 {
		if epoch == 0 {
			epoch = defaultNode.epoch
		}
		var (
			tim, _ = time.ParseInLocation(timeFormat, s[:14], lCST)
			seq, _ = strconv.ParseInt(s[14:19], 10, 64)
			nod, _ = strconv.ParseInt(s[19:], 10, 64)
		)
		return ID{
			epoch:    epoch,
			time:     tim,
			sequence: seq,
			node:     nod,
		}
	}
	return ID{}
}
