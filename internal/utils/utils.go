package utils

import (
	"errors"
	"fmt"
	"github.com/jinzhu/copier"
	"time"
)

var StringToTimeConverter = copier.TypeConverter{
	SrcType: "",
	DstType: time.Time{},
	Fn: func(src interface{}) (interface{}, error) {
		// 将字符串转换为 time.Time
		srcStr, ok := src.(string)
		if !ok {
			return nil, fmt.Errorf("expected string, got %T", src)
		}
		parsedTime, err := time.Parse("2006-01-02 15:04:05", srcStr)
		if err != nil {
			return nil, err
		}
		return parsedTime, nil
	},
}
var TimetoStringConverter = copier.TypeConverter{
	SrcType: time.Time{},
	DstType: "",
	Fn: func(src interface{}) (interface{}, error) {
		// 将字符串转换为 time.Time
		t, ok := src.(time.Time)
		if !ok {
			return nil, errors.New("")
		}
		return t.Format("2006-01-02 15:04:05"), nil
	},
}
var CopierOption = copier.Option{
	Converters: []copier.TypeConverter{StringToTimeConverter, TimetoStringConverter},
}
