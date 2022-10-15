package main

import (
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		expiration := time.Now().Add(time.Hour)
		cookie := http.Cookie{Name: "cookie", Value: "value", Expires: expiration}
		http.SetCookie(writer, &cookie)
	})

	http.HandleFunc("/read", func(writer http.ResponseWriter, request *http.Request) {
		cookie, err := request.Cookie("cookie")
		if err != nil {
			writer.WriteHeader(http.StatusNotFound)
			return
		}
		writer.Write([]byte(cookie.Value))
	})

	http.ListenAndServe(":3000", nil)
}