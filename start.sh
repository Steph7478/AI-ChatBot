#!/bin/bash

# ============================================
# NEURAL CHATBOT - INTERACTIVE MENU
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

# ============================================
# Menu
# ============================================

show_menu() {
    clear
    echo -e "${BOLD}${MAGENTA}┌─────────────────────────────────────────────────┐${NC}"
    echo -e "${BOLD}${MAGENTA}│           🧠  NEURAL CHATBOT  🧠                │${NC}"
    echo -e "${BOLD}${MAGENTA}└─────────────────────────────────────────────────┘${NC}"
    echo
    echo -e "  ${CYAN}📌 Main Menu${NC}"
    echo
    echo -e "    ${GREEN}1)${NC} 🚀 Start Chatbot"
    echo -e "    ${GREEN}2)${NC} 📊 Statistics"
    echo -e "    ${GREEN}3)${NC} 🗑️ Delete Model"
    echo -e "    ${GREEN}4)${NC} 🧠 Train Network"
    echo -e "    ${GREEN}5)${NC} 🔧 Recompile"
    echo -e "    ${GREEN}6)${NC} ❓ Help"
    echo -e "    ${GREEN}7)${NC} 👋 Exit"
    echo
    echo -ne "${YELLOW}➜ Choose: ${NC}"
}

# ============================================
# Statistics
# ============================================

show_stats() {
    clear
    echo -e "${BOLD}${CYAN}┌─────────────────────────────────────────────────┐${NC}"
    echo -e "${BOLD}${CYAN}│                    📊 STATS                     │${NC}"
    echo -e "${BOLD}${CYAN}└─────────────────────────────────────────────────┘${NC}"
    echo
    
    if [ -f "data/conversations.txt" ]; then
        CONV_COUNT=$(grep -c "|" data/conversations.txt 2>/dev/null || echo "0")
        CONV_SIZE=$(du -h data/conversations.txt 2>/dev/null | cut -f1)
        echo -e "  ${GREEN}✓${NC} conversations.txt"
        echo -e "    📝 ${CONV_COUNT} entries  |  💾 ${CONV_SIZE}"
    fi
    
    if [ -f "data/training_data.txt" ]; then
        TRAIN_COUNT=$(grep -c "|" data/training_data.txt 2>/dev/null || echo "0")
        echo -e "  ${GREEN}✓${NC} training_data.txt"
        echo -e "    📝 ${TRAIN_COUNT} entries"
    fi
    
    if [ -f "data/model.gob" ]; then
        MODEL_SIZE=$(du -h data/model.gob 2>/dev/null | cut -f1)
        echo -e "  ${GREEN}✓${NC} model.gob"
        echo -e "    💾 ${MODEL_SIZE}"
    fi
    
    if [ -f "data/model.gob.vocab" ]; then
        VOCAB_SIZE=$(du -h data/model.gob.vocab 2>/dev/null | cut -f1)
        echo -e "  ${GREEN}✓${NC} vocabulary"
        echo -e "    📖 ${VOCAB_SIZE}"
    fi
    
    echo
    read -p "  Press Enter..."
}

# ============================================
# Delete Model
# ============================================

delete_model() {
    clear
    echo -e "${BOLD}${RED}┌─────────────────────────────────────────────────┐${NC}"
    echo -e "${BOLD}${RED}│                  🗑️ DELETE MODEL                │${NC}"
    echo -e "${BOLD}${RED}└─────────────────────────────────────────────────┘${NC}"
    echo
    echo -e "  ${YELLOW}⚠️  This will erase all learning!${NC}"
    echo
    read -p "  Are you sure? (yes/no): " confirm
    
    if [ "$confirm" = "yes" ]; then
        rm -f data/model.gob data/model.gob.vocab
        echo -e "  ${GREEN}✅ Deleted${NC}"
    else
        echo -e "  ${BLUE}❌ Cancelled${NC}"
    fi
    
    echo
    read -p "  Press Enter..."
}

# ============================================
# Train
# ============================================

train_model() {
    clear
    echo -e "${BOLD}${CYAN}┌─────────────────────────────────────────────────┐${NC}"
    echo -e "${BOLD}${CYAN}│                  🧠 TRAIN NETWORK               │${NC}"
    echo -e "${BOLD}${CYAN}└─────────────────────────────────────────────────┘${NC}"
    echo
    
    if [ ! -f "data/training_data.txt" ]; then
        echo -e "  ${RED}❌ No training_data.txt found!${NC}"
        read -p "  Press Enter..."
        return
    fi
    
    TRAIN_COUNT=$(grep -c "|" data/training_data.txt 2>/dev/null || echo "0")
    echo -e "  📚 ${TRAIN_COUNT} training examples"
    echo
    read -p "  Epochs: " epochs
    
    if ! [[ "$epochs" =~ ^[0-9]+$ ]]; then
        echo -e "  ${RED}❌ Invalid${NC}"
        read -p "  Press Enter..."
        return
    fi
    
    echo
    echo -e "  ${BLUE}Training...${NC}"
    echo
    
    ./chatbot -train "$epochs"
    
    echo
    read -p "  Press Enter..."
}

# ============================================
# Recompile
# ============================================

recompile() {
    clear
    echo -e "${BOLD}${CYAN}┌─────────────────────────────────────────────────┐${NC}"
    echo -e "${BOLD}${CYAN}│                  🔧 RECOMPILE                   │${NC}"
    echo -e "${BOLD}${CYAN}└─────────────────────────────────────────────────┘${NC}"
    echo
    
    if ! command -v go &> /dev/null; then
        echo -e "  ${RED}❌ Go not installed${NC}"
        read -p "  Press Enter..."
        return
    fi
    
    echo -e "  ${BLUE}Building...${NC}"
    rm -f chatbot
    go build -o chatbot ./cmd
    
    if [ $? -eq 0 ]; then
        SIZE=$(du -h chatbot 2>/dev/null | cut -f1)
        echo -e "  ${GREEN}✅ Success (${SIZE})${NC}"
    else
        echo -e "  ${RED}❌ Failed${NC}"
    fi
    
    echo
    read -p "  Press Enter..."
}

# ============================================
# Help
# ============================================

show_help() {
    clear
    echo -e "${BOLD}${CYAN}┌─────────────────────────────────────────────────┐${NC}"
    echo -e "${BOLD}${CYAN}│                    ❓ HELP                      │${NC}"
    echo -e "${BOLD}${CYAN}└─────────────────────────────────────────────────┘${NC}"
    echo
    echo -e "  ${GREEN}Chat Commands:${NC}"
    echo -e "    /quit      Exit"
    echo -e "    /save      Save model"
    echo -e "    /stats     Show stats"
    echo -e "    /temp X    Set temperature"
    echo
    echo -e "  ${GREEN}Files:${NC}"
    echo -e "    conversations.txt  → Quick responses"
    echo -e "    training_data.txt  → Neural training"
    echo
    echo -e "  ${GREEN}Tips:${NC}"
    echo -e "    Lower temp = more predictable"
    echo -e "    Higher temp = more creative"
    echo
    read -p "  Press Enter..."
}

# ============================================
# Start Chatbot
# ============================================

start_chatbot() {
    clear
    echo -e "${BOLD}${GREEN}┌─────────────────────────────────────────────────┐${NC}"
    echo -e "${BOLD}${GREEN}│                  🚀 STARTING                    │${NC}"
    echo -e "${BOLD}${GREEN}└─────────────────────────────────────────────────┘${NC}"
    echo
    
    if ! command -v go &> /dev/null; then
        echo -e "  ${RED}❌ Go not installed${NC}"
        read -p "  Press Enter..."
        return
    fi
    
    # Create data dir if needed
    if [ ! -d "data" ]; then
        mkdir -p data
    fi
    
    # Create default conversations if needed
    if [ ! -f "data/conversations.txt" ]; then
        echo "hi|Hello!" > data/conversations.txt
        echo "what is java|Java is a programming language" >> data/conversations.txt
    fi
    
    # Compile if needed
    if [ ! -f "chatbot" ] || [ "cmd/main.go" -nt "chatbot" ]; then
        echo -e "  ${BLUE}Compiling...${NC}"
        go build -o chatbot ./cmd
        if [ $? -ne 0 ]; then
            echo -e "  ${RED}❌ Compilation failed${NC}"
            read -p "  Press Enter..."
            return
        fi
    fi
    
    clear
    echo -e "${BOLD}${MAGENTA}┌─────────────────────────────────────────────────┐${NC}"
    echo -e "${BOLD}${MAGENTA}│              🤖 CHATBOT READY 🤖                │${NC}"
    echo -e "${BOLD}${MAGENTA}└─────────────────────────────────────────────────┘${NC}"
    echo
    
    ./chatbot
    
    echo
    read -p "  Press Enter..."
}

# ============================================
# Main
# ============================================

while true; do
    show_menu
    read choice
    
    case $choice in
        1) start_chatbot ;;
        2) show_stats ;;
        3) delete_model ;;
        4) train_model ;;
        5) recompile ;;
        6) show_help ;;
        7) 
            echo
            echo -e "  ${GREEN}👋 Bye!${NC}"
            exit 0
            ;;
        *) 
            echo -e "  ${RED}Invalid${NC}"
            sleep 1
            ;;
    esac
done