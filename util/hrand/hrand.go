package hrand

import (
	"github.com/bamcop/kit"
	"github.com/zlabwork/snowflake"
)

var (
	snow *snowflake.Node
)

func init() {
	var err error
	snow, err = snowflake.NewNode(0)
	kit.Try(err)
}

// SnowId 生成雪花ID, 53bits 版本
func SnowId() int64 {
	return snow.Generate().Int64()
}
