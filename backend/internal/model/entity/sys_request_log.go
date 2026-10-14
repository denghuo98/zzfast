// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// SysRequestLog is the golang structure for table sys_request_log.
type SysRequestLog struct {
	Id        int64       `json:"id"        orm:"id"         description:"日志ID"`
	ReqId     string      `json:"reqId"     orm:"req_id"     description:"请求ID"`
	Method    string      `json:"method"    orm:"method"     description:"请求方法"`
	Path      string      `json:"path"      orm:"path"       description:"请求路径"`
	MemberId  int64       `json:"memberId"  orm:"member_id"  description:"用户ID"`
	Username  string      `json:"username"  orm:"username"   description:"用户名"`
	Ip        string      `json:"ip"        orm:"ip"         description:"客户端IP"`
	UserAgent string      `json:"userAgent" orm:"user_agent" description:"UA信息"`
	Status    int         `json:"status"    orm:"status"     description:"响应状态码"`
	ErrorMsg  string      `json:"errorMsg"  orm:"error_msg"  description:"错误信息"`
	CostTime  int64       `json:"costTime"  orm:"cost_time"  description:"耗时(毫秒)"`
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" description:"创建时间"`
}
