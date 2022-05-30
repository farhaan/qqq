package utility

import "testing"

func TestStringInArray(t *testing.T) {
	type args struct {
		str string
		arr []string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "success",
			args: args{
				str: "hehe",
				arr: []string{"wkwkwk", "haha", "hihi", "hehe"},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StringInArray(tt.args.str, tt.args.arr); got != tt.want {
				t.Errorf("StringInArray() = %v, want %v", got, tt.want)
			}
		})
	}
}
