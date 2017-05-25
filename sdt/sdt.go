package sdt

import (
	"net/http"
	"bytes"
	"crypto/tls"
	"strings"
	"io"
	"encoding/xml"
	"time"
)

type Sdt struct {
	ApiUrl    string
	PartnerId string
	Password  string
}

type resultReader interface {
	decodeXml(apiResult) error
}

type apiResults struct {
	Results []apiResult
}

type apiResult struct {
	Type      int
	Body      io.Reader
	Err       error
	ServiceId int
}

type requests struct {
	Request []request
}

type request struct {
	Partner_id  string      `xml:"partner_id,attr"`
	Password    string      `xml:"password,attr"`
	Parcel      interface{} `xml:"parcel"`
	RequestType int         `xml:"request_type,attr"`
	ServiceId   int         `xml:"-"`
}

func Init(url string, id string, password string) Sdt {
	return Sdt{
		ApiUrl: url,
		PartnerId: id,
		Password: password,
	}
}

func (sdt *Sdt)sendRequests(req requests) (apiResults, error) {
	var results apiResults

	ch := make(chan apiResult)

	for _, r := range req.Request {
		tmp := struct {
			request
			XMLName struct{}  `xml:"request"`
		}{request: r}

		res, err := xml.MarshalIndent(tmp, "", "    ")
		if err != nil {
			return results, err
		}

		go makePost(r.RequestType, r.ServiceId, string(res), sdt.ApiUrl, ch)
	}

	for range req.Request {
		r := <-ch
		results.Results = append(results.Results, r)
	}

	return results, nil
}

func makePost(rt int, sid int, postData string, urlPath string, ch chan <-apiResult) {
	request, err := http.NewRequest("POST", urlPath, bytes.NewBuffer([]byte(postData)))

	if err != nil {
		ch <- apiResult{Err: err, Type: rt, ServiceId: sid}
		return
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	timeout := time.Duration(30 * time.Second)
	client := &http.Client{Transport: tr, Timeout: timeout}
	resp, err := client.Do(request)

	if err != nil {
		ch <- apiResult{Err: err, Type: rt, ServiceId: sid}
		return
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	s := buf.String()
	s = strings.Replace(s, "WINDOWS-1251", "UTF-8", 1)
	b := bytes.NewReader([]byte(s))

	ch <- apiResult{Body: b, Err: nil, Type: rt, ServiceId: sid}
}
