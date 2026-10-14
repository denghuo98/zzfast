package middleware

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gtime"

	"github.com/denghuo98/zfast/internal/model/entity"
	"github.com/denghuo98/zfast/web/zcontext"
	"github.com/denghuo98/zfast/web/zqueue"
	"github.com/denghuo98/zfast/zconsts"
)

// AccessLog 访问日志中间件
// 在请求处理完成后记录请求信息（方法、路径、用户、IP、状态码、耗时、错误），异步推送至队列落库
func (s *sMiddleware) AccessLog(r *ghttp.Request) {
	start := gtime.Now()
	r.Middleware.Next()

	var (
		user     = zcontext.GetUser(r.Context())
		memberId int64
		username string
	)
	if user != nil {
		memberId = user.Id
		username = user.Username
	}

	errMsg := ""
	if err := r.GetError(); err != nil {
		errMsg = err.Error()
	}

	log := &entity.SysRequestLog{
		ReqId:     gctx.CtxId(r.Context()),
		Method:    r.Method,
		Path:      r.URL.Path,
		MemberId:  memberId,
		Username:  username,
		Ip:        r.GetClientIp(),
		UserAgent: r.UserAgent(),
		Status:    r.Response.Status,
		ErrorMsg:  errMsg,
		CostTime:  gtime.Now().Sub(start).Milliseconds(),
		CreatedAt: gtime.Now(),
	}

	if err := zqueue.Push(zconsts.QueueLogTopic, log); err != nil {
		g.Log().Warningf(r.Context(), "push request log failed: %v", err)
	}
}
