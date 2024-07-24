# Lexipets
_Golang + Cassandra + My poor understanding of genetics = virtual pet site_

## What?
This is a "virtual pet site" (think neopets, wajas) developed in Golang with Cassandra for persistence. Right now it's just a basic genetics/animal breeding simulator. Once I have this hosted somewhere with a front end, there will be art that illustrates the concept better.

Just think punnet squares. Dominant genes, recessive genes, and random 50/50 chances.

## Why?

I want my daughter to have the virtual pet site experience I had in the 00's. Maybe she'll learn something about genetics. I'm not really qualified to teach genetics.

## How?

Currently this is an API written in Go using Gin for the REST API, Cassandra (on Docker) for data persistence, and Testify for unit testing.

When I add a front end, it will be written in React or Vue. Probably. 

The API generates "pets" which have a species, genes, and an "img" string which will eventually correlate to an image filepath. The image string is based on the "genes" present on each pet, so each image will be slightly different based on the pets' dominant and recessive gene makeup.

## How do I run the API locally?

You'll need to run a Cassandra instance and input the fixture data in `fixtures.cql` via CQLSH. 

Your Cassandra instance must also be accessible on `127.0.0.1:9042` as of writing. If you run Cassandra elsewhere, you'll need to edit the host and port in `utils.go` to match where your instance can be found.

`go run .` in the `lexipets` directory (y'know, after you've cloned the repo) will run the API on localhost.
