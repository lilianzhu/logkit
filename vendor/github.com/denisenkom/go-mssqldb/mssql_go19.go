// +build go1.9

package mssql

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	// "github.com/cockroachdb/apd"
)

// Type alias provided for compibility.
//
// Deprecated: users should transition to the new names when possible.
type MssqlDriver = Driver
type MssqlBulk = Bulk
type MssqlBulkOptions = BulkOptions
type MssqlConn = Conn
type MssqlResult = Result
type MssqlRows = Rows
type MssqlStmt = Stmt

var _ driver.NamedValueChecker = &Conn{}

func (c *Conn) CheckNamedValue(nv *driver.NamedValue) error {
	switch v := nv.Value.(type) {
	case sql.Out:
		if c.outs == nil {
			c.outs = make(map[string]interface{})
		}
		c.outs[nv.Name] = v.Dest

		// Unwrap the Out value and check the inner value.
		lnv := *nv
		lnv.Value = v.Dest
		err := c.CheckNamedValue(&lnv)
		if err != nil {
			if err != driver.ErrSkip {
				return err
			}
			lnv.Value, err = driver.DefaultParameterConverter.ConvertValue(lnv.Value)
			if err != nil {
				return err
			}
		}
		nv.Value = sql.Out{Dest: lnv.Value}
		return nil
	// case *apd.Decimal:
	// 	return nil
	default:
		return driver.ErrSkip
	}
}

func (s *Stmt) makeParamExtra(val driver.Value) (res Param, err error) {
	switch val := val.(type) {
	case sql.Out:
		res, err = s.makeParam(val.Dest)
		res.Flags = fByRevValue
	default:
		err = fmt.Errorf("mssql: unknown type for %T", val)
	}
	return
}
