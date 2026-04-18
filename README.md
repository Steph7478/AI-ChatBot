# ☕ Java Chatbot - AI-Powered Java Assistant

A lightweight neural network chatbot specialized in the Java programming language. Built with pure Go, this chatbot uses a Transformer with Multi-Head Attention to generate responses and learns from conversations.

## Features

- **Neural Transformer Architecture** - 2-layer Transformer with Multi-Head Attention (4 heads)
- Lightweight - Runs on any PC including old hardware with minimal RAM usage (~100MB)
- Persistent Memory - Saves neural network weights to model.gob
- Multiple File Support - Loads from conversations.txt, prompts.txt, and examples.txt
- Semantic Matching - Word similarity matching (30% threshold) when neural net needs help
- Synonym Support - Maps variations like "hi|hello|hey" to the same response
- Fallback Learning - Learns new answers directly from users when it doesn't know something
- Zero Dependencies - Pure Go implementation no external libraries or GPU required
- Comprehensive Java Knowledge - Covers history syntax JVM frameworks testing microservices and career
- Fun Personality - Enthusiastic Java fanatic who absolutely loves Java and playfully teases Python (all in good fun!)
- Java vs Python Jokes - Expect lots of good-natured jokes about why Java is the king and Python is... well, also a language that exists
- Interactive CLI - Simple command line interface with useful commands
- Configurable - Adjust similarity threshold temperature and response length
- Multi-platform - Works on Linux macOS and Windows

## Commands

- /quit - Save checkpoint and exit
- /save - Save memory checkpoint manually
- /reload - Reload conversations from file
- /stats - Show statistics including conversation count and vocabulary size
- /temp [0.1-1.5] - Adjust response creativity lower is more predictable

## How It Works

The chatbot reads conversations from a text file where each conversation has a User question and a Bot answer. It also reads prompts.txt for synonyms and examples.txt for training examples. It uses a Transformer neural network to generate responses. When you ask something, it first tries exact match, then neural generation, then semantic similarity. If no match is found with high confidence, the chatbot asks you to teach the correct answer, which is then saved for future use. All learned data is saved to a checkpoint file for instant loading on next startup.

### Matching Strategy

1. **Exact Match** - Returns the exact answer if the question exists
2. **Neural Generation** - Uses Transformer neural network to generate response
3. **Semantic Similarity** - Finds similar questions based on word overlap (threshold > 0.3)
4. **Fallback Learning** - Asks the user to teach when no match is found

## Topics Covered

- Java history from Oak to modern versions
- JVM internals memory management and garbage collection
- Object oriented programming concepts
- Collections framework and data structures
- Java 8 features lambdas streams and optional
- Modern Java with virtual threads records and pattern matching
- Spring Boot framework and microservices
- Build tools Maven and Gradle
- Testing with JUnit Mockito and Testcontainers
- Databases JPA Hibernate and JDBC
- **Java vs Python jokes and friendly roasts** - Why Java is like espresso strong fast reliable and Python is like chamomile tea tasty but slow
- Career advice certifications and interview questions

## Personality & Humor

This chatbot has a strong personality! It absolutely LOVES Java and playfully teases Python in a cute, lighthearted way. Expect responses like:

- "Java is like espresso - strong fast reliable. Python is like chamomile tea - tasty but slow! ☕⚡"
- "Java is like a tank - strong typing compiled runs anywhere. Python is like a bicycle - fun but don't take it to war! 😂"
- "Why did Python go to the doctor? Unhealthy indentation! Java just checked its exceptions 🤣"
- "Write once run anywhere - Python can't say that! 🎯☕"

The bot is designed to be fun and engaging while still providing accurate technical information about Java. It's perfect for learning Java with a smile on your face!