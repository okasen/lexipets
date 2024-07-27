package main

import (
	"github.com/gocql/gocql"
	"os"
)

func cassandra() (*gocql.Session, error) {
	cluster := gocql.NewCluster(os.Getenv("CASSANDRA"))
	cluster.Keyspace = os.Getenv("KEYSPACE")

	session, err := cluster.CreateSession()
	if err != nil {
		return nil, err
	}
	return session, nil
}
