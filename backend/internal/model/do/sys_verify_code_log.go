// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// SysVerifyCodeLog is the golang structure of table zz_sys_verify_code_log for DAO operations like Where/Data.
type SysVerifyCodeLog struct {
	g.Meta    `orm:"table:zz_sys_verify_code_log, do:true"`
	Id        any         // 日志ID
	Scene     any         // 业务场景
	Channel   any         // 发送通道
	Mobile    any         // 手机号
	Ip        any         // 客户端IP
	Status    any         // 状态
	ErrMsg    any         // 错误信息
	CreatedAt *gtime.Time // 创建时间
}
