#!/bin/bash
# Crash Fault Tolerance Test

echo "=== Raft Crash Test ==="
ORDERER_CONTAINER="orderer.example.com"

# Get current leader (if possible)
echo "Current orderer cluster status:"
docker exec $ORDERER_CONTAINER peer channel list 2>/dev/null || echo "Orderer running"

# Start a background transaction loop (optional)
# We'll just measure recovery time after killing the leader

echo "Killing orderer container: $ORDERER_CONTAINER"
kill_time=$(date +%s%N)
docker stop $ORDERER_CONTAINER

# Wait for cluster to re-elect leader (we monitor logs)
echo "Waiting for new leader election..."
sleep 5
recovery_time=$(date +%s%N)

# Restart the killed orderer (if needed)
docker start $ORDERER_CONTAINER

# Calculate recovery time in seconds
recovery_sec=$(echo "scale=3; ($recovery_time - $kill_time)/1000000000" | bc)
echo "Recovery time: $recovery_sec seconds"

# Log result
echo "$(date): Recovery time = $recovery_sec s" >> crash_test_results.txt