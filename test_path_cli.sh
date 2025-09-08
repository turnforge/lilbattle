#!/bin/bash

# Test the path display in CLI
cd cmd/weewar-cli

# Try with the test game
echo "Testing path display with game 32112070..."
echo -e "options -1,-2\nquit" | ./weewar-cli 32112070