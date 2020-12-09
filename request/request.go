package request

import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strings"

	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
	"golang.org/x/text/transform"
)

type Request struct {
	baseUrl *url.URL

	header http.Header

	defaultClient *http.Client
}

// Get请求
func (r *Request) Get(path string, query url.Values) (resp *http.Response, err error) {
	if strings.HasPrefix(path, "http") {
		r.baseUrl, _ = url.Parse(path)
	} else {
		r.baseUrl.Path = path
	}

	req, err := http.NewRequest("GET", r.baseUrl.String(), nil)
	if err != nil {
		return nil, err
	}

	if r.header != nil {
		req.Header = r.header
	}

	req.URL.RawQuery = query.Encode()

	return r.defaultClient.Do(req)
}

// Post请求
func (r *Request) Post(path string, body io.Reader) (resp *http.Response, err error) {
	if strings.HasPrefix(path, "http") {
		r.baseUrl, _ = url.Parse(path)
	} else {
		r.baseUrl.Path = path
	}

	req, err := http.NewRequest("POST", r.baseUrl.String(), nil)
	if err != nil {
		return nil, err
	}

	if r.header != nil {
		req.Header = r.header
	}

	rc, ok := body.(io.ReadCloser)
	if !ok && body != nil {
		rc = ioutil.NopCloser(body)
	}
	req.Body = rc

	return r.defaultClient.Do(req)
}

func NewRequest(opts ...Option) *Request {
	var request Request

	// 设置基础地址
	urlParse, _ := url.Parse("http://127.0.0.1:80")
	request.baseUrl = urlParse
	// 设置默认的client
	request.defaultClient = http.DefaultClient

	for _, o := range opts {
		o(&request)
	}

	return &request
}

// 获取网页内容的编码
func GetCharset(resp *http.Response) (encoding string, err error) {
	reader := bufio.NewReader(resp.Body)
	content, err := reader.Peek(1024)
	if err != nil {
		return "", err
	}
	_, encoding, _ = charset.DetermineEncoding(content, resp.Header.Get("content-type"))
	return encoding, nil
}

// 转码，比如：gbk -> utf-8，返回一个解码后的reader
func Transform(source io.Reader, e encoding.Encoding) io.Reader {
	// source: 原始的body
	return transform.NewReader(source, e.NewDecoder())
}

func UploadFile() error {
	var body bytes.Buffer

	writer := multipart.NewWriter(&body)
	if err := writer.WriteField("name", "liuxuech"); err != nil {
		return err
	}

	// 第一个参数：form表单的field名称，第二个参数：上传的文件的名称
	// 文件1
	fileWriter1, err := writer.CreateFormFile("file1", "README.md")
	if err != nil {
		return err
	}
	fd, err := os.Open("README.md")
	if err != nil {
		return err
	}
	if _, err := io.Copy(fileWriter1, fd); err != nil {
		return err
	}

	// 文件2
	fileWriter2, err := writer.CreateFormFile("file2", "README.md")
	fd, err = os.Open("README.md")
	if err != nil {
		return err
	}
	if _, err := io.Copy(fileWriter2, fd); err != nil {
		return err
	}

	if err := writer.Close(); err != nil {
		return err
	}

	contentType := writer.FormDataContentType()

	resp, err := http.Post("http://127.0.0.1:3500/upload", contentType, &body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	io.Copy(os.Stdout, resp.Body)

	return nil
}
