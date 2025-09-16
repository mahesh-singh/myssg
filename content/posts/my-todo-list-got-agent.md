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

Here’s what we need:

- Go
- Anthropic API key that you set as an environment variable, ANTHROPIC_API_KEY

Setup a new Go project

```
mkdir code-editing-agent
cd code-editing-agent
go mod init agent
touch main.go
```

https://ampcode.com/how-to-build-an-agent