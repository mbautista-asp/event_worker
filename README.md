## event_worker

### Install

* Get a .env file from Charlie
* Make sure you have dbmate installed, and run `dbmate up`
* Make sure RabbitMQ is running on port 5672, and its management console is running on 15672. Running ...

```
docker-compose up -d rabbitmq
```

   ... will run rabbitmq in a Docker container.
* Create a fanout exchange on your RabbitMQ install called `events_development_fanout`, and create a queue called `events_development`. Make sure you bind your queue to your exchange.

* Make sure you have Postgres installed, at the URL in your .env file

* Run the worker like `binaryName worker`
* Run the publisher like `binaryName publisher <count>`, where <count> is the number of fake events you want to publish to your fanout exchange