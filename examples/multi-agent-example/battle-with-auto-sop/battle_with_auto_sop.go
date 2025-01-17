package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/antgroup/aievo/agent"
	"github.com/antgroup/aievo/aievo"
	"github.com/antgroup/aievo/callback"
	"github.com/antgroup/aievo/environment"
	"github.com/antgroup/aievo/feedback"
	"github.com/antgroup/aievo/llm"
	"github.com/antgroup/aievo/llm/openai"
	"github.com/antgroup/aievo/memory"
	"github.com/antgroup/aievo/schema"
)

func main() {
	// init llm
	client, err := openai.New(
		openai.WithToken(os.Getenv("OPENAI_API_KEY")),
		openai.WithModel(os.Getenv("OPENAI_MODEL")),
		openai.WithBaseURL(os.Getenv("OPENAI_BASE_URL")))
	if err != nil {
		log.Fatal(err)
	}
	// init callback
	cb := &callback.LogHandler{}

	// init environment && mem
	mem := memory.NewBufferWindowMemory(25)
	env := environment.NewEnv()
	env.Memory = mem

	// init team and leader
	team, leader := getTeam(client, cb, env)

	// init sop agent
	sopFeedback, _ := feedback.NewSopFeedback(client)
	sop, _ := agent.NewSopAgent(agent.WithLLM(client),
		agent.WithEnv(env),
		agent.WithCallback(cb),
		agent.WithFeedbacks(sopFeedback))

	// init aievo running options
	opts := make([]aievo.Option, 0, 5)
	opts = append(opts,
		aievo.WithTeam(team),
		aievo.WithMaxTurn(50),
		aievo.WithCallback(cb),
		aievo.WithSubScribeMode(environment.ALLSubMode),
		aievo.WithLLM(client),
		aievo.WithEnvironment(env),
		aievo.WithTeamLeader(leader),
		aievo.WithUserProxy(nil),
		aievo.WithSopExpert(sop),
	)

	// start aievo
	evo, err := aievo.NewAIEvo(opts...)
	if err != nil {
		panic(err)
	}
	run, err := evo.Run(context.Background(), "Topic: climate change. Under 80 words per message",
		llm.WithTemperature(0.5),
	)
	if err != nil {
		panic(err)
	}
	fmt.Println(run)
}

// getTeam return team all members and leader member
func getTeam(client llm.LLM, cb callback.Handler,
	env schema.Environment) ([]schema.Agent, schema.Agent) {
	// init host agent
	host, _ := agent.NewBaseAgent(
		agent.WithPrompt(_defaultHostPrompt),
		agent.WithInstruction(_defaultHostInstruction),
		agent.WithSuffix(_defaultHostSuffix),
		agent.WithName("Host"),
		agent.WithDesc(`The Host of debate`),
		agent.WithLLM(client),
		agent.WithCallback(cb),
		agent.WithEnv(env))

	// init agent
	Alice, _ := agent.NewBaseAgent(
		agent.WithPrompt(_defaultAffirmativePrompt),
		agent.WithInstruction(_defaultPlayerInstruction),
		agent.WithSuffix(_defaultPlayerSuffix),
		agent.WithName("Alice"),
		agent.WithDesc(`The Affirmative Side player of debate`),
		agent.WithLLM(client),
		agent.WithCallback(cb),
		agent.WithEnv(env))

	// init agent
	Bob, _ := agent.NewBaseAgent(
		agent.WithPrompt(_defaultNegativePrompt),
		agent.WithInstruction(_defaultPlayerInstruction),
		agent.WithSuffix(_defaultPlayerSuffix),
		agent.WithName("Bob"),
		agent.WithDesc(`The Negative Side player of debate`),
		agent.WithLLM(client),
		agent.WithCallback(cb),
		agent.WithEnv(env))
	team := make([]schema.Agent, 0, 6)
	team = append(team, host, Alice, Bob)
	experts := []string{"Expert1", "Expert2", "Expert3"}
	for _, expert := range experts {
		tmp, _ := agent.NewBaseAgent(
			agent.WithPrompt(_defaultExpertPrompt),
			agent.WithInstruction(_defaultExpertInstruction),
			agent.WithSuffix(_defaultExpertSuffix),
			agent.WithName(expert),
			agent.WithDesc(`The Expert of debate`),
			agent.WithLLM(client),
			agent.WithCallback(cb),
			agent.WithEnv(env))
		team = append(team, tmp)
	}
	return team, host
}
