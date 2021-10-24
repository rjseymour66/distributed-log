package server

import (
	"fmt"
	"sync"
)

// Struct with mutual exclusion lock for reading and appending
// records: A slice of Records
type Log struct {
	mu	sync.Mutex
	records	[]Record
}

// Returns a pointer to a new Log
func NewLog() *Log {
	return &Log{}
}

// Appends a record to a log. 
// Returns the new offset and an error
func (c *Log) Append(record Record) (uint64, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	record.Offset = uint64(len(c.records))
	c.records = append(c.records, record)
	return record.Offset, nil
}

// Returns the log record at the provided offset and an error.
// Returns an erroriIf the offset is larger than the highest 
// record offset.
func (c *Log) Read(offset uint64) (Record, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if offset >= uint64(len(c.records)) {
		return Record{}, ErrOffsetNotFound
	}
	return c.records[offset], nil
}

// Represents a Record object consisting of a slice of bytes with 
// the record information, and the last offset.
type Record struct {
	Value	[]byte	`json:"value"`
	Offset	uint64	`json:"offset"`
}

var ErrOffsetNotFound = fmt.Errorf("offset not found")
