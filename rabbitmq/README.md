# RabbitMQ Example

This is a basic example of [RabbitMQ](https://www.rabbitmq.com/) usage.

The publisher `cmd/pub` sends 10,000 messages to the queue.
The subscriber `cmd/sub` consumes all the messages from the queue and print the elapsed time.

## Dependencies

It is necessary to have a RabbitMQ instance running on your computer.
The easiest way to do it is running it on Docker:

```
docker run -d --rm --name rabbitmq -p 15672:15672 -p 5672:5672 rabbitmq:management
```

Note that this RabbitMQ image already has the [Managemant Plugin](https://www.rabbitmq.com/management.html) enabled.