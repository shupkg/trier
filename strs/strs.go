package strs

import (
	"bytes"
	"math/bits"
	"math/rand"
	"reflect"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
	"unicode/utf8"
)

const (
	Num    Spec = 1 << iota //包含数字
	Symbol                  //包含非字母数字符号
	Upper                   //包含大写字母
	Lower                   //喊喊小写字母

	Alphabet = Upper | Lower     //字母（大写和小写）
	NoSymbol = Alphabet | Num    //字母和数字
	All      = NoSymbol | Symbol //所有字符串

	maxSpec = 4 //枚举种类数量
)

var (
	rnd  = rand.New(rand.NewSource(time.Now().UnixNano()))
	cMap = map[Spec]string{
		Num:    "0123456789",
		Symbol: ":;<=>?@[\\]^_`{|}~",
		Upper:  "ABCDEFGHIJKLMNOPQRSTUVWXYZ",
		Lower:  "abcdefghijklmnopqrstuvwxyz",
	}
)

//Spec 随机字符串枚举类型
type Spec uint

func (s Spec) String() string {
	return strconv.FormatUint(uint64(s), 2)
}

//Len 枚举包含的字符种类数量
func (s Spec) Len() int {
	return bits.OnesCount(uint(s))
}

// RandString 产生随机字符串
// @auth 树哥<shu@aryao.com> 2020/07/17 14:03:42
// @param	size		string		随机字符串长度
// @param	spec		Spec枚举	类型
// @param	distinct	boolean		不允许出现重复字符
// @return  随机字符串
func RandString(size int, spec Spec, distinct bool) string {
	if size <= 0 {
		return ""
	}
	tl := spec.Len()
	if tl == 0 {
		return strings.Repeat("0", size)
	}

	b := make([]byte, size)
	c := make([]byte, 0, 80)
	var j int
	for i := 0; i < maxSpec; i++ {
		p := Spec(1 << i)
		if spec&p != 0 {
			s := []byte(cMap[p])
			r := rnd.Intn(len(s))
			s[0], s[r] = s[r], s[0]
			b[j] = s[0]
			c = append(c, s[1:]...)
			j++
			if j == size {
				return UnsafeString(b)
			}
		}
	}

	if distinct && size-j > len(c) {
		distinct = false
	}

	if !distinct {
		c = append(c, b[:j]...)
	}

	for ; j < size; j++ {
		r := rnd.Intn(len(c))
		c[0], c[r] = c[r], c[0]
		b[j] = c[0]
		if distinct {
			c = c[1:]
		}
	}

	for i := 0; i < tl; i++ {
		if r := rnd.Intn(size - 1); i != r {
			b[i], b[r] = b[r], b[i]
		}
	}

	return UnsafeString(b)
}

// RandHs 产生随机汉字 汉字范围: 0x4e00 - 0x9fa5
// @auth 树哥<shu@aryao.com> 2020/07/17 14:03:42
// @param	size		string		随机字符串长度
// @param	distinct	boolean		不允许出现重复字符
// @return  随机汉字
func RandHs(size int, distinct bool) string {
	b := make([]byte, size*3)
	s := rune(0x4e00)
	l := rune(0x9fa5) - s

	if distinct && rune(size) > l {
		distinct = false
	}

	t := make([]byte, 3)
	var n int
	for i := 0; i < size; i++ {
		n = utf8.EncodeRune(t, rand.Int31n(l)+s)
		if distinct {
			for bytes.Contains(b, t) {
				n = utf8.EncodeRune(t, rand.Int31n(l)+s)
			}
		}
		copy(b[i*3:i*3+n], t)
	}
	return UnsafeString(b)
}

func RandSlice(s []interface{}) {
	for i := 0; i < len(s); i++ {
		r := rnd.Intn(len(s))
		s[i], s[r] = s[r], s[i]
	}
}

func RandSlice2(s interface{}) {
	if s == nil {
		return
	}

	v := reflect.Indirect(reflect.ValueOf(s))
	if v.Kind() != reflect.Slice && v.Kind() != reflect.Array {
		return
	}

	l := v.Len()
	for i := 0; i < l; i++ {
		j := rnd.Intn(l)
		vi := reflect.ValueOf(v.Index(i).Interface())
		vj := reflect.ValueOf(v.Index(j).Interface())
		v.Index(i).Set(vj)
		v.Index(j).Set(vi)
	}
}

var idx int32

func Incr() int {
	if atomic.LoadInt32(&idx) > 100000 {
		atomic.StoreInt32(&idx, 1)
	}
	return int(atomic.AddInt32(&idx, 3+rnd.Int31n(7)))
}

func Incrs() string {
	return strconv.Itoa(Incr())
}
