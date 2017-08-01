// githubSearch project main.go
package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/bitly/go-simplejson"
)

func index(w http.ResponseWriter, req *http.Request) {

	io.WriteString(w, "<HTML><form method=\"POST\" action=\"search\">GitHubSearch:<input name=\"target\">")
}

func Search(w http.ResponseWriter, req *http.Request) {
	target := req.PostFormValue("target")
	url := fmt.Sprintf("https://api.github.com/search/code?access_token=xxxx&q=%s", target)
	fmt.Print(url + "\n")
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
		return
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	js, err := simplejson.NewJson(body)
	if err != nil {
		panic(err.Error())
	}
	data, _ := js.Map()
	result := fmt.Sprintf("总数:%s", data["total_count"])
	items := js.Get("items")
	io.WriteString(w, "<html>")
	io.WriteString(w, result)

	for i := 0; i < 100; i++ {
		io.WriteString(w, "<br>")
		fmt.Println(items.GetIndex(i).Get("repository").Get("full_name"))
		rres, err := items.GetIndex(i).Get("repository").Get("full_name").String()
		rurl, _ := items.GetIndex(i).Get("path").String()
		if err != nil {
			break
		}
		io.WriteString(w, "<a href=\"https://github.com/")
		io.WriteString(w, rres)
		io.WriteString(w, "/tree/master/")
		io.WriteString(w, rurl)
		io.WriteString(w, "\" target=\"_blank\">")
		io.WriteString(w, rres)
		io.WriteString(w, "/")
		io.WriteString(w, rurl)
		io.WriteString(w, "</a>")
	}

}

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/search", Search)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
