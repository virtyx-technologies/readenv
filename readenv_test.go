package readenv

import (
	"fmt"
	"log"
	"os"
	"testing"
)

func TestNotWriteable(t *testing.T) {
	type testOpts struct {
		notWriteable string `env:"NOT_WRITEABLE"`
	}
	os.Setenv("NOT_WRITEABLE", "hello")
	opts := &testOpts{}
	err := ReadEnv(opts)
	if err == nil {
		t.Error("readenv with unexported field should fail")
	}
}

func TestFloats(t *testing.T) {
	type testOpts struct {
		Float64Type float64 `env:"E_FLOAT64"`
		Float32Type float32 `env:"E_FLOAT32"`
	}
	os.Setenv("E_FLOAT64", "3.14")
	os.Setenv("E_FLOAT32", "3.14")
	opts := &testOpts{}
	err := ReadEnv(opts)
	if err != nil {
		t.Errorf("readenv failed: %v", err)
	}
	if opts.Float64Type != 3.14 {
		t.Errorf("float64 should have been 3.14 but was %f", opts.Float64Type)
	}
	if opts.Float32Type != 3.14 {
		t.Errorf("float32 should have been 3.14 but was %f", opts.Float32Type)
	}
}

func TestInts(t *testing.T) {
	type testOpts struct {
		Int64Type int64 `env:"E_INT64"`
		Int32Type int32 `env:"E_INT32"`
		Int16Type int16 `env:"E_INT16"`
		Int8Type  int8  `env:"E_INT8"`
		IntType   int   `env:"E_INT"`
	}
	os.Setenv("E_INT64", "42")
	os.Setenv("E_INT32", "42")
	os.Setenv("E_INT16", "42")
	os.Setenv("E_INT8", "42")
	os.Setenv("E_INT", "42")
	opts := &testOpts{}
	err := ReadEnv(opts)
	if err != nil {
		t.Errorf("readenv failed: %v", err)
	}
	if opts.Int64Type != 42 {
		t.Errorf("int64 should have been 42 but was %d", opts.Int64Type)
	}
	if opts.Int32Type != 42 {
		t.Errorf("int32 should have been 42 but was %d", opts.Int32Type)
	}
	if opts.Int16Type != 42 {
		t.Errorf("int16 should have been 42 but was %d", opts.Int16Type)
	}
	if opts.Int8Type != 42 {
		t.Errorf("int8 should have been 42 but was %d", opts.Int8Type)
	}
	if opts.IntType != 42 {
		t.Errorf("int should have been 42 but was %d", opts.IntType)
	}
}

func TestString(t *testing.T) {
	type testOpts struct {
		StringType string `env:"E_STRING"`
	}
	os.Setenv("E_STRING", "string")
	opts := &testOpts{}
	err := ReadEnv(opts)
	if err != nil {
		t.Errorf("readenv failed: %v", err)
	}
	if opts.StringType != "string" {
		t.Errorf("string should have been 'string' but was %s", opts.StringType)
	}
}

func TestBool(t *testing.T) {
	type testOpts struct {
		BoolEmpty bool `env:"E_EMPTY"`
		BoolNo    bool `env:"E_NO"`
		BoolOff   bool `env:"E_OFF"`
		BoolZero  bool `env:"E_ZERO"`
		BoolElse  bool `env:"E_TRUE"`
	}
	os.Setenv("E_EMPTY", "")
	os.Setenv("E_NO", "NO")
	os.Setenv("E_OFF", "off")
	os.Setenv("E_ZERO", "0")
	os.Setenv("E_TRUE", "1")
	opts := &testOpts{}
	err := ReadEnv(opts)
	if err != nil {
		t.Errorf("readenv failed: %v", err)
	}
	if opts.BoolEmpty != false {
		t.Error("empty environment variable was not set to false")
	}
	if opts.BoolNo != false {
		t.Error("environment variable 'no' was not set to false")
	}
	if opts.BoolOff != false {
		t.Error("environment variable 'off' was not set to false")
	}
	if opts.BoolZero != false {
		t.Error("environment variable '0' was not set to false")
	}
	if opts.BoolElse != true {
		t.Error("environment variable '1' was not set to true")
	}
}

func TestReadToStruct(t *testing.T) {
	type testOpts struct {
		F string `env:"F"`
	}
	opts := &testOpts{}
	if err := ReadEnv(&opts); err == nil {
		t.Error("reading into **struct did not return an error")
	}
	if err := ReadEnv(*opts); err == nil {
		t.Error("reading into struct did not return an error")
	}
	if err := ReadEnv("asdf"); err == nil {
		t.Error("reading into string did not return an error")
	}
}

func TestReadNoTags(t *testing.T) {
	type testOpts struct {
		F string
	}
	opts := &testOpts{}
	if err := ReadEnv(opts); err != nil {
		t.Errorf("got an error reading opts: %v", err)
	}
	if opts.F != "" {
		t.Errorf("field was set to %s despite not being tagged", opts.F)
	}
}

func TestReadBadInt(t *testing.T) {
	type testOpts struct {
		I int `env:"I"`
	}
	opts := &testOpts{}
	if err := ReadEnv(opts); err == nil {
		t.Error("should have gotten an error but did not")
	}
}

func TestReadBadFloat(t *testing.T) {
	type testOpts struct {
		F float64 `env:"F"`
	}
	opts := &testOpts{}
	if err := ReadEnv(opts); err == nil {
		t.Error("should have gotten an error but did not")
	}
}

func TestReadBadString(t *testing.T) {
	type testOpts struct {
		S string `env:"S"`
	}
	opts := &testOpts{}
	if err := ReadEnv(opts); err == nil {
		t.Error("should have gotten an error but did not")
	}
}

func ExampleReadEnv() {
	os.Setenv("PORT", "8000")
	os.Setenv("DEBUG", "1")
	type MyOpts struct {
		Port  int  `env:"PORT"`
		Debug bool `env:"DEBUG"`
	}
	opts := &MyOpts{}
	if err := ReadEnv(opts); err != nil {
		log.Fatalf("could not read config: %v", err)
	}
	fmt.Printf("port=%d, debug=%t", opts.Port, opts.Debug)
	// Output: port=8000, debug=true
}
