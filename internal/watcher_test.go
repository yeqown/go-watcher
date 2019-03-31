package internal

import (
	"testing"
)

func Test_NewWatcher(t *testing.T) {
	exit := make(chan bool)
	watchingFiletypes := []string{"go"}
	unwatchingRegular := []string{}

	opt := &WatcherOption{
		D:                 2000,
		IncludedFiletypes: watchingFiletypes,
		ExcludedRegexps:   unwatchingRegular,
	}
	if w, err := NewWatcher([]string{"."}, exit, opt); err != nil {
		t.Error(err)
		t.FailNow()
	} else {
		w.Exit()
	}
}

func Test_filetypeIncludePre(t *testing.T) {
	type args struct {
		filetype string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "case 0",
			args: args{
				filetype: ".go",
			},
			want:    "go",
			wantErr: false,
		},
		{
			name: "case 1",
			args: args{
				filetype: "tar.gz",
			},
			want:    "gz",
			wantErr: false,
		},
		{
			name: "case 2",
			args: args{
				filetype: "main.py",
			},
			want:    "py",
			wantErr: false,
		},
		{
			name: "case 3",
			args: args{
				filetype: "demo-bin",
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := filetypeIncludePre(tt.args.filetype)
			if (err != nil) != tt.wantErr {
				t.Errorf("filetypeIncludePre() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("filetypeIncludePre() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_filetypeIncludeJudge(t *testing.T) {
	type args struct {
		s string
		m map[string]bool
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "case 0",
			args: args{
				s: "main.go",
				m: map[string]bool{
					"go": true,
					"py": true,
				},
			},
			want: true,
		},
		{
			name: "case 1",
			args: args{
				s: "main.tar.gz",
				m: map[string]bool{
					"go": true,
					"gz": true,
				},
			},
			want: true,
		},
		{
			name: "case 2",
			args: args{
				s: "main-tar-gz",
				m: map[string]bool{
					"go": true,
					"gz": true,
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := filetypeIncludeJudge(tt.args.s, tt.args.m); got != tt.want {
				t.Errorf("filetypeIncludeJudge() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_regularExcludePre(t *testing.T) {
	type args struct {
		regular string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "case 0",
			args: args{
				regular: "^api.go$",
			},
			want:    "^api.go$",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := regularExcludePre(tt.args.regular)
			if (err != nil) != tt.wantErr {
				t.Errorf("regularExcludePre() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("regularExcludePre() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_regularExcludeJudge(t *testing.T) {
	type args struct {
		s string
		m map[string]bool
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "case 0",
			args: args{
				s: "api.go",
				m: map[string]bool{
					"^api.go$": true,
					"*.go":     true,
				},
			},
			want: true,
		},
		{
			name: "case 0",
			args: args{
				s: "main.py",
				m: map[string]bool{
					"^api.go$": true,
					"*.go":     true,
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := regularExcludeJudge(tt.args.s, tt.args.m); got != tt.want {
				t.Errorf("regularExcludeJudge() = %v, want %v", got, tt.want)
			}
		})
	}
}
