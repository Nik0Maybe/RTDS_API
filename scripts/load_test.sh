API_KEY="your-secure-api-key"
SERVER_URL="http://localhost:8080"
NUM_STREAMS=1000
DATA="sample data"

for i in $(seq 1 $NUM_STREAMS); do
    RESPONSE=$(curl -s -X POST "$SERVER_URL/stream/start" -H "X-API-Key: $API_KEY")
    STREAM_ID=$(echo $RESPONSE | jq -r '.message' | awk '{print $NF}')
    
    curl -s -X POST "$SERVER_URL/stream/$STREAM_ID/send" -d "$DATA" -H "X-API-Key: $API_KEY" &
    
    # (Optional) Connect to the results WebSocket (requires WebSocket client like `wscat`)
    # wscat -c ws://localhost:8080/stream/$STREAM_ID/results &
done

echo "Load test initiated with $NUM_STREAMS concurrent streams."
