/*
MIT License

Copyright (c) 2018 Danny van Kooten

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/
package respond

import (
	"html/template"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

type Greeting struct {
	One string `json:"one"`
	Two string `json:"two"`
}

func TestJSON(t *testing.T) {

	w := httptest.NewRecorder()
	err := JSON(w, 299, Greeting{"hello", "world"})
	res := w.Result()

	if err != nil {
		t.Errorf("expected %#v, got %#v", nil, err)
	}

	if res.StatusCode != 299 {
		t.Errorf("invalid status code: expected %#v, got %#v", 299, res.StatusCode)
	}

	e := ContentJSON + "; charset=" + Charset
	if v := res.Header.Get(ContentType); v != e {
		t.Errorf("invalid content type: expected %#v, got %#v", e, v)
	}

	body, _ := ioutil.ReadAll(res.Body)
	if v := string(body); v != "{\"one\":\"hello\",\"two\":\"world\"}\n" {
		t.Errorf("invalid response body: expected %#v, got %#v", "{\"one\":\"hello\",\"two\":\"world\"}\n", v)
	}
}

func TestJSONError(t *testing.T) {
	w := httptest.NewRecorder()
	err := JSON(w, 200, make(chan int))
	if err == nil {
		t.Error("Expected err, got nil")
	}
}

func TestXMLError(t *testing.T) {
	w := httptest.NewRecorder()
	err := XML(w, 200, make(chan int))
	if err == nil {
		t.Error("Expected err, got nil")
	}
}

func TestTemplate(t *testing.T) {
	w := httptest.NewRecorder()
	tmpl := template.Must(template.New("").Parse("Hello {{.}}"))
	data := "world"
	err := Template(w, http.StatusOK, tmpl, data)
	res := w.Result()
	if err != nil {
		t.Errorf("expected %#v, got %#v", nil, err)
	}

	if res.StatusCode != http.StatusOK {
		t.Errorf("invalid status code: expected %#v, got %#v", http.StatusOK, res.StatusCode)
	}

	e := ContentHTML + "; charset=" + Charset
	if v := res.Header.Get(ContentType); v != e {
		t.Errorf("invalid content type: expected %#v, got %#v", e, v)
	}

	body, _ := ioutil.ReadAll(res.Body)
	if v := string(body); v != "Hello world" {
		t.Errorf("invalid body: expected %#v, got %#v", string(data), v)
	}
}

func TestTemplateError(t *testing.T) {
	w := httptest.NewRecorder()
	tmpl := template.Must(template.New("").Parse("Hello <span class=\">world</span>"))
	err := Template(w, http.StatusOK, tmpl)
	if err == nil {
		t.Error("expected err, got nil")
	}
}

func TestHTML(t *testing.T) {
	w := httptest.NewRecorder()
	data := []byte("Hello <strong>world</strong>!")
	err := HTML(w, http.StatusOK, data)
	res := w.Result()
	if err != nil {
		t.Errorf("expected %#v, got %#v", nil, err)
	}

	if res.StatusCode != http.StatusOK {
		t.Errorf("invalid status code: expected %#v, got %#v", http.StatusOK, res.StatusCode)
	}

	e := ContentHTML + "; charset=" + Charset
	if v := res.Header.Get(ContentType); v != e {
		t.Errorf("invalid content type: expected %#v, got %#v", e, v)
	}

	body, _ := ioutil.ReadAll(res.Body)
	if v := string(body); v != string(data) {
		t.Errorf("invalid body: expected %#v, got %#v", string(data), v)
	}
}

func TestXML(t *testing.T) {
	w := httptest.NewRecorder()
	err := XML(w, 299, Greeting{"hello", "world"})
	res := w.Result()

	if err != nil {
		t.Errorf("expected %#v, got %#v", nil, err)
	}

	if res.StatusCode != 299 {
		t.Errorf("invalid status code: expected %#v, got %#v", 299, res.StatusCode)
	}

	e := ContentXML + "; charset=" + Charset
	if v := res.Header.Get(ContentType); v != e {
		t.Errorf("invalid content type: expected %#v, got %#v", e, v)
	}

	body, _ := ioutil.ReadAll(res.Body)
	expected := "<Greeting><One>hello</One><Two>world</Two></Greeting>"
	if v := string(body); v != expected {
		t.Errorf("invalid response body: expected %#v, got %#v", expected, v)
	}
}

func TestText(t *testing.T) {
	w := httptest.NewRecorder()
	err := Text(w, 200, []byte("Hello world!"))
	res := w.Result()

	if err != nil {
		t.Errorf("expected %#v, got %#v", nil, err)
	}

	if res.StatusCode != 200 {
		t.Errorf("invalid status code: expected %#v, got %#v", 200, res.StatusCode)
	}

	e := ContentText + "; charset=" + Charset
	if v := res.Header.Get(ContentType); v != e {
		t.Errorf("invalid content type: expected %#v, got %#v", e, v)
	}

	body, _ := ioutil.ReadAll(res.Body)
	expected := "Hello world!"
	if v := string(body); v != expected {
		t.Errorf("invalid response body: expected %#v, got %#v", expected, v)
	}
}
