# ☕ Java Chatbot - AI-Powered Java Assistant

A lightweight retrieval-based chatbot specialized in the Java programming language. Built with pure Go, this chatbot learns from conversations and uses TF-IDF similarity to find the most relevant answers.

## Features

- Lightweight - Runs on any PC including old hardware with minimal RAM usage
- Persistent Memory - Saves learned conversations to checkpoint files for instant loading
- TF-IDF Similarity - Intelligent answer matching based on word importance
- Zero Dependencies - Pure Go implementation no external libraries or GPU required
- Comprehensive Java Knowledge - Covers history syntax JVM frameworks testing microservices and career
- Fun Personality - Enthusiastic Java fanatic with a sense of humor and jokes
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
- Java jokes and humor
- Career advice certifications and interview questions

## Configuration

Edit internal/config/config.go to adjust behavior including similarity threshold for match precision temperature for response randomness and top K candidates for answer selection.