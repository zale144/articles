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
			if err != nil {
				t.Error(err)
			}

			if rec.Code != tt.wantCode {
				t.Errorf("Login() = %v, wantCode %v", rec.Code, tt.wantCode)
			}

			resp := rec.Result()
			buf, _ := ioutil.ReadAll(resp.Body)

			rsp := &dto.LoginResponse{}
			_ = json.Unmarshal(buf, rsp)

			if rsp.Token != tt.want {
				t.Errorf("Login() = %v, want %v", rsp.Token, tt.want)
			}
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
			if err != nil {
				t.Error(err)
			}

			if rec.Code != tt.wantCode {
				t.Errorf("Register() = %v, wantCode %v", rec.Code, tt.wantCode)
			}

			resp := rec.Result()
			buf, _ := ioutil.ReadAll(resp.Body)

			rsp := &dto.ResponseMessage{}
			_ = json.Unmarshal(buf, rsp)

			if rsp.Message != tt.want {
				t.Errorf("Register() = %v, want %v", rsp.Message, tt.want)
			}
		})
	}
}
