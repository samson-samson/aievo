package aievo

import (
	"context"

	"github.com/antgroup/aievo/environment"
	"github.com/antgroup/aievo/llm"
	"github.com/antgroup/aievo/memory"
)

var (
	_defaultSender = "User"
)

func NewAIEvo(opts ...Option) (*AIEvo, error) {
	e := &AIEvo{}
	o := &options{
		maxTurn:  20,
		maxToken: 20 * 4000,
		subMode:  environment.DefaultSubMode,
		env:      environment.NewEnv(),
	}
	for _, opt := range opts {
		opt(o)
	}

	e.Environment = o.env
	if e.Environment.Team == nil {
		e.Environment.Team = environment.NewTeam()
	}
	e.Team.SubMode = o.subMode
	if e.Environment.Memory == nil {
		e.Environment.Memory = memory.NewBufferMemory()
	}
	e.Callback = o.callback
	e.MaxToken = o.maxToken
	e.MaxTurn = o.maxTurn
	e.Team.AddMembers(o.team...)
	if o.user != nil {
		e.Team.AddMembers(o.user)
	}
	e.Team.Subscribes = o.subscribes
	e.Team.Leader = o.leader
	e.Sop = o.sop
	e.SopExpert = o.sopExpert
	e.Planner = o.planner
	e.Watcher = o.watcher
	e.WatchCondition = o.watchCondition
	e.Handler = Chain(e.BuildPlan, e.BuildSOP, e.Watch, e.Scheduler)

	// 填充各个agent与环境交互的方式
	for _, agent := range e.GetTeam() {
		agent.WithEnv(e.Environment)
	}

	if e.GetTeamLeader() == nil {
		return nil, ErrMissingLeader
	}
	if len(e.GetTeam()) == 0 {
		return nil, ErrMissTeam
	}
	return e, nil
}

func (e *AIEvo) Run(ctx context.Context, prompt string, opts ...llm.GenerateOption) (string, error) {
	return e.Handler(ctx, prompt, opts...)
}
