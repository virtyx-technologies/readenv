# readenv

> easily read environment variables into a struct

[![GoDoc](https://godoc.org/github.com/virtyx-technologies/readenv?status.png)](https://godoc.org/github.com/virtyx-technologies/readenv)

Example:

```go
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
```
