# go-resque

Simple [Resque](https://github.com/resque/resque) queue client for [Go](http://golang.org).

## Introduction

This is a fork of [go-resque](https://github.com/kavu/go-resque) (one of many).

Differences from the original are:
- it is a Go module
- travis config is updated to use newer Go versions
- all drivers except one removed
- this driver is [go-redis](https://github.com/go-redis/redis) v9
- driver update required to change some signatures to include context
- also, in newer Redis' versions some commands are deprecated, and I used suggested replacements (namely, ZRANGE instead of ZRANGEBYSCORE)

## Installation

Installation is simple and familiar for Go programmers:

```
go get github.com/skaurus/go-resque
```

## Usage

Let's assume that you have such Resque Job (taken from Resque examples):

```ruby
module Demo
  class Job
    def self.perform(params)
      puts "Processed a job!"
    end
  end
end
```

So, we can enqueue this job from Go.

```go
package main

import (
	"context"
	"github.com/skaurus/go-resque"            // Import this package
	_ "github.com/skaurus/go-resque/redis.v9" // Use go-redis v9
	"github.com/go-redis/redis"               // Redis client from go-redis package
)

func main() {
	var err error

	// Create new Redis client to use for enqueuing
	client := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
		DB:   0,
	})
	// Create enqueuer instance
	enqueuer := resque.NewRedisEnqueuer(context.Background(), "redis.v9", client, "resque:")

	// Enqueue the job into the "go" queue with appropriate client
	_, err = enqueuer.Enqueue(context.Background(), "go", "Demo::Job")
	if err != nil {
		panic(err)
	}

	// Enqueue into the "default" queue with passing one parameter to the Demo::Job.perform
	_, err = enqueuer.Enqueue(context.Background(), "default", "Demo::Job", 1)
	if err != nil {
		panic(err)
	}

	// Enqueue into the "extra" queue with passing multiple
	// parameters to the Demo::Job.perform so it will fail
	_, err = enqueuer.Enqueue(context.Background(), "extra", "Demo::Job", 1, 2, "woot")
	if err != nil {
		panic(err)
	}

}
```

Simple enough? I hope so.

## Contributing

Just open pull request or ping me directly on e-mail, if you want to discuss some ideas.
