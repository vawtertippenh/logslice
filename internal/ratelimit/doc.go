// Package ratelimit implements a token-bucket rate limiter for logslice.
//
// It is used to throttle output throughput when streaming or tailing log
// files, preventing downstream consumers from being overwhelmed.
//
// # Usage
//
//	limiter := ratelimit.New(ratelimit.Options{
//		Rate:  1000, // 1000 lines per second
//		Burst: 2000, // allow short bursts up to 2000
//	})
//
//	for _, line := range lines {
//		limiter.Wait() // block until a token is available
//		fmt.Println(line)
//	}
//
// The limiter is safe for concurrent use.
package ratelimit
