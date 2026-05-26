Created go module using:

```
go mod init golang_linter_agent
```

Get ADK package to use with  `golang_linter_agent`

```
go get google.golang.org/adk
```

Process to create a runner:

```

// Create runner
runner, err := runner.New(runner.Config{
    AppName:        appName,
    Agent:          timeAgent, // TODO: add the root agent here
    SessionService: service,
})

fmt.Println("####### Starting golang_linter_agent #######")
session, events, err := runner.RunLive(ctx, resp.Session.UserID(), resp.Session.ID(), agent.LiveRunConfig{
    ResponseModalities: []genai.Modality{"TEXT"},
})
if err != nil {
    log.Fatalf("Failed to run agent: %v", err)
}
defer session.Close()

fmt.Printf("Event for session: %+v\n", events)

```