# NATS Publisher Example

This is a simple example of how to use NATS Publisher.

## How to run?

Firtly, you must have your NATS Server running. The easiest way is by running it on Docker:

```bash
$ docker run -d --name nats-main --rm -p 4222:4222 -p 6222:6222 -p 8222:8222 nats
```

Now that your NATS Server is running, it's time to build this code:

```bash
$ go build .
```

And finally, you can run your publisher:

```bash
$ ./nats-pub
```

It will send to NATS a "Hello World" message.
