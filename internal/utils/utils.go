package utils

import (
	"github.com/jinzhu/copier"
	"time"
)

var TimeToStringConverter = copier.TypeConverter{
	SrcType: "",
	DstType: time.Time{},
	Fn: func(src interface{}) (interface{}, error) {
		return time.Parse("2006-01-02 15:04:05", src.(string))
	},
}
var StringToTimeConverter = copier.TypeConverter{
	SrcType: time.Time{},
	DstType: "",
	Fn: func(src interface{}) (interface{}, error) {
		return src.(time.Time).Format("2006-01-02 15:04:05"), nil
	},
}

var CopierOption = copier.Option{
	IgnoreEmpty: true,
	DeepCopy:    true,
	Converters: []copier.TypeConverter{
		TimeToStringConverter,
		StringToTimeConverter,
	},
}
