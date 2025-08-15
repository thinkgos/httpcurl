package httpcurl

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
)

func Test_IntoCurl(t *testing.T) {
	form := url.Values{}
	form.Add("age", "10")
	form.Add("name", "John")
	body := form.Encode()

	req, _ := http.NewRequest(http.MethodPost, "http://example/cats", io.NopCloser(bytes.NewBufferString(body))) // nolint: noctx
	req.Header.Add("API_KEY", "123")

	got, _ := IntoCurl(req)
	want := `curl 'http://example/cats' -X 'POST' -H 'Api_key: 123' -d 'age=10&name=John' --compressed`
	if got != want {
		t.Errorf("%s\ngot: %v\nwant: %v\n", t.Name(), got, want)
	}
}

func Test_IntoCurl_JSON(t *testing.T) {
	req, _ := http.NewRequestWithContext(context.Background(), "PUT", "http://www.example.com/abc/def.png?jlk=mno&pqr=stu", bytes.NewBufferString(`{"hello":"world","answer":42}`))
	req.Header.Set("Content-Type", "application/json")

	got, _ := IntoCurl(req)
	want := `curl 'http://www.example.com/abc/def.png?jlk=mno&pqr=stu' -X 'PUT' -H 'Content-Type: application/json' -d '{"hello":"world","answer":42}' --compressed`
	if got != want {
		t.Errorf("%s\ngot: %v\nwant: %v\n", t.Name(), got, want)
	}
}

func Test_IntoCurl_Separator(t *testing.T) {
	req, _ := http.NewRequest("PUT", "http://www.example.com/abc/def.png?jlk=mno&pqr=stu", bytes.NewBufferString(`{"hello":"world","answer":42}`)) // nolint: noctx
	req.Header.Set("Content-Type", "application/json")

	got, _ := New(WithSeparator(" \\\n"), WithDumpRequestBody(dumpRequestBody)).IntoCurl(req)
	want := `curl \
 'http://www.example.com/abc/def.png?jlk=mno&pqr=stu' \
 -X 'PUT' \
 -H 'Content-Type: application/json' \
 -d '{"hello":"world","answer":42}' \
 --compressed`
	if got != want {
		t.Errorf("%s\ngot: %v\nwant: %v\n", t.Name(), got, want)
	}
}

func Test_IntoCurl_NoBody(t *testing.T) {
	req, _ := http.NewRequest("GET", "http://www.example.com/abc/def.png?jlk=mno&pqr=stu", nil) // nolint: noctx
	req.Header.Set("Content-Type", "application/json")

	got, _ := IntoCurl(req)
	want := `curl 'http://www.example.com/abc/def.png?jlk=mno&pqr=stu' -X 'GET' -H 'Content-Type: application/json' --compressed`
	if got != want {
		t.Errorf("%s\ngot: %v\nwant: %v\n", t.Name(), got, want)
	}
}

func Test_IntoCurl_EmptyStringBody(t *testing.T) {
	req, _ := http.NewRequestWithContext(context.Background(), "PUT", "http://www.example.com/abc/def.png?jlk=mno&pqr=stu", bytes.NewBufferString(""))
	req.Header.Set("Content-Type", "application/json")

	got, _ := IntoCurl(req)
	want := `curl 'http://www.example.com/abc/def.png?jlk=mno&pqr=stu' -X 'PUT' -H 'Content-Type: application/json' --compressed`
	if got != want {
		t.Errorf("%s\ngot: %v\nwant: %v\n", t.Name(), got, want)
	}
}

func Test_IntoCurl_NewlineInBody(t *testing.T) {
	req, _ := http.NewRequestWithContext(context.Background(), "POST", "http://www.example.com/abc/def.png?jlk=mno&pqr=stu", bytes.NewBufferString("hello\nworld"))
	req.Header.Set("Content-Type", "application/json")

	got, _ := IntoCurl(req)
	want := `curl 'http://www.example.com/abc/def.png?jlk=mno&pqr=stu' -X 'POST' -H 'Content-Type: application/json' -d 'hello
world' --compressed`
	if got != want {
		t.Errorf("%s\ngot: %v\nwant: %v\n", t.Name(), got, want)
	}
}

func Test_IntoCurl_SpecialCharsInBody(t *testing.T) {
	req, _ := http.NewRequest("POST", "http://www.example.com/abc/def.png?jlk=mno&pqr=stu", bytes.NewBufferString(`Hello $123 o'neill -"-`)) // nolint: noctx
	req.Header.Set("Content-Type", "application/json")

	got, _ := IntoCurl(req)
	want := `curl 'http://www.example.com/abc/def.png?jlk=mno&pqr=stu' -X 'POST' -H 'Content-Type: application/json' -d 'Hello $123 o'\''neill -"-' --compressed`
	if got != want {
		t.Errorf("%s\ngot: %v\nwant: %v\n", t.Name(), got, want)
	}
}

func Test_IntoCurl_Other(t *testing.T) {
	uri := "http://www.example.com/abc/def.png?jlk=mno&pqr=stu"
	payload := new(bytes.Buffer)
	payload.Write([]byte(`{"hello":"world","answer":42}`))
	req, err := http.NewRequestWithContext(context.Background(), "PUT", uri, payload)
	if err != nil {
		panic(err)
	}
	req.Header.Set("X-Auth-Token", "private-token")
	req.Header.Set("Content-Type", "application/json")

	got, _ := IntoCurl(req)
	want := `curl 'http://www.example.com/abc/def.png?jlk=mno&pqr=stu' -X 'PUT' -H 'Content-Type: application/json' -H 'X-Auth-Token: private-token' -d '{"hello":"world","answer":42}' --compressed`
	if got != want {
		t.Errorf("%s\ngot: %v\nwant: %v\n", t.Name(), got, want)
	}
}

func Test_IntoCurl_Https(t *testing.T) {
	uri := "https://www.example.com/abc/def.png?jlk=mno&pqr=stu"
	payload := new(bytes.Buffer)
	payload.Write([]byte(`{"hello":"world","answer":42}`))
	req, err := http.NewRequest("PUT", uri, payload) // nolint: noctx
	if err != nil {
		panic(err)
	}
	req.Header.Set("X-Auth-Token", "private-token")
	req.Header.Set("Content-Type", "application/json")

	got, _ := IntoCurl(req)
	want := `curl 'https://www.example.com/abc/def.png?jlk=mno&pqr=stu' -k -X 'PUT' -H 'Content-Type: application/json' -H 'X-Auth-Token: private-token' -d '{"hello":"world","answer":42}' --compressed`
	if got != want {
		t.Errorf("%s\ngot: %v\nwant: %v\n", t.Name(), got, want)
	}
}

func Test_ServerSide(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s, err := IntoCurl(r)
		if err != nil {
			t.Error(err)
		}
		_, _ = fmt.Fprint(w, s)
	}))
	defer svr.Close()

	url := svr.URL + "/?a=b"
	resp, err := http.Get(url) // nolint: noctx
	if err != nil {
		t.Error(err)
	}
	defer resp.Body.Close() // nolint: errcheck
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	want := fmt.Sprintf("curl '%s' -X 'GET' -H 'Accept-Encoding: gzip' -H 'User-Agent: Go-http-client/1.1' --compressed", url)
	if out := string(data); out != want {
		t.Errorf("got: %s, want: %s", out, want)
	}
}

func Benchmark_IntoCurl(b *testing.B) {
	form := url.Values{}

	for i := 0; i <= b.N; i++ {
		form.Add("number", strconv.Itoa(i))
		body := form.Encode()
		req, _ := http.NewRequest(http.MethodPost, "http://example", io.NopCloser(bytes.NewBufferString(body))) // nolint: noctx
		_, err := IntoCurl(req)
		if err != nil {
			panic(err)
		}
	}
}
