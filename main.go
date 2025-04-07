package main

import (
	"encoding/json"
	"fmt"
	_ "github.com/hashicorp/go-msgpack/codec"
	"io"
	"net/http"
)

//TIP To run your code, right-click the code and select <b>Run</b>. Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.

func main() {
	http.HandleFunc("/", handler)
	fmt.Println(http.ListenAndServe(":80", nil))
}

type js struct {
	Age  int    `json:"age"`
	Name string `json:"name"`
}

type answer struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	bytes, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, "pizdaparse:"+err.Error(), http.StatusInternalServerError)
		return
	}
	j := js{}
	err = json.Unmarshal(bytes, &j)
	if err != nil {
		http.Error(w, "pizdaunmarshal:"+err.Error(), http.StatusInternalServerError)
	}

	answ := answer{
		Name: j.Name,
	}

	switch j.Name {
	case "Artur":
		answ.Status = "Mobile"
	case "Ivan":
		answ.Status = "Backend"
	default:
		answ.Status = "Gay"
	}

	bytes, err = json.Marshal(answ)
	if err != nil {
		http.Error(w, "marshal:"+err.Error(), http.StatusInternalServerError)
	}

	_, err = w.Write(bytes)
	if err != nil {
		http.Error(w, "write:"+err.Error(), http.StatusInternalServerError)
	}
}

//TIP See GoLand help at <a href="https://www.jetbrains.com/help/go/">jetbrains.com/help/go/</a>.
// Also, you can try interactive lessons for GoLand by selecting 'Help | Learn IDE Features' from the main menu.
