#!/bin/bash

# ============================================
# 🧠 Neural Chatbot - Interactive Menu
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
    echo "    ║     🧠  NEURAL CHATBOT - MAIN MENU  🧠       ║"
    echo "    ╚══════════════════════════════════════════════╝"
    echo -e "${NC}"
    echo
    echo -e "  ${CYAN}📌${NC} Choose an option:"
    echo
    echo -e "    ${GREEN}1)${NC} 🚀 Start Chatbot"
    echo -e "    ${GREEN}2)${NC} 📊 Show Statistics"
    echo -e "    ${GREEN}3)${NC} ✏️  Edit conversations.txt"
    echo -e "    ${GREEN}4)${NC} 🗑️  Delete Model (reset learning)"
    echo -e "    ${GREEN}5)${NC} 🧠 Train Neural Network"
    echo -e "    ${GREEN}6)${NC} 🔧 Recompile"
    echo -e "    ${GREEN}7)${NC} ❓ Show Help"
    echo -e "    ${GREEN}8)${NC} 👋 Exit"
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
        CONV_COUNT=$(grep -c "|" data/conversations.txt 2>/dev/null || echo "0")
        FILE_SIZE=$(du -h data/conversations.txt 2>/dev/null | cut -f1)
        
        echo -e "  ${GREEN}✓${NC} ${BOLD}conversations.txt${NC}"
        echo -e "    📝 Conversations: ${CYAN}$CONV_COUNT${NC}"
        echo -e "    💾 Size: ${CYAN}$FILE_SIZE${NC}"
        
        if [ $CONV_COUNT -gt 0 ]; then
            echo
            echo -e "  ${CYAN}📋 Last conversations:${NC}"
            tail -5 data/conversations.txt | grep "|" | while IFS='|' read -r input response; do
                echo -e "    ${GREEN}→${NC} ${input:0:40}"
                echo -e "    ${BLUE}←${NC} ${response:0:40}..."
            done
        fi
    else
        echo -e "  ${RED}✗${NC} conversations.txt ${RED}not found${NC}"
        echo -e "    ${YELLOW}💡 Create data/conversations.txt with:${NC}"
        echo -e "       what is java|Java is a programming language"
    fi
    
    echo
    
    if [ -f "data/model.gob" ]; then
        MODEL_SIZE=$(du -h data/model.gob 2>/dev/null | cut -f1)
        echo -e "  ${GREEN}✓${NC} ${BOLD}model.gob${NC}"
        echo -e "    💾 Size: ${CYAN}$MODEL_SIZE${NC}"
    fi
    
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
        echo -e "  ${YELLOW}💡 Creating with example...${NC}"
        mkdir -p data
        echo "hi|Hello! How can I help you?" > data/conversations.txt
        echo "what is java|Java is a programming language" >> data/conversations.txt
    fi
    
    echo -e "  ${CYAN}📋 Last conversations:${NC}\n"
    tail -5 data/conversations.txt | grep "|" | while IFS='|' read -r input response; do
        echo -e "    ${GREEN}Q:${NC} ${input:0:50}"
        echo -e "    ${BLUE}A:${NC} ${response:0:50}"
        echo
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
    
    echo -e "\n  ${GREEN}✅ Done!${NC}"
    read -p "  Press Enter to continue..."
}

# Delete model
delete_model() {
    echo
    echo -e "${BOLD}${RED}╔════════════════════════════════════════════╗${NC}"
    echo -e "${BOLD}${RED}║           🗑️  DELETE MODEL                 ║${NC}"
    echo -e "${BOLD}${RED}╚════════════════════════════════════════════╝${NC}"
    echo
    echo -e "  ${YELLOW}⚠️  This will delete all learned neural weights!${NC}"
    echo -e "  ${CYAN}The model will restart with random weights${NC}"
    echo
    read -p "  Are you sure? (yes/no): " confirm
    
    if [ "$confirm" = "yes" ]; then
        rm -f data/model.gob
        echo -e "  ${GREEN}✅ Model deleted!${NC}"
    else
        echo -e "  ${BLUE}Cancelled${NC}"
    fi
    
    read -p "  Press Enter to continue..."
}

# Train model
# Train model
train_model() {
    echo
    echo -e "${BOLD}${CYAN}╔════════════════════════════════════════════╗${NC}"
    echo -e "${BOLD}${CYAN}║           🧠  TRAIN NEURAL NETWORK         ║${NC}"
    echo -e "${BOLD}${CYAN}╚════════════════════════════════════════════╝${NC}"
    echo
    
    if [ ! -f "data/conversations.txt" ]; then
        echo -e "  ${RED}❌ No conversations.txt found!${NC}"
        echo -e "  ${YELLOW}💡 Please add some conversations first${NC}"
        read -p "  Press Enter to continue..."
        return
    fi
    
    CONV_COUNT=$(grep -c "|" data/conversations.txt 2>/dev/null || echo "0")
    if [ "$CONV_COUNT" -eq 0 ]; then
        echo -e "  ${RED}❌ No valid conversations found!${NC}"
        echo -e "  ${YELLOW}💡 Add lines with format: question|answer${NC}"
        read -p "  Press Enter to continue..."
        return
    fi
    
    echo -e "  ${GREEN}✓${NC} Found ${CYAN}$CONV_COUNT${NC} conversations"
    echo
    echo -e "  ${YELLOW}How many epochs? (recommended: 10-100)${NC}"
    echo -e "  ${CYAN}💡 More epochs = better learning but slower${NC}"
    read -p "  ➜ " epochs
    
    if ! [[ "$epochs" =~ ^[0-9]+$ ]]; then
        echo -e "  ${RED}❌ Invalid number!${NC}"
        read -p "  Press Enter to continue..."
        return
    fi
    
    echo
    echo -e "  ${BLUE}🧠 Training for $epochs epochs...${NC}"
    echo -e "  ${YELLOW}This may take a while...${NC}"
    echo
    
    ./chatbot -train "$epochs"

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
        echo -e "  ${YELLOW}Install: https://golang.org/dl/${NC}"
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
    echo -e "    /quit      ${CYAN}→${NC} Exit chatbot"
    echo -e "    /save      ${CYAN}→${NC} Save model"
    echo -e "    /stats     ${CYAN}→${NC} Show model statistics"
    echo -e "    /temp X    ${CYAN}→${NC} Set temperature (0.1-1.5)"
    echo
    echo -e "  ${BOLD}${GREEN}Training Options:${NC}"
    echo -e "    ${CYAN}Option 5 in menu${NC} → Train neural network with X epochs"
    echo -e "    ${CYAN}./chatbot -train 50${NC} → Train from command line"
    echo
    echo -e "  ${BOLD}${GREEN}Model Architecture:${NC}"
    echo -e "    • Transformer with Multi-Head Attention"
    echo -e "    • Vocab: 10,000 | Embedding: 128-dim"
    echo -e "    • Hidden: 256-dim | Heads: 4 | Layers: 2"
    echo
    echo -e "  ${BOLD}${GREEN}File Format (conversations.txt):${NC}"
    echo -e "    ${YELLOW}question|answer${NC}"
    echo -e "    ${CYAN}Example: what is Java|Java is a programming language${NC}"
    echo
    echo -e "  ${BOLD}${GREEN}How Learning Works:${NC}"
    echo -e "    1. Neural network generates response"
    echo -e "    2. If confidence is low, asks for teaching"
    echo -e "    3. Teaching is saved to conversations.txt (pipe format)"
    echo -e "    4. Run training to update neural weights"
    echo
    echo -e "  ${BOLD}${GREEN}Tips:${NC}"
    echo -e "    • More conversations = better responses"
    echo -e "    • Train after adding new conversations"
    echo -e "    • Lower temperature = more predictable"
    echo -e "    • Higher temperature = more creative"
    echo
    read -p "  Press Enter to continue..."
}

# Start chatbot
start_chatbot() {
    echo
    echo -e "${BOLD}${GREEN}╔════════════════════════════════════════════╗${NC}"
    echo -e "${BOLD}${GREEN}║           🚀  STARTING CHATBOT             ║${NC}"
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
        echo -e "  ${RED}❌ conversations.txt not found!${NC}"
        echo -e "  ${YELLOW}💡 Creating default file...${NC}"
        echo "# Format: question|answer" > data/conversations.txt
        echo "hi|Hello! How can I help you?" >> data/conversations.txt
        echo "what is your name|I'm Neural Chatbot!" >> data/conversations.txt
        echo -e "  ${GREEN}✅ Created with examples${NC}"
    fi
    
    CONV_COUNT=$(grep -c "|" data/conversations.txt 2>/dev/null | grep -v "^0$" || echo "0")
    echo -e "  ${GREEN}✓${NC} Found ${CYAN}$CONV_COUNT${NC} conversations"
    
    NEED_COMPILE=0
    if [ ! -f "chatbot" ]; then
        NEED_COMPILE=1
    else
        if [ -f "cmd/main.go" ] && [ "cmd/main.go" -nt "chatbot" ] 2>/dev/null; then
            NEED_COMPILE=1
        fi
        for file in internal/*/*.go; do
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
    
    clear
    echo -e "${BOLD}${MAGENTA}"
    echo "    ╔══════════════════════════════════════════════╗"
    echo "    ║    ☕  NEURAL CHATBOT - READY TO CHAT  ☕    ║"
    echo "    ╚══════════════════════════════════════════════╝"
    echo -e "${NC}"
    echo
    echo -e "  ${CYAN}💡 Commands:${NC} /quit ${GREEN}|${NC} /save ${GREEN}|${NC} /stats ${GREEN}|${NC} /temp [0.1-1.5]"
    echo
    echo -e "${BOLD}${MAGENTA}════════════════════════════════════════════════${NC}"
    echo -e "${BOLD}${GREEN}  💬 Start typing your messages below${NC}"
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
        4) delete_model ;;
        5) train_model ;;
        6) recompile ;;
        7) show_help ;;
        8) 
            echo
            echo -e "  ${GREEN}👋 Goodbye! Keep coding! 🧠☕${NC}"
            exit 0
            ;;
        *) 
            echo -e "  ${RED}Invalid option!${NC}"
            sleep 1
            ;;
    esac
done