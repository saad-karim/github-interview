# Setup & Code Segment

This exercise is meant to be a paired programming activity
that's open book, open notes. So feel free to have a conversation with your interviewers, google commands, or use any libraries you think will best solve this task.

Scenario:
> We have a CSV file containing movie metadata that we want to "import" into a web service using a supplied endpoint. This endpoint only accepts one object at a time.

> We don't know how accurate the data is so we'll need to validate it. But we still want to send as much valid data as possible into the server.

> Our priority is to get a working import solution first. If we have time, we can optimize process. We'll have plenty of time to discuss improvements during a retrospective.

## Setup

- Please clone this repo and share your screen
- Run `bin/setup`
- Make a test POST call:

  ```bash
  curl localhost:9009/movies -d '{"year":1997, "length": 123, "title": "Face Off", "subject": "action", "actor": "Cage, Nicholas", "actress": "Allen, Joan", "director": "Woo, John", "popularity": 82, "awards": "No", "image": "NicholasCage.png"}'
  ```

- (optional) Check monitoring endpoint: `curl localhost:9009/metrics`
  - You should see 1 count for `sink_get_total` and `sink_post_total`
  - Reset the metrics count by calling `bin/reset` in the same window as the log

## Helper Commands

| Command | Description |
| --- | --- |
| `bin/setup` | Builds and start the service that you'll be writing a script against |
| `bin/reset` | Restarts the service to reset prometheus metrics |
| `bin/log` | Display logs from server |
| `bin/clean` | Stops and remove docker containers |
| `bin/destroy` | Destroy all local docker artifacts. *Use with caution* |
