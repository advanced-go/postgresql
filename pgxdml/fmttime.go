package pgxdml

import (
	"strconv"
	"strings"
	"time"
)

// FmtTimestamp - format time.Time into the following string format : 2023-04-14 14:14:45.522460
func FmtTimestamp(t time.Time) string {
	var buf []byte
	t = t.UTC()
	year, month, day := t.Date()
	itoa(&buf, year, 4)
	buf = append(buf, '-')
	itoa(&buf, int(month), 2)
	buf = append(buf, '-')
	itoa(&buf, day, 2)
	buf = append(buf, ' ')

	hour, min, sec := t.Clock()
	itoa(&buf, hour, 2)
	buf = append(buf, ':')
	itoa(&buf, min, 2)
	buf = append(buf, ':')
	itoa(&buf, sec, 2)
	//if l.flag&Lmicroseconds != 0 {
	buf = append(buf, '.')
	itoa(&buf, t.Nanosecond()/1e3, 6)
	//}
	//buf = append(buf, ' ')
	return string(buf)
}

func itoa(buf *[]byte, i int, wid int) {
	// Assemble decimal in reverse order.
	var b [20]byte
	bp := len(b) - 1
	for i >= 10 || wid > 1 {
		wid--
		q := i / 10
		b[bp] = byte('0' + i - q*10)
		bp--
		i = q
	}
	// i < 10
	b[bp] = byte('0' + i)
	*buf = append(*buf, b[bp:]...)
}

// ParseTimestamp - parse a string into a time.Time, using the following string : 2023-04-14 14:14:45.522460
func ParseTimestamp(s string) (time.Time, error) {
	if len(s) == 0 {
		return time.Now().UTC(), nil
	}
	year, err := strconv.Atoi(s[0:4])
	if err != nil {
		return time.Now().UTC(), err
	}
	month, er1 := strconv.Atoi(s[5:7])
	if er1 != nil {
		return time.Now().UTC(), er1
	}
	day, er2 := strconv.Atoi(s[8:10])
	if er2 != nil {
		return time.Now().UTC(), er2
	}
	hour, er3 := strconv.Atoi(s[11:13])
	if er3 != nil {
		return time.Now().UTC(), er3
	}
	min, er4 := strconv.Atoi(s[14:16])
	if er4 != nil {
		return time.Now().UTC(), er4
	}
	sec, er5 := strconv.Atoi(s[17:19])
	if er5 != nil {
		return time.Now().UTC(), er5
	}
	ns, er6 := strconv.Atoi(s[20:26])
	if er6 != nil {
		return time.Now().UTC(), er6
	}
	return time.Date(year, time.Month(month), day, hour, min, sec, ns*1000, time.UTC), nil
}

// ParseDuration - parse a duration string which contains the time unit abbreviation: m, s, ms, µs
func ParseDuration(s string) (time.Duration, error) {
	if s == "" {
		return 0, nil
	}
	tokens := strings.Split(s, "ms")
	if len(tokens) == 2 {
		val, err := strconv.Atoi(tokens[0])
		if err != nil {
			return 0, err
		}
		return time.Duration(val) * time.Millisecond, nil
	}
	tokens = strings.Split(s, "µs")
	if len(tokens) == 2 {
		val, err := strconv.Atoi(tokens[0])
		if err != nil {
			return 0, err
		}
		return time.Duration(val) * time.Microsecond, nil
	}
	tokens = strings.Split(s, "m")
	if len(tokens) == 2 {
		val, err := strconv.Atoi(tokens[0])
		if err != nil {
			return 0, err
		}
		return time.Duration(val) * time.Minute, nil
	}
	// Assume seconds
	tokens = strings.Split(s, "s")
	if len(tokens) == 2 {
		s = tokens[0]
	}
	val, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	return time.Duration(val) * time.Second, nil
}
