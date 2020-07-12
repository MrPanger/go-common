package http

import (
	"crypto/tls"
	"errors"
	"go/common/log"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

//HTTP Client工具
func SendHttpReq(host, uri, method, body string, isHttps bool, headers map[string]string, timeout int) ([]byte, int, error) {
	tr := &http.Transport{
		TLSClientConfig:   &tls.Config{InsecureSkipVerify: isHttps},
		DisableKeepAlives: true,
	}
	httpClient := &http.Client{Transport: tr, Timeout: time.Duration(timeout) * time.Second}
	req, err := buildHttpReq(host, uri, method, strings.NewReader(body))
	if err != nil {
		log.Errorf("failed to create speed request, error: %s", err.Error())
		return []byte{0}, 0, err
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	sendTime := time.Now()
	log.Debugf("the http request method[%s], url is: %s", req.Method, req.URL.String())
	log.Debugf("the http request body is: %s", body)
	resp, err := httpClient.Do(req)
	if err != nil {
		log.Errorf("failed to post request to speed server, error: %s", err.Error())
		return []byte{0}, 0, err
	}
	processTime := time.Now().Sub(sendTime)
	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{0}, resp.StatusCode, err
	}
	defer resp.Body.Close()
	log.Debugf("Receive response after %v, the body is :%s, the code is: %d.", processTime, buf, resp.StatusCode)

	return buf, resp.StatusCode, nil
}

func buildHttpReq(speedSrvAddr, path, method string, body io.Reader) (*http.Request, error) {
	if len(path) == 0 || len(method) == 0 {
		return nil, errors.New("invalid path or method")
	}
	host := speedSrvAddr
	if strings.HasPrefix(host, "http://") {
		host = strings.TrimPrefix(host, "http://")
	}
	if strings.HasPrefix(host, "https://") {
		host = strings.TrimPrefix(host, "https://")
	}

	sep := strings.Index(host, "/")
	if sep != -1 {
		host = strings.Split(host, "/")[0]
	}

	req, err := http.NewRequest(method, speedSrvAddr+path, body)
	if err != nil {
		return nil, err
	}

	return req, nil
}
