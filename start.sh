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
    echo -e "${GREEN}☕ Java Chatbot - Main Menu${NC}"
    echo "========================================="
    echo
    echo -e "${CYAN}1)${NC} Start Chatbot"
    echo -e "${CYAN}2)${NC} Show Statistics"
    echo -e "${CYAN}3)${NC} Edit conversations.txt"
    echo -e "${CYAN}4)${NC} Edit knowledge.txt"
    echo -e "${CYAN}5)${NC} Delete Checkpoint (reset learning)"
    echo -e "${CYAN}6)${NC} Delete Learned (reset N-Gram)"
    echo -e "${CYAN}7)${NC} Recompile"
    echo -e "${CYAN}8)${NC} Show Help"
    echo -e "${CYAN}9)${NC} Exit"
    echo
    echo -n "Choose an option [1-9]: "
}

# Function to show stats
show_stats() {
    echo
    echo "========================================="
    echo -e "${GREEN}📊 Statistics${NC}"
    echo "========================================="
    
    # conversations.txt stats
    if [ -f "data/conversations.txt" ]; then
        CONV_COUNT=$(grep -c "^User:" data/conversations.txt 2>/dev/null || echo "0")
        BOT_COUNT=$(grep -c "^Bot:" data/conversations.txt 2>/dev/null || echo "0")
        TOTAL_LINES=$(wc -l < data/conversations.txt 2>/dev/null || echo "0")
        FILE_SIZE=$(du -h data/conversations.txt 2>/dev/null | cut -f1)
        
        echo -e "${GREEN}✅ conversations.txt:${NC}"
        echo "   Conversations: $CONV_COUNT"
        echo "   Bot responses: $BOT_COUNT"
        echo "   Total lines: $TOTAL_LINES"
        echo "   Size: $FILE_SIZE"
    else
        echo -e "${RED}❌ conversations.txt not found${NC}"
    fi
    
    echo
    
    # knowledge.txt stats
    if [ -f "data/knowledge.txt" ]; then
        KNOW_LINES=$(grep -v "^#" data/knowledge.txt | grep -v "^$" | wc -l)
        KNOW_SIZE=$(du -h data/knowledge.txt 2>/dev/null | cut -f1)
        echo -e "${GREEN}✅ knowledge.txt:${NC}"
        echo "   Phrases: $KNOW_LINES"
        echo "   Size: $KNOW_SIZE"
    fi
    
    echo
    
    # learned.txt stats
    if [ -f "data/learned.txt" ]; then
        LEARN_LINES=$(grep -v "^#" data/learned.txt | grep -v "^$" | wc -l)
        LEARN_SIZE=$(du -h data/learned.txt 2>/dev/null | cut -f1)
        echo -e "${GREEN}✅ learned.txt:${NC}"
        echo "   Learned phrases: $LEARN_LINES"
        echo "   Size: $LEARN_SIZE"
    fi
    
    echo
    
    # Checkpoint stats
    if [ -f "data/checkpoint.gob" ]; then
        CHECK_SIZE=$(du -h data/checkpoint.gob 2>/dev/null | cut -f1)
        CHECK_DATE=$(stat -c %y data/checkpoint.gob 2>/dev/null || stat -f %Sm data/checkpoint.gob 2>/dev/null)
        echo -e "${GREEN}✅ Checkpoint (TF-IDF):${NC}"
        echo "   Size: $CHECK_SIZE"
        echo "   Last modified: $CHECK_DATE"
    else
        echo -e "${YELLOW}⚠️  No checkpoint found${NC}"
    fi
    
    echo
    
    # Binary stats
    if [ -f "chatbot" ]; then
        BIN_SIZE=$(du -h chatbot 2>/dev/null | cut -f1)
        echo -e "${GREEN}✅ Binary:${NC}"
        echo "   Size: $BIN_SIZE"
    fi
    
    echo
    echo -n "Press Enter to return to menu..."
    read
}

# Function to edit conversations.txt
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
    tail -20 data/conversations.txt | grep -E "User:|Bot:" | tail -10
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

# Function to edit knowledge.txt
edit_knowledge() {
    echo
    echo "========================================="
    echo -e "${GREEN}✏️  Edit knowledge.txt${NC}"
    echo "========================================="
    echo
    
    if [ ! -f "data/knowledge.txt" ]; then
        echo -e "${RED}❌ data/knowledge.txt not found!${NC}"
        echo -n "Press Enter to return to menu..."
        read
        return
    fi
    
    # Show first 5 lines
    echo -e "${CYAN}First 5 phrases in knowledge base:${NC}"
    echo
    head -10 data/knowledge.txt | grep -v "^#" | head -5
    echo
    echo "Choose editor:"
    echo "  1) nano"
    echo "  2) vim"
    echo "  3) VS Code"
    echo "  4) Cancel"
    echo
    read -p "Choose editor [1-4]: " editor_choice
    
    case $editor_choice in
        1) nano data/knowledge.txt ;;
        2) vim data/knowledge.txt ;;
        3) code data/knowledge.txt ;;
        *) echo "Cancelled" ;;
    esac
    
    echo
    echo -e "${GREEN}✅ Knowledge edited! Restart chatbot to apply${NC}"
    echo
    echo -n "Press Enter to return to menu..."
    read
}

# Function to delete checkpoint
delete_checkpoint() {
    echo
    echo "========================================="
    echo -e "${RED}⚠️  DELETE CHECKPOINT (TF-IDF)${NC}"
    echo "========================================="
    echo
    echo -e "${YELLOW}This will delete all learned TF-IDF memory!${NC}"
    echo "The chatbot will restart learning from conversations.txt"
    echo
    read -p "Are you sure? (yes/no): " confirm
    
    if [ "$confirm" = "yes" ]; then
        rm -f data/checkpoint.gob
        echo -e "${GREEN}✅ Checkpoint deleted!${NC}"
    else
        echo -e "${BLUE}Cancelled${NC}"
    fi
    
    echo
    echo -n "Press Enter to return to menu..."
    read
}

# Function to delete learned.txt
delete_learned() {
    echo
    echo "========================================="
    echo -e "${RED}⚠️  DELETE LEARNED (N-Gram)${NC}"
    echo "========================================="
    echo
    echo -e "${YELLOW}This will delete all N-Gram learned phrases!${NC}"
    echo "The chatbot will only have knowledge.txt base"
    echo
    read -p "Are you sure? (yes/no): " confirm
    
    if [ "$confirm" = "yes" ]; then
        rm -f data/learned.txt
        touch data/learned.txt
        echo -e "${GREEN}✅ learned.txt cleared!${NC}"
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
    
    # Check if Go is installed
    if ! command -v go &> /dev/null; then
        echo -e "${RED}❌ Go is not installed!${NC}"
        echo -n "Press Enter to continue..."
        read
        return
    fi
    
    # Clean old binary
    rm -f chatbot
    
    # Compile
    go build -o chatbot ./cmd
    
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✅ Compilation successful!${NC}"
        BIN_SIZE=$(du -h chatbot 2>/dev/null | cut -f1)
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
    echo -e "${CYAN}File Format (conversations.txt):${NC}"
    echo "  User: your question here"
    echo "  Bot: answer here"
    echo "  (blank line between conversations)"
    echo
    echo -e "${CYAN}File Format (knowledge.txt):${NC}"
    echo "  # Comments start with #"
    echo "  just write phrases one per line"
    echo "  the bot uses these to generate new responses"
    echo
    echo -e "${CYAN}Example (conversations.txt):${NC}"
    echo "  User: what is java"
    echo "  Bot: java is a programming language"
    echo "  "
    echo "  User: is java good"
    echo "  Bot: yes java is wonderful"
    echo
    echo -e "${CYAN}How learning works:${NC}"
    echo "  1. Bot tries to find answer in conversations.txt"
    echo "  2. If not found, tries to generate with N-Gram"
    echo "  3. If still nothing, asks you to teach"
    echo "  4. Your teaching is saved to learned.txt and checkpoint"
    echo
    echo -e "${CYAN}Tips:${NC}"
    echo "  - More conversations = better answers"
    echo "  - Use /reload after editing conversations.txt"
    echo "  - Lower temperature = more predictable"
    echo "  - Higher temperature = more creative"
    echo "  - Delete checkpoint to reset TF-IDF memory"
    echo "  - Delete learned.txt to reset N-Gram memory"
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
            
            # Check if data directory exists
            if [ ! -d "data" ]; then
                mkdir -p data
                echo -e "${GREEN}✅ Created data directory${NC}"
            fi
            
            # Check if conversations.txt exists
            if [ ! -f "data/conversations.txt" ]; then
                echo -e "${YELLOW}⚠️  data/conversations.txt not found! Creating example...${NC}"
                cat > data/conversations.txt << 'EOF'
User: hello
Bot: hello welcome to the java chatbot

User: what is java
Bot: java is a programming language created by james gosling in 1991

User: who created java
Bot: james gosling created java at sun microsystems
EOF
                echo -e "${GREEN}✅ Created example conversations.txt${NC}"
            fi
            
            # Check if knowledge.txt exists
            if [ ! -f "data/knowledge.txt" ]; then
                echo -e "${YELLOW}⚠️  data/knowledge.txt not found! Creating example...${NC}"
                cat > data/knowledge.txt << 'EOF'
# Java Knowledge Base
java is the most amazing programming language
java runs on over 3 billion devices
james gosling created java in 1991
the jvm is the heart of java
spring boot makes java web development easy
EOF
                echo -e "${GREEN}✅ Created example knowledge.txt${NC}"
            fi
            
            # Check if learned.txt exists
            if [ ! -f "data/learned.txt" ]; then
                touch data/learned.txt
                echo -e "${GREEN}✅ Created learned.txt${NC}"
            fi
            
            # Count conversations
            CONV_COUNT=$(grep -c "^User:" data/conversations.txt 2>/dev/null || echo "0")
            echo -e "${GREEN}✅ Found $CONV_COUNT conversations${NC}"
            
            # Check if binary needs compilation
            NEED_COMPILE=0
            
            if [ ! -f "chatbot" ]; then
                NEED_COMPILE=1
            else
                # Check source files
                for file in cmd/main.go internal/app/*.go internal/dataset/*.go internal/memory/*.go internal/ngram/*.go internal/config/*.go; do
                    if [ -f "$file" ] && [ "$file" -nt "chatbot" ] 2>/dev/null; then
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
            edit_knowledge
            ;;
        5)
            delete_checkpoint
            ;;
        6)
            delete_learned
            ;;
        7)
            recompile
            ;;
        8)
            show_help
            ;;
        9)
            echo
            echo -e "${GREEN}👋 Goodbye! Keep coding in Java! ☕${NC}"
            exit 0
            ;;
        *)
            echo -e "${RED}Invalid option!${NC}"
            sleep 1
            ;;
    esac
done