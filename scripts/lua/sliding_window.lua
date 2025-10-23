--- Sliding window algorithm lua
local key = KEYS[1]
local window = tonumber(ARGV[1])
local limit = tonumber(ARGV[2])
local now = tonumber(ARGV[3])
local request_id = ARGV[4]

-- Remove old timestamps
redis.call('ZREMRANGEBYSCORE', key, 0, now - window)

-- Count current window
local count = redis.call('ZCARD', key)

if count < limit then
    redis.call('ZADD', key, now, request_id)
    redis.call('EXPIRE', key, math.ceil(window / 1000))
    return {1, count + 1}
else
    return {0, count}
end