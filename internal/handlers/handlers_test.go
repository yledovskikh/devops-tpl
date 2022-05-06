package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/yledovskikh/devops-tpl/internal/storage"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	type args struct {
		storage storage.Storage
	}
	tests := []struct {
		name string
		args args
		want *Server
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.storage); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServer_GetMetric(t *testing.T) {
	type fields struct {
		storage storage.Storage
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				storage: tt.fields.storage,
			}
			s.GetMetric(tt.args.w, tt.args.r)
		})
	}
}

func TestServer_PostMetric(t *testing.T) {
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
			//s := &Server{
			//	storage: tt.fields.storage,
			//}
			//s.PostMetric(tt.args.w, tt.args.r)
			s := storage.NewMetricStore()
			h := New(s)

			path := "/update/" + tt.metric.metricType + "/" + tt.metric.metricName + "/" + tt.metric.metricValue
			req, err := http.NewRequest("POST", path, nil)
			if err != nil {
				t.Fatal(err)
			}
			tr := httptest.NewRecorder()
			r := chi.NewRouter()
			r.HandleFunc("/update/{metricType}/{metricName}/{metricValue}", h.PostMetric)
			r.ServeHTTP(tr, req)
			res := tr.Result()
			assert.Equal(t, tt.want.statusCode, res.StatusCode)
			res.Body.Close()
		})
	}
}

func Test_storageErrToStatus(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := storageErrToStatus(tt.args.err); got != tt.want {
				t.Errorf("storageErrToStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}
