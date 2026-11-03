// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// SysVerifyCodeLog is the golang structure for table sys_verify_code_log.
type SysVerifyCodeLog struct {
	Id        int64       `json:"id"        orm:"id"         description:"日志ID"`
	Scene     string      `json:"scene"     orm:"scene"      description:"业务场景"`
	Channel   string      `json:"channel"   orm:"channel"    description:"发送通道"`
	Mobile    string      `json:"mobile"    orm:"mobile"     description:"手机号"`
	Ip        string      `json:"ip"        orm:"ip"         description:"客户端IP"`
	Status    int         `json:"status"    orm:"status"     description:"状态"`
	ErrMsg    string      `json:"errMsg"    orm:"err_msg"    description:"错误信息"`
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" description:"创建时间"`
}
