package flag

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/pflag"
)

type Set struct {
	name string
	args []string
	set  *pflag.FlagSet
}

func New(args ...string) *Set {
	if len(args) == 0 {
		args = os.Args
	}
	name := filepath.Base(args[0])
	pflag.ErrHelp = errors.New("")
	set := pflag.NewFlagSet(name, pflag.ExitOnError)
	set.SortFlags = false
	fs := &Set{args: args[1:], set: set, name: name}
	fs.set.Usage = fs.Usage
	return fs
}

func (f *Set) Usage() {
	_, _ = fmt.Fprintf(os.Stderr, "Usage of %s:\n", f.name)
	f.set.PrintDefaults()
}

func (f *Set) Parse() {
	_ = f.set.Parse(f.args)
}

func (f *Set) VarP(param interface{}, name, shorthand string, value interface{}, usage string) {
	usage = usage + " [" + f.fixEnvKey(name) + "]"
	switch out := param.(type) {
	case *string:
		f.set.StringVarP(out, name, shorthand, f.getEnvString(name, value.(string)), usage)
	case *[]string:
		f.set.StringSliceVarP(out, name, shorthand, f.getEnvStrings(name, value.([]string)), usage)

	case *int64:
		f.set.Int64VarP(out, name, shorthand, f.getEnvInt(name, value.(int64)), usage)
	case *[]int64:
		f.set.Int64SliceVarP(out, name, shorthand, f.getEnvInts(name, value.([]int64)), usage)

	case *float64:
		f.set.Float64VarP(out, name, shorthand, f.getEnvFloat(name, value.(float64)), usage)
	case *[]float64:
		f.set.Float64SliceVarP(out, name, shorthand, f.getEnvFloats(name, value.([]float64)), usage)

	case *time.Duration:
		f.set.DurationVarP(out, name, shorthand, f.getEnvDuration(name, value.(time.Duration)), usage)
	case *[]time.Duration:
		f.set.DurationSliceVarP(out, name, shorthand, f.getEnvDurations(name, value.([]time.Duration)), usage)

	case *bool:
		f.set.BoolVarP(out, name, shorthand, f.getEnvBool(name, value.(bool)), usage)
	default:
		f.Usage()
	}
}

/*Env*/

func (f *Set) fixEnvKey(name string) string {
	name = f.name + "_" + name
	name = strings.ReplaceAll(name, ".", "_")
	name = regexp.MustCompile(`([A-Z]+)`).ReplaceAllString(name, "_$1")
	name = regexp.MustCompile(`([^\w]+)`).ReplaceAllString(name, "_")
	name = strings.ToUpper(strings.Trim(name, "_"))
	return name
}

func (f *Set) getEnv(name string) string {
	return os.Getenv(f.fixEnvKey(name))
}

func (f *Set) getEnvString(name string, def string) string {
	val := f.getEnv(name)
	if val == "" {
		return def
	}
	return val
}

func (f *Set) getEnvStrings(name string, def []string) []string {
	if val := f.getEnv(name); val != "" {
		return strings.Split(val, ",")
	}
	return def
}

func (f *Set) getEnvFloat(name string, def float64) float64 {
	val := f.getEnv(name)
	v, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return def
	}
	return v
}

func (f *Set) getEnvFloats(name string, def []float64) []float64 {
	if val := f.getEnv(name); val != "" {
		valArray := strings.Split(val, ",")
		var ret []float64
		for _, s := range valArray {
			v, err := strconv.ParseFloat(s, 64)
			if err == nil {
				ret = append(ret, v)
			}
		}
		if len(ret) > 0 {
			return ret
		}
	}

	return def
}

func (f *Set) getEnvInt(name string, def int64) int64 {
	val := f.getEnv(name)
	v, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return def
	}
	return v
}

func (f *Set) getEnvInts(name string, def []int64) []int64 {
	if val := f.getEnv(name); val != "" {
		valArray := strings.Split(val, ",")
		var ret []int64
		for _, s := range valArray {
			v, err := strconv.ParseInt(s, 10, 64)
			if err == nil {
				ret = append(ret, v)
			}
		}
		if len(ret) > 0 {
			return ret
		}
	}

	return def
}

func (f *Set) getEnvBool(name string, def bool) bool {
	val := f.getEnv(name)
	v, err := strconv.ParseBool(val)
	if err != nil {
		return def
	}
	return v
}

func (f *Set) getEnvDuration(name string, def time.Duration) time.Duration {
	val := f.getEnv(name)
	v, err := time.ParseDuration(val)
	if err != nil {
		return def
	}
	return v
}

func (f *Set) getEnvDurations(name string, def []time.Duration) []time.Duration {
	if val := f.getEnv(name); val != "" {
		valArray := strings.Split(val, ",")
		var ret []time.Duration
		for _, s := range valArray {
			v, err := time.ParseDuration(s)
			if err == nil {
				ret = append(ret, v)
			}
		}
		if len(ret) > 0 {
			return ret
		}
	}

	return def
}
