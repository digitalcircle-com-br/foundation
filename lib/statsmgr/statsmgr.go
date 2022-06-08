package statsmgr

import (
	"fmt"
	"strings"

	"github.com/digitalcircle-com-br/foundation/lib/redismgr"
)

func SetI(k string, v int64) error {
	return redismgr.Set("stats:"+k, fmt.Sprintf("%v", v))
}

func GetI(k string, v int64) (int64, error) {
	return redismgr.GetI("stats:" + k)
}

func Incr(k string) (int64, error) {
	return redismgr.Incr("stats:" + k)
}

func Decr(k string) (int64, error) {
	return redismgr.Decr("stats:" + k)
}

func GetStats(pattern string) (map[string]int64, error) {
	ks, err := redismgr.Keys("stats:" + pattern)
	if err != nil {
		return nil, err
	}
	var ret map[string]int64
	for _, k := range ks {
		v, err := redismgr.GetI(k)
		if err != nil {
			return nil, err
		}
		nk := strings.Replace(k, "stats:", "", 1)
		ret[nk] = v
	}
	return ret, nil
}
