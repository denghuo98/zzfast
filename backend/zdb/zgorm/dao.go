package zgorm

import (
	"fmt"
	"sync"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

type daoInstance interface {
	Table() string
	Ctx(ctx g.Ctx) *gdb.Model
}

// pkFieldCache 缓存各表主键字段名，避免每次唯一性校验都反射查询表结构
var pkFieldCache sync.Map // tableName -> pkFieldName(string)

// GetPkField 获取dao实例中的主键名称
func GetPkField(ctx g.Ctx, dao daoInstance) (string, error) {
	table := dao.Table()
	if v, ok := pkFieldCache.Load(table); ok {
		return v.(string), nil
	}

	// 获取表字段
	fields, err := dao.Ctx(ctx).TableFields(table)
	if err != nil {
		return "", err
	}
	if len(fields) == 0 {
		return "", gerror.New("field not found")
	}

	for _, field := range fields {
		// 获取主键，支持大写和小写的 PRI
		if field.Key == "PRI" || field.Key == "pri" {
			pkFieldCache.Store(table, field.Name)
			return field.Name, nil
		}
	}
	return "", gerror.New("no primary key")
}

// IsUnique 判断属性是否唯一
// dao 数据库操作实例
// where 查询条件
// message 错误消息
// pkId 主键ID, 用于排除当前记录,默认不传则不排除
func IsUnique(ctx g.Ctx, dao daoInstance, where g.Map, message string, pkId ...interface{}) error {
	if len(where) == 0 {
		return gerror.New("查询条件不能为空")
	}

	m := dao.Ctx(ctx).Where(where)
	if len(pkId) > 0 {
		field, err := GetPkField(ctx, dao)
		if err != nil {
			return err
		}
		m = m.WhereNot(field, pkId)
	}

	count, err := m.Count(1)
	if err != nil {
		return err
	}
	if count > 0 {
		if message == "" {
			for k := range where {
				message = fmt.Sprintf("数据表%s中%s已经存在", dao.Table(), where[k])
			}
		}
		return gerror.New(message)
	}
	return nil
}
