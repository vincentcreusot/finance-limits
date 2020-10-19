package fileutils

import (
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
	"time"
	"fmt"
)

func Test_ReadLines(t *testing.T) {
	_, testFileName, _, _ := runtime.Caller(0)
	baseFolder := filepath.Dir(testFileName)
	type args struct {
		filename string
	}
	type output struct {
		lines []string
	}
	tests := []struct {
		name string
		args args
		want output
	}{
		{
			name: "existingFile",
			args: args{
				filename: baseFolder + "/../test/filetest.txt",
			},
			want: output{
				lines: []string{"first line", "second line"},
			},
		},
		{
			name: "notExistingFile",
			args: args{
				filename: "notexisting.txt",
			},
			want: output{
				lines: make([]string, 0),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stringChannel := make(chan string)
			go ReadLines(tt.args.filename, stringChannel)
			linesParsed := make([]string, 0)
			for line := range stringChannel {
				linesParsed = append(linesParsed, line)
			}
			if !reflect.DeepEqual(linesParsed, tt.want.lines) {
				t.Errorf("ReadFile = %v, want %v", linesParsed, tt.want)
			}
		})
	}
}

func Test_WriteLines(t *testing.T) {
	_, testFileName, _, _ := runtime.Caller(0)
	baseFolder := filepath.Dir(testFileName)

	type args struct {
		filename string
		lines []string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "newFile",
			args: args{
				filename: baseFolder + "/../test/newFile_"+ fmt.Sprintf("%d", time.Now().Unix()) + ".txt",
				lines: []string{"first line", "second line"},
			},
			want: false,
		},
		{
			name: "existingFile",
			args: args{
				filename: baseFolder + "/../test/existingfiletowrite",
				lines: []string{"first line", "second line"},
			},
			want: false,
		},
		{
			name: "notWritable",
			args: args{
				filename: "notExistingFolder/notexisting.txt",
				lines: make([]string,0),
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WriteLines(tt.args.filename, tt.args.lines) ; (got != nil) != tt.want {
				t.Errorf("WriteLines = %v, want %v", got, tt.want)
			}
		})
	}
}
