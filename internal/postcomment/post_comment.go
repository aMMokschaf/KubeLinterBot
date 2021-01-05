package postcomment

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func PostComment() {
	fmt.Println("post comment method")
	var input string = "testmsg"
	fmt.Println(input)
	testfunc()
}

func testfunc() {
	requestBody, err := json.Marshal(map[string]string{
		"name": "Marcel",
		"mail": "blabla@bla.de",
	})
	if err != nil {
		log.Fatalln(err)
	}

	resp, err := http.Post("https://api.github.com/repos/aMMokschaf/yamls/commits/6c052db5bb1ce419adb725ad4e726a3548149452/comments", "application/vnd.github.v3+json", bytes.NewBuffer(requestBody))
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(string(body))
}
