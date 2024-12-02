package read

import "reflect"

func ToGenericSlice(elements ...any) []any {
	var res []interface{}
	for _, e := range elements {
		val := reflect.ValueOf(e)
		if val.Kind() == reflect.Slice {
			for i := 0; i < val.Len(); i++ {
				res = append(res, val.Index(i).Interface())
			}
		} else {
			res = append(res, e)
		}
	}
	return res
}
