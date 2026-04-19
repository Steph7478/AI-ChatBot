# ☕ Neural Java Chatbot - AI-Powered Java Assistant

A lightweight neural network chatbot specialized in the Java programming language. Built with pure Go, this chatbot uses a Transformer with Multi-Head Attention to generate responses and learns from conversations.

## Features

- **Neural Transformer Architecture** - Multi-layer Transformer with Multi-Head Attention
- **Word-based Tokenization** - Dynamic vocabulary building from conversations
- **Persistent Memory** - Saves neural network weights AND vocabulary to disk
- **Multiple File Support** - Loads from conversations.txt and prompts.txt
- **Intelligent Matching** - Three-tier matching system (exact → fuzzy → neural)
- **Synonym Resolution** - Maps variations like "hi hello hey sup" to the same response
- **Fuzzy String Matching** - Word overlap similarity with configurable threshold
- **Early Stopping** - Automatically stops training when loss stops improving
- **Repetition Penalty** - Prevents the model from repeating the same word
- **Interactive Learning** - Learns new answers directly from users when it doesn't know something
- **Colored Output** - Cyan bot responses with gray status indicators
- **Zero External Dependencies** - Pure Go implementation, no GPU required
- **Comprehensive Java Knowledge** - Covers history, syntax, JVM, frameworks, and career
- **Fun Personality** - Enthusiastic Java fanatic who playfully teases Python

## Commands

| Command | Description |
|---------|-------------|
| /quit | Save checkpoint and exit |
| /save | Save memory checkpoint manually |
| /stats | Show statistics including conversation count |
| /temp [0.1-2.0] | Adjust response creativity (lower = more predictable) |

## How It Works

The chatbot reads conversations from data/conversations.txt where each line follows the format question|answer. It also reads data/prompts.txt for synonym mappings. The matching strategy works in this order:

1. Exact Match - Returns exact answer if question exists in conversations
2. Synonym Resolution - Maps input to main phrase using prompts.txt
3. Fuzzy Match - Finds similar questions based on word overlap (threshold greater than 35 percent)
4. Neural Generation - Uses Transformer neural network to generate response
5. Fallback Learning - Asks user to teach when no match is found

When the neural network generates a response, you can confirm if it was useful or provide the correct answer. All learned data is saved to data/model.gob and data/model.gob.vocab for instant loading on next startup.

## Training

To train the neural network on your conversations:

./chatbot -train 20

The training process includes:
- Early stopping when loss stops improving (patience of 5 epochs)
- Repetition penalty to prevent mode collapse
- Progress display with loss values
- Automatic vocabulary building from conversations

## File Formats

data/conversations.txt:
hi|Hello! How can I help you?
what is java|Java is a programming language

data/prompts.txt:
hi|hi|hello|hey|hola|sup|howdy
goodbye|goodbye|bye|see you later

## Topics Covered

- Java history from Oak to modern versions
- JVM internals, memory management and garbage collection
- Object oriented programming concepts
- Collections framework and data structures
- Java 8 features including lambdas, streams, and optional
- Modern Java with virtual threads, records, and pattern matching
- Spring Boot framework and microservices
- Build tools like Maven and Gradle
- Testing with JUnit, Mockito, and Testcontainers
- Databases, JPA, Hibernate, and JDBC
- Java vs Python friendly roasts and jokes
- Career advice, certifications, and interview questions

## Personality

This chatbot has a strong personality. It absolutely loves Java and playfully teases Python in a lighthearted way. Expect responses like:

Java is like espresso - strong, fast, reliable. Python is like chamomile tea - tasty but slow.

Why did Python go to the doctor? Unhealthy indentation. Java just checked its exceptions.

Write once, run anywhere - Python can't say that.