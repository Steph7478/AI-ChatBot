#!/bin/bash

# Colors
GREEN='\033[0;32m'
CYAN='\033[0;36m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
BOLD='\033[1m'
DIM='\033[2m'
NC='\033[0m'

show_menu() {
    clear
    echo -e "${BOLD}${CYAN}"
    echo "    ┌─────────────────────────────────┐"
    echo "    │      🧠  NEURAL CHATBOT  🧠     │"
    echo "    └─────────────────────────────────┘"
    echo -e "${NC}"
    echo
    echo -e "  ${BOLD}${GREEN}1${NC}  🚀 Start Chatbot"
    echo -e "  ${BOLD}${GREEN}2${NC}  📊 Statistics"
    echo -e "  ${BOLD}${GREEN}3${NC}  🗑️ Delete Model"
    echo -e "  ${BOLD}${GREEN}4${NC}  🧠 Train Network"
    echo -e "  ${BOLD}${GREEN}5${NC}  🔧 Recompile"
    echo -e "  ${BOLD}${GREEN}6${NC}  ❓ Help"
    echo -e "  ${BOLD}${GREEN}7${NC}  👋 Exit"
    echo
    echo -ne "  ${BOLD}${YELLOW}➜  ${NC}"
}

show_stats() {
    clear
    echo -e "${BOLD}${CYAN}    📊 STATISTICS${NC}\n"
    
    [ -f "data/conversations.txt" ] && echo -e "  ${GREEN}✓${NC} conversations: ${CYAN}$(grep -c "|" data/conversations.txt 2>/dev/null || echo 0)${NC}"
    [ -f "data/training_data.txt" ] && echo -e "  ${GREEN}✓${NC} training: ${CYAN}$(grep -c "|" data/training_data.txt 2>/dev/null || echo 0)${NC}"
    [ -f "data/model.gob" ] && echo -e "  ${GREEN}✓${NC} model: ${CYAN}$(du -h data/model.gob 2>/dev/null | cut -f1)${NC}"
    [ -f "data/model.gob.vocab" ] && echo -e "  ${GREEN}✓${NC} vocab: ${CYAN}$(du -h data/model.gob.vocab 2>/dev/null | cut -f1)${NC}"
    
    echo
    read -p "  Press Enter..."
}

delete_model() {
    clear
    echo -e "${BOLD}${CYAN}    🗑️ DELETE MODEL${NC}\n"
    echo -e "  ${YELLOW}⚠️  This will erase all learning${NC}\n"
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

train_model() {
    clear
    echo -e "${BOLD}${CYAN}    🧠 TRAIN NETWORK${NC}\n"
    
    if [ ! -f "data/training_data.txt" ]; then
        echo -e "  ${RED}❌ No training_data.txt${NC}"
        read -p "  Press Enter..."
        return
    fi
    
    TRAIN_COUNT=$(grep -c "|" data/training_data.txt 2>/dev/null || echo "0")
    echo -e "  📚 ${TRAIN_COUNT} examples\n"
    read -p "  Epochs: " epochs
    
    [[ ! "$epochs" =~ ^[0-9]+$ ]] && echo -e "  ${RED}❌ Invalid${NC}" && read -p "  Press Enter..." && return
    
    echo -e "\n  ${BLUE}Training...${NC}\n"
    ./chatbot -train "$epochs"
    echo
    read -p "  Press Enter..."
}

recompile() {
    clear
    echo -e "${BOLD}${CYAN}    🔧 RECOMPILE${NC}\n"
    
    if ! command -v go &> /dev/null; then
        echo -e "  ${RED}❌ Go not installed${NC}"
        read -p "  Press Enter..."
        return
    fi
    
    echo -e "  ${BLUE}Building...${NC}"
    go build -o chatbot ./cmd
    
    if [ $? -eq 0 ]; then
        echo -e "  ${GREEN}✅ Success${NC}"
    else
        echo -e "  ${RED}❌ Failed${NC}"
    fi
    echo
    read -p "  Press Enter..."
}

show_help() {
    clear
    echo -e "${BOLD}${CYAN}    ❓ HELP${NC}\n"
    echo -e "  ${GREEN}Commands:${NC}"
    echo "    /quit      Exit"
    echo "    /save      Save model"
    echo "    /stats     Show stats"
    echo "    /temp X    Set temperature (0.1-2.0)"
    echo
    echo -e "  ${GREEN}Files:${NC}"
    echo "    conversations.txt  → Quick responses"
    echo "    training_data.txt  → Neural training"
    echo
    echo -e "  ${GREEN}Tips:${NC}"
    echo "    Lower temp = more predictable"
    echo "    Higher temp = more creative"
    echo
    read -p "  Press Enter..."
}

start_chatbot() {
    if ! command -v go &> /dev/null; then
        echo -e "  ${RED}❌ Go not installed${NC}"
        read -p "  Press Enter..."
        return
    fi
    
    mkdir -p data
    [ ! -f "data/conversations.txt" ] && echo "hi|Hello!" > data/conversations.txt
    
    if [ ! -f "chatbot" ] || [ "cmd/main.go" -nt "chatbot" ]; then
        echo -e "  ${BLUE}Compiling...${NC}"
        go build -o chatbot ./cmd
        [ $? -ne 0 ] && echo -e "  ${RED}❌ Failed${NC}" && read -p "  Press Enter..." && return
    fi
    
    clear
    echo -e "${BOLD}${CYAN}    🤖 CHATBOT READY${NC}\n"
    ./chatbot
    echo
    read -p "  Press Enter..."
}

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
        7) echo -e "\n  ${GREEN}👋 Bye!${NC}" && exit 0 ;;
        *) echo -e "  ${RED}❌ Invalid${NC}" && sleep 1 ;;
    esac
done