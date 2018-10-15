package jsonmuncher

import (
	"errors"
	"strconv"
	"strings"
)

// EndOfValue denotes that the end of an object or array has already been
// reached, and no new elements can be read. This is returned by the NextKey()
// and NextValue() functions.
var EndOfValue = errors.New("End of value reached")

// ErrIncomplete is returned when a call is made to a value with an "Incomplete"
// status, meaning a read error had occurred during a previous operation.
var ErrIncomplete = errors.New("Status incomplete denotes failed read")

// ErrWorkingChild is returned when NextKey(), NextValue(), or Close() is
// called on an object or array, but one of its elements is only partially read.
// The child element must be fully read or closed first.
var ErrWorkingChild = errors.New("Unable to consume when child element is partially read")

// ErrNoParamsSpecified is returned from the Compare() and FindKey() functions
// when no arguments are passed. These variadic functions expect at least one
// argument.
var ErrNoParamsSpecified = errors.New("At least one argument should be provided")

// ErrTypeMismatch is returned when a JsonValue method specific to a particular
// JSON type is called on a different JSON type. For example, ValueNum() will
// return this error if called on any JsonValue that isn't a Number.
type ErrTypeMismatch struct {
	Provided JsonType
	Expected []JsonType
}

var showType = [...]string{
	Null:   "Null",
	Bool:   "Bool",
	Number: "Number",
	String: "String",
	Array:  "Array",
	Object: "Object",
}

func newErrTypeMismatch(p JsonType, e ...JsonType) ErrTypeMismatch {
	return ErrTypeMismatch{p, e}
}

// Error implements error for ErrTypeMismatch.
func (e ErrTypeMismatch) Error() string {
	var bld strings.Builder
	bld.WriteString("Method cannot be called on type ")
	bld.WriteString(showType[e.Provided])
	bld.WriteString(", only on ")
	for i := 0; i < len(e.Expected); i++ {
		if i > 0 {
			bld.WriteString(" or ")
		}
		bld.WriteString(showType[e.Expected[i]])
	}
	return bld.String()
}

// ErrUnexpectedChar is returned whenever a syntactic parse error is
// encountered: an illegal character or an unexpected EOF.
type ErrUnexpectedChar struct {
	Offset      uint64
	ProvidedEOF bool
	Provided    byte
	Expected    []byte
	CustomMsg   string
}

func newErrUnexpectedChar(off uint64, p byte, e ...byte) ErrUnexpectedChar {
	return ErrUnexpectedChar{off, false, p, e, ""}
}

func newErrUnexpectedEOF(off uint64, e ...byte) ErrUnexpectedChar {
	return ErrUnexpectedChar{off, true, 0, e, ""}
}

// Error implements error for ErrUnexpectedChar.
func (e ErrUnexpectedChar) Error() string {
	var bld strings.Builder
	bld.WriteString("Unexpected ")
	if e.ProvidedEOF {
		bld.WriteString("EOF")
	} else {
		bld.WriteString(strconv.QuoteRuneToGraphic(rune(e.Provided)))
	}
	bld.WriteString(" at file offset ")
	bld.WriteString(strconv.FormatUint(e.Offset, 10))
	if e.CustomMsg != "" {
		bld.WriteString(": ")
		bld.WriteString(e.CustomMsg)
		return bld.String()
	}
	bld.WriteString(", expected ")
	if len(e.Expected) == 1 {
		bld.WriteString(strconv.QuoteRuneToGraphic(rune(e.Expected[0])))
		return bld.String()
	}
	bld.WriteString("one of ")
	for i := 0; i < len(e.Expected); i++ {
		if i > 0 {
			bld.WriteString(", ")
		}
		if i+1 < len(e.Expected) && e.Expected[i]+1 == e.Expected[i+1] {
			fst := e.Expected[i]
			for i+1 < len(e.Expected) && e.Expected[i]+1 == e.Expected[i+1] {
				i++
			}
			bld.WriteString(strconv.QuoteRuneToGraphic(rune(fst)))
			bld.WriteString("-")
			bld.WriteString(strconv.QuoteRuneToGraphic(rune(e.Expected[i])))
		} else {
			bld.WriteString(strconv.QuoteRuneToGraphic(rune(e.Expected[i])))
		}
	}
	return bld.String()
}
