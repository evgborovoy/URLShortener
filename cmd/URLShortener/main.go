package main

import "URLShortener/internal/config"

func main() {
	config := config.MustLoad()
	_ = config

	// init logger: slog
	// init storage: sqlite
	// init router: chi, "chi render"
	// run server

}
