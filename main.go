package main

import (
	"fmt"
	"io"
	"joubertredrat/go-env-management/pkg"
	"net/http"
)

func main() {
	config, err := pkg.GetConfig()
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, dsn(
			config.DatabaseHost,
			config.DatabasePort,
			config.DatabaseUsername,
			config.DatabasePassword,
			config.DatabaseDBName,
			string(config.DatabaseSSLMode),
		))
	})

	listen := fmt.Sprintf("%s:%s", config.ApiHost, config.ApiPort)
	fmt.Printf("Running app %s at %s\n", config.ApiEnv, listen)

	if err := http.ListenAndServe(listen, nil); err != nil {
		panic(err)
	}
}

func dsn(host, port, user, pass, name, sslmode string) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", user, pass, host, port, name, sslmode)
}
