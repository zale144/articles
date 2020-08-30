package users_test

import (
	"articles/usertags/internal/api/handler/users"
	"articles/usertags/internal/dto"
	"articles/usertags/internal/pkg/middleware"
	"bytes"
	"encoding/json"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

type mockUserService struct {
	fail bool
}

func (m mockUserService) Login(dto.LoginPayload) (string, error) {
	if m.fail {
		return "", errors.New("")
	}
	return "token", nil
}

func TestLogin(t *testing.T) {
	type args struct {
		svc  users.UserService
		body dto.LoginPayload
	}
	tests := []struct {
		name     string
		args     args
		wantCode int
		want     string
	}{
		{
			name: "success",
			args: args{
				svc: mockUserService{},
				body: dto.LoginPayload{
					Email:    "user@test.com",
					Password: "secret",
				},
			},
			want:     "token",
			wantCode: http.StatusOK,
		}, {
			name: "fail service",
			args: args{
				svc: mockUserService{true},
				body: dto.LoginPayload{
					Email:    "user@test.com",
					Password: "secret",
				},
			},
			wantCode: http.StatusUnauthorized,
		}, {
			name: "fail validation",
			args: args{
				svc: mockUserService{},
				body: dto.LoginPayload{
					Password: "secret",
				},
			},
			wantCode: http.StatusBadRequest,
		},
	}

	e := echo.New()
	e.Validator = middleware.NewCustomValidator(validator.New())

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			b, _ := json.Marshal(tt.args.body)

			body := ioutil.NopCloser(bytes.NewBuffer(b))
			h := http.Header{}
			h.Add("Content-Type", "application/json")

			rec := httptest.NewRecorder()

			ctx := e.NewContext(&http.Request{
				Method:        echo.POST,
				URL:           &url.URL{},
				Body:          body,
				ContentLength: int64(len(b)),
				Header:        h,
			}, rec)

			err := users.Login(tt.args.svc)(ctx)
			require.Nil(t, err, "error getting card", err)
			assert.Equal(t, tt.wantCode, rec.Code, "status code is not equal")

			resp := rec.Result()
			buf, _ := ioutil.ReadAll(resp.Body)

			rsp := dto.LoginResponse{}
			_ = json.Unmarshal(buf, &rsp)

			assert.Equal(t, tt.want, rsp.Token, "response does not match expected output")
		})
	}
}

func (m mockUserService) Register(dto.RegisterPayload) error {
	if m.fail {
		return errors.New("internal")
	}
	return nil
}

func TestRegister(t *testing.T) {
	type args struct {
		svc  users.UserService
		body dto.RegisterPayload
	}
	tests := []struct {
		name     string
		args     args
		wantCode int
		want     string
	}{
		{
			name: "success",
			args: args{
				svc: mockUserService{},
				body: dto.RegisterPayload{
					Name: "User1",
					LoginPayload: dto.LoginPayload{
						Email:    "user@test.com",
						Password: "secret",
					},
				},
			},
			want:     "registration successful",
			wantCode: http.StatusOK,
		}, {
			name: "fail service",
			args: args{
				svc: mockUserService{true},
				body: dto.RegisterPayload{
					Name: "User1",
					LoginPayload: dto.LoginPayload{
						Email:    "user@test.com",
						Password: "secret",
					},
				},
			},
			want:     "could not register user: internal",
			wantCode: http.StatusBadRequest,
		}, {
			name: "fail validation",
			args: args{
				svc: mockUserService{true},
				body: dto.RegisterPayload{
					Name: "User1",
					LoginPayload: dto.LoginPayload{
						Password: "secret",
					},
				},
			},
			want:     "Key: 'RegisterPayload.LoginPayload.Email' Error:Field validation for 'Email' failed on the 'required' tag",
			wantCode: http.StatusBadRequest,
		},
	}

	e := echo.New()
	e.Validator = middleware.NewCustomValidator(validator.New())

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			b, _ := json.Marshal(tt.args.body)

			body := ioutil.NopCloser(bytes.NewBuffer(b))
			h := http.Header{}
			h.Add("Content-Type", "application/json")

			rec := httptest.NewRecorder()

			ctx := e.NewContext(&http.Request{
				Method:        echo.POST,
				URL:           &url.URL{},
				Body:          body,
				ContentLength: int64(len(b)),
				Header:        h,
			}, rec)

			err := users.Register(tt.args.svc)(ctx)
			require.Nil(t, err, "error getting card", err)
			assert.Equal(t, tt.wantCode, rec.Code, "status code is not equal")

			resp := rec.Result()
			buf, _ := ioutil.ReadAll(resp.Body)

			rsp := &dto.ResponseMessage{}
			_ = json.Unmarshal(buf, rsp)

			assert.Equal(t, tt.want, rsp.Message, "response does not match expected output")
		})
	}
}
