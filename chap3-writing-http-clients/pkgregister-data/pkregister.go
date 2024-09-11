// chap3/pkgregister-data/pkgregister.go
package pkgregister

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

func registerPackageData(client *http.Client, url string, data pkgData) (pkgRegisterResult, error) {
	p := pkgRegisterResult{}
	payload, contentType, err := createMultiPartMessage(data)
	if err != nil {
		return p, nil
	}

	reader := bytes.NewReader(payload)
	r, err := http.Post(url, contentType, reader)
	if err != nil {
		return p, nil
	}
	defer r.Body.Close()

	respData, err := io.ReadAll(r.Body)
	if err != nil {
		return p, nil
	}

	err = json.Unmarshal(respData, &p)
	return p, err
}

func createHTTPClientWithTimeout(d time.Duration) *http.Client {
	client := http.Client{Timeout: d}
	return &client
}