package blog

import (
	"fmt"
	"testing"
)

func TestexecAdvQuery(t *testing.T) {

	t.Log(execAdvQuery(map[string]interface{}{
		"table":     "test1",
		"cols":      "id,body,ss",
		"condition": "id=1",
		"order":     "DESC",
		"ordercol":  "id",
		"limit":     "1,5"}))
}

func execAdvQuery(args map[string]interface{}, datas ...interface{}) string {
	sql := "SElECT "
	// 查询柱头
	if args["cols"] != nil {
		sql += args["cols"].(string) + " FROM "
	} else {
		sql += "* FROM "
	}
	// 表名
	if args["table"] == nil {
		return ""
	}
	sql += args["table"].(string)
	// 条件
	if args["condition"] != nil {
		sql += " WHERE " + args["condition"].(string)
	}
	// 排序
	if args["order"] != nil && args["ordercol"] != nil {
		sql += " ORDER BY " + args["ordercol"].(string) + " " + args["order"].(string)
	}
	// 记录选择
	if args["limit"] != nil {
		sql += " LIMIT " + args["limit"].(string)
	}
	fmt.Println("test sql=", sql)
	return sql
}
