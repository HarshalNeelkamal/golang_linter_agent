package main

import (
	"context"
	"log"
	"os"

	"golang_linter_agent/agents/root_agent"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/cmd/launcher"
	"google.golang.org/adk/cmd/launcher/full"
	"google.golang.org/adk/model/gemini"
	"google.golang.org/adk/session"
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

	// Create loop agent with the Gemini model
	rootAgent := root_agent.NewAgent(model)

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
		AgentLoader:    agent.NewSingleLoader(rootAgent),
	}
	l := full.NewLauncher()
	if err = l.Execute(ctx, config, os.Args[1:]); err != nil {
		log.Fatalf("Run failed: %v\n\n%s", err, l.CommandLineSyntax())
	}
}
