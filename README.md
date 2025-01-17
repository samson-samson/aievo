# AIEvo: a **multi-agent** framework open-sourced by **Ant Group**

---

## Introduction

AIEVO is a **multi-agent** framework open-sourced by **Ant Group**, designed to efficiently create multi-agent applications.

## Core Advantages

- **High SOP Compliance**: Strictly adheres to the SOP Graph of tasks to execute complex tasks.
- **High Execution Success Rate**: Improves the success rate of complex tasks through feedback and watcher mechanisms.
- **High Flexibility**: Supports dynamic task creation and subscription settings based on task attributes.
- **Extensibility**: Provides numerous extension points for custom enhancements to the framework.
- **Enterprise-Proven**: Validated in production within Ant Group, ensuring the framework's stability and reliability.

## Usage Examples

### Multi-Agent Examples

- [Debate Competition](/examples/multi-agent-example/battle): Multiple agents engage in a debate competition.
- [Werewolf Game](/examples/multi-agent-example/werewolf): Multiple agents play a game of Werewolf.
- [Undercover Game](/examples/multi-agent-example/undercover): Multiple agents play a game of Undercover.
- [Paper Writing](/examples/multi-agent-example/paper_write): Multiple agents collaborate to write a paper.

### Single-Agent Examples

- [Engineer](/examples/single-agent-example/engineer): A single agent that can complete code writing tasks.

## Architecture Overview

![](docs/static/arch.jpg)

## Module Introduction

### Agent Module

The primary function of this module is to facilitate the construction of Agents. We employ the ReAct approach for building Agents, and compared to using LangChain for Agent construction, we support autonomous interaction with other Agents.

For example, to create a programmer Agent:
```go
// 实例化基座模型
client, err := openai.New(
    openai.WithToken(os.Getenv("OPENAI_API_KEY")),
    openai.WithModel(os.Getenv("OPENAI_MODEL")),
    openai.WithBaseURL(os.Getenv("OPENAI_BASE_URL")))
if err != nil {
    log.Fatal(err)
}// 文件操作相关工具
fileTools, _ := file.GetFileRelatedTools(workspace)
// 命令执行
bashTool, _ := bash.New()
// 构建工具集
engineerTools := make([]tool.Tool, 0)
engineerTools = append(engineerTools, fileTools...)
engineerTools = append(engineerTools, bashTool)
// 回调处理器，主要用于分析Agent的执行过程
callbackHandler := &CallbackHandler{}
// 定义Env，用于和其他Agent进行交互以及记忆存储
env := environment.NewEnv()
// 构建Agent
engineer, _ := agent.NewBaseAgent( 
    // Agent的名称 
    agent.WithName("engineer"),
    // Agent的描述 
    agent.WithDesc(EngineerDescription),
    // Agent的Prompt（支持设置动态参数）
    agent.WithPrompt(EngineerPrompt),
    // Agent的指令（支持设置动态参数）
    agent.WithInstruction(SingleAgentInstructions),
    // Agent的动态参数
	// 1. 当前工作流程
    agent.WithVars("sop", Workflow),
	// 2. 当前工作区
    agent.WithVars("workspace", workspace),
    // Agent的工具集
    agent.WithTools(engineerTools),
    // Agent的基座模型
    agent.WithLLM(client),
    // Agent的回调处理器
    agent.WithCallback(callbackHandler),
	// Agent的环境
    aievo.WithEnvironment(env),
)

// 运行Agent
gen, _ := engineer.Run(
    context.Background(), 
	// 用户消息
    []schema.Message
    {
        {
            Type:     schema.MsgTypeMsg,
            Content:  "写一个终端版本的贪吃蛇",
            Sender:   "User",
            Receiver: "engineer",
        },
    }, 
	// 基座模型参数
    llm.WithTemperature(0.1)
)

// 打印Agent的回复
fmt.Println(gen.Messages[0].Content)
```

### Env Module

In multi-agent systems, this module is primarily used to store information such as team members, subscription relationships, historical messages, and the SopGraph of tasks. It also serves as an intermediary pool for interactions between Agents. Agents send messages to the Env, and the driver module distributes them to the corresponding Agents.

Below is the interface definition of the Env module:
```go
type Environment interface {
	// Produce 生产消息
	Produce(ctx context.Context, msgs ...Message) error
	// Consume 消费消息
	Consume(ctx context.Context) *Message
	// SOP 获取任务的SopGraph
	SOP() string
	// GetTeam 获取团队里面的所有Agent
	GetTeam() []Agent
	// GetTeamLeader 获取团队Leader
	GetTeamLeader() Agent
	// LoadMemory 获取Agent的历史消息
	LoadMemory(ctx context.Context, receiver Agent) []Message
	// GetSubscribeAgents 获取某个Agent的订阅者
	GetSubscribeAgents(_ context.Context, subscribed Agent) []Agent
}
```

In addition to storing the basic information mentioned above, the Env module can also store some control information, such as token (maximum token consumption for a task), max_turns (maximum number of Agent runs), and so on.

Since we mentioned that Env is the message relay pool, let's delve into the details of Message in AIEvo.

Below is the structure definition of Message:
```go
type Message struct {
	// 消息类型
	Type      string       `json:"cate"`
	// 产出该消息的思考过程
	Thought   string       `json:"thought"`
	// 消息内容
	Content   string       `json:"content"`
	// 发送方
	Sender    string       `json:"sender"`
	// 接收方
	Receiver  string       `json:"receiver"`
	// 接受条件
	Condition string       `json:"condition"`
	// 生成该消息过程中的工具调用记录
	Steps     []StepAction `json:"steps"`
	// 相关日志存储
	Log       string
	// 控制信息，用于剔除和更新Agent
    MngInfo   *MngInfo
	// 所有的可以接收该消息的Agent
	AllReceiver []string
}
```

Currently supported message types include:
1. MsgTypeMsg: Regular interaction messages
2. MsgTypeEnd: Session end messages
3. MsgTypeCreative: Messages generated when CreativeAgent creates a Team
4. MsgTypeSOP: Messages generated when SopAgent creates a SopGraph based on tasks

Different message types have different processing strategies when delivered to the Env:
- MsgTypeMsg -> msgStrategy: Store the message in Memory
- MsgTypeCreative -> creativeStrategy: Modify the Team (used to remove and update agents)
- MsgTypeSOP -> sopStrategy: Store the SopGraph in the Env

Currently supported team modes include:
1. DefaultSubMode: Default mode. If a LeaderAgent exists, LeaderSubMode is used; otherwise, ALLSubMode is used.
2. LeaderSubMode: All agents subscribe only to the LeaderAgent, while the LeaderAgent subscribes to all agents. The LeaderAgent drives the execution of the entire task.
3. ALLSubMode: All agents subscribe to each other, and everyone collectively drives the execution of the task.
4. CustomSubMode: Custom subscription relationships, where the user specifies the subscription relationships.

> How to choose a team mode?
> - LeaderSubMode: Suitable for scenarios where the task Sop is complex and requires a high success rate.
> - ALLSubMode: Suitable for scenarios where the task Sop is relatively simple and you want to fully leverage the autonomy of the agents.
> - CustomSubMode: Suitable for scenarios where the subscription relationships between agents are well-defined and the Sop is relatively simple.

### Feedback Module

This module is used to review and provide feedback on the content generated by the Agent.

Before the introduction of this module, the following issues might have occurred:
1. Content was not generated in a fixed format.
2. Due to hallucinations in LLM, the Agent might repeatedly call a certain tool.
3. Task progression messages did not follow the SopGraph.
4. Generated responses contained sensitive information.

After introducing Feedback, these issues can be resolved. For example, in Feedback, we can validate the format of the output message. If it meets the requirements, suggestions can be provided, and a retry can be initiated.

Of course, the scenarios where Feedback is needed are far more extensive. Users can customize the Feedback as required.

Below is the interface definition of the Feedback module:
```go
type Feedback interface {
	Feedback(ctx context.Context, agent schema.Agent, messages []schema.Message, actions []schema.StepAction,
		steps []schema.StepAction, prompt string) *FeedbackInfo
}
```
Below is the structure definition of FeedbackInfo:
```go
type FeedbackInfo struct {
	// 反馈类型: 通过/不通过
	Type  FeedbackType `json:"type"`
	// 反馈建议
	Msg   string       `json:"msg"`
	// 消耗的Token
	Token int          `json:"token"`
}
```

Multiple Feedbacks can form a FeedbackChain as follows:
```go
// feedbackChain Define a struct contains slice of Feedback
type feedbackChain struct {
	chains []Feedback
}

// Feedback Implement the Feedback method for the feedbackChain
func (fc *feedbackChain) Feedback(ctx context.Context, agent schema.Agent, messages []schema.Message, actions []schema.StepAction,
	steps []schema.StepAction, prompt string) *FeedbackInfo {

	info := &FeedbackInfo{
		Type: Approved,
	}

	for _, feedback := range fc.chains {
		if feedback == nil {
			continue
		}
		info = feedback.Feedback(ctx, agent, messages, actions, steps, prompt)
		if info.Type == NotApproved {
			return info
		}
	}

	return info
}

// Chain function to create a new Feedback that chains multiple Feedback
func Chain(chains ...Feedback) Feedback {
	return &feedbackChain{chains: chains}
}
```

### Watcher Module

This module is used to monitor the operation of the entire multi-agent system and intervene in the process when appropriate.

For example: In a Werewolf game scenario, if an Agent is killed, the watcher will remove that Agent from the Team to prevent it from receiving further messages, which could disrupt the entire operation.

The operation process of the watcher is as follows:

1. Start a watcher to observe all execution processes and generate Team change messages for removing and updating agents.
    ```go
    func (e *AIEvo) Watch(ctx context.Context, _ string, opts ...llm.GenerateOption) (string, error) {
        if e.Watcher != nil {
            e.WatchChan = make(chan schema.Message)
            e.WatchChanDone = make(chan struct{})
            go func() {
                for message := range e.WatchChan {
                    if e.WatchCondition != nil && e.WatchCondition(message) == false {
                        e.WatchChanDone <- struct{}{}
                        continue
                    }
                    generation, err := e.Watcher.Run(ctx, e.LoadMemory(ctx, e.Watcher))
                    if err != nil {
                        e.WatchChanDone <- struct{}{}
                        continue
                    }
                    e.Produce(ctx, generation.Messages...)
                    e.WatchChanDone <- struct{}{}
                }
            }()
        }
        return "", nil
    }
    ```

2. apply change messages
   ```go
   func (e *Environment) mngInfoStrategy(ctx context.Context, msg *schema.Message) error {
       if msg.MngInfo == nil {
           return nil
       }
       // 当前仅支持移除
       if msg.MngInfo.Remove != nil {
           e.Team.RemoveMembers(msg.MngInfo.Remove)
       }
       _ = e.Memory.Save(ctx, *msg)
       return nil
   }
   ```


### Driver Module

This module is used to drive the operation of the entire multi-agent system. The entire driving logic is message-driven.

The operation steps are as follows:
```go
e.Handler = Chain(e.BuildPlan, e.BuildSOP, e.Watch, e.Scheduler)
```
1. Build Team (can be manually specified): Construct the team based on the current task attributes.
2. Build Sop (can be manually specified): Construct the Sop based on the task document.
3. Start Watcher: Launch the Watcher to monitor the entire operation process.
4. Start Scheduling: Schedule each Agent based on the messages in the Env.

The scheduling logic is as follows:

```go
func (e *AIEvo) Scheduler(ctx context.Context, prompt string, opts ...llm.GenerateOption) (string, error) {
	// 首先由用户消息进行驱动
	e.Produce(ctx, schema.Message{
		Type:     schema.MsgTypeMsg,
		Content:  prompt,
		Sender:   _defaultSender,
		Receiver: e.Leader.Name(),
	})
	// 循环从Env中消费消息，直到消息为空
	for msg := e.Consume(ctx); msg != nil; msg = e.Consume(ctx) {
		if msg.IsEnd() {
			return msg.Content, nil
		}
		// 获取消息的接受者
		receivers := msg.Receivers()
		for _, rec := range receivers {
			receiver := e.Agent(rec)
			if receiver == nil {
				if len(receivers) == 1 {
					return msg.Content, fmt.Errorf(
						"get unexcept agent %s", msg.Receiver)
				}
				continue
			}
			messages := e.LoadMemory(ctx, receiver)
			if e.Callback != nil {
				e.Callback.HandleAgentStart(ctx, receiver, messages)
			}
			// 调度对应的Agent运行
			gen, err := receiver.Run(ctx, messages, opts...)
			if err != nil {
				return "", err
			}
			if e.Callback != nil {
				e.Callback.HandleAgentEnd(ctx, receiver, gen)
			}

			if gen.Messages == nil {
				return "", fmt.Errorf("gen messages is nil for agent %s", msg.Receiver)
			}

			// 将消息投递到Env中
			e.Produce(ctx, gen.Messages...)
			e.broadcast(gen.Messages...)
		}
	}
	return "", nil
}
```

## Communication
<table>
  <tr>
    <td>
      DingTalk Group:<br>
      <img src="docs/static/qr.png" alt="DingTalk QR Code" style="width: 150px; height: 150px;">
    </td>
    <td>
      WeChat Official Account:<br>
      <img src="docs/static/account.jpg" alt="WeChat QR Code" style="width: 150px; height: 150px;">
    </td>
  </tr>
</table>


## License

AIEvo is licensed under the Apache 2.0 License. For more details, please read [LICENSE]()