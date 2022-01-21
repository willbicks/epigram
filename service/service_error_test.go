package service

import "testing"

func TestServiceError_Error(t *testing.T) {
	tests := []struct {
		name string
		e    ServiceError
		want string
	}{
		{
			name: "Empty",
			e:    ServiceError{},
			want: "",
		},
		{
			name: "One Issue",
			e: ServiceError{
				Issues: []string{"ErrMsg Test"},
			},
			want: "ErrMsg Test",
		},
		{
			name: "Issue with Status",
			e: ServiceError{
				StatusCode: 400,
				Issues:     []string{"ErrMsg Test"},
			},
			want: "400: ErrMsg Test",
		},
		{
			name: "Multiple Issue with Status",
			e: ServiceError{
				StatusCode: 400,
				Issues:     []string{"ErrMsg 1.", "ErrMsg 2."},
			},
			want: "400: ErrMsg 1. ErrMsg 2.",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.Error(); got != tt.want {
				t.Errorf("ServiceError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServiceError_HasIssues(t *testing.T) {
	tests := []struct {
		name string
		e    ServiceError
		want bool
	}{
		{
			name: "Empty",
			e:    ServiceError{},
			want: false,
		},
		{
			name: "Status Code Only",
			e: ServiceError{
				StatusCode: 412,
			},
			want: false,
		},
		{
			name: "One Issue",
			e: ServiceError{
				Issues: []string{"ErrMsg Test"},
			},
			want: true,
		},
		{
			name: "Issue with Status",
			e: ServiceError{
				StatusCode: 400,
				Issues:     []string{"ErrMsg Test"},
			},
			want: true,
		},
		{
			name: "Multiple Issue with Status",
			e: ServiceError{
				StatusCode: 400,
				Issues:     []string{"ErrMsg 1.", "ErrMsg 2."},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.HasIssues(); got != tt.want {
				t.Errorf("ServiceError.HasIssues() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServiceError_addIssue(t *testing.T) {
	type args struct {
		iss string
	}
	tests := []struct {
		name string
		e    *ServiceError
		args args
	}{
		{
			name: "Empty",
			e:    &ServiceError{},
			args: args{iss: "test error 1"},
		},
		{
			name: "One Issue",
			e: &ServiceError{
				Issues: []string{"test error 1"},
			},
			args: args{iss: "test error 2"},
		},
		{
			name: "Issue with Status",
			e: &ServiceError{
				StatusCode: 400,
				Issues:     []string{"ErrMsg Test"},
			},
			args: args{iss: "test error 2"},
		},
		{
			name: "Multiple Issue with Status",
			e: &ServiceError{
				StatusCode: 400,
				Issues:     []string{"ErrMsg 1.", "ErrMsg 2."},
			},
			args: args{iss: "test error 3"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.e.addIssue(tt.args.iss)
			if tt.e.Issues[len(tt.e.Issues)-1] != tt.args.iss {
				t.Errorf("ServiceError.AddIssue() did not add %v", tt.args.iss)
			}
		})
	}
}
