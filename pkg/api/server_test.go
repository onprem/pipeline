package api

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/prmsrswt/pipeline/pkg/task"
)

const (
	sampleCSV = `id,name
1,x
2,y
3,z`
	maxProcessingSec = 2
)

func constructFileUpload(csv string, t *testing.T) (bytes.Buffer, string) {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	fw, err := w.CreateFormFile("file", "test.csv")
	if err != nil {
		t.Fatal(err)
	}

	fw.Write([]byte(csv))
	w.Close()

	return b, w.FormDataContentType()
}

func getID(body io.ReadCloser, t *testing.T) string {
	var res map[string]interface{}
	err := json.NewDecoder(body).Decode(&res)
	if err != nil {
		t.Fatal(err)
	}

	data, ok := (res["data"]).(map[string]interface{})
	if !ok {
		t.Fatalf("decoding response id: %v", res)
	}

	id, ok := data["id"]

	if !ok {
		t.Fatal("ID not found")
	}

	return id.(string)
}

func getStatus(body io.ReadCloser, t *testing.T) task.Status {
	var res map[string]interface{}
	err := json.NewDecoder(body).Decode(&res)
	if err != nil {
		t.Fatal(err)
	}

	data, ok := res["data"].(map[string]interface{})
	if !ok {
		t.Fatal("decoding response status")
	}

	status, ok := data["status"]

	if !ok {
		t.Fatal("Status not found")
	}

	return task.Status(status.(string))
}

func checkStatus(id string, status task.Status, ts *httptest.Server, t *testing.T) {
	respStatus, err := ts.Client().PostForm(ts.URL+"/status", url.Values{"id": []string{id}})
	if err != nil {
		t.Fatal(err)
	}
	defer respStatus.Body.Close()

	if respStatus.StatusCode != http.StatusOK {
		t.Fatalf("bad http status: %s", respStatus.Status)
	}

	statusP := getStatus(respStatus.Body, t)

	if statusP != status {
		t.Fatalf("incorrrect status. expected: %s; got: %s", status, statusP)
	}
}

func requestAndCheckStatus(id, path string, status task.Status, ts *httptest.Server, t *testing.T) {
	resp, err := ts.Client().PostForm(ts.URL+path, url.Values{"id": []string{id}})
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("bad status %s: %s", path, resp.Status)
	}

	checkStatus(id, status, ts, t)
}

func uploadSampleCSV(ts *httptest.Server, t *testing.T) string {
	b, contentType := constructFileUpload(sampleCSV, t)

	resp, err := ts.Client().Post(ts.URL+"/upload", contentType, &b)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("bad status: %s", resp.Status)
	}

	return getID(resp.Body, t)
}

func setupServer(t *testing.T) *httptest.Server {
	dir, err := ioutil.TempDir("", "pipeline-test")
	if err != nil {
		t.Fatal(err)
	}

	mux := http.NewServeMux()

	api := NewAPI(dir)
	api.Register(mux)

	ts := httptest.NewServer(mux)
	t.Cleanup(func() {
		ts.Close()
		os.RemoveAll(dir)
	})

	return ts
}

func TestOverAll(t *testing.T) {
	ts := setupServer(t)

	id := uploadSampleCSV(ts, t)

	checkStatus(id, task.TaskRunning, ts, t)

	time.Sleep(time.Duration(maxProcessingSec * time.Second / 4))

	requestAndCheckStatus(id, "/pause", task.TaskPaused, ts, t)

	time.Sleep(time.Duration(maxProcessingSec * time.Second))

	requestAndCheckStatus(id, "/resume", task.TaskRunning, ts, t)

	time.Sleep(time.Duration(maxProcessingSec * time.Second / 4))

	requestAndCheckStatus(id, "/terminate", task.TaskTerminated, ts, t)
}

func TestUpload(t *testing.T) {
	ts := setupServer(t)

	b, contentType := constructFileUpload(sampleCSV, t)

	resp, err := ts.Client().Post(ts.URL+"/upload", contentType, &b)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("bad status: %s", resp.Status)
	}
}

func TestStatusRunning(t *testing.T) {
	ts := setupServer(t)

	id := uploadSampleCSV(ts, t)

	checkStatus(id, task.TaskRunning, ts, t)
}

func TestStatusFinished(t *testing.T) {
	ts := setupServer(t)

	id := uploadSampleCSV(ts, t)

	time.Sleep(time.Duration(maxProcessingSec * time.Second))

	checkStatus(id, task.TaskFinished, ts, t)
}
