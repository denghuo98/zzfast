package zqueues

import (
	"context"
	"encoding/json"

	"github.com/denghuo98/zzfast/internal/dao"
	"github.com/denghuo98/zzfast/internal/model/entity"
	"github.com/denghuo98/zzfast/web/zqueue"
	"github.com/denghuo98/zzfast/zconsts"
)

func init() {
	zqueue.RegisterConsumer(RequestLog)
}

// RequestLog 请求访问日志
var RequestLog = &qRequestLog{}

type qRequestLog struct{}

// GetTopic 主题
func (q *qRequestLog) GetTopic() string {
	return zconsts.QueueLogTopic
}

// Handle 处理消息
func (q *qRequestLog) Handle(ctx context.Context, mqMsg zqueue.MqMsg) (err error) {
	var data entity.SysRequestLog
	if err = json.Unmarshal(mqMsg.Body, &data); err != nil {
		return err
	}
	_, err = dao.SysRequestLog.Ctx(ctx).Data(data).OmitEmpty().Insert()
	return err
}
