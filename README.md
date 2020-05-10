[![GoDoc](https://godoc.org/github.com/MakotoE/go-fahapi?status.svg)](https://godoc.org/github.com/MakotoE/go-fahapi)

# Folding@home client API wrapper for Go

```
package main

import "github.com/MakotoE/go-fahapi"

func main() {
	api, err := fahapi.NewAPI(fahapi.DefaultAddr)
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