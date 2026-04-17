# ☕ Java Chatbot - AI-Powered Java Assistant

A lightweight retrieval-based chatbot specialized in the Java programming language. Built with pure Go, this chatbot learns from conversations and uses TF-IDF similarity to find the most relevant answers.

## Features

- Lightweight - Runs on any PC including old hardware with minimal RAM usage
- Persistent Memory - Saves learned conversations to checkpoint files for instant loading
- TF-IDF Similarity - Intelligent answer matching based on word importance
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

The chatbot reads conversations from a text file where each conversation has a User question and a Bot answer. It then builds a vocabulary and calculates TF-IDF vectors for each question. When you ask something, it converts your question to a vector and finds the most similar question in memory using cosine similarity, then returns the corresponding answer. All learned data is saved to a checkpoint file for instant loading on next startup.

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
- **Java vs Python jokes and friendly roasts** - Why Java is like a warm hug from a friendly robot and Python is like a slow bicycle with training wheels
- Career advice certifications and interview questions

## Personality & Humor

This chatbot has a strong personality! It absolutely LOVES Java and playfully teases Python in a cute, lighthearted way. Expect responses like:

- "java is like a warm hug from a friendly robot ☕💕"
- "python is okay i guess but java is the king"
- "python is like a bicycle java is like a spaceship 🚀"
- "java writes once runs anywhere python writes once cries everywhere because indentation broke again"

The bot is designed to be fun and engaging while still providing accurate technical information about Java. It's perfect for learning Java with a smile on your face!

## Configuration

Edit internal/config/config.go to adjust behavior including similarity threshold for match precision temperature for response randomness and top K candidates for answer selection.