#!/bin/bash
echo "Starting Discussion Platform..."

# Start backend
cd backend && go run cmd/auth/main.go &
BACKEND_PID=$!

# Start frontend
cd frontend && npm run dev &
FRONTEND_PID=$!

# Handle ctrl+c — kill both
trap "kill $BACKEND_PID $FRONTEND_PID" SIGINT SIGTERM

wait
