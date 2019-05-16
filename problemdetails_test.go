package problemdetails

import (
	"encoding/xml"
	"reflect"
	"testing"
)

func TestNewHTTP(t *testing.T) {
	type args struct {
		statusCode int
	}
	tests := []struct {
		name string
		args args
		want *ProblemDetails
	}{
		{
			name: "http_400",
			args: args{
				statusCode: 400,
			},
			want: &ProblemDetails{
				Type:     "about:blank",
				Title:    "Bad Request",
				Status:   400,
				Detail:   "",
				Instance: "",
			},
		},
		{
			name: "zero",
			args: args{
				statusCode: 0,
			},
			want: &ProblemDetails{
				Type:     "about:blank",
				Title:    "",
				Status:   0,
				Detail:   "",
				Instance: "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewHTTP(tt.args.statusCode); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewHTTP() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNew(t *testing.T) {
	type args struct {
		statusCode  int
		problemType string
		title       string
		detail      string
		instance    string
	}
	tests := []struct {
		name string
		args args
		want *ProblemDetails
	}{
		{
			name: "http_400_blank",
			args: args{
				statusCode:  400,
				problemType: "",
				title:       "",
				detail:      "",
				instance:    "",
			},
			want: &ProblemDetails{
				Type:     "about:blank",
				Title:    "Bad Request",
				Status:   400,
				Detail:   "",
				Instance: "",
			},
		},
		{
			name: "http_404",
			args: args{
				statusCode:  404,
				problemType: "https://example.net/problem/not_found",
				title:       "Not found",
				detail:      "Object with id was not found, another id should be given.",
				instance:    "https://api.example.net/objects/1234",
			},
			want: &ProblemDetails{
				Type:     "https://example.net/problem/not_found",
				Title:    "Not found",
				Status:   404,
				Detail:   "Object with id was not found, another id should be given.",
				Instance: "https://api.example.net/objects/1234",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.statusCode, tt.args.problemType, tt.args.title, tt.args.detail, tt.args.instance); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProblemDetailsError(t *testing.T) {
	type fields struct {
		XMLName  xml.Name
		Type     string
		Title    string
		Status   int
		Detail   string
		Instance string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "404",
			fields: fields{
				Type:     "https://example.net/problem/not_found",
				Title:    "Not found",
				Status:   404,
				Detail:   "Object with id was not found, another id should be given.",
				Instance: "https://api.example.net/objects/1234",
			},
			want: "Not found: Object with id was not found, another id should be given.",
		},
		{
			name: "404_title",
			fields: fields{
				Status: 404,
				Title:  "Not found",
			},
			want: "Not found",
		},
		{
			name: "404_status",
			fields: fields{
				Status: 404,
			},
			want: "Status 404",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pd := &ProblemDetails{
				XMLName:  tt.fields.XMLName,
				Type:     tt.fields.Type,
				Title:    tt.fields.Title,
				Status:   tt.fields.Status,
				Detail:   tt.fields.Detail,
				Instance: tt.fields.Instance,
			}
			if got := pd.Error(); got != tt.want {
				t.Errorf("ProblemDetails.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}
