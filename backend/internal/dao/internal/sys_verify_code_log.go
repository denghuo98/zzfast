// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// SysVerifyCodeLogDao is the data access object for the table zz_sys_verify_code_log.
type SysVerifyCodeLogDao struct {
	table    string                  // table is the underlying table name of the DAO.
	group    string                  // group is the database configuration group name of current DAO.
	columns  SysVerifyCodeLogColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler      // handlers for customized model modification.
}

// SysVerifyCodeLogColumns defines and stores column names for the table zz_sys_verify_code_log.
type SysVerifyCodeLogColumns struct {
	Id        string // 日志ID
	Scene     string // 业务场景
	Channel   string // 发送通道
	Mobile    string // 手机号
	Ip        string // 客户端IP
	Status    string // 状态
	ErrMsg    string // 错误信息
	CreatedAt string // 创建时间
}

// sysVerifyCodeLogColumns holds the columns for the table zz_sys_verify_code_log.
var sysVerifyCodeLogColumns = SysVerifyCodeLogColumns{
	Id:        "id",
	Scene:     "scene",
	Channel:   "channel",
	Mobile:    "mobile",
	Ip:        "ip",
	Status:    "status",
	ErrMsg:    "err_msg",
	CreatedAt: "created_at",
}

// NewSysVerifyCodeLogDao creates and returns a new DAO object for table data access.
func NewSysVerifyCodeLogDao(handlers ...gdb.ModelHandler) *SysVerifyCodeLogDao {
	return &SysVerifyCodeLogDao{
		group:    "default",
		table:    "zz_sys_verify_code_log",
		columns:  sysVerifyCodeLogColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *SysVerifyCodeLogDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current DAO.
func (dao *SysVerifyCodeLogDao) Table() string {
	return dao.table
}

// Columns returns all column names of current DAO.
func (dao *SysVerifyCodeLogDao) Columns() SysVerifyCodeLogColumns {
	return dao.columns
}

// Group returns the database configuration group name of current DAO.
func (dao *SysVerifyCodeLogDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for current operation.
func (dao *SysVerifyCodeLogDao) Ctx(ctx context.Context) *gdb.Model {
	model := dao.DB().Model(dao.table)
	for _, handler := range dao.handlers {
		model = handler(model)
	}
	return model.Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
func (dao *SysVerifyCodeLogDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
