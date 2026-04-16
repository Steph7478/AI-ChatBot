#!/bin/bash

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Function to show menu
show_menu() {
    clear
    echo "========================================="
    echo -e "${GREEN}🤖 Java Chatbot - Main Menu${NC}"
    echo "========================================="
    echo
    echo -e "${CYAN}1)${NC} Start Chatbot"
    echo -e "${CYAN}2)${NC} Show Statistics"
    echo -e "${CYAN}3)${NC} Edit conversations.txt"
    echo -e "${CYAN}4)${NC} Delete Checkpoint (reset learning)"
    echo -e "${CYAN}5)${NC} Recompile"
    echo -e "${CYAN}6)${NC} Show Help"
    echo -e "${CYAN}7)${NC} Exit"
    echo
    echo -n "Choose an option [1-7]: "
}

# Function to show stats
show_stats() {
    echo
    echo "========================================="
    echo -e "${GREEN}📊 Statistics${NC}"
    echo "========================================="
    
    if [ -f "data/conversations.txt" ]; then
        CONV_COUNT=$(grep -c "^User:" data/conversations.txt 2>/dev/null || echo "0")
        LINE_COUNT=$(wc -l < data/conversations.txt)
        WORD_COUNT=$(wc -w < data/conversations.txt)
        FILE_SIZE=$(du -h data/conversations.txt | cut -f1)
        
        echo -e "${GREEN}✅ conversations.txt:${NC}"
        echo "   Conversations: $CONV_COUNT"
        echo "   Lines: $LINE_COUNT"
        echo "   Words: $WORD_COUNT"
        echo "   Size: $FILE_SIZE"
    else
        echo -e "${RED}❌ conversations.txt not found${NC}"
    fi
    
    echo
    
    if [ -f "data/memory_checkpoint.gob" ]; then
        CHECK_SIZE=$(du -h data/memory_checkpoint.gob | cut -f1)
        CHECK_DATE=$(stat -c %y data/memory_checkpoint.gob 2>/dev/null || stat -f %Sm data/memory_checkpoint.gob 2>/dev/null)
        echo -e "${GREEN}✅ Checkpoint:${NC}"
        echo "   Size: $CHECK_SIZE"
        echo "   Last modified: $CHECK_DATE"
    else
        echo -e "${YELLOW}⚠️  No checkpoint found${NC}"
    fi
    
    echo
    
    if [ -f "chatbot" ]; then
        BIN_SIZE=$(du -h chatbot | cut -f1)
        echo -e "${GREEN}✅ Binary:${NC}"
        echo "   Size: $BIN_SIZE"
    fi
    
    echo
    echo -n "Press Enter to return to menu..."
    read
}

# Function to edit conversations
edit_conversations() {
    echo
    echo "========================================="
    echo -e "${GREEN}✏️  Edit conversations.txt${NC}"
    echo "========================================="
    echo
    
    if [ ! -f "data/conversations.txt" ]; then
        echo -e "${RED}❌ data/conversations.txt not found!${NC}"
        echo -n "Press Enter to return to menu..."
        read
        return
    fi
    
    # Show last 5 conversations
    echo -e "${CYAN}Last 5 conversations in file:${NC}"
    echo
    tail -10 data/conversations.txt | grep -E "User:|Bot:" | tail -5
    echo
    echo "Choose editor:"
    echo "  1) nano"
    echo "  2) vim"
    echo "  3) VS Code"
    echo "  4) Cancel"
    echo
    read -p "Choose editor [1-4]: " editor_choice
    
    case $editor_choice in
        1) nano data/conversations.txt ;;
        2) vim data/conversations.txt ;;
        3) code data/conversations.txt ;;
        *) echo "Cancelled" ;;
    esac
    
    echo
    echo -e "${GREEN}✅ File edited! Use /reload in chatbot to reload${NC}"
    echo
    echo -n "Press Enter to return to menu..."
    read
}

# Function to delete checkpoint
delete_checkpoint() {
    echo
    echo "========================================="
    echo -e "${RED}⚠️  DELETE CHECKPOINT${NC}"
    echo "========================================="
    echo
    echo -e "${YELLOW}This will delete all learned memory!${NC}"
    echo "The chatbot will restart learning from conversations.txt"
    echo
    read -p "Are you sure? (yes/no): " confirm
    
    if [ "$confirm" = "yes" ]; then
        rm -f data/memory_checkpoint.gob
        echo -e "${GREEN}✅ Checkpoint deleted!${NC}"
    else
        echo -e "${BLUE}Cancelled${NC}"
    fi
    
    echo
    echo -n "Press Enter to return to menu..."
    read
}

# Function to recompile
recompile() {
    echo
    echo "========================================="
    echo -e "${GREEN}📦 Recompiling...${NC}"
    echo "========================================="
    
    # Clean old binary
    rm -f chatbot
    
    # Compile
    go build -o chatbot ./cmd
    
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✅ Compilation successful!${NC}"
        BIN_SIZE=$(du -h chatbot | cut -f1)
        echo "   Binary size: $BIN_SIZE"
    else
        echo -e "${RED}❌ Compilation failed!${NC}"
    fi
    
    echo
    echo -n "Press Enter to return to menu..."
    read
}

# Function to show help
show_help() {
    echo
    echo "========================================="
    echo -e "${GREEN}📖 Help Guide${NC}"
    echo "========================================="
    echo
    echo -e "${CYAN}Chatbot Commands:${NC}"
    echo "  /quit      - Save and exit chatbot"
    echo "  /save      - Save checkpoint manually"
    echo "  /reload    - Reload from conversations.txt"
    echo "  /stats     - Show statistics"
    echo "  /temp X    - Set temperature (0.1-1.5)"
    echo
    echo -e "${CYAN}File Format:${NC}"
    echo "  User: your question here"
    echo "  Bot: answer here"
    echo "  (blank line between conversations)"
    echo
    echo -e "${CYAN}Example:${NC}"
    echo "  User: what is java"
    echo "  Bot: java is a programming language"
    echo "  "
    echo "  User: is java good"
    echo "  Bot: yes java is wonderful"
    echo
    echo -e "${CYAN}Tips:${NC}"
    echo "  - Add as many conversations as you want"
    echo "  - More conversations = better answers"
    echo "  - Use /reload after editing file"
    echo "  - Lower temperature = more predictable"
    echo "  - Higher temperature = more creative"
    echo
    echo -n "Press Enter to return to menu..."
    read
}

# Main loop
while true; do
    show_menu
    read choice
    
    case $choice in
        1)
            echo
            echo "========================================="
            echo -e "${GREEN}🚀 Starting Chatbot...${NC}"
            echo "========================================="
            echo
            
            # Check if Go is installed
            if ! command -v go &> /dev/null; then
                echo -e "${RED}❌ Go is not installed!${NC}"
                echo "Install from: https://golang.org/dl/"
                echo -n "Press Enter to continue..."
                read
                continue
            fi
            
            # Check if conversations.txt exists
            if [ ! -f "data/conversations.txt" ]; then
                echo -e "${RED}❌ data/conversations.txt not found!${NC}"
                echo -n "Press Enter to continue..."
                read
                continue
            fi
            
            # Count conversations
            CONV_COUNT=$(grep -c "^User:" data/conversations.txt 2>/dev/null || echo "0")
            echo -e "${GREEN}✅ Found $CONV_COUNT conversations${NC}"
            
            if [ "$CONV_COUNT" -eq 0 ]; then
                echo -e "${YELLOW}⚠️  No conversations found!${NC}"
                echo -n "Press Enter to continue..."
                read
                continue
            fi
            
            # Check if source files have changed
            NEED_COMPILE=0
            
            if [ ! -f "chatbot" ]; then
                NEED_COMPILE=1
            else
                for file in cmd/main.go internal/dataset/dataset.go internal/memory/memory.go; do
                    if [ -f "$file" ] && [ "$file" -nt "chatbot" ]; then
                        NEED_COMPILE=1
                        break
                    fi
                done
            fi
            
            # Compile if needed
            if [ $NEED_COMPILE -eq 1 ]; then
                echo -e "${BLUE}📦 Compiling...${NC}"
                go build -o chatbot ./cmd
                if [ $? -ne 0 ]; then
                    echo -e "${RED}❌ Compilation failed${NC}"
                    echo -n "Press Enter to continue..."
                    read
                    continue
                fi
                echo -e "${GREEN}✅ Compiled successfully${NC}"
            fi
            
            echo
            echo "========================================="
            echo -e "${GREEN}💬 Starting chatbot...${NC}"
            echo "Type /quit to return to menu"
            echo "========================================="
            echo
            
            # Run the chatbot
            ./chatbot
            
            echo
            echo -e "${GREEN}✅ Returned to menu${NC}"
            echo -n "Press Enter to continue..."
            read
            ;;
        2)
            show_stats
            ;;
        3)
            edit_conversations
            ;;
        4)
            delete_checkpoint
            ;;
        5)
            recompile
            ;;
        6)
            show_help
            ;;
        7)
            echo
            echo -e "${GREEN}👋 Goodbye!${NC}"
            exit 0
            ;;
        *)
            echo -e "${RED}Invalid option!${NC}"
            sleep 1
            ;;
    esac
done