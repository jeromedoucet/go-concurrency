package database

import "time"

type DbClient interface {
	Set(string, interface{}, time.Duration) error
	Get(string) (interface{}, error)
	Remove(string) error
	Close()
}
