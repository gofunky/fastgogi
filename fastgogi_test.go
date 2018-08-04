package fastgogi

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/valyala/fasthttp"
	"sort"
	"strings"
)

const (
	testAgent = "test/agent"
	testHost  = "https://github.com/gofunky/gogi"
)

func TestNewClientWithOptions(t *testing.T) {
	type args struct {
		options FastGogiOptions
	}
	tests := []struct {
		name        string
		args        args
		wantOptions FastGogiOptions
	}{
		{
			name:        "Default options",
			args:        args{FastGogiOptions{}},
			wantOptions: FastGogiOptions{UserAgent: defaultUserAgent, Host: defaultHost},
		},
		{
			name:        "Custom options",
			args:        args{FastGogiOptions{UserAgent: testAgent}},
			wantOptions: FastGogiOptions{UserAgent: testAgent, Host: defaultHost},
		},
		{
			name:        "Custom options",
			args:        args{FastGogiOptions{Host: testHost}},
			wantOptions: FastGogiOptions{UserAgent: defaultUserAgent, Host: testHost},
		},
		{
			name:        "Custom options",
			args:        args{FastGogiOptions{UserAgent: testAgent, Host: testHost}},
			wantOptions: FastGogiOptions{UserAgent: testAgent, Host: testHost},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotClient := NewClientWithOptions(tt.args.options); !reflect.DeepEqual(gotClient.FastGogiOptions, &tt.wantOptions) {
				t.Errorf("NewClientWithOptions().FastGogiOptions = %v, want %v",
					gotClient.FastGogiOptions, tt.wantOptions)
			}
		})
	}
}

func TestNewClient(t *testing.T) {
	tests := []struct {
		name        string
		wantOptions FastGogiOptions
	}{
		{
			name:        "Default options",
			wantOptions: FastGogiOptions{UserAgent: defaultUserAgent, Host: defaultHost},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotClient := NewClient(); !reflect.DeepEqual(gotClient.FastGogiOptions, &tt.wantOptions) {
				t.Errorf("NewClientWithOptions().FastGogiOptions = %v, want %v",
					gotClient.FastGogiOptions, tt.wantOptions)
			}
		})
	}
}

func Test_gogiClient_List(t *testing.T) {
	type fields struct {
		client          *fasthttp.Client
		FastGogiOptions *FastGogiOptions
	}
	tests := []struct {
		name      string
		fields    fields
		wantTypes []string
		wantErr   bool
	}{
		{
			name: "Check some patterns",
			fields: fields{
				client:          &fasthttp.Client{},
				FastGogiOptions: &FastGogiOptions{UserAgent: defaultUserAgent, Host: defaultHost},
			},
			wantTypes: []string{"java", "go"},
		},
		{
			name: "Invalid host",
			fields: fields{
				client:          &fasthttp.Client{},
				FastGogiOptions: &FastGogiOptions{UserAgent: testAgent, Host: testHost},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &fastGogiClient{
				client:          tt.fields.client,
				FastGogiOptions: tt.fields.FastGogiOptions,
			}
			gotTypes, err := c.List()
			if (err != nil) != tt.wantErr {
				t.Errorf("fastGogiClient.List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				for _, want := range tt.wantTypes {
					if !contains(gotTypes, want) {
						t.Errorf("List() result is missing %v", want)
					}
				}
			}
		})
	}
}

func Test_gogiClient_Get(t *testing.T) {
	endSignature := []byte("End of")
	allTypes, _ := NewClient().List()
	allTypesURI := NewClient().GetPath(allTypes...)
	type fields struct {
		client          *fasthttp.Client
		FastGogiOptions *FastGogiOptions
	}
	type args struct {
		includedTypes []string
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		wantContent []byte
		wantErr     bool
	}{
		{
			name: "Check some patterns",
			fields: fields{
				client:          &fasthttp.Client{},
				FastGogiOptions: &FastGogiOptions{UserAgent: defaultUserAgent, Host: defaultHost},
			},
			args:        args{[]string{"java", "go"}},
			wantContent: []byte("go,java"),
		},
		{
			name: "Check all patterns",
			fields: fields{
				client:          &fasthttp.Client{},
				FastGogiOptions: &FastGogiOptions{UserAgent: defaultUserAgent, Host: defaultHost},
			},
			args:        args{allTypes},
			wantContent: []byte(allTypesURI),
		},
		{
			name: "Invalid host",
			fields: fields{
				client:          &fasthttp.Client{},
				FastGogiOptions: &FastGogiOptions{UserAgent: testAgent, Host: testHost},
			},
			args:    args{[]string{"java", "go"}},
			wantErr: true,
		},
		{
			name: "Invalid arguments",
			fields: fields{
				client:          &fasthttp.Client{},
				FastGogiOptions: &FastGogiOptions{UserAgent: defaultUserAgent, Host: defaultHost},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &fastGogiClient{
				client:          tt.fields.client,
				FastGogiOptions: tt.fields.FastGogiOptions,
			}
			gotContent, err := c.Get(tt.args.includedTypes...)
			if (err != nil) != tt.wantErr {
				t.Errorf("fastGogiClient.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if !bytes.Contains(gotContent, tt.wantContent) {
					t.Errorf("fastGogiClient.Get() = %s, want %s", gotContent, tt.wantContent)
				}
				if !bytes.Contains(gotContent, endSignature) {
					t.Errorf("fastGogiClient.Get() = %s, want %s", gotContent, endSignature)
				}
			}
		})
	}
}

func Test_fastGogiClient_GetPath(t *testing.T) {
	allTypes, _ := NewClient().List()
	sort.Strings(allTypes)
	type fields struct {
		client          *fasthttp.Client
		FastGogiOptions *FastGogiOptions
	}
	type args struct {
		includedTypes []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantURI string
	}{
		{
			name: "Check some patterns",
			fields: fields{
				client:          &fasthttp.Client{},
				FastGogiOptions: &FastGogiOptions{UserAgent: defaultUserAgent, Host: defaultHost},
			},
			args:    args{[]string{"java", "go"}},
			wantURI: "https://www.gitignore.io/api/go,java",
		},
		{
			name: "Check all patterns",
			fields: fields{
				client:          &fasthttp.Client{},
				FastGogiOptions: &FastGogiOptions{UserAgent: defaultUserAgent, Host: defaultHost},
			},
			args:    args{allTypes},
			wantURI: "https://www.gitignore.io/api/" + strings.Join(allTypes, comma),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &fastGogiClient{
				client:          tt.fields.client,
				FastGogiOptions: tt.fields.FastGogiOptions,
			}
			if gotURI := c.GetPath(tt.args.includedTypes...); gotURI != tt.wantURI {
				t.Errorf("fastGogiClient.GetPath() = %v, want %v", gotURI, tt.wantURI)
			}
		})
	}
}
