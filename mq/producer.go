package mq

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/succko/hera/global"
	"go.uber.org/zap"
)

type producer struct {
}

var Producer = new(producer)

func (p *producer) SendSync(topic string, body interface{}) {
	p.SendSyncWithTag(topic, body, "")
}

func (p *producer) SendSyncWithTag(topic string, body interface{}, tag string) {
	data, _ := json.Marshal(body)
	//实例化消息
	msg := &primitive.Message{
		Topic: topic,
		Body:  data,
	}
	msg.WithDelayTimeLevel(1)
	if tag != "" {
		msg.WithTag(tag)
	}
	msg.WithKeys([]string{"DEFAULT"})
	//msg.WithDelayTimeLevel(1)
	//同步发送
	res, err := global.App.RocketMqProducer.SendSync(context.Background(), msg)
	if err != nil {
		zap.L().Error(fmt.Sprintf("send message error: %s", err))
	} else {
		zap.L().Info(fmt.Sprintf("send message success: result=%s", res.String()))
	}
}
