# Order Book System

This project is a simple order book system written in **golang** that can be used to place orders and view the order
book.

## Technologies used

- [Golang](https://golang.org/), Programming language.
- [Gin](https://github.com/gin-gonic/gin), HTTP web framework.
- [Gorm](https://gorm.io/), ORM library for Golang.
- [PostgreSQL](https://www.postgresql.org/), Database.
- [Docker](https://www.docker.com/), Containerization platform.
- [Apache Kafka](https://kafka.apache.org/), Message broker and streaming platform.

## Run project

First, you need to have docker and postgreSQL installed on your machine. Then, you can run the following command to
start the project:

```shell
$ docker-compose up
```

Create topic by this command:

```shell
$ docker-compose exec kafka kafka-topics --create --topic yourTopic --bootstrap-server localhost:9092
```

To test if the topic is created, you can run the following command:

```shell
$ docker-compose exec broker kafka-topics --bootstrap-server broker:29092 --list
```

Then run your postgreSQL server on port 5432

At last, you can run the following command to start the project:

```shell
$ start.sh
```

## Test project

You can run the following command to run the tests:

```shell
$ go test ./...
```
