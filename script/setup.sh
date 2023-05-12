#!/usr/bin/env bash

SCRIPT_PATH=$(realpath "$0")
cd $SCRIPT_PATH
cd ../../chat-apps
echo "Script path is: $(pwd)"

# Define the function for setting up the application
setup() {
    destroy
    mkdir -p logs
    docker-compose -f $(pwd)/docker/docker-compose.yml -p docker up  > logs/docker.log 2>&1 &
    echo "Application has been set up."

    echo "Run:"
    echo "tail -f $(pwd)/logs/docker.log"
}

# Define the function for destroying the application
destroy() {
    docker-compose -f $(pwd)/docker/docker-compose.yml -p docker down
    docker rmi -f chat-apps:latest
    echo "Application has been destroyed."
}

# Call the appropriate function based on the first argument passed to the script
# Determine which function to call based on the first argument
case $1 in
  setup)
    setup
    ;;
  destroy)
    destroy
    ;;
  *)
    echo "Invalid argument. Usage: ./script.sh [setup|destroy]"
    exit 1
    ;;
esac
