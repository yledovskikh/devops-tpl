package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

//func testRequest(t *testing.T, ts *httptest.Server, method, path string) (*http.Response, string) {
//	req, err := http.NewRequest(method, ts.URL+path, nil)
//	require.NoError(t, err)
//
//	resp, err := http.DefaultClient.Do(req)
//	require.NoError(t, err)
//
//	respBody, err := ioutil.ReadAll(resp.Body)
//	require.NoError(t, err)
//
//	defer resp.Body.Close()
//
//	return resp, string(respBody)
//}

//func TestGetMetric(t *testing.T) {
//	//type args struct {
//	//	w http.ResponseWriter
//	//	r *http.Request
//	//}
//	type want struct {
//		code int
//		//response string
//		//contentType string
//	}
//	type metric struct {
//		metricType  string
//		metricName  string
//		metricValue string
//	}
//	tests := []struct {
//		name   string
//		metric metric
//		want   want
//	}{
//		{
//			name:   "simple test #1",
//			metric: metric{"goage", "metric1", "111.111"},
//			want:   want{200},
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//
//			r := chi.NewRouter()
//			ts := httptest.NewServer(r)
//			defer ts.Close()
//			path := "/update/{metricType}/{metricName}/{metricValue}"
//			response, _ := testRequest(t, ts, "POST", path)
//
//			if response.StatusCode != 200 {
//				t.Fatal("err")
//			}
//
//			//assert.Equal(t, http.StatusOK, resp.StatusCode)
//			//assert.Equal(t, "brand:renault", body)
//
//			//
//			//
//			//
//			//m := tt.metric
//			////path := fmt.Sprintf("/update/%s/%s/%s",
//			////	m.metricType, m.metricName, m.metricValue)
//			//req, err := http.NewRequest("POST", path, nil)
//			//if err != nil {
//			//	t.Fatal(err)
//			//}
//			//rr := httptest.NewRecorder()
//			//router := mux.NewRouter()
//			//router.HandleFunc("/update/{metricType}/{metricName}/{metricValue}", handlers.UpdateMeticHandler)
//			//router.ServeHTTP(rr, req)
//			//res := rr.Result()
//			//
//			//if res.StatusCode != tc.want.code {
//			//	t.Errorf("Expected status code %d, got %d", tc.want.code, rr.Code)
//			//}
//			//
//			//defer res.Body.Close()
//			//resBody, err := io.ReadAll(res.Body)
//			//if err != nil {
//			//	t.Fatal(err)
//			//}
//			//if string(resBody) != tc.want.response {
//			//	t.Errorf("Expected body %s, got %s", tc.want.response, rr.Body.String())
//			//}
//		})
//	}
//}

func TestPostMetric(t *testing.T) {
	type want struct {
		statusCode int
		//response string
		//contentType string
	}
	type metric struct {
		metricType  string
		metricName  string
		metricValue string
	}
	tests := []struct {
		name   string
		metric metric
		want   want
	}{
		{
			name:   "simple test upload metric gauge",
			metric: metric{"gauge", "testm1", "1111.1"},
			want:   want{200},
		},
		{
			name:   "incorrect upload metric gauge",
			metric: metric{"gauge", "testm1", "1.111.1"},
			want:   want{400},
		},
		{
			name:   "simple test upload metric counter",
			metric: metric{"counter", "testm1", "1"},
			want:   want{200},
		},
		{
			name:   "incorrect upload metric counter",
			metric: metric{"counter", "testm1", "1.1"},
			want:   want{400},
		},
		{
			name:   "incorrect metric type",
			metric: metric{"incorrect", "testm1", "1.1"},
			want:   want{501},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			path := "/update/" + tt.metric.metricType + "/" + tt.metric.metricName + "/" + tt.metric.metricValue
			req, err := http.NewRequest("POST", path, nil)
			if err != nil {
				t.Fatal(err)
			}
			tr := httptest.NewRecorder()
			r := chi.NewRouter()
			r.HandleFunc("/update/{metricType}/{metricName}/{metricValue}", PostMetric)
			r.ServeHTTP(tr, req)
			res := tr.Result()
			assert.Equal(t, tt.want.statusCode, res.StatusCode)
			res.Body.Close()
		})
	}
}
