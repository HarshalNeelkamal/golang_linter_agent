package main

import (
	"context"
	"log"
	"os"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/cmd/launcher"
	"google.golang.org/adk/cmd/launcher/full"
	"google.golang.org/adk/model/gemini"
	"google.golang.org/adk/session"
	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/geminitool"
	"google.golang.org/genai"
)

const (
	appName   = "golang_linter_agent"
	userID    = "default-user"
	sessionID = "default-session"
)

var (
	initialState map[string]any
)

func init() {
	// Set up the initial state for the session
	initialState = map[string]any{
		"initial_key": "initial_value",
	}
}

func main() {
	ctx := context.Background()

	// Initialize the Gemini model
	clientConfig := &genai.ClientConfig{
		APIKey: os.Getenv("GOOGLE_AI_API_KEY"),
	}
	model, err := gemini.NewModel(ctx, "gemini-2.5-flash", clientConfig)
	if err != nil {
		log.Fatalf("Failed to initialize Gemini model: %v", err)
	}

	// Create agent
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

	// create session with initial state
	service := session.InMemoryService()
	_, err = service.Create(ctx, &session.CreateRequest{
		SessionID: sessionID,
		UserID:    userID,
		AppName:   appName,
		State:     initialState,
	})
	if err != nil {
		log.Fatalf("Failed to create session: %v", err)
	}

	// Launch the agent using the full launcher
	config := &launcher.Config{
		SessionService: service,
		AgentLoader:    agent.NewSingleLoader(timeAgent),
	}
	l := full.NewLauncher()
	if err = l.Execute(ctx, config, os.Args[1:]); err != nil {
		log.Fatalf("Run failed: %v\n\n%s", err, l.CommandLineSyntax())
	}
}
