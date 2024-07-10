package rke2

import "testing"

func Test_replaceVersionLink(t *testing.T) {
	type args struct {
		version string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{
			name: "Test replaceVersionLink with right version",
			args: args{
				version: "v1.21.3+rke2r1",
			},
			want: "v1.21.3%2Brke2r1",
		},
		// TODO: Add test cases.
		{
			name: "Test replaceVersionLink with right version",
			args: args{
				version: "v1.21.3",
			},
			want: "v1.21.3",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := replaceVersionLink(tt.args.version); got != tt.want {
				t.Errorf("replaceVersionLink() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_ensureTrailingSlash(t *testing.T) {
	type args struct {
		dir string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{
			name: "Test ensureTrailingSlash with right dir",
			args: args{
				dir: "test",
			},
			want: "test/",
		},
		{
			name: "Test ensureTrailingSlash with right dir but including slash",
			args: args{
				dir: "test/",
			},
			want: "test/",
		},
		{
			name: "Test ensureTrailingSlash with right dir but including slash at the beginning",
			args: args{
				dir: "./test/",
			},
			want: "./test/",
		},
		{
			name: "Test ensureTrailingSlash with right dir but including slash at the beginning",
			args: args{
				dir: "./aa/test/",
			},
			want: "./aa/test/",
		},
		{
			name: "Test ensureTrailingSlash with right dir but including slash at the beginning but not in the end",
			args: args{
				dir: "./aa/test",
			},
			want: "./aa/test/",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ensureTrailingSlash(tt.args.dir); got != tt.want {
				t.Errorf("ensureTrailingSlash() = %v, want %v", got, tt.want)
			}
		})
	}
}
