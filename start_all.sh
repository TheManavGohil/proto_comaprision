#!/bin/bash

# Configuration
GO_PORT=8081
JAVA_PORT=8082
NODE_PORT=8083
FRONTEND_PORT=8080

echo "============================================="
echo "🚀 Starting Performance Comparison Servers..."
echo "============================================="

# Keep track of background process PIDs
PIDS=()

# Cleanup function on Ctrl+C (SIGINT)
cleanup() {
    echo -e "\n🛑 Stopping all servers..."
    for pid in "${PIDS[@]}"; do
        if kill -0 "$pid" 2>/dev/null; then
            kill "$pid" 2>/dev/null
        fi
    done
    echo "Done. All servers stopped."
    exit 0
}

trap cleanup SIGINT

# 1. Compile Java Server
echo "☕ Compiling Java Server..."
javac JavaServer.java
if [ $? -ne 0 ]; then
    echo "❌ Java compilation failed. Exiting."
    exit 1
fi

# 2. Start Go Server
echo "🚀 Starting Go Server on port $GO_PORT..."
go run go_server.go > go_server.log 2>&1 &
PIDS+=($!)

# 3. Start Java Server
echo "☕ Starting Java Server on port $JAVA_PORT..."
java JavaServer > java_server.log 2>&1 &
PIDS+=($!)

# 4. Start Node.js Server
echo "⬢ Starting Node.js Server on port $NODE_PORT..."
node node_server.js > node_server.log 2>&1 &
PIDS+=($!)

# 5. Start Frontend HTTP Server
echo "🌐 Starting Frontend Server on port $FRONTEND_PORT..."
python3 -m http.server $FRONTEND_PORT > frontend_server.log 2>&1 &
PIDS+=($!)

# Wait a moment for servers to initialize
sleep 2

# Check if servers are running
echo "---------------------------------------------"
echo "📊 Checking Server Status:"

if curl -s http://localhost:$GO_PORT/api/stats >/dev/null; then
    echo "  [OK] Go Server: http://localhost:$GO_PORT"
else
    echo "  [FAIL] Go Server failed to start! Check go_server.log"
fi

# Send options request or small post request to Java server since it only accepts POST/OPTIONS on api endpoints
if curl -s -X OPTIONS http://localhost:$JAVA_PORT/api/factorial >/dev/null; then
    echo "  [OK] Java Server: http://localhost:$JAVA_PORT"
else
    echo "  [FAIL] Java Server failed to start! Check java_server.log"
fi

if curl -s http://localhost:$NODE_PORT/api/stats >/dev/null; then
    echo "  [OK] Node.js Server: http://localhost:$NODE_PORT"
else
    echo "  [FAIL] Node.js Server failed to start! Check node_server.log"
fi

# Check frontend
if curl -s http://localhost:$FRONTEND_PORT >/dev/null; then
    echo "  [OK] Frontend Dashboard: http://localhost:$FRONTEND_PORT"
else
    echo "  [FAIL] Frontend Server failed to start! Check frontend_server.log"
fi

echo "---------------------------------------------"
echo "🎉 All servers started! Access the comparison at: http://localhost:$FRONTEND_PORT"
echo "📝 Logs are saved to: go_server.log, java_server.log, node_server.log, frontend_server.log"
echo "Press [Ctrl+C] to stop all servers."
echo "============================================="

# Keep script running to trap Ctrl+C
wait
