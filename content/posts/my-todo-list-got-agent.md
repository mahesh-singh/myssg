+++
title = 'My TODO List Finally Got Agent'
date = "2025-9-9"
tags = ['llm', 'agent', 'go']
slug = 'my-todo-list-got-agent'
draft = false
+++

# My TODO List Finally Got Agent (Hello World Agent)

Remember when your biggest productivity problem was remembering where you put that sticky note with "do the code review" written on it? Well, I've accidentally created something that makes that look like child's play. I built a todo list that talks back.

I was hearing all about agents, I was thinking what could be a smallest agent hands-on I could built. I partially manage my running todo list and notes in a text file. Sometime list get too long and take long time to look for something specific or sometime I return after couple of days back.

I was thinking what if I can have a conversation with it with the help of LLM. It took me about 300 lines of code to make this happen.


## The Magic 

Here's the weird thing about AI agents - they're basically just really smart middlemen. Think of it like this:

- You say something in human language
- The AI thinks "hmm, they probably want me to read/write a file"
- code calls a function to do that
- It tells you what happened in a nice way

It's like having a translator who speaks both Human and Computer, except the translator is suspiciously good at understanding what you actually meant, not just what you said.

## The Two Tools (function) That Rule Them All

My todo agent only knows how to do two things:

- Read todo files (because how else would it know what you wrote?)
- Edit todo files (because that's literally the point)

That's it. Two functions. LLM figures out which tool (function) to use and when. Ask it "What's on my list?" and it knows to read the file. Say "Add wash dishes to my todos" and it knows to edit the file.

## The Code

1. Need an anthropic key, go and get the key. I used $5
2. First have a conversation 
3. Read the file 
4. Edit the file 
5. Create the file 

Lets start 

We are going to use Golang, but you can use what you know, concept will be the same. We need one addition thing (Anthropic API key)[https://console.anthropic.com/settings/keys]. Once you get the API Key set on the env variable `ANTHROPIC_API_KEY`

```
export ANTHROPIC_API_KEY="this is the last time i will tell you to set this"
```

Every time you restart your terminal or machine, you need to set it again. This shouldn't be committed in git repo. It will cerate risk of exposing your API keys to others.


Now setup a new Go project

```
mkdir todo-agent
cd todo-agent
go mod init agent
touch main.go
```

Here is high level structure of tasks we would like to achieve 

- Establish back and fourth communication with Claude via integrating the API in our code
- Setting up first tool for reading the todo.txt file 
- If todo.txt file doesn't exists, create the file
- Modify the todo.txt file based on users request

Now, let’s open main.go and, as a first step

```
package main

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"github.com/anthropics/anthropic-sdk-go"
)

func main() {
	client := anthropic.NewClient()

	scanner := bufio.NewScanner(os.Stdin)
	getUserMessage := func() (string, bool) {
		if !scanner.Scan() {
			return "", false
		}
		return scanner.Text(), true
	}

	agent := NewAgent(&client, getUserMessage)
	err := agent.Run(context.TODO())
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
	}
}

func NewAgent(client *anthropic.Client, getUserMessage func() (string, bool)) *Agent {
	return &Agent{
		client:         client,
		getUserMessage: getUserMessage,
	}
}

type Agent struct {
	client         *anthropic.Client
	getUserMessage func() (string, bool)
}
```

Now let’s add the missing Run() method:

```
// main.go

func (a *Agent) Run(ctx context.Context) error {
	conversation := []anthropic.MessageParam{}

	fmt.Println("🤖 Todo Agent Ready! (Press Ctrl+C to quit)")
    fmt.Println("Try: ''what's on my todo list?'")

	for {
		fmt.Print("You: ")
		userInput, ok := a.getUserMessage()
		if !ok {
			break
		}

		userMessage := anthropic.NewUserMessage(anthropic.NewTextBlock(userInput))
		conversation = append(conversation, userMessage)

		message, err := a.askClaude(ctx, conversation)
		if err != nil {
			return err
		}
		conversation = append(conversation, message.ToParam())

		for _, content := range message.Content {
			switch content.Type {
			case "text":
				fmt.Printf("Claude: %s\n", content.Text)
			}
		}
	}

	return nil
}

func (a *Agent) askClaude(ctx context.Context, conversation []anthropic.MessageParam) (*anthropic.Message, error) {
	message, err := a.client.Messages.New(ctx, anthropic.MessageNewParams{
		Model:     anthropic.ModelClaude3_7SonnetLatest,
		MaxTokens: int64(1024),
		Messages:  conversation,
	})
	return message, err
}

```

Above code every AI chat application you’ve ever used, except it’s written by you and you are using it in the terminal.

Now time to run it!

```
export ANTHROPIC_API_KEY="this is the last time i will tell you to set this"

# Download the dependencies
go mod tidy
# Run it
go run main.go
```


```
🤖 Todo Agent Ready! (Press Ctrl+C to quit)
Try: what's on my todo list?
You: Hey I am Mahesh, How are you?
Claude: Hello Mahesh! I'm doing well, thank you for asking. How are you today? Is there something specific related to todo files that you'd like help with? I can help you read or edit todo files if you need assistance with that.
```

https://ampcode.com/how-to-build-an-agent