package httpcurl

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
)

var defaultHttpCurl = New()

// IntoCurl returns a curl corresponding to an http.Request use default option.
func IntoCurl(req *http.Request) (string, error) {
	return defaultHttpCurl.IntoCurl(req)
}

// Option HttpCurl option
type Option func(h *HttpCurl)

// HttpCurl http curl instance.
type HttpCurl struct {
	sep             string
	dumpRequestBody func(io.ReadCloser) (io.ReadCloser, string, error)
}

// WithSeparator set separator, default is no separator.
func WithSeparator(sep string) Option {
	return func(h *HttpCurl) {
		h.sep = sep
	}
}

// WithDumpRequestBody dump request body. default dump.
func WithDumpRequestBody(dumpRequestBody func(b io.ReadCloser) (io.ReadCloser, string, error)) Option {
	return func(h *HttpCurl) {
		if dumpRequestBody != nil {
			h.dumpRequestBody = dumpRequestBody
		}
	}
}

// New a new HttpCurl.
func New(opts ...Option) *HttpCurl {
	h := &HttpCurl{
		sep:             "",
		dumpRequestBody: dumpRequestBody,
	}
	for _, f := range opts {
		f(h)
	}
	return h
}

// IntoCurl returns a curl corresponding to an http.Request
func (h *HttpCurl) IntoCurl(req *http.Request) (string, error) {
	if req.URL == nil {
		return "", fmt.Errorf("httpcurl(IntoCurl): invalid request, req.URL is nil")
	}
	b := builder{
		buf: strings.Builder{},
		sep: h.sep,
	}
	b.buf.Grow(256)
	b.WriteLine("curl")

	schema := req.URL.Scheme
	url := req.URL.String()
	if schema == "" {
		schema = "http"
		if req.TLS != nil {
			schema = "https"
		}
		url = schema + "://" + req.Host + req.URL.RequestURI()
	}
	b.WriteLine(bashEscape(url))
	if schema == "https" {
		b.WriteLine("-k")
	}
	b.WriteLine("-X", bashEscape(req.Method))

	keys := make([]string, 0, len(req.Header))
	for k := range req.Header {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		b.WriteLine("-H", bashEscape(fmt.Sprintf("%s: %s", k, strings.Join(req.Header[k], " "))))
	}

	if req.Body != nil {
		rc, reqBody, err := h.dumpRequestBody(req.Body)
		if err != nil {
			return "", err
		}
		req.Body = rc
		if reqBody != "" {
			b.WriteLine("-d", bashEscape(reqBody))
		}
	}
	b.WriteLine("--compressed")

	return b.String(), nil
}

func dumpRequestBody(b io.ReadCloser) (io.ReadCloser, string, error) {
	if b == nil || b == http.NoBody {
		// No copying needed. Preserve the magic sentinel meaning of NoBody.
		return http.NoBody, "", nil
	}
	var buff bytes.Buffer
	_, err := buff.ReadFrom(b)
	if err != nil {
		return nil, "", fmt.Errorf("httpcurl: buffer read from request body error, %w", err)
	}
	if err = b.Close(); err != nil {
		return nil, "", err
	}
	return io.NopCloser(bytes.NewReader(buff.Bytes())), buff.String(), nil
}

type builder struct {
	buf strings.Builder
	sep string
}

func (b *builder) String() string {
	return b.buf.String()
}

func (b *builder) WriteLine(vs ...string) *builder {
	length := b.buf.Len()
	if length > 0 && b.sep != "" {
		b.buf.WriteString(b.sep)
	}
	for _, v := range vs {
		if length > 0 {
			b.buf.WriteString(" ")
		}
		b.buf.WriteString(v)
	}
	return b
}

func bashEscape(str string) string {
	return `'` + strings.ReplaceAll(str, `'`, `'\''`) + `'`
}
