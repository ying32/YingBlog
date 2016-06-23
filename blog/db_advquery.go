package blog

import (
	"database/sql"
	"errors"
)

type QueryResultFunc func(*sql.Rows)

/**
高级查询
	"table":     "test1",
	"cols":      "id,body,ss",
	"condition": "id=1",
	"order":     "DESC",
	"ordercol":  "id",
	"limit":     "1,5"
* **/
func (self *MDB) execAdvQuery(args map[string]interface{}, queryResult QueryResultFunc, datas ...interface{}) (err error) {
	sql := "SElECT "
	// 查询柱头
	if args["cols"] != nil {
		sql += args["cols"].(string) + " FROM "
	} else {
		sql += "* FROM "
	}
	// 表名
	if args["table"] == nil {
		return errors.New("table is empty")
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
	rows, err := self.Query(sql, datas...)
	if err != nil {
		return err
	}
	defer rows.Close()
	if queryResult != nil {
		queryResult(rows)
	}
	return nil
}
