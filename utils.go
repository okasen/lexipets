package main

import "github.com/gocql/gocql"

func cassandra() (*gocql.Session, error) {
	cluster := gocql.NewCluster("127.0.0.1:9042")
	cluster.Keyspace = "lexipets"

	session, err := cluster.CreateSession()
	if err != nil {
		return nil, err
	}
	return session, nil
}
