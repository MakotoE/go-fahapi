[![GoDoc](https://godoc.org/github.com/MakotoE/go-fahapi?status.svg)](https://pkg.go.dev/github.com/MakotoE/go-fahapi)

# Folding@home client API wrapper for Go

This is a work in progress!

```
func ExampleAPI_PauseAll() {
	api, err := NewAPI()
	if err != nil {
		panic(err)
	}

	defer api.Close()

	if err := api.PauseAll(); err != nil {
		panic(err)
	}

	if err := api.UnpauseAll(); err != nil {
		panic(err)
	}
}
```