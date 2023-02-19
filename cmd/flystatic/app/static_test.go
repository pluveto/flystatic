package app

import "testing"

func Test_buildHttpPathPrefix(t *testing.T) {
	type args struct {
		path       string
		userPrefix string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "empty", args: args{path: "", userPrefix: ""}, want: ""},
		{name: "empty path", args: args{path: "", userPrefix: "a"}, want: "a"},
		{name: "empty use prefix", args: args{path: "a", userPrefix: ""}, want: "a"},
		{name: "normal", args: args{path: "a", userPrefix: "b"}, want: "a/b"},
		{name: "path with slash", args: args{path: "a/", userPrefix: "b"}, want: "a/b"},
		{name: "user prefix with slash", args: args{path: "a", userPrefix: "b/"}, want: "a/b"},
		{name: "path and user prefix with slash", args: args{path: "a/", userPrefix: "b/"}, want: "a/b"},
		{name: "duplicate slash", args: args{path: "a/", userPrefix: "/b/"}, want: "a/b"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := buildHttpPathPrefix(tt.args.path, tt.args.userPrefix); got != tt.want {
				t.Errorf("buildPathPrefix() = %v, want %v", got, tt.want)
			}
		})
	}
}
