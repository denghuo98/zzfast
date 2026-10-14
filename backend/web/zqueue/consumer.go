package zqueue

import (
	"context"
	"sync"
)

// Consumer 消费者接口，实现该接口即可加入到消费队列中
type Consumer interface {
	GetTopic() string                                    // 获取消费主题
	Handle(ctx context.Context, mqMsg MqMsg) (err error) // 处理消息的方法
}

// consumerManager 消费者管理
type consumerManager struct {
	sync.Mutex
	list map[string]Consumer // 维护的消费者列表
}

var consumers = &consumerManager{
	list: make(map[string]Consumer),
}

// RegisterConsumer 注册任务到消费者队列
func RegisterConsumer(cs Consumer) {
	consumers.Lock()
	defer consumers.Unlock()
	topic := cs.GetTopic()
	if _, ok := consumers.list[topic]; ok {
		Logger().Warningf(ctx, "queue.RegisterConsumer topic:%v duplicate registration.", topic)
		return
	}
	Logger().Infof(ctx, "queue.RegisterConsumer topic:%v register success.", topic)
	consumers.list[topic] = cs
}

// StartConsumersListener 启动所有已注册的消费者监听
func StartConsumersListener(ctx context.Context) {
	for _, c := range consumers.list {
		go func(c Consumer) {
			consumerListen(ctx, c)
		}(c)
	}
}

// consumerListen 消费者监听
func consumerListen(ctx context.Context, job Consumer) {
	var (
		topic  = job.GetTopic()
		c, err = InstanceConsumer()
	)

	if err != nil {
		Logger().Fatalf(ctx, "InstanceConsumer %s err:%+v", topic, err)
		return
	}

	if listenErr := c.ListenReceiveMsgDo(topic, func(mqMsg MqMsg) {
		err = job.Handle(ctx, mqMsg)

		if err != nil {
			// 消费失败时记录错误日志，包含消息ID和topic以便追踪
			// 注意：当前磁盘队列暂未实现死信/重试机制，失败消息会被丢弃
			Logger().Errorf(ctx, "消息消费失败,topic:%s,msgId:%s,err:%+v", topic, mqMsg.MsgId, err)
		}

		// 记录消费队列日志
		ConsumerLog(ctx, topic, mqMsg, err)
	}); listenErr != nil {
		Logger().Fatalf(ctx, "消费队列：%s 监听失败, err:%+v", topic, listenErr)
	}
}
