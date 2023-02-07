package manage

import "testing"

func Test_trim(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "测试无\n",
			args: args{
				s: "desc app_info;",
			},
			want: "desc app_info;",
		},
		{
			name: "测试\n开头",
			args: args{
				s: "\ndesc app_info;",
			},
			want: "desc app_info;",
		},
		{
			name: "测试\n末尾",
			args: args{
				s: "desc app_info;\n",
			},
			want: "desc app_info;",
		},
		{
			name: "测试\n末尾",
			args: args{
				s: "desc app_info;\n\n",
			},
			want: "desc app_info;",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := trim(tt.args.s); got != tt.want {
				t.Errorf("trim() = %v, want %v", got, tt.want)
			}
		})
	}
}
