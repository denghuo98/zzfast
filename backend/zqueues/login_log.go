package zqueues

import (
	"context"
	"encoding/json"

	"github.com/denghuo98/zzfast/internal/model/entity"
	"github.com/denghuo98/zzfast/web/zqueue"
	"github.com/denghuo98/zzfast/zconsts"
	"github.com/denghuo98/zzfast/zservice"
)

func init() {
	zqueue.RegisterConsumer(LoginLog)
}

// LoginLog 登录日志
var LoginLog = &qLoginLog{}

type qLoginLog struct{}

// GetTopic 主题
func (q *qLoginLog) GetTopic() string {
	return zconsts.QueueLoginLogTopic
}

// Handle 处理消息
func (q *qLoginLog) Handle(ctx context.Context, mqMsg zqueue.MqMsg) (err error) {
	var data entity.SysLoginLog
	if err = json.Unmarshal(mqMsg.Body, &data); err != nil {
		return err
	}
	return zservice.SysLoginLog().RealWrite(ctx, data)
}
