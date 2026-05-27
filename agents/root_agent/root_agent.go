// Package root_agent implements a root agent for the golang_linter.
package root_agent

import (
	"log"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/agent/workflowagents/loopagent"
	"google.golang.org/adk/model"
	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/geminitool"
)

func NewAgent(model model.LLM) agent.Agent {
	// Create a placeholder time agent
	timeAgent, err := llmagent.New(llmagent.Config{
		Name:        "hello_time_agent",
		Model:       model,
		Description: "Tells the current time in a specified city.",
		Instruction: "You are a helpful assistant that tells the current time in a city.",
		Tools: []tool.Tool{
			geminitool.GoogleSearch{},
		},
	})
	if err != nil {
		log.Fatalf("Failed to create agent: %v", err)
	}

	cfg := loopagent.Config{
		MaxIterations: 3, // Loop indefinitely
		AgentConfig: agent.Config{
			Name:        "linter_root_agent",
			Description: "Root agent for golang linter. Delegates linting tasks to sub-agents and compiles results.",
			SubAgents:   []agent.Agent{timeAgent},
			// Absense of custom run forces loop agent to auto delegate to sub-agents.
			// Run:         run,
		},
	}

	agent, err := loopagent.New(cfg)
	if err != nil {
		panic("Failed to create loop agent: " + err.Error())
	}

	return agent
}
