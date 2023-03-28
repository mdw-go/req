// Package should info: github.com/mdwhatcott/testing/should@v0.21.3 (a little copy-paste is better than a little dependency)
// AUTO-GENERATED: 2023-03-28 07:48:58.930639 -0600 MDT m=+0.003550197
package should

import (
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"reflect"
	"runtime/debug"
	"sort"
	"strconv"
	"strings"
	"testing"
	"time"
)

// FILE: be_chronological.go

func BeChronological(actual any, expected ...any) error {
	err := validateExpected(0, expected)
	if err != nil {
		return err
	}

	var t []time.Time
	err = validateType(actual, t)
	if err != nil {
		return err
	}

	times := actual.([]time.Time)
	if sort.SliceIsSorted(times, func(i, j int) bool { return times[i].Before(times[j]) }) {
		return nil
	}
	return failure("expected to be chronological: %v", times)
}
func (negated) BeChronological(actual any, expected ...any) error {
	err := BeChronological(actual, expected...)
	if errors.Is(err, ErrAssertionFailure) {
		return nil
	}
	if err != nil {
		return err
	}
	return failure("want non-chronological times, got chronological times:", actual)
}

// FILE: be_empty.go

func BeEmpty(actual any, expected ...any) error {
	err := validateExpected(0, expected)
	if err != nil {
		return err
	}

	err = validateKind(actual, kindsWithLength...)
	if err != nil {
		return err
	}

	length := reflect.ValueOf(actual).Len()
	if length == 0 {
		return nil
	}

	TYPE := reflect.TypeOf(actual).String()
	return failure("got len(%s) == %d, want empty %s", TYPE, length, TYPE)
}
func (negated) BeEmpty(actual any, expected ...any) error {
	err := BeEmpty(actual, expected...)
	if errors.Is(err, ErrAssertionFailure) {
		return nil
	}
	if err != nil {
		return err
	}
	TYPE := reflect.TypeOf(actual).String()
	return failure("got empty %s, want non-empty %s", TYPE, TYPE)
}

// FILE: be_false.go

func BeFalse(actual any, expected ...any) error {
	err := validateExpected(0, expected)
	if err != nil {
		return err
	}

	err = validateType(actual, *new(bool))
	if err != nil {
		return err
	}

	boolean := actual.(bool)
	if boolean {
		return failure("got <true>, want <false>")
	}

	return nil
}

// FILE: be_greater_than.go

func BeGreaterThan(actual any, EXPECTED ...any) error {
	lessThanErr := BeLessThan(actual, EXPECTED...)
	if errors.Is(lessThanErr, ErrTypeMismatch) || errors.Is(lessThanErr, ErrExpectedCountInvalid) {
		return lessThanErr
	}
	equalityErr := Equal(actual, EXPECTED...)
	if lessThanErr == nil || equalityErr == nil {
		return failure("%v was not greater than %v", actual, EXPECTED[0])
	}
	return nil
}
func (negated) BeGreaterThan(actual any, expected ...any) error {
	err := BeGreaterThan(actual, expected...)
	if errors.Is(err, ErrAssertionFailure) {
		return nil
	}

	if err != nil {
		return err
	}

	return failure("\n"+
		"  expected:               %#v\n"+
		"  to not be greater than: %#v\n"+
		"  (but it was)",
		expected[0],
		actual,
	)
}

// FILE: be_greater_than_or_equal_to.go

func BeGreaterThanOrEqualTo(actual any, expected ...any) error {
	err := Equal(actual, expected...)
	if err == nil {
		return nil
	}
	err = BeGreaterThan(actual, expected...)
	if errors.Is(err, ErrAssertionFailure) {
		return failure("%v was not greater than or equal to %v", actual, expected)
	}

	if err != nil {
		return err
	}
	return nil
}
func (negated) BeGreaterThanOrEqualTo(actual any, expected ...any) error {
	err := BeGreaterThanOrEqualTo(actual, expected...)
	if errors.Is(err, ErrAssertionFailure) {
		return nil
	}

	if err != nil {
		return err
	}

	return failure("\n"+
		"  expected:                           %#v\n"+
		"  to not be greater than or equal to: %#v\n"+
		"  (but it was)",
		expected[0],
		actual,
	)
}

// FILE: be_in.go

func BeIn(actual any, expected ...any) error {
	err := validateExpected(1, expected)
	if err != nil {
		return err
	}

	err = Contain(expected[0], actual)
	if err != nil {
		return err
	}

	return nil
}
func (negated) BeIn(actual any, expected ...any) error {
	err := validateExpected(1, expected)
	if err != nil {
		return err
	}

	return NOT.Contain(expected[0], actual)
}

// FILE: be_less_than.go

func BeLessThan(actual any, EXPECTED ...any) error {
	err := validateExpected(1, EXPECTED)
	if err != nil {
		return err
	}

	expected := EXPECTED[0]
	failed := false

	for _, spec := range lessThanSpecs {
		if !spec.assertable(actual, expected) {
			continue
		}
		if spec.passes(actual, expected) {
			return nil
		}
		failed = true
		break
	}

	if failed {
		return failure("%v was not less than %v", actual, expected)
	}
	return wrap(ErrTypeMismatch, "could not compare [%v] and [%v]",
		reflect.TypeOf(actual), reflect.TypeOf(expected))
}
func (negated) BeLessThan(actual any, expected ...any) error {
	err := BeLessThan(actual, expected...)
	if errors.Is(err, ErrAssertionFailure) {
		return nil
	}

	if err != nil {
		return err
	}

	return failure("\n"+
		"  expected:            %#v\n"+
		"  to not be less than: %#v\n"+
		"  (but it was)",
		expected[0],
		actual,
	)
}

var lessThanSpecs = []specification{
	bothStringsLessThan{},
	bothSignedIntegersLessThan{},
	bothUnsignedIntegersLessThan{},
	bothFloatsLessThan{},
	signedAndUnsignedLessThan{},
	unsignedAndSignedLessThan{},
	floatAndIntegerLessThan{},
	integerAndFloatLessThan{},
	bothTimesLessThan{},
}

type bothStringsLessThan struct{}

func (bothStringsLessThan) assertable(a, b any) bool {
	return reflect.ValueOf(a).Kind() == reflect.String && reflect.ValueOf(b).Kind() == reflect.String
}
func (bothStringsLessThan) passes(a, b any) bool {
	return reflect.ValueOf(a).String() < reflect.ValueOf(b).String()
}

type bothSignedIntegersLessThan struct{}

func (bothSignedIntegersLessThan) assertable(a, b any) bool {
	return isSignedInteger(a) && isSignedInteger(b)
}
func (bothSignedIntegersLessThan) passes(a, b any) bool {
	return reflect.ValueOf(a).Int() < reflect.ValueOf(b).Int()
}

type bothUnsignedIntegersLessThan struct{}

func (bothUnsignedIntegersLessThan) assertable(a, b any) bool {
	return isUnsignedInteger(a) && isUnsignedInteger(b)
}
func (bothUnsignedIntegersLessThan) passes(a, b any) bool {
	return reflect.ValueOf(a).Uint() < reflect.ValueOf(b).Uint()
}

type bothFloatsLessThan struct{}

func (bothFloatsLessThan) assertable(a, b any) bool {
	return isFloat(a) && isFloat(b)
}
func (bothFloatsLessThan) passes(a, b any) bool {
	return reflect.ValueOf(a).Float() < reflect.ValueOf(b).Float()
}

type signedAndUnsignedLessThan struct{}

func (signedAndUnsignedLessThan) assertable(a, b any) bool {
	return isSignedInteger(a) && isUnsignedInteger(b)
}
func (signedAndUnsignedLessThan) passes(a, b any) bool {
	A := reflect.ValueOf(a)
	B := reflect.ValueOf(b)
	if A.Int() < 0 {
		return true
	}
	return uint64(A.Int()) < B.Uint()
}

type unsignedAndSignedLessThan struct{}

func (unsignedAndSignedLessThan) assertable(a, b any) bool {
	return isUnsignedInteger(a) && isSignedInteger(b)
}
func (unsignedAndSignedLessThan) passes(a, b any) bool {
	A := reflect.ValueOf(a)
	B := reflect.ValueOf(b)
	if A.Uint() > math.MaxInt64 {
		return false
	}
	return int64(A.Uint()) < B.Int()
}

type floatAndIntegerLessThan struct{}

func (floatAndIntegerLessThan) assertable(a, b any) bool {
	return isFloat(a) && isInteger(b)
}
func (floatAndIntegerLessThan) passes(a, b any) bool {
	return asFloat(a) < asFloat(b)
}

type integerAndFloatLessThan struct{}

func (integerAndFloatLessThan) assertable(a, b any) bool {
	return isInteger(a) && isFloat(b)
}
func (integerAndFloatLessThan) passes(a, b any) bool {
	return asFloat(a) < asFloat(b)
}

type bothTimesLessThan struct{}

func (bothTimesLessThan) assertable(a, b any) bool {
	return isTime(a) && isTime(b)
}
func (bothTimesLessThan) passes(a, b any) bool {
	return a.(time.Time).Before(b.(time.Time))
}

// FILE: be_less_than_or_equal_to.go

func BeLessThanOrEqualTo(actual any, expected ...any) error {
	err := Equal(actual, expected...)
	if err == nil {
		return nil
	}
	err = BeLessThan(actual, expected...)
	if errors.Is(err, ErrAssertionFailure) {
		return failure("%v was not less than or equal to %v", actual, expected)
	}

	if err != nil {
		return err
	}
	return nil
}
func (negated) BeLessThanOrEqualTo(actual any, expected ...any) error {
	err := BeLessThanOrEqualTo(actual, expected...)
	if errors.Is(err, ErrAssertionFailure) {
		return nil
	}

	if err != nil {
		return err
	}

	return failure("\n"+
		"  expected:                        %#v\n"+
		"  to not be less than or equal to: %#v\n"+
		"  (but it was)",
		expected[0],
		actual,
	)
}

// FILE: be_nil.go

func BeNil(actual any, expected ...any) error {
	err := validateExpected(0, expected)
	if err != nil {
		return err
	}

	if actual == nil || interfaceHasNilValue(actual) {
		return nil
	}

	return failure("got %#v, want <nil>", actual)
}
func interfaceHasNilValue(actual any) bool {
	value := reflect.ValueOf(actual)
	kind := value.Kind()
	nillable := kind == reflect.Slice ||
		kind == reflect.Chan ||
		kind == reflect.Func ||
		kind == reflect.Ptr ||
		kind == reflect.Map

	return nillable && value.IsNil()
}
func (negated) BeNil(actual any, expected ...any) error {
	err := BeNil(actual, expected...)
	if errors.Is(err, ErrAssertionFailure) {
		return nil
	}

	if err != nil {
		return err
	}

	return failure("got nil, want non-<nil>")
}

// FILE: be_true.go

func BeTrue(actual any, expected ...any) error {
	err := validateExpected(0, expected)
	if err != nil {
		return err
	}

	err = validateType(actual, *new(bool))
	if err != nil {
		return err
	}

	boolean := actual.(bool)
	if !boolean {
		return failure("got <false>, want <true>")
	}
	return nil
}

// FILE: contain.go

func Contain(actual any, expected ...any) error {
	err := validateExpected(1, expected)
	if err != nil {
		return err
	}

	err = validateKind(actual, containerKinds...)
	if err != nil {
		return err
	}

	actualValue := reflect.ValueOf(actual)
	EXPECTED := expected[0]

	switch reflect.TypeOf(actual).Kind() {
	case reflect.Map:
		expectedValue := reflect.ValueOf(EXPECTED)
		value := actualValue.MapIndex(expectedValue)
		if value.IsValid() {
			return nil
		}
	case reflect.Array, reflect.Slice:
		for i := 0; i < actualValue.Len(); i++ {
			item := actualValue.Index(i).Interface()
			if Equal(EXPECTED, item) == nil {
				return nil
			}
		}
	case reflect.String:
		err = validateKind(EXPECTED, reflect.String, reflectRune)
		if err != nil {
			return err
		}

		expectedRune, ok := EXPECTED.(rune)
		if ok {
			EXPECTED = string(expectedRune)
		}

		full := actual.(string)
		sub := EXPECTED.(string)
		if strings.Contains(full, sub) {
			return nil
		}
	}

	return failure("\n"+
		"   item absent: %#v\n"+
		"   within:      %#v",
		EXPECTED,
		actual,
	)
}
func (negated) Contain(actual any, expected ...any) error {
	err := Contain(actual, expected...)
	if errors.Is(err, ErrAssertionFailure) {
		return nil
	}

	if err != nil {
		return err
	}

	return failure("\n"+
		"item found: %#v\n"+
		"within:     %#v",
		expected[0],
		actual,
	)
}

const reflectRune = reflect.Int32

// FILE: doc.go

// FILE: end_with.go

func EndWith(actual any, expected ...any) error {
	err := validateExpected(1, expected)
	if err != nil {
		return err
	}

	err = validateKind(actual, orderedContainerKinds...)
	if err != nil {
		return err
	}

	actualValue := reflect.ValueOf(actual)
	EXPECTED := expected[0]

	switch reflect.TypeOf(actual).Kind() {
	case reflect.Array, reflect.Slice:
		if actualValue.Len() == 0 {
			break
		}
		last := actualValue.Index(actualValue.Len() - 1).Interface()
		if Equal(EXPECTED, last) == nil {
			return nil
		}
	case reflect.String:
		err = validateKind(EXPECTED, reflect.String, reflectRune)
		if err != nil {
			return err
		}

		expectedRune, ok := EXPECTED.(rune)
		if ok {
			EXPECTED = string(expectedRune)
		}

		full := actual.(string)
		prefix := EXPECTED.(string)
		if strings.HasSuffix(full, prefix) {
			return nil
		}
	}

	return failure("\n"+
		"   proposed prefix: %#v\n"+
		"   not a prefix of: %#v",
		EXPECTED,
		actual,
	)
}

// FILE: equal.go

func Equal(actual any, EXPECTED ...any) error {
	err := validateExpected(1, EXPECTED)
	if err != nil {
		return err
	}

	expected := EXPECTED[0]

	for _, spec := range equalitySpecs {
		if !spec.assertable(actual, expected) {
			continue
		}
		if spec.passes(actual, expected) {
			return nil
		}
		break
	}
	return failure(report(actual, expected))
}
func (negated) Equal(actual any, expected ...any) error {
	err := Equal(actual, expected...)
	if errors.Is(err, ErrAssertionFailure) {
		return nil
	}

	if err != nil {
		return err
	}

	return failure("\n"+
		"  expected:     %#v\n"+
		"  to not equal: %#v\n"+
		"  (but it did)",
		expected[0],
		actual,
	)
}

var equalitySpecs = []specification{
	numericEquality{},
	timeEquality{},
	deepEquality{},
}

func report(a, b any) string {
	aType := fmt.Sprintf("(%v)", reflect.TypeOf(a))
	bType := fmt.Sprintf("(%v)", reflect.TypeOf(b))
	longestType := int(math.Max(float64(len(aType)), float64(len(bType))))
	aType += strings.Repeat(" ", longestType-len(aType))
	bType += strings.Repeat(" ", longestType-len(bType))
	aFormat := fmt.Sprintf(format(a), a)
	bFormat := fmt.Sprintf(format(b), b)
	typeDiff := diff(bType, aType)
	valueDiff := diff(bFormat, aFormat)

	builder := new(strings.Builder)
	_, _ = fmt.Fprintf(builder, "\n")
	_, _ = fmt.Fprintf(builder, "Expected: %s %s\n", bType, bFormat)
	_, _ = fmt.Fprintf(builder, "Actual:   %s %s\n", aType, aFormat)
	_, _ = fmt.Fprintf(builder, "          %s %s", typeDiff, valueDiff)

	return builder.String()
}
func format(v any) string {
	if isNumeric(v) || isTime(v) {
		return "%v"
	} else {
		return "%#v"
	}
}
func diff(a, b string) string {
	result := new(strings.Builder)
	for x := 0; x < len(a) && x < len(b); x++ {
		if x >= len(a) || x >= len(b) || a[x] != b[x] {
			result.WriteString("^")
		} else {
			result.WriteString(" ")
		}
	}
	return result.String()
}

type deepEquality struct{}

func (deepEquality) assertable(a, b any) bool {
	return reflect.TypeOf(a) == reflect.TypeOf(b)
}
func (deepEquality) passes(a, b any) bool {
	return reflect.DeepEqual(a, b)
}

type numericEquality struct{}

func (numericEquality) assertable(a, b any) bool {
	return isNumeric(a) && isNumeric(b)
}
func (numericEquality) passes(a, b any) bool {
	aValue := reflect.ValueOf(a)
	bValue := reflect.ValueOf(b)
	if isUnsignedInteger(a) && isSignedInteger(b) && aValue.Uint() >= math.MaxInt64 {
		return false
	}
	if isSignedInteger(a) && isUnsignedInteger(b) && bValue.Uint() >= math.MaxInt64 {
		return false
	}
	aAsB := aValue.Convert(bValue.Type()).Interface()
	bAsA := bValue.Convert(aValue.Type()).Interface()
	return a == bAsA && b == aAsB
}

type timeEquality struct{}

func (timeEquality) assertable(a, b any) bool {
	return isTime(a) && isTime(b)
}
func (timeEquality) passes(a, b any) bool {
	return a.(time.Time).Equal(b.(time.Time))
}
func isTime(v any) bool {
	_, ok := v.(time.Time)
	return ok
}

// FILE: errors.go

var (
	ErrExpectedCountInvalid = errors.New("expected count invalid")
	ErrTypeMismatch         = errors.New("type mismatch")
	ErrKindMismatch         = errors.New("kind mismatch")
	ErrAssertionFailure     = errors.New("assertion failure")
)

func failure(format string, args ...any) error {
	trace := stack()
	if len(trace) > 0 {
		format += "\nStack: (filtered)\n%s"
		args = append(args, trace)
	}
	return wrap(ErrAssertionFailure, format, args...)
}
func stack() string {
	lines := strings.Split(string(debug.Stack()), "\n")
	var filtered []string
	for x := 1; x < len(lines)-1; x += 2 {
		fileLineRaw := lines[x+1]
		if strings.Contains(fileLineRaw, "_test.go:") {
			filtered = append(filtered, lines[x], fileLineRaw)
			line, ok := readSourceCodeLine(fileLineRaw)
			if ok {
				filtered = append(filtered, "  "+line)
			}

		}
	}
	if len(filtered) == 0 {
		return ""
	}
	return "> " + strings.Join(filtered, "\n> ")
}
func readSourceCodeLine(fileLineRaw string) (string, bool) {
	fileLineJoined := strings.Fields(strings.TrimSpace(fileLineRaw))[0]
	fileLine := strings.Split(fileLineJoined, ":")
	sourceCode, _ := os.ReadFile(fileLine[0])
	sourceCodeLines := strings.Split(string(sourceCode), "\n")
	lineNumber, _ := strconv.Atoi(fileLine[1])
	lineNumber--
	if len(sourceCodeLines) <= lineNumber {
		return "", false
	}
	return sourceCodeLines[lineNumber], true
}
func wrap(inner error, format string, args ...any) error {
	return fmt.Errorf("%w: "+fmt.Sprintf(format, args...), inner)
}

// FILE: expected.go

func validateExpected(count int, expected []any) error {
	length := len(expected)
	if length == count {
		return nil
	}

	s := pluralize(length)
	return wrap(ErrExpectedCountInvalid, "got %d value%s, want %d", length, s, count)
}
func pluralize(count int) string {
	if count == 1 {
		return ""
	}
	return "s"
}
func validateType(actual, expected any) error {
	ACTUAL := reflect.TypeOf(actual)
	EXPECTED := reflect.TypeOf(expected)
	if ACTUAL == EXPECTED {
		return nil
	}
	return wrap(ErrTypeMismatch, "got %s, want %s", ACTUAL, EXPECTED)
}
func validateKind(actual any, kinds ...reflect.Kind) error {
	value := reflect.ValueOf(actual)
	kind := value.Kind()
	for _, k := range kinds {
		if k == kind {
			return nil
		}
	}
	return wrap(ErrKindMismatch, "got %s, want one of %v", kind, kinds)
}

// FILE: happen_after.go

func HappenAfter(actual any, expected ...any) error {
	err := validateExpected(1, expected)
	if err != nil {
		return err
	}
	err = validateType(actual, time.Time{})
	if err != nil {
		return err
	}
	err = validateType(expected[0], time.Time{})
	if err != nil {
		return err
	}
	return BeGreaterThan(actual, expected[0])
}

// FILE: happen_before.go

func HappenBefore(actual any, expected ...any) error {
	err := validateExpected(1, expected)
	if err != nil {
		return err
	}
	err = validateType(actual, time.Time{})
	if err != nil {
		return err
	}
	err = validateType(expected[0], time.Time{})
	if err != nil {
		return err
	}
	return BeLessThan(actual, expected[0])
}

// FILE: happen_on.go

func HappenOn(actual any, expected ...any) error {
	err := validateExpected(1, expected)
	if err != nil {
		return err
	}
	err = validateType(actual, time.Time{})
	if err != nil {
		return err
	}
	err = validateType(expected[0], time.Time{})
	if err != nil {
		return err
	}
	return Equal(actual, expected...)
}
func (negated) HappenOn(actual any, expected ...any) error {
	err := HappenOn(actual, expected...)
	if errors.Is(err, ErrAssertionFailure) {
		return nil
	}
	if err != nil {
		return err
	}

	return failure("\n"+
		"  expected:     %#v\n"+
		"  to not equal: %#v\n"+
		"  (but it did)",
		expected[0],
		actual,
	)
}

// FILE: have_length.go

func HaveLength(actual any, expected ...any) error {
	err := validateExpected(1, expected)
	if err != nil {
		return err
	}

	err = validateKind(actual, kindsWithLength...)
	if err != nil {
		return err
	}

	err = validateKind(expected[0], kindSlice(signedIntegerKinds)...)
	if err != nil {
		return err
	}

	expectedLength := reflect.ValueOf(expected[0]).Int()
	actualLength := int64(reflect.ValueOf(actual).Len())
	if actualLength == expectedLength {
		return nil
	}

	return failure("got length of %d, want %d", actualLength, expectedLength)
}

// FILE: kinds.go

var floatTypes = map[reflect.Kind]struct{}{
	reflect.Float32: {},
	reflect.Float64: {},
}

func isFloat(v any) bool {
	_, found := floatTypes[reflect.TypeOf(v).Kind()]
	return found
}
func asFloat(a any) float64 {
	v := reflect.ValueOf(a)
	if isSignedInteger(a) {
		return float64(v.Int())
	}
	if isUnsignedInteger(a) {
		return float64(v.Uint())
	}
	return v.Float()
}

var unsignedIntegerKinds = map[reflect.Kind]struct{}{
	reflect.Uint:    {},
	reflect.Uint8:   {},
	reflect.Uint16:  {},
	reflect.Uint32:  {},
	reflect.Uint64:  {},
	reflect.Uintptr: {},
}

func isUnsignedInteger(v any) bool {
	_, found := unsignedIntegerKinds[reflect.TypeOf(v).Kind()]
	return found
}

var signedIntegerKinds = map[reflect.Kind]struct{}{
	reflect.Int:   {},
	reflect.Int8:  {},
	reflect.Int16: {},
	reflect.Int32: {},
	reflect.Int64: {},
}

func isSignedInteger(v any) bool {
	_, found := signedIntegerKinds[reflect.TypeOf(v).Kind()]
	return found
}
func isInteger(v any) bool {
	return isSignedInteger(v) || isUnsignedInteger(v)
}

var numericKinds = map[reflect.Kind]struct{}{
	reflect.Int:     {},
	reflect.Int8:    {},
	reflect.Int16:   {},
	reflect.Int32:   {},
	reflect.Int64:   {},
	reflect.Uint:    {},
	reflect.Uint8:   {},
	reflect.Uint16:  {},
	reflect.Uint32:  {},
	reflect.Uint64:  {},
	reflect.Float32: {},
	reflect.Float64: {},
}

func isNumeric(v any) bool {
	of := reflect.TypeOf(v)
	if of == nil {
		return false
	}
	_, found := numericKinds[of.Kind()]
	return found
}

var kindsWithLength = []reflect.Kind{
	reflect.Map,
	reflect.Chan,
	reflect.Array,
	reflect.Slice,
	reflect.String,
}
var containerKinds = []reflect.Kind{
	reflect.Map,
	reflect.Array,
	reflect.Slice,
	reflect.String,
}
var orderedContainerKinds = []reflect.Kind{
	reflect.Array,
	reflect.Slice,
	reflect.String,
}

func kindSlice(kinds map[reflect.Kind]struct{}) (result []reflect.Kind) {
	for kind := range kinds {
		result = append(result, kind)
	}
	return result
}

// FILE: not.go

var NOT negated

type negated struct{}

// FILE: panic.go

func Panic(actual any, expected ...any) (err error) {
	err = NOT.Panic(actual, expected...)
	if errors.Is(err, ErrAssertionFailure) {
		return nil
	}

	if err != nil {
		return err
	}

	return failure("provided func did not panic as expected")
}
func (negated) Panic(actual any, expected ...any) (err error) {
	err = validateExpected(0, expected)
	if err != nil {
		return err
	}

	err = validateType(actual, func() {})
	if err != nil {
		return err
	}

	panicked := true
	defer func() {
		r := recover()
		if panicked {
			err = failure(""+
				"provided func should not have"+
				"panicked but it did with: %s", r,
			)
		}
	}()

	actual.(func())()
	panicked = false
	return nil
}

// FILE: reporter.go

type Func func(actual any, expected ...any) error
type Reporter interface {
	Helper()
	Report(error)
	io.Writer
}
type T struct{ Reporter }

func New(t *testing.T) *T {
	return &T{Reporter: NewTestingReporter(t)}
}
func Report(reporters ...Reporter) *T {
	if len(reporters) == 0 {
		reporters = append(reporters, NewWriterReporter(os.Stdout))
	}
	return &T{Reporter: NewCompositeReporter(reporters...)}
}
func (this *T) So(actual any, assertion Func, expected ...any) (ok bool) {
	this.Helper()
	err := assertion(actual, expected...)
	this.Reporter.Report(err)
	return err == nil
}
func (this *T) Print(v ...any) {
	this.Reporter.Helper()
	_, _ = this.Write([]byte(fmt.Sprint(v...)))
}
func (this *T) Printf(f string, v ...any) {
	this.Reporter.Helper()
	_, _ = this.Write([]byte(fmt.Sprintf(f, v...)))
}
func (this *T) Println(v ...any) {
	this.Reporter.Helper()
	_, _ = this.Write([]byte(fmt.Sprintln(v...)))
}

type TestingReporter struct{ *testing.T }

func NewTestingReporter(t *testing.T) *TestingReporter {
	return &TestingReporter{T: t}
}
func (this *TestingReporter) Report(err error) {
	if err != nil {
		this.Helper()
		this.Error(err)
	}
}
func (this *TestingReporter) Write(p []byte) (n int, err error) {
	this.Helper()
	this.Log(string(p))
	return len(p), nil
}

type CompositeReporter struct{ reporters []Reporter }

func (this *CompositeReporter) Helper() {
	for _, reporter := range this.reporters {
		reporter.Helper()
	}
}
func NewCompositeReporter(reporters ...Reporter) *CompositeReporter {
	return &CompositeReporter{reporters: reporters}
}
func (this *CompositeReporter) Report(err error) {
	for _, reporter := range this.reporters {
		reporter.Report(err)
	}
}
func (this *CompositeReporter) Write(p []byte) (n int, err error) {
	for _, reporter := range this.reporters {
		n, err = reporter.Write(p)
		if err != nil {
			break
		}
	}
	return n, err
}

type WriterReporter struct{ io.Writer }

func (this *WriterReporter) Helper() {}
func NewWriterReporter(writer io.Writer) *WriterReporter {
	return &WriterReporter{Writer: writer}
}
func (this *WriterReporter) Report(err error) {
	if err != nil {
		_, _ = fmt.Fprintln(this, err.Error())
	}
}

type LogReporter struct{ logger *log.Logger }

func NewLogReporter(logger *log.Logger) *LogReporter {
	return &LogReporter{logger: logger}
}
func (this LogReporter) Report(err error) {
	if err != nil {
		this.logger.Print(err.Error())
	}
}
func (this LogReporter) Write(p []byte) (n int, err error) {
	this.logger.Print(string(p))
	return len(p), nil
}
func (this LogReporter) Helper() {}

// FILE: run.go

func Run(fixture any, options ...Option) {
	config := new(config)
	for _, option := range options {
		option(config)
	}

	fixtureValue := reflect.ValueOf(fixture)
	fixtureType := reflect.TypeOf(fixture)
	t := fixtureValue.Elem().FieldByName("T").Elem().FieldByName("Reporter").Interface().(*TestingReporter)

	var (
		testNames        []string
		skippedTestNames []string
		focusedTestNames []string
	)
	for x := 0; x < fixtureType.NumMethod(); x++ {
		name := fixtureType.Method(x).Name
		method := fixtureValue.MethodByName(name)
		_, isNiladic := method.Interface().(func())
		if !isNiladic {
			continue
		}

		if strings.HasPrefix(name, "Test") {
			testNames = append(testNames, name)
		} else if strings.HasPrefix(name, "LongTest") {
			testNames = append(testNames, name)

		} else if strings.HasPrefix(name, "SkipLongTest") {
			skippedTestNames = append(skippedTestNames, name)
		} else if strings.HasPrefix(name, "SkipTest") {
			skippedTestNames = append(skippedTestNames, name)

		} else if strings.HasPrefix(name, "FocusLongTest") {
			focusedTestNames = append(focusedTestNames, name)
		} else if strings.HasPrefix(name, "FocusTest") {
			focusedTestNames = append(focusedTestNames, name)
		}
	}

	if len(focusedTestNames) > 0 {
		testNames = focusedTestNames
	}

	if len(testNames) == 0 {
		t.Skip("NOT IMPLEMENTED (no test cases defined, or they are all marked as skipped)")
		return
	}

	if config.parallelFixture {
		t.Parallel()
	}

	setup, hasSetup := fixture.(setupSuite)
	if hasSetup {
		setup.SetupSuite()
	}

	teardown, hasTeardown := fixture.(teardownSuite)
	if hasTeardown {
		defer teardown.TeardownSuite()
	}

	for _, name := range skippedTestNames {
		testCase{t: t, manualSkip: true, name: name}.Run()
	}

	for _, name := range testNames {
		testCase{t, name, config, false, fixtureType, fixtureValue}.Run()
	}
}

type testCase struct {
	t            *TestingReporter
	name         string
	config       *config
	manualSkip   bool
	fixtureType  reflect.Type
	fixtureValue reflect.Value
}

func (this testCase) Run() {
	_ = this.t.Run(this.name, this.decideRun())
}
func (this testCase) decideRun() func(*testing.T) {
	if this.manualSkip {
		return this.skipFunc("Skipping: " + this.name)
	}

	if isLongRunning(this.name) && testing.Short() {
		return this.skipFunc("Skipping long-running test in -test.short mode: " + this.name)
	}

	return this.runTest
}
func (this testCase) skipFunc(message string) func(*testing.T) {
	return func(t *testing.T) { t.Skip(message) }
}
func (this testCase) runTest(t *testing.T) {
	if this.config.parallelTests {
		t.Parallel()
	}

	fixtureValue := this.fixtureValue
	if this.config.freshFixture {
		fixtureValue = reflect.New(this.fixtureType.Elem())
	}
	fixtureValue.Elem().FieldByName("T").Set(reflect.ValueOf(New(t)))

	setup, hasSetup := fixtureValue.Interface().(setupTest)
	if hasSetup {
		setup.Setup()
	}

	teardown, hasTeardown := fixtureValue.Interface().(teardownTest)
	if hasTeardown {
		defer teardown.Teardown()
	}

	fixtureValue.MethodByName(this.name).Call(nil)
}
func isLongRunning(name string) bool {
	return strings.HasPrefix(name, "Long") ||
		strings.HasPrefix(name, "FocusLong")
}

type (
	setupSuite    interface{ SetupSuite() }
	setupTest     interface{ Setup() }
	teardownTest  interface{ Teardown() }
	teardownSuite interface{ TeardownSuite() }
)

// FILE: run_options.go

type config struct {
	freshFixture    bool
	parallelFixture bool
	parallelTests   bool
}
type Option func(*config)
type Opt struct{}

var Options Opt

func (Opt) FreshFixture() Option {
	return func(c *config) {
		c.freshFixture = true
	}
}
func (Opt) SharedFixture() Option {
	return func(c *config) {
		c.freshFixture = false
		c.parallelTests = false
		c.parallelFixture = false
	}
}
func (Opt) ParallelFixture() Option {
	return func(c *config) {
		c.parallelFixture = true
	}
}
func (Opt) ParallelTests() Option {
	return func(c *config) {
		c.parallelTests = true
		c.freshFixture = true
		Options.FreshFixture()(c)
	}
}
func (Opt) UnitTests() Option {
	return func(c *config) {
		Options.ParallelTests()(c)
		Options.ParallelFixture()(c)
	}
}
func (Opt) IntegrationTests() Option {
	return func(c *config) {
		Options.SharedFixture()(c)
	}
}

// FILE: so.go

func So(t *testing.T, actual any, assertion Func, expected ...any) {
	_ = New(t).So(actual, assertion, expected...)
}

// FILE: spec.go

type specification interface {
	assertable(a, b any) bool
	passes(a, b any) bool
}

// FILE: start_with.go

func StartWith(actual any, expected ...any) error {
	err := validateExpected(1, expected)
	if err != nil {
		return err
	}

	err = validateKind(actual, orderedContainerKinds...)
	if err != nil {
		return err
	}

	actualValue := reflect.ValueOf(actual)
	EXPECTED := expected[0]

	switch reflect.TypeOf(actual).Kind() {
	case reflect.Array, reflect.Slice:
		if actualValue.Len() == 0 {
			break
		}
		first := actualValue.Index(0).Interface()
		if Equal(EXPECTED, first) == nil {
			return nil
		}
	case reflect.String:
		err = validateKind(EXPECTED, reflect.String, reflectRune)
		if err != nil {
			return err
		}

		expectedRune, ok := EXPECTED.(rune)
		if ok {
			EXPECTED = string(expectedRune)
		}

		full := actual.(string)
		prefix := EXPECTED.(string)
		if strings.HasPrefix(full, prefix) {
			return nil
		}
	}

	return failure("\n"+
		"   proposed prefix: %#v\n"+
		"   not a prefix of: %#v",
		EXPECTED,
		actual,
	)
}

// FILE: wrap_error.go

func WrapError(actual any, expected ...any) error {
	err := validateExpected(1, expected)
	if err != nil {
		return err
	}

	inner, ok := expected[0].(error)
	if !ok {
		return errTypeMismatch(expected[0])
	}

	outer, ok := actual.(error)
	if !ok {
		return errTypeMismatch(actual)
	}

	if errors.Is(outer, inner) {
		return nil
	}

	return fmt.Errorf("%w:\n"+
		"\t            outer err: (%s)\n"+
		"\tshould wrap inner err: (%s)",
		ErrAssertionFailure,
		outer,
		inner,
	)
}
func errTypeMismatch(v any) error {
	return fmt.Errorf(
		"%w: got %s, want error",
		ErrTypeMismatch,
		reflect.TypeOf(v),
	)
}
