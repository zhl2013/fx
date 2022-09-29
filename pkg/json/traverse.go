package json

import (
	"encoding/json"
	"fmt"
	"strconv"

	. "github.com/antonmedv/fx/pkg/dict"
)

type Iterator struct {
	Object       interface{}
	Path, Parent string
}

func Dfs(object interface{}, f func(it Iterator)) {
	sub(Iterator{Object: object}, f)
}

func sub(it Iterator, f func(it Iterator)) {
	f(it)
	switch it.Object.(type) {
	case *Dict:
		keys := it.Object.(*Dict).Keys
		for _, k := range keys {
			subpath := it.Path + "." + k
			value, _ := it.Object.(*Dict).Get(k)
			sub(Iterator{
				Object: value,
				Path:   subpath,
				Parent: it.Path,
			}, f)
		}

	case Array:
		slice := it.Object.(Array)
		for i, value := range slice {
			subpath := accessor(it.Path, i)
			sub(Iterator{
				Object: value,
				Path:   subpath,
				Parent: it.Path,
			}, f)
		}
	}
}

func GetVal(obj interface{}, path string) string {
	var res interface{}
	Dfs(obj, func(it Iterator) {
		if it.Path == path {
			res = it.Object
		}
	})

	switch v := res.(type) {
	case nil:
		return "can not find val"
	case bool:
		return strconv.FormatBool(v)
	case Number:
		return "num"
	case string:
		return res.(string)
	case *Dict, Array:
		v2, _ := json.Marshal(v)
		return string(v2)
	default:
		return "unknown type"
	}
}

func accessor(path string, to interface{}) string {
	return fmt.Sprintf("%v[%v]", path, to)
}
