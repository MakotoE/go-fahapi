package fahapi

func Example() {
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
