#!/bin/bash

# ============================================
# ☕ Java Chatbot - Interactive Menu
# ============================================

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
MAGENTA='\033[0;35m'
BOLD='\033[1m'
NC='\033[0m'

# Show menu
show_menu() {
    clear
    echo -e "${BOLD}${MAGENTA}"
    echo "    ╔══════════════════════════════════════════════╗"
    echo "    ║     ☕  JAVA CHATBOT - MAIN MENU  ☕         ║"
    echo "    ╚══════════════════════════════════════════════╝"
    echo -e "${NC}"
    echo
    echo -e "  ${CYAN}📌${NC} Choose an option:"
    echo
    echo -e "    ${GREEN}1)${NC} 🚀 Start Chatbot"
    echo -e "    ${GREEN}2)${NC} 📊 Show Statistics"
    echo -e "    ${GREEN}3)${NC} ✏️  Edit conversations.txt"
    echo -e "    ${GREEN}4)${NC} 🗑️  Delete Checkpoint (reset learning)"
    echo -e "    ${GREEN}5)${NC} 🔧 Recompile"
    echo -e "    ${GREEN}6)${NC} ❓ Show Help"
    echo -e "    ${GREEN}7)${NC} 👋 Exit"
    echo
    echo -ne "${YELLOW}➜${NC} "
}

# Show stats
show_stats() {
    echo
    echo -e "${BOLD}${CYAN}╔════════════════════════════════════════════╗${NC}"
    echo -e "${BOLD}${CYAN}║              📊 STATISTICS                 ║${NC}"
    echo -e "${BOLD}${CYAN}╚════════════════════════════════════════════╝${NC}"
    echo
    
    if [ -f "data/conversations.txt" ]; then
        CONV_COUNT=$(grep -c "^User:" data/conversations.txt 2>/dev/null || echo "0")
        FILE_SIZE=$(du -h data/conversations.txt 2>/dev/null | cut -f1)
        
        echo -e "  ${GREEN}✓${NC} ${BOLD}conversations.txt${NC}"
        echo -e "    📝 Conversations: ${CYAN}$CONV_COUNT${NC}"
        echo -e "    💾 Size: ${CYAN}$FILE_SIZE${NC}"
    else
        echo -e "  ${RED}✗${NC} conversations.txt ${RED}not found${NC}"
    fi
    
    echo
    
    if [ -f "data/checkpoint.gob" ]; then
        CHECK_SIZE=$(du -h data/checkpoint.gob 2>/dev/null | cut -f1)
        echo -e "  ${GREEN}✓${NC} ${BOLD}checkpoint.gob${NC}"
        echo -e "    💾 Size: ${CYAN}$CHECK_SIZE${NC}"
    else
        echo -e "  ${YELLOW}⚠${NC} No checkpoint found"
    fi
    
    echo
    
    if [ -f "chatbot" ]; then
        BIN_SIZE=$(du -h chatbot 2>/dev/null | cut -f1)
        echo -e "  ${GREEN}✓${NC} ${BOLD}Binary${NC}"
        echo -e "    💾 Size: ${CYAN}$BIN_SIZE${NC}"
    fi
    
    echo
    read -p "  Press Enter to continue..."
}

# Edit conversations
edit_conversations() {
    echo
    echo -e "${BOLD}${CYAN}╔════════════════════════════════════════════╗${NC}"
    echo -e "${BOLD}${CYAN}║           ✏️  EDIT CONVERSATIONS           ║${NC}"
    echo -e "${BOLD}${CYAN}╚════════════════════════════════════════════╝${NC}"
    echo
    
    if [ ! -f "data/conversations.txt" ]; then
        echo -e "  ${RED}❌ conversations.txt not found!${NC}"
        read -p "  Press Enter to continue..."
        return
    fi
    
    echo -e "  ${CYAN}📋 Last conversations:${NC}\n"
    tail -20 data/conversations.txt | grep -E "User:|Bot:" | tail -6 | while read line; do
        if [[ $line == User:* ]]; then
            echo -e "    ${GREEN}→${NC} $line"
        else
            echo -e "    ${BLUE}←${NC} $line"
        fi
    done
    
    echo
    echo -e "  ${YELLOW}Choose editor:${NC}"
    echo -e "    ${GREEN}1)${NC} nano"
    echo -e "    ${GREEN}2)${NC} vim"
    echo -e "    ${GREEN}3)${NC} VS Code"
    echo -e "    ${GREEN}4)${NC} Cancel"
    echo
    read -p "  ➜ " editor_choice
    
    case $editor_choice in
        1) nano data/conversations.txt ;;
        2) vim data/conversations.txt ;;
        3) code data/conversations.txt ;;
        *) echo -e "  ${BLUE}Cancelled${NC}" ;;
    esac
    
    echo -e "\n  ${GREEN}✅ Done! Use /reload in chatbot to apply changes${NC}"
    read -p "  Press Enter to continue..."
}

# Delete checkpoint
delete_checkpoint() {
    echo
    echo -e "${BOLD}${RED}╔════════════════════════════════════════════╗${NC}"
    echo -e "${BOLD}${RED}║           🗑️  DELETE CHECKPOINT            ║${NC}"
    echo -e "${BOLD}${RED}╚════════════════════════════════════════════╝${NC}"
    echo
    echo -e "  ${YELLOW}⚠️  This will delete all learned memory!${NC}"
    echo -e "  ${CYAN}The chatbot will restart learning from conversations.txt${NC}"
    echo
    read -p "  Are you sure? (yes/no): " confirm
    
    if [ "$confirm" = "yes" ]; then
        rm -f data/checkpoint.gob
        echo -e "  ${GREEN}✅ Checkpoint deleted!${NC}"
    else
        echo -e "  ${BLUE}Cancelled${NC}"
    fi
    
    read -p "  Press Enter to continue..."
}

# Recompile
recompile() {
    echo
    echo -e "${BOLD}${CYAN}╔════════════════════════════════════════════╗${NC}"
    echo -e "${BOLD}${CYAN}║            🔧  RECOMPILING                 ║${NC}"
    echo -e "${BOLD}${CYAN}╚════════════════════════════════════════════╝${NC}"
    echo
    
    if ! command -v go &> /dev/null; then
        echo -e "  ${RED}❌ Go is not installed!${NC}"
        read -p "  Press Enter to continue..."
        return
    fi
    
    echo -e "  ${BLUE}📦 Building...${NC}"
    rm -f chatbot
    go build -o chatbot ./cmd
    
    if [ $? -eq 0 ]; then
        BIN_SIZE=$(du -h chatbot 2>/dev/null | cut -f1)
        echo -e "  ${GREEN}✅ Success!${NC}"
        echo -e "  💾 Binary size: ${CYAN}$BIN_SIZE${NC}"
    else
        echo -e "  ${RED}❌ Compilation failed!${NC}"
    fi
    
    read -p "  Press Enter to continue..."
}

# Show help
show_help() {
    echo
    echo -e "${BOLD}${CYAN}╔════════════════════════════════════════════╗${NC}"
    echo -e "${BOLD}${CYAN}║              ❓  HELP GUIDE                ║${NC}"
    echo -e "${BOLD}${CYAN}╚════════════════════════════════════════════╝${NC}"
    echo
    echo -e "  ${BOLD}${GREEN}Chatbot Commands:${NC}"
    echo -e "    /quit      ${CYAN}→${NC} Save and exit"
    echo -e "    /save      ${CYAN}→${NC} Save checkpoint manually"
    echo -e "    /reload    ${CYAN}→${NC} Reload conversations from file"
    echo -e "    /stats     ${CYAN}→${NC} Show statistics"
    echo -e "    /temp X    ${CYAN}→${NC} Set temperature (0.1-1.5)"
    echo
    echo -e "  ${BOLD}${GREEN}File Format:${NC}"
    echo -e "    ${YELLOW}User:${NC} your question here"
    echo -e "    ${YELLOW}Bot:${NC} answer here"
    echo -e "    ${CYAN}(blank line between conversations)${NC}"
    echo
    echo -e "  ${BOLD}${GREEN}How Learning Works:${NC}"
    echo -e "    1. Bot searches for exact match in conversations.txt"
    echo -e "    2. If not found, uses LCS similarity (threshold > 0.6)"
    echo -e "    3. If still nothing, asks you to teach"
    echo -e "    4. Your teaching is saved to checkpoint"
    echo
    echo -e "  ${BOLD}${GREEN}Tips:${NC}"
    echo -e "    • More conversations = better answers"
    echo -e "    • Use /reload after editing conversations.txt"
    echo -e "    • Lower temperature = more predictable"
    echo -e "    • Higher temperature = more creative"
    echo
    read -p "  Press Enter to continue..."
}

# Start chatbot
start_chatbot() {
    echo
    echo -e "${BOLD}${GREEN}╔════════════════════════════════════════════╗${NC}"
    echo -e "${BOLD}${GREEN}║           🚀  STARTING CHATBOT            ║${NC}"
    echo -e "${BOLD}${GREEN}╚════════════════════════════════════════════╝${NC}"
    echo
    
    if ! command -v go &> /dev/null; then
        echo -e "  ${RED}❌ Go is not installed!${NC}"
        read -p "  Press Enter to continue..."
        return
    fi
    
    if [ ! -d "data" ]; then
        mkdir -p data
        echo -e "  ${GREEN}✅ Created data directory${NC}"
    fi
    
    if [ ! -f "data/conversations.txt" ]; then
        echo -e "  ${YELLOW}⚠️  Creating example conversations.txt...${NC}"
        cat > data/conversations.txt << 'EOF'
User: hello
Bot: hello welcome to the java chatbot

User: what is java
Bot: java is a programming language created by james gosling in 1991

User: who created java
Bot: james gosling created java at sun microsystems
EOF
        echo -e "  ${GREEN}✅ Created${NC}"
    fi
    
    CONV_COUNT=$(grep -c "^User:" data/conversations.txt 2>/dev/null || echo "0")
    echo -e "  ${GREEN}✓${NC} Found ${CYAN}$CONV_COUNT${NC} conversations"
    
    NEED_COMPILE=0
    if [ ! -f "chatbot" ]; then
        NEED_COMPILE=1
    else
        for file in cmd/main.go internal/app/*.go internal/dataset/*.go internal/model/*.go internal/config/*.go; do
            if [ -f "$file" ] && [ "$file" -nt "chatbot" ] 2>/dev/null; then
                NEED_COMPILE=1
                break
            fi
        done
    fi
    
    if [ $NEED_COMPILE -eq 1 ]; then
        echo -e "  ${BLUE}📦 Compiling...${NC}"
        go build -o chatbot ./cmd
        if [ $? -ne 0 ]; then
            echo -e "  ${RED}❌ Compilation failed${NC}"
            read -p "  Press Enter to continue..."
            return
        fi
        echo -e "  ${GREEN}✅ Compiled${NC}"
    fi
    
    echo
    echo -e "${BOLD}${MAGENTA}════════════════════════════════════════════════${NC}"
    echo -e "${BOLD}${GREEN}  💬 Type /quit to return to menu${NC}"
    echo -e "${BOLD}${MAGENTA}════════════════════════════════════════════════${NC}"
    echo
    
    ./chatbot
    
    echo
    echo -e "  ${GREEN}✅ Returned to menu${NC}"
    read -p "  Press Enter to continue..."
}

# Main loop
while true; do
    show_menu
    read choice
    
    case $choice in
        1) start_chatbot ;;
        2) show_stats ;;
        3) edit_conversations ;;
        4) delete_checkpoint ;;
        5) recompile ;;
        6) show_help ;;
        7) 
            echo
            echo -e "  ${GREEN}👋 Goodbye! Keep coding in Java! ☕${NC}"
            exit 0
            ;;
        *) 
            echo -e "  ${RED}Invalid option!${NC}"
            sleep 1
            ;;
    esac
done