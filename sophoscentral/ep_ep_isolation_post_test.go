package sophoscentral

import (
	"context"
	"reflect"
	"testing"
)

func TestEndpointService_ToggleIsolation(t *testing.T) {
	type args struct {
		ctx       context.Context
		tenantID  string
		tenantURL BaseURL
		ti        ToggleIsolations
	}
	tests := []struct {
		name    string
		e       EndpointService
		args    args
		want    *ToggleIsolationSettings
		want1   *Response
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := tt.e.ToggleIsolation(tt.args.ctx, tt.args.tenantID, tt.args.tenantURL, tt.args.ti)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToggleIsolation() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToggleIsolation() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("ToggleIsolation() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
