+++
title = 'My TODO List Finally Got Agent'
date = "2025-9-9"
tags = ['llm', 'agent', 'go']
slug = 'my-todo-list-got-agent'
draft = false
+++

# My TODO List Finally Got Agent (Hello World Agent)

You know that feeling when your todo list becomes so long you spend more time searching for tasks than actually doing them? I had that exact problem. My running todo list lived in a simple text file, which worked fine until it didn't. Coming back after a few days away meant scrolling through endless bullet points trying to remember what was urgent.

Then I heard everyone talking about AI agents. While others were building complex systems, I wondered: what's the smallest possible agent I could create that would actually solve my real problem? 

The answer surprised me. In just 300 lines of Go code, I built something that lets me have natural conversations with my todo list. No more hunting through text files—I just ask what's due tomorrow or tell it to add something new.


## What Makes This Work

The secret isn't complex algorithms or fancy frameworks. An AI agent is fundamentally about bridging two worlds: human language and computer operations. When I say "add wash dishes to my todos," the AI understands my intent and translates it into a file modification. When I ask "what's on my list?", it knows to read the file and format the response nicely.

This translation layer is where the magic happens. The LLM acts as an intelligent interpreter that can understand context, intent, and natural language, then decide which programmatic actions to take.

## Two Simple Tools, Infinite Possibilities

My todo agent operates with just two capabilities:

- **Reading todo files** - How else would it know what I've written down?
- **Editing todo files** - The whole point of any todo system

That's the complete toolkit. The AI figures out when to use each tool based on our conversation. Ask "What's on my list?" and it reads. Say "Add buy milk to my todos" and it edits. The intelligence lies in this decision-making process.


## Building It Step By Step

The implementation starts simple and grows incrementally. Here's how I approached it:

### Prerequisites

Lets start 

You'll need an Anthropic API key from [console.anthropic.com](https://console.anthropic.com/settings/keys). Set it as an environment variable:

```
export ANTHROPIC_API_KEY="this is the last time i will tell you to set this"
```

This key gives you access to Claude's capabilities. I purchase token worth $5 to experimenting and building this project.

Every time you restart your terminal or machine, you need to set it again. This shouldn't be committed in git repo. It will cerate risk of exposing your API keys to others.


### Project Setup

```
mkdir todo-agent
cd todo-agent
go mod init agent
touch main.go
```

### Core Architecture

- Establish back and fourth communication with Claude via integrating the API in our code
- Setting up first tool for reading the todo.txt file 
- If todo.txt file doesn't exists, create the file
- Modify the todo.txt file based on users request

The foundation is a conversation loop that maintains context across multiple exchanges:

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

### The Conversation Engine

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

This creates the basic chat interface. Each conversation turn builds on the previous context, allowing for natural, flowing interactions.

### First Test Run

```
export ANTHROPIC_API_KEY="this is the last time i will tell you to set this"

# Download the dependencies
go mod tidy
# Run it
go run main.go
```

Here's what a basic conversation looks like:

```
🤖 Todo Agent Ready! (Press Ctrl+C to quit)
Try: what's on my todo list?
You: Hey I am Mahesh, How are you?
Claude: Hello Mahesh! I'm doing well, thank you for asking. How are you today? Is there something specific related to todo files that you'd like help with? I can help you read or edit todo files if you need assistance with that.
You: What is my  Name
Claude: Your name is Mahesh, as you mentioned in your introduction. Is there something specific you'd like help with regarding todo files? I can help you read or edit todo files using the available tools if needed.
You: List three famous person with this name
Claude: I'd be happy to tell you about three famous people with the name Mahesh, though I should note that I don't have access to a specific tool to look up this information. Based on general knowledge:

1. Mahesh Babu - A very popular Indian film actor who works primarily in Telugu cinema. He's known for films like "Pokiri," "Dookudu," and "Bharat Ane Nenu."

2. Mahesh Bhatt - A well-known Indian film director, producer, and screenwriter who has directed critically acclaimed films like "Arth," "Saaransh," and "Zakhm."

3. Mahesh Bhupathi - A former professional tennis player from India who became the first Indian to win a Grand Slam tournament. He's known for his doubles expertise, having won numerous Grand Slam doubles titles.

Is there anything specific about todo management you'd like help with? I can assist with reading or editing todo files using the available tools.
```


The AI maintains conversational memory and stays focused on its primary purpose—helping with todo management.

## Implementing the Read Tool


Tools are the bridge between conversation and action. Each tool needs four components:

- **Name and description** - What the tool does and when to use it
- **Input schema** - What parameters it expects
- **Execution function** - The actual implementation

Here's the complete tool infrastructure:


```
type ToolDefinition struct {
	Name        string                         `json:"name"`
	Description string                         `json:"description"`
	InputSchema anthropic.ToolInputSchemaParam `json:"input_schema"`
	Function    func(input json.RawMessage) (string, error)
}
```

Update the Agent struct to include tools:

```
// `tools` is added here:
type Agent struct {
	client         *anthropic.Client
	getUserMessage func() (string, bool)
	tools          []ToolDefinition
}

// And here:
func NewAgent(client *anthropic.Client,
	getUserMessage func() (string, bool),
	tools []ToolDefinition) *Agent {
	return &Agent{
		client:         client,
		getUserMessage: getUserMessage,
		tools:          tools,
	}
}

// And here:
func main() {
	client := anthropic.NewClient()

	scanner := bufio.NewScanner(os.Stdin)

	getUserMessage := func() (string, bool) {
		if !scanner.Scan() {
			return "", false
		}

		return scanner.Text(), true
	}

	tools := []ToolDefinition{}

	agent := NewAgent(&client, getUserMessage, tools)
	err := agent.Run(context.TODO())
	if err != nil {
		fmt.Printf("error: %s\n", err.Error())
	}

}
```

Integrate tools into the API call:

```
func (a *Agent) askClaude(ctx context.Context, conversation []anthropic.MessageParam) (*anthropic.Message, error) {

	anthropicTools := []anthropic.ToolUnionParam{}

	for _, tool := range a.tools {
		anthropicTools = append(anthropicTools, anthropic.ToolUnionParam{
			OfTool: &anthropic.ToolParam{
				Name:        tool.Name,
				Description: anthropic.String(tool.Description),
				InputSchema: tool.InputSchema,
			},
		})
	}

	message, err := a.client.Messages.New(ctx, anthropic.MessageNewParams{
		Model:     anthropic.ModelClaude3_7SonnetLatest,
		MaxTokens: int64(1024),
		Messages:  conversation,
		Tools:     anthropicTools,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create message: %w", err)
	}

	return message, err
}

```

Now implement the read_todo tool:

```
var ReadToDoDefinition = ToolDefinition{
	Name:        "read_todo",
	Description: "Read the contents of a todo file (todo.txt) . Use this to see what's in the todo file.",
	InputSchema: ReadToDoInputSchema,
	Function:    ReadToDo,
}

type ReadToDoInput struct {
	Path string `json:"path" jsonschema_description: "Path to the todo file in current working directory (e.g., 'todo.txt')"`
}

var ReadToDoInputSchema = GenerateSchema[ReadToDoInput]()

func ReadToDo(input json.RawMessage) (string, error) {
	readToDoInput := ReadToDoInput{}
	err := json.Unmarshal(input, &readToDoInput)

	if err != nil {
		panic(err)
	}

	content, err := os.ReadFile(readToDoInput.Path)

	if err != nil {
		return "", fmt.Errorf("could not read file: %v", err)
	}

	return string(content), nil
}

func GenerateSchema[T any]() anthropic.ToolInputSchemaParam {
	reflector := jsonschema.Reflector{
		AllowAdditionalProperties: false,
		DoNotReference:            true,
	}

	var v T
	schema := reflector.Reflect(v)

	return anthropic.ToolInputSchemaParam{
		Properties: schema.Properties,
	}
}
```

Add the required import:

```
// main.go

package main

import (
	"bufio"
	"context"
    // Add this:
	"encoding/json"
	"fmt"
	"os"

	"github.com/anthropics/anthropic-sdk-go"
    // Add this:
	"github.com/invopop/jsonschema"
)
```

Install dependencies:

```
go mod tidy
```

Enable the tool in main:

```
func main() {
    // [... previous code ...]
    tools := []ToolDefinition{ReadToDoDefinition}
    // [... previous code ...]
}
```


### Testing the Read Functionality

Create a sample todo file:

```
touch todo.txt
```

and add few tasks in `todo.txt` like
```
- buy pencil from store
- do code review of pull request #456
- send email note to someone 
```

Test it:

```
go run main.go
🤖 Todo Agent Ready! (Press Ctrl+C to quit)
Try: 'add buy milk to my todo list' or 'what's on my todo list?'

You: Whats in my todo list
Claude: I'll help you check what's in your todo list. Let me retrieve that information for you.
tool: read_todo({"path":"todo.txt"})
Claude: Here is what's currently in your todo list:
1. Buy pencil from store
2. Do code review of pull request #456
3. Send email note to someone

Is there anything specific you'd like to add, remove, or modify in your todo list?
```
The AI automatically decided to use the read_todo tool based on my natural language request. This demonstrates the power of letting the LLM determine when and how to use tools.

## Adding Edit Capabilities

Reading is only half the story. The real power comes from being able to modify the todo list through conversation:


We are going to use the same pattern as `read_todo`. 

- Define  EditToDoScheme definition and EditToDoInput schema
- `EditToDo` function which do the action 
- `createNewFile` create new todo.txt if needed 

```

var EditToDoDefinition = ToolDefinition{
	Name: "edit_todo",
	Description: `Edit a todo file by replacing old text with new text.
	Replaces 'old_str' with 'new_str' in the given todo file. 'old_str' and 'new_str' MUST be different from each other.

	If the file doesn't exist, it will be created.
	`,
	InputSchema: EditToDoInputSchema,
	Function:    EditToDo,
}

type EditToDoInput struct {
	Path   string `json:"path" jsonschema_description:"The path to the todo file"`
	OldStr string `json:"old_str" jsonschema_description:"Text to search for - must match exactly and must only have one match exactly"`
	NewStr string `json:"new_str" jsonschema_description:"Text to replace old_str with"`
}

var EditToDoInputSchema = GenerateSchema[EditToDoInput]()

func EditToDo(input json.RawMessage) (string, error) {
	editToDoInput := EditToDoInput{}

	err := json.Unmarshal(input, &editToDoInput)

	if err != nil {
		return "", err
	}

	if editToDoInput.Path == "" || editToDoInput.OldStr == editToDoInput.NewStr {
		return "", fmt.Errorf("invalid input params")
	}

	content, err := os.ReadFile(editToDoInput.Path)

	if err != nil {
		if os.IsNotExist(err) && editToDoInput.OldStr == "" {
			return createNewFile(editToDoInput.Path, editToDoInput.NewStr)
		}
		return "", err
	}

	oldContent := string(content)

	newContent := strings.Replace(oldContent, editToDoInput.OldStr, editToDoInput.NewStr, -1)

	if oldContent == newContent && editToDoInput.OldStr != "" {
		return "", fmt.Errorf("old_str not found in file")
	}

	err = os.WriteFile(editToDoInput.Path, []byte(newContent), 0644)
	if err != nil {
		return "", err
	}

	return "OK", nil
}

func createNewFile(filePath, content string) (string, error) {
	dir := path.Dir(filePath)
	if dir != "." {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return "", fmt.Errorf("failed to create dir: %w", err)
		}
	}

	err := os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}

	return fmt.Sprintf("Successfully created file %s", filePath), nil
}
```

Add both tools to the agent:


```
// main.go

func main() {
	// [... previous code ...]

	tools := []ToolDefinition{ReadToDoDefinition, EditToDoDefinition}

	// [... previous code ...]
}
```

### Testing the Complete System

Now the real magic happens:

```
go run main.go
🤖 Todo Agent Ready! (Press Ctrl+C to quit)
Try: 'add buy milk to my todo list' or 'what's on my todo list?'

You: What's in my todo
Claude: I'll help you check what's in your todo list. Let me retrieve the contents of your todo file.
tool: read_todo({"path":"todo.txt"})
Claude: Here's what's in your todo list:
- buy pencil from store
- do code review of pull request #456
- send email note to someone
You: Add reply to matt's email
Claude: I'll add "reply to Matt's email" to your todo list. Let me do that for you.

First, I'll read the current contents of your todo file to know where to add the new item:
tool: read_todo({"path":"todo.txt"})
Claude: Now I'll add your new todo item to the list:
tool: edit_todo({"path":"todo.txt","old_str":"- buy pencil from store\n- do code review of pull request #456\n- send email note to someone ","new_str":"- buy pencil from store\n- do code review of pull request #456\n- send email note to someone \n- reply to Matt's email"})
Claude: I've successfully added "reply to Matt's email" to your todo list. Your updated todo list now includes this new item.
You: Now what's in my todo
Claude: Let me check what's currently in your todo list:
tool: read_todo({"path":"todo.txt"})
Claude: Here's what's in your todo list now:

- buy pencil from store
- do code review of pull request #456
- send email note to someone
- reply to Matt's email
```

The AI understood my intent to add a task, read the current file to understand the structure, then edited it appropriately. This multi-step reasoning happens automatically.

### Advanced Interactions

The system handles complex requests naturally:

```
You: schedule by pencil due tomorrow
Claude: I understand you want to set a due date for the pencil item on your todo list. Let me update it to indicate it's due tomorrow.

First, let me read your current todo list:
tool: read_todo({"path":"todo.txt"})
Claude: Now, let me update the pencil item to include the due date:
tool: edit_todo({"path":"todo.txt","old_str":"- buy pencil from store","new_str":"- buy pencil from store (due tomorrow)"})
Claude: I've updated your todo list to indicate that buying a pencil is due tomorrow. The item now reads "- buy pencil from store (due tomorrow)".
You: now what's in my todo
Claude: Let me check your current todo list for you:
tool: read_todo({"path":"todo.txt"})
Claude: Here's what's currently in your todo list:

- buy pencil from store (due tomorrow)
- do code review of pull request #456
- send email note to someone
- reply to Matt's email
```

Even with a typo in my request ("schedule by pencil"), the AI understood my intent and correctly modified the task to include a due date.

### Query-Based Interactions

The agent can answer specific questions about todo items:

```
You: What is due tomorrow
Claude: Let me check your todo list and find what's due tomorrow:
tool: read_todo({"path":"todo.txt"})
Claude: Based on your todo list, you have one item due tomorrow:
- buy pencil from store (due tomorrow)
```

## The Power of Simplicity

What strikes me most about this project is how simple yet powerful it is. With just 300 lines of code and two basic operations (read and edit), I've created something that fundamentally changes how I interact with my todo list.

The key insight is that the intelligence doesn't live in complex algorithms or elaborate frameworks. It lives in the LLM's ability to understand intent and translate it into appropriate actions. The agent layer just provides the mechanisms for that translation to happen.

This opens up incredible possibilities. I could easily extend this to integrate with other todo apps through their APIs, add calendar functionality, or create entirely new interfaces. The conversational layer removes the friction between human intention and computer action.

The models available today make this kind of natural language interface not just possible, but surprisingly straightforward to implement. We're witnessing a fundamental shift in how humans and computers can interact, and tools like this are just the beginning.

Building this agent taught me that the future isn't necessarily about more complex software—it's about software that understands us better.




