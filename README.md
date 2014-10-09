# go-resque

Simple [Resque](https://github.com/defunkt/resque) queue client for [Go](http://golang.org).

## Installation

Installation is simple and familiar for Go programmers:

```
go get github.com/kavu/go-resque
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
  "github.com/kavu/go-resque" // Import this package
  _ "github.com/kavu/go-resque/godis" // Use godis driver
  "github.com/simonz05/godis/redis" // Redis client from godis package
)

func main() {
  var err error

  client := redis.New("tcp:127.0.0.1:6379", 0, "") // Create new Redis client to use for enqueuing
  enqueuer := resque.NewRedisEnqueuer("godis", client) // Create enqueuer instance

  // Enqueue the job into the "go" queue with appropriate client
  _, err = enqueuer.Enqueue("resque:queue:go", "Demo::Job")
  if err != nil {
    panic(err)
  }

  // Enqueue into the "default" queue with passing one parameter to the Demo::Job.perform
  _, err = enqueuer.Enqueue("resque:queue:default", "Demo::Job", 1)
  if err != nil {
    panic(err)
  }

  // Enqueue into the "extra" queue with passing multiple
  // parameters to the Demo::Job.perform so it will fail
  _, err = enqueuer.Enqueue("resque:queue:extra", "Demo::Job", 1, 2, "woot")
  if err != nil {
    panic(err)
  }

}
```

Simple enough? I hope so.

## Contributing

Just open pull request or ping me directly on e-mail, if you want to discuss some ideas.
