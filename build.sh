#!/bin/bash

# ANSI escape codes for colors and styles
RED='\033[1;31m'       # Bold Red
GREEN='\033[1;32m'     # Bold Green
BOLD='\033[1m'         # Bold
NC='\033[0m'           # No Color

# Array of common operating systems
os=(linux windows darwin)

# Array of common architectures for each OS
linux_arch=(amd64 386 arm arm64)
windows_arch=(amd64 386)
darwin_arch=(amd64 arm64)

# Function to get the valid architectures for a given OS
get_valid_architectures() {
  local os=$1
  local architectures_var="${os}_arch[@]"
  local architectures=("${!architectures_var}")
  echo "${architectures[@]}"
}

get_parent_folder_name() {
    local parent_folder=$(dirname $PWD)
    local parent_folder_name=$(basename $PWD)
    echo "$parent_folder_name"
}

# Function to compile the Go application for a given OS/Architecture
compile_for_platform() {
  local os=$1
  local arch=$2
  local output="build/$3.${os}_${arch}"
  
  # Compile the application
  mkdir -p build
  GOOS=$os GOARCH=$arch go build -o "$output" *.go
  # Show the completion message
  echo -e "\r[ ${GREEN}$(tput bold)DONE$(tput sgr0)${NC} ]  $os/$arch"
}

# Function to show the status of the compilation with an ASCII animation
show_status() {
  local os=$1
  local arch=$2
  local animation="    "  # Initial animation with 4 spaces
  local index=0
  while true; do
    local frame="${animation:0:index}*${animation:index+1}"
    echo -ne "\r[ ${RED}$(tput bold)${frame}${NC} ]  $os/$arch"
    index=$(( (index + 1) % 4 ))  # Ensure animation is contained within 4 characters
    sleep 0.1
  done
}

# Function to print a banner
print_banner() {
    local message="TTKBuilder v1.0"
    echo -e "\033[1;36m$message\033[0m"
}

# Main function to compile for all common platforms
compile_all() {
  for os in "${os[@]}"; do
    architectures=($(get_valid_architectures "$os"))

    for arch in "${architectures[@]}"; do
      compile_for_platform "$os" "$arch" "$1" &
      local pid=$!
      # Show the status with an ASCII animation
      show_status "$os" "$arch" &
      local animation_pid=$!
      wait $pid 2>/dev/null
      kill $animation_pid 2>/dev/null
      wait $animation_pid 2>/dev/null
    done
  done
}

# Check if any Go file exists
if [[ -z $(ls *.go 2>/dev/null) ]]; then
  echo -e "${RED}[ FAIL ]${NC} No *.go files found in the current directory!"
  exit 1
fi

# Run the compilation for all common platforms
print_banner
compile_all ttkagent
