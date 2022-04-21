package main

import (
	//"github.com/stretchr/testify/assert"
	"testing"
)

func TestRunTimeMetrics_UpdateRTMetric(t *testing.T) {
	type fields struct {
		counter map[string]int64
		gauge   map[string]float64
	}
	type args struct {
		mtype  string
		mname  string
		mvalue string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "Simple test gauge #1",
			args:    args{"gauge", "allocate", "44412"},
			wantErr: false,
		},
		{
			name:    "Negative test gauge #2",
			args:    args{"gauge", "allocate", "a44412"},
			wantErr: true,
		},
		{
			name:    "Simple test counter #2",
			args:    args{"counter", "testField", "1"},
			wantErr: false,
			fields:  fields{counter: map[string]int64{"testField": 1}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rtm := &RunTimeMetrics{
				counter: tt.fields.counter,
				gauge:   tt.fields.gauge,
			}
			if err := rtm.UpdateRTMetric(tt.args.mtype, tt.args.mname, tt.args.mvalue); (err != nil) != tt.wantErr {
				t.Errorf("UpdateRTMetric() error = %v, wantErr %v", err, tt.wantErr)
			} else if tt.fields.counter != nil && tt.fields.counter["testField"] != 2 {
				t.Errorf("UpdateRTMetric() Counter not increase - %v", tt.fields.counter["testField"])
			}
		})
	}
}
