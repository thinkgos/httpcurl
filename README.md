# httpcurl

Transform Golang's http.Request to cURL command line.

[![Go.Dev reference](https://img.shields.io/badge/go.dev-reference-blue?logo=go&logoColor=white)](https://pkg.go.dev/github.com/thinkgos/httpcurl?tab=doc)
[![codecov](https://codecov.io/gh/thinkgos/httpcurl/graph/badge.svg?token=aHu5wq1m6i)](https://codecov.io/gh/thinkgos/httpcurl)
[![Tests](https://github.com/thinkgos/httpcurl/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/thinkgos/httpcurl/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/thinkgos/httpcurl)](https://goreportcard.com/report/github.com/thinkgos/httpcurl)
[![License](https://img.shields.io/github/license/thinkgos/httpcurl)](https://raw.githubusercontent.com/thinkgos/httpcurl/main/LICENSE)
[![Tag](https://img.shields.io/github/v/tag/thinkgos/httpcurl)](https://github.com/thinkgos/httpcurl/tags)

## Usage

Use go get.

```bash
go get github.com/thinkgos/httpcurl
```

Then import the package into your own code.

```go
import "github.com/thinkgos/httpcurl"
```

## Example

```go
import (
    "fmt"
    "github.com/thinkgos/httpcurl"
)

func main() {
	form := url.Values{}
	form.Add("age", "10")
	form.Add("name", "John")
	body := form.Encode()

	req, _ := http.NewRequest(http.MethodPost, "http://example/cats", io.NopCloser(bytes.NewBufferString(body)))
	req.Header.Add("API_KEY", "123")

	curl, _ := httpcurl.IntoCurl(req)
    fmt.Println(curl) // curl 'http://example/cats' -X 'POST' -H 'Api_key: 123' -d 'age=10&name=John' --compressed
}
```

## Reference

## License

This project is under MIT License. See the [LICENSE](LICENSE) file for the full license text.
