package main

import (
	"fmt"
	"time"

	"github.com/jinzhu/copier"
)

func main() {
	// 定义自定义转换器
	stringToTimeConverter := copier.TypeConverter{
		SrcType: "",
		DstType: time.Time{},
		Fn: func(src interface{}) (interface{}, error) {
			// 将字符串转换为 time.Time
			srcStr, ok := src.(string)
			if !ok {
				return nil, fmt.Errorf("expected string, got %T", src)
			}
			parsedTime, err := time.Parse(time.RFC3339, srcStr)
			if err != nil {
				return nil, err
			}
			return parsedTime, nil
		},
	}

	// 定义源和目标结构体
	type Source struct {
		Born string
	}
	type Target struct {
		Born time.Time
	}

	source := Source{
		Born: "2024-12-01T08:42:12.759Z",
	}
	var target Target

	// 使用 CopyWithOption 进行复制
	err := copier.CopyWithOption(&target, &source, copier.Option{
		Converters: []copier.TypeConverter{stringToTimeConverter},
	})
	if err != nil {
		fmt.Println("Error copying:", err)
		return
	}

	// 输出结果
	fmt.Printf("Source: %+v\n", source)
	fmt.Printf("Target: %+v\n", target)
}
