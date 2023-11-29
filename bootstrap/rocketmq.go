package bootstrap

import (
	"context"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"github.com/gin-gonic/gin"
	"github.com/succko/hera/global"
	"go.uber.org/zap"
)

type mq struct{}

var Mq = new(mq)

func InitializeRocketMqConsumers(c map[string]func(message []byte)) []rocketmq.PushConsumer {
	var consumers []rocketmq.PushConsumer
	for k, v := range c {
		consumers = append(consumers, initializeRocketMqConsumer(k, v))
	}
	return consumers
}

func initializeRocketMqConsumer(topic string, f func(message []byte)) rocketmq.PushConsumer {
	return initializeRocketMqConsumerWithTag(topic, "", f)
}

func initializeRocketMqConsumerWithTag(topic string, tag string, f func(message []byte)) rocketmq.PushConsumer {
	c, err := rocketmq.NewPushConsumer(
		consumer.WithGroupName(topic+"-"+global.App.Config.App.AppName+"-"+gin.Mode()),
		consumer.WithNsResolver(primitive.NewPassthroughResolver([]string{global.App.Config.Rokcetmq.Addr})),
		consumer.WithConsumerModel(consumer.Clustering),
		consumer.WithConsumeFromWhere(consumer.ConsumeFromLastOffset),
	)

	if err != nil {
		zap.L().Error("Failed to create consumer", zap.Error(err))
		return nil
	}

	var selector consumer.MessageSelector
	if tag != "" {
		selector = consumer.MessageSelector{
			Type:       consumer.TAG,
			Expression: tag,
		}
	} else {
		selector = consumer.MessageSelector{}
	}

	// 订阅的消费者列表
	subscribe(c, topic, selector, f)

	if err = c.Start(); err != nil {
		zap.L().Error("Failed to start consumer", zap.Error(err))
		return nil
	}

	return c
}

func subscribe(c rocketmq.PushConsumer, topic string, selector consumer.MessageSelector, f func(message []byte)) {
	// 订阅消息
	err := c.Subscribe(topic, selector, func(ctx context.Context,
		msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
		for _, msg := range msgs {
			zap.L().Info(fmt.Sprintf("Consumer subscribe callback, topic: %s, msg: %s", topic, msg))
			go f(msg.Body)
			// retry TODO
		}
		return consumer.ConsumeSuccess, nil
	})
	if err != nil {
		zap.L().Error("Failed to subscribe", zap.Error(err))
	}
}

func InitializeRocketMqProducer() rocketmq.Producer {
	p, err := rocketmq.NewProducer(
		producer.WithNsResolver(primitive.NewPassthroughResolver([]string{global.App.Config.Rokcetmq.Addr})),
		producer.WithRetry(16),
		producer.WithGroupName(global.App.Config.App.AppName),
	)

	if err != nil {
		zap.L().Error("Failed to create producer", zap.Error(err))
		return nil
	}

	err = p.Start()
	if err != nil {
		zap.L().Error("Failed to start producer", zap.Error(err))
		return nil
	}

	return p
}
