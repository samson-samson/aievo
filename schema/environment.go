package schema

import (
	"context"
)

type Environment interface {
	// Produce 生产消息
	Produce(ctx context.Context, msgs ...Message) error
	// Consume 消费消息
	Consume(ctx context.Context) *Message

	// SOP 任务SOP
	SOP() string
	// GetTeam 团队所有成员
	GetTeam() []Agent
	// GetTeamLeader 团队Leader
	GetTeamLeader() Agent

	// LoadMemory 获取Agent的历史消息
	LoadMemory(ctx context.Context, receiver Agent) []Message

	GetSubscribeAgents(_ context.Context, subscribed Agent) []Agent
}

// Memory is the interface for memory in chains.
type Memory interface {
	Load(ctx context.Context, filter func(index, consumption int, message Message) bool) []Message

	// LoadNext 加载下一条消息，filter 对下一条消息进行检查，不符合的话，则不会返回
	LoadNext(ctx context.Context, filter func(message Message) bool) *Message

	Save(ctx context.Context, msg Message) error
	// Clear memory contents.
	Clear(ctx context.Context) error
}

const (
	MsgTypeMsg      = "MSG"
	MsgTypeCreative = "CREATIVE"
	MsgTypeSOP      = "SOP"
	MsgTypeEnd      = "END"
)

const (
	MsgAllReceiver = "ALL"
)

type MngInfo struct {
	Create []struct {
		Name        string   `json:"name"`
		Description string   `json:"description"`
		Tools       []string `json:"tools"`
		Prompt      string   `json:"prompt"`
	} `json:"create"`
	Select []string `json:"select"`
	Remove []string `json:"remove"`
}

type Subscribe struct {
	Subscribed Agent
	Subscriber Agent
	Condition  string
}
