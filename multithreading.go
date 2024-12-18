package main

import (
	"context"
	"io"
	"net/http"
	"fmt"
	"time"
)

func make_request(ctx context.Context, chanel chan string, url string) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		fmt.Println(err)
		return 
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	chanel <- string(body)
}

func main() {
	ctxreq, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	cep := "01153000"

	URLBrasilAPI := "https://brasilapi.com.br/api/cep/v1/01153000" + cep
	URLViaCEP := "http://viacep.com.br/ws/" + cep + "/json/"

	brasilapi := make(chan string)
	defer close(brasilapi)
	viacep := make(chan string)
	defer close(viacep)


	go make_request(ctxreq, brasilapi, URLBrasilAPI)
	go make_request(ctxreq, viacep, URLViaCEP)

	select{
		case <- ctxreq.Done():
			fmt.Println("Timeout")
		case res1 := <- brasilapi:
			fmt.Println(res1)
		case res2 := <- viacep:
			fmt.Println(res2)
	}
}
