local counter = 0

function setup()
    local res = wrk.format("POST", "/stream/start", {["X-API-Key"] = "your-secure-api-key"})
    return res
end

request = function()
    counter = counter + 1
    local stream_id = "<your_stream_id>"
    local path = "/stream/" .. stream_id .. "/send"
    local headers = {["X-API-Key"] = "your-secure-api-key"}
    local body = "sample data " .. counter
    return wrk.format("POST", path, headers, body)
end

done = function(summary, latency, requests)
    print("Complete")
    print("Requests: " .. requests)
    print("Latency: " .. latency.mean / 1000 .. "ms")
end
