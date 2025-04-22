package models

import (
  "math/big"
  "strconv"
  "time"
  "fmt"
)

// BigInt is a wrapper over big.Int to implement only unmarshalText
// for json decoding.
type BigInt big.Int

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (b *BigInt) UnmarshalText(text []byte) (err error) {
	var bigInt = new(big.Int)
	err = bigInt.UnmarshalText(text)
	if err != nil {
		return
	}

	*b = BigInt(*bigInt)
	return nil
}

// MarshalText implements the encoding.TextMarshaler
func (b *BigInt) MarshalText() (text []byte, err error) {
	return []byte(b.Int().String()), nil
}

// Int returns b's *big.Int form
func (b *BigInt) Int() *big.Int {
	return (*big.Int)(b)
}

func (b *BigInt) String() string {
    return b.Int().String()
}

func (b *BigInt) SetString(s string, base int) error {
    i, ok := new(big.Int).SetString(s, base)
    if !ok {
        return fmt.Errorf("invalid big.Int string: %s", s)
    }
    *b = BigInt(*i)
    return nil
}

func (b *BigInt) ToInt64() (int64, error) {
    return b.Int().Int64(), nil // optional: check for overflow
}

func (b *BigInt) Clone() *BigInt {
    cloned := new(big.Int).Set(b.Int())
    bb := BigInt(*cloned)
    return &bb
}

// Time is a wrapper over big.Int to implement only unmarshalText
// for json decoding.
type Time time.Time

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (t *Time) UnmarshalText(text []byte) (err error) {
	input, err := strconv.ParseInt(string(text), 10, 64)
	if err != nil {
		err = wrapErr(err, "strconv.ParseInt")
		return
	}

	var timestamp = time.Unix(input, 0)
	*t = Time(timestamp)

	return nil
}

// Time returns t's time.Time form
func (t Time) Time() time.Time {
	return time.Time(t)
}

// MarshalText implements the encoding.TextMarshaler
func (t Time) MarshalText() (text []byte, err error) {
	return []byte(strconv.FormatInt(t.Time().Unix(), 10)), nil
}

func (t *Time) SetUnix(unix int64) {
    *t = Time(time.Unix(unix, 0))
}

func (t Time) After(other Time) bool {
    return t.Time().After(other.Time())
}

func (t Time) Before(other Time) bool {
    return t.Time().Before(other.Time())
}

func (t Time) IsZero() bool {
    return t.Time().IsZero()
}
