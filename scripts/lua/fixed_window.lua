-- LUA SCRIPT FOR FIXED WINDOW REDIS
local key = KEYS[1]
local window = tonumber(ARGV[1]) -- window time
local limit = tonumber(ARGV[2]) -- limit count

-- increment counter
local current = redis.call("INCR", key)

-- if it was not set then set expire for key, ttl: window time
if current == 1 then
    redis.call("EXPIRE", key, window)
end

-- check for limit
if current > limit then
    return 0 -- return fail
else
    return 1 -- return success
end
