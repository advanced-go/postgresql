package pgxdml

import (
	"errors"
	"fmt"
	"strings"
)

// TODO : create a way to detect and remove SQL inject attacks
// DROP TABLE, DELETE FROM, SELECT * FROM, a double-dashed sequence ‘--’, or a semicolon ;
// quotes /*

// SanitizeString - verify that a string does not contain any text associated with a SQL injection
func SanitizeString(s string) error {
	trimmed := TrimDoubleSpace(strings.ToLower(s))
	for _, t := range tokens {
		index := strings.Index(trimmed, t)
		if index != -1 {
			return errors.New(fmt.Sprintf("SQL injection embedded in string [%v] : %v", trimmed, t))
		}
	}
	return nil
}
