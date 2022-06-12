package middleware

import "testing"

func Test_phaseJwtToken(t *testing.T) {
	PhaseJwtToken("")
	type args struct {
		token string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			PhaseJwtToken(tt.args.token)
		})
	}
}
