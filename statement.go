package neoql

import (
	"database/sql/driver"
	"strconv"
)

// stmt implements the sql/driver.Stmt interface.
type stmt struct {
	conn  *conn
	query string
}

// Exec implements the Exec() method of the sql/driver.Stmt interface.
func (stm *stmt) Exec(args []driver.Value) (driver.Result, error) {
	if _, err := stm.conn.run(stm.query, makeArgsMap(args)); err != nil {
		return nil, err
	}
	return &result{}, nil
}

// Query implements the Query() method of the sql/driver.Stmt interface.
func (stm *stmt) Query(args []driver.Value) (driver.Rows, error) {
	return stm.conn.run(stm.query, makeArgsMap(args))
}

// NumInput implements the NumInput() method of the sql/driver.Stmt interface.
// It is not supported by this driver.
func (stm *stmt) NumInput() int {
	return -1
}

// Exec implements the Exec() method of the sql/driver.Stmt interface.
func (stm *stmt) Close() error {
	stm.query = ""
	stm.conn = nil
	return nil
}

// makeArgsMap converts a driver.Value slices into a named parameters map.
// It also converts driver.Value if necessary to fit in the neo4j database.
func makeArgsMap(args []driver.Value) map[string]interface{} {
	params := make(map[string]interface{})
	for i, v := range args {
		// Neo4j databases using Bolt protocol do not support binary data. If we receive a byte slice,
		// it probably comes from a custom type which implements MarshalPS so we don't want to encode it again.
		// If some day Neo4j handles binary data, we might want to unmarshal the bytes, then if it's valid
		// assume it must not be encoded, if it is not, assume it should be encoded. Or maybe there is a better
		// way.
		if b, ok := v.([]byte); ok {
			v = rawBytes(b)
		}
		params[strconv.Itoa(i)] = v
	}
	return params
}
