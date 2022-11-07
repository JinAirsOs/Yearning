package model

import "testing"

func TestParseNameFromDN(t *testing.T) {
	type args struct {
		userdn string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test userdn",
			args: args{
				userdn: "CN=Zhendong Pan 潘振东,OU=CN5,OU=China,OU=Staff,OU=China,DC=fareast,DC=nevint,DC=com",
			},
			want: "Zhendong Pan 潘振东",
		},
		{
			name: "test cornercase",
			args: args{
				userdn: "OU=China,OU=Staff,OU=China,DC=fareast,DC=nevint,DC=com",
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseNameFromDN(tt.args.userdn); got != tt.want {
				t.Errorf("ParseNameFromDN() = %v, want %v", got, tt.want)
			}
		})
	}
}
