package jwt

import (
	"testing"
)

func TestToken_SigningString(t1 *testing.T) {
	type fields struct {
		Raw       string
		Method    SigningMethod
		Header    map[string]interface{}
		Claims    Claims
		Signature string
		Valid     bool
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{
			name: "hmac_signed_token",
			fields: fields{
				Raw:    "",
				Method: SigningMethodHS256,
				Header: map[string]interface{}{
					"typ": "JWT",
					"alg": SigningMethodHS256.Alg(),
				},
				Claims:    RegisteredClaims{},
				Signature: "",
				Valid:     false,
			},
			want:    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.e30",
			wantErr: false,
		},
		{
			name: "rsa_signed_token",
			fields: fields{
				Raw:    "",
				Method: SigningMethodRS256,
				Header: map[string]interface{}{
					"typ": "JWT",
					"alg": SigningMethodRS256.Alg(),
				},
				Claims:    RegisteredClaims{},
				Signature: "",
				Valid:     false,
			},
			want:    "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.e30",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t1.Run(
			tt.name, func(t1 *testing.T) {
				t := &Token{
					Raw:       tt.fields.Raw,
					Method:    tt.fields.Method,
					Header:    tt.fields.Header,
					Claims:    tt.fields.Claims,
					Signature: tt.fields.Signature,
					Valid:     tt.fields.Valid,
				}
				got, err := t.Sign()
				if (err != nil) != tt.wantErr {
					t1.Errorf("SigningString() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if got != tt.want {
					t1.Errorf("SigningString() got = %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func BenchmarkToken_SigningString(b *testing.B) {
	t := &Token{
		Method: SigningMethodHS256,
		Header: map[string]interface{}{
			"typ": "JWT",
			"alg": SigningMethodHS256.Alg(),
		},
		Claims: RegisteredClaims{},
	}
	b.Run(
		"BenchmarkToken_SigningString", func(b *testing.B) {
			b.ResetTimer()
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				t.Sign()
			}
		},
	)
}
