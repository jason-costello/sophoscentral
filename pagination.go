package sophoscentral

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"strconv"
)
type Fields struct{
	ctx                   context.Context
	httpClient            *http.Client
	lastSuccessfulRequest *http.Request
	pages                 Pages
}

func NewPaginationFields(ctx context.Context, httpClient *http.Client, lastRequest *http.Request, pages Pages) (Fields, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	if lastRequest == nil {
		return Fields{}, errors.New("must include last request")
	}

	if pages.Total == 0 || pages.Items == 0 {
		return Fields{}, errors.New("must include pageTotal=true in initial request query args")
	}
return Fields{
	ctx:                   ctx,
	httpClient:            httpClient,
	lastSuccessfulRequest: lastRequest,
	pages:                pages,
}, nil
}
func  GetRemainingPages(pf Fields) ([][]byte, error) {
	defer fmt.Printf("Done")
	var pageBytes [][]byte

	urls := extractAllURLs(pf.lastSuccessfulRequest.URL.String(), pf.pages)

	// func generateRequest(ctx context.Context, cURL string, headers map[string]string, method string, body []byte )(*http.Request, error){

  generator := func(done <-chan interface{}, urls ...string) <-chan string{
  	urlStream := make(chan string)
  	go func(){
  		defer close (urlStream)
  		for _, s := range urls{
  			select {
			case <-done:
				return
				case urlStream <- s:
			}
			}

	}()
  	return urlStream
	}

	buildRequest := func(done <- chan interface{},
		urlStream <- chan string,
		) <- chan *http.Request{
		reqStream := make(chan *http.Request)
		go func(){
			defer close(reqStream)
			for reqUrl := range urlStream{
				select{
				case <- done:
				return
				case reqStream <- generateRequest(pf.ctx, pf.lastSuccessfulRequest, reqUrl, nil):
				}
			}
		}()
		return reqStream
	}


	makeRequest := func(done <- chan interface{},
		reqStream <- chan *http.Request,
		) <- chan []byte{
		byteStream := make(chan []byte)
		go func(){
			defer close(byteStream)
			for  req := range reqStream{
				select {

				case <-done:
					return
				case byteStream <- MustMakeRequest(pf.httpClient, req):
				}
			}
		}()
		return byteStream
	}



	done := make(chan interface{})
	defer close(done)
	urlStream := generator(done, urls...)
	pipeline := makeRequest(done, buildRequest(done, urlStream))

	for v := range pipeline{
		fmt.Printf("%s\n", v)
		pageBytes = append(pageBytes, v)
	}


return pageBytes, nil
}
func extractAllURLs(currentURL string, p Pages) []string{

	if currentURL == ""{
		return nil
	}

	cURL, err := url.Parse(currentURL)
	if err != nil{
		return nil
	}

	remainingPages := getRemainingPageCount(p.Total, p.Current, p.MaxSize)
	var urls []string

	for i := p.Current +1; i < remainingPages; i++{
		q := cURL.Query()

		q.Set("page", strconv.Itoa(i))
		cURL.RawQuery = q.Encode()
		urls = append(urls, cURL.String())
	}

	return urls


}

func generateRequest(ctx context.Context, lastSuccessfulRequest *http.Request, reqUrl string, body []byte )*http.Request{

	var req *http.Request
	var err error
	if body == nil {
		if len(body) < 1{
			req, err = http.NewRequestWithContext(ctx, lastSuccessfulRequest.Method, reqUrl, nil)
			}
	} else {
		req, err = http.NewRequestWithContext(ctx, lastSuccessfulRequest.Method,reqUrl, bytes.NewBuffer(body))
	}
	if err != nil{
		return nil
	}
		req.Header = lastSuccessfulRequest.Header

	return req
}

func getRemainingPageCount(ttlItems, currItems, maxReturn int) int{
	if currItems >= ttlItems{
		return 0
	}
	remainingItemsCount := ttlItems - currItems
	if math.Mod(float64(remainingItemsCount), float64(maxReturn))==0{
		return int(float64(remainingItemsCount)/float64(maxReturn))
	}
		return  int(math.Ceil(float64(remainingItemsCount)/float64(maxReturn)))
}
