#!/bin/bash

# Run the 'seed' main.go in the background
go run /cmd/seed/main.go &
# Store the process ID (PID) of the 'seed' process
seed_pid=$!
# Wait for the 'seed' process to finish
wait $seed_pid
go run /cmd/consume/main.go &
go run /cmd/api/main.go
