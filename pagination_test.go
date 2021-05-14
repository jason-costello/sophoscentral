package sophoscentral

import "testing"

func Test_getRemainingPages(t *testing.T) {
	type args struct {
		ttlItems  int
		currItems int
		maxReturn int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name:"one",
			args: args{
				ttlItems:  10,
				currItems: 8,
				maxReturn: 5,
			},
			want: 1,
		},
		{
			name:"40",
			args: args{
				ttlItems:  4000,
				currItems: 50,
				maxReturn: 100,
			},
			want: 40,
		},
		{
			name:"0",
			args: args{
				ttlItems:  4000,
				currItems: 4000,
				maxReturn: 100,
			},
			want: 0,
		},
		{
			name:"-1",
			args: args{
				ttlItems:  3999,
				currItems: 4000,
				maxReturn: 100,
			},
			want: 0,
		},
	}
		for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getRemainingPageCount(tt.args.ttlItems, tt.args.currItems, tt.args.maxReturn); got != tt.want {
				t.Errorf("getRemainingPages() = %v, want %v", got, tt.want)
			}
		})
	}
}
