package wallet

import "testing"

func TestService_Reject(t *testing.T) {
	type args struct {
		paymentID string
	}
	tests := []struct {
		name    string
		s       *Service
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.s.Reject(tt.args.paymentID); (err != nil) != tt.wantErr {
				t.Errorf("Service.Reject() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
