-- Token Bucket Rate Limiter

local key = KEYS[1]
local capacity = tonumber(ARGV[1])          -- Maximum number of tokens in the bucket
local refill_rate = tonumber(ARGV[2])       -- Tokens added per second
local tokens_to_consume = tonumber(ARGV[3]) -- Tokens required for this request
local now = tonumber(ARGV[4])               -- Current timestamp (in seconds)
local bucket_ttl = tonumber(ARGV[5])        -- TTL for Redis key

-- Get current bucket state
local bucket = redis.call('HMGET', key, 'tokens', 'last_refill')
local tokens = tonumber(bucket[1]) or capacity
local last_refill = tonumber(bucket[2]) or now

-- Refill tokens based on elapsed time
local time_passed = now - last_refill
local refill_amount = time_passed * refill_rate
tokens = math.min(capacity, tokens + refill_amount)
last_refill = now

-- Check if enough tokens are available
if tokens < tokens_to_consume then
	redis.call('HSET', key, 'tokens', tokens, 'last_refill', last_refill)
	redis.call('EXPIRE', key, bucket_ttl)
	return {0, tokens}
end

-- Consume tokens and save state
tokens = tokens - tokens_to_consume
redis.call('HSET', key, 'tokens', tokens, 'last_refill', last_refill)
redis.call('EXPIRE', key, bucket_ttl)

return {1, tokens}