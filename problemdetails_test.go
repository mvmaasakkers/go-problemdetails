package problemdetails

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestMarshal(t *testing.T) {
	tests := []struct {
		name     string
		args     *ProblemDetails
		wantXML  string
		wantJSON string
	}{
		{
			name: "http_400_blank",
			args: &ProblemDetails{
				Type:     "about:blank",
				Title:    "Bad Request",
				Status:   400,
				Detail:   "",
				Instance: "",
			},
			wantXML:  `<problem xmlns="urn:ietf:rfc:7807"><type>about:blank</type><title>Bad Request</title><status>400</status></problem>`,
			wantJSON: `{"type":"about:blank","title":"Bad Request","status":400}`,
		},
		{
			name: "http_type_title",
			args: &ProblemDetails{
				Type:     "https://example.net/problem/issue",
				Title:    "Issue title",
				Status:   0,
				Detail:   "",
				Instance: "",
			},
			wantXML:  `<problem xmlns="urn:ietf:rfc:7807"><type>https://example.net/problem/issue</type><title>Issue title</title></problem>`,
			wantJSON: `{"type":"https://example.net/problem/issue","title":"Issue title"}`,
		},
		{
			name: "http_404",
			args: &ProblemDetails{
				Type:     "https://example.net/problem/issue",
				Title:    "Issue title",
				Status:   404,
				Detail:   "Object with id was not found, another id should be given.",
				Instance: "https://api.example.net/objects/1234",
			},
			wantXML:  `<problem xmlns="urn:ietf:rfc:7807"><type>https://example.net/problem/issue</type><title>Issue title</title><status>404</status><detail>Object with id was not found, another id should be given.</detail><instance>https://api.example.net/objects/1234</instance></problem>`,
			wantJSON: `{"type":"https://example.net/problem/issue","title":"Issue title","status":404,"detail":"Object with id was not found, another id should be given.","instance":"https://api.example.net/objects/1234"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bXML, err := xml.Marshal(tt.args)
			if err != nil {
				t.Errorf("xml got error %v", err)
				t.FailNow()
			}
			if tt.wantXML != string(bXML) {
				t.Errorf("xml got %s, want %s", string(bXML), tt.wantXML)
			}

			bJSON, err := json.Marshal(tt.args)
			if err != nil {
				t.Errorf("json got error %v", err)
				t.FailNow()
			}
			if tt.wantJSON != string(bJSON) {
				t.Errorf("json got %s, want %s", string(bJSON), tt.wantJSON)
			}
		})
	}
}

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
		{
			name: "http_relative_type",
			args: args{
				statusCode:  404,
				problemType: "/problem/not_found",
				title:       "Not found",
				detail:      "Object with id was not found, another id should be given.",
				instance:    "https://api.example.net/objects/1234",
			},
			want: &ProblemDetails{
				Type:     "/problem/not_found",
				Title:    "Not found",
				Status:   404,
				Detail:   "Object with id was not found, another id should be given.",
				Instance: "https://api.example.net/objects/1234",
			},
		},
		{
			name: "http_relative_instance",
			args: args{
				statusCode:  404,
				problemType: "/problem/not_found",
				title:       "Not found",
				detail:      "Object with id was not found, another id should be given.",
				instance:    "objects-1234",
			},
			want: &ProblemDetails{
				Type:     "/problem/not_found",
				Title:    "Not found",
				Status:   404,
				Detail:   "Object with id was not found, another id should be given.",
				Instance: "objects-1234",
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

func TestProblemDetailsServeJSON(t *testing.T) {

	tests := []struct {
		name     string
		args     *ProblemDetails
		wantXML  string
		wantJSON string
	}{
		{
			name: "http_400_blank",
			args: &ProblemDetails{
				Type:     "about:blank",
				Title:    "Bad Request",
				Status:   400,
				Detail:   "",
				Instance: "",
			},
			wantXML:  `<problem xmlns="urn:ietf:rfc:7807"><type>about:blank</type><title>Bad Request</title><status>400</status></problem>`,
			wantJSON: `{"type":"about:blank","title":"Bad Request","status":400}`,
		},
		{
			name: "http_type_title",
			args: &ProblemDetails{
				Type:     "https://example.net/problem/issue",
				Title:    "Issue title",
				Status:   0,
				Detail:   "",
				Instance: "",
			},
			wantXML:  `<problem xmlns="urn:ietf:rfc:7807"><type>https://example.net/problem/issue</type><title>Issue title</title></problem>`,
			wantJSON: `{"type":"https://example.net/problem/issue","title":"Issue title"}`,
		},
		{
			name: "http_404",
			args: &ProblemDetails{
				Type:     "https://example.net/problem/issue",
				Title:    "Issue title",
				Status:   404,
				Detail:   "Object with id was not found, another id should be given.",
				Instance: "https://api.example.net/objects/1234",
			},
			wantXML:  `<problem xmlns="urn:ietf:rfc:7807"><type>https://example.net/problem/issue</type><title>Issue title</title><status>404</status><detail>Object with id was not found, another id should be given.</detail><instance>https://api.example.net/objects/1234</instance></problem>`,
			wantJSON: `{"type":"https://example.net/problem/issue","title":"Issue title","status":404,"detail":"Object with id was not found, another id should be given.","instance":"https://api.example.net/objects/1234"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/", nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()

			if err := tt.args.ServeJSON(rr, req); err != nil {
				t.Errorf("ProblemDetails.ServeJSON() = %v, want nil", err)
				t.FailNow()
			}

			pd := &ProblemDetails{}

			if err := json.NewDecoder(rr.Body).Decode(pd); err != nil {
				t.Errorf("json got error %v", err)
				t.FailNow()
			}
			if !reflect.DeepEqual(pd, tt.args) {
				t.Errorf("json got %v, want %v", pd, tt.args)
			}
		})
	}
}

func TestProblemDetailsServeXML(t *testing.T) {
	tests := []struct {
		name     string
		args     *ProblemDetails
		wantXML  string
		wantJSON string
	}{
		{
			name: "http_400_blank",
			args: &ProblemDetails{
				Type:     "about:blank",
				Title:    "Bad Request",
				Status:   400,
				Detail:   "",
				Instance: "",
			},
			wantXML:  `<problem xmlns="urn:ietf:rfc:7807"><type>about:blank</type><title>Bad Request</title><status>400</status></problem>`,
			wantJSON: `{"type":"about:blank","title":"Bad Request","status":400}`,
		},
		{
			name: "http_type_title",
			args: &ProblemDetails{
				Type:     "https://example.net/problem/issue",
				Title:    "Issue title",
				Status:   0,
				Detail:   "",
				Instance: "",
			},
			wantXML:  `<problem xmlns="urn:ietf:rfc:7807"><type>https://example.net/problem/issue</type><title>Issue title</title></problem>`,
			wantJSON: `{"type":"https://example.net/problem/issue","title":"Issue title"}`,
		},
		{
			name: "http_404",
			args: &ProblemDetails{
				Type:     "https://example.net/problem/issue",
				Title:    "Issue title",
				Status:   404,
				Detail:   "Object with id was not found, another id should be given.",
				Instance: "https://api.example.net/objects/1234",
			},
			wantXML:  `<problem xmlns="urn:ietf:rfc:7807"><type>https://example.net/problem/issue</type><title>Issue title</title><status>404</status><detail>Object with id was not found, another id should be given.</detail><instance>https://api.example.net/objects/1234</instance></problem>`,
			wantJSON: `{"type":"https://example.net/problem/issue","title":"Issue title","status":404,"detail":"Object with id was not found, another id should be given.","instance":"https://api.example.net/objects/1234"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/", nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			if err := tt.args.ServeXML(rr, req); err != nil {
				t.Errorf("ProblemDetails.ServeXML() = %v, want nil", err)
				t.FailNow()
			}

			pd := &ProblemDetails{}
			if err := xml.NewDecoder(rr.Body).Decode(pd); err != nil {
				t.Errorf("xml got error %v", err)
				t.FailNow()
			}
			// Can't use deepequal because of XMLName
			if pd.Status != tt.args.Status ||
				pd.Detail != tt.args.Detail ||
				pd.Title != tt.args.Title ||
				pd.Instance != tt.args.Instance ||
				pd.Type != tt.args.Type {
				t.Errorf("xml got %v, want %v", pd, tt.args)
			}
		})
	}
}
