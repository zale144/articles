package tags

import (
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
	"reflect"
	"testing"
)

type mockTagService struct {
	fail bool
}

func (m mockTagService) Add(string, dto.AddTagsPayload) error {
	if m.fail {
		return errors.New("internal")
	}
	return nil
}

func TestAdd(t *testing.T) {
	type args struct {
		svc  TagService
		body dto.AddTagsPayload
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
				svc: mockTagService{},
				body: dto.AddTagsPayload{
					Tags: []string{"tag1", "tag2"},
				},
			},
			want:     "tags added successfully",
			wantCode: http.StatusOK,
		}, {
			name: "fail service",
			args: args{
				svc: mockTagService{true},
				body: dto.AddTagsPayload{
					Tags: []string{"tag1", "tag2"},
				},
			},
			want: "could not add tag: internal",
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

			ctx.Set("email", "user@test.com")

			err := Add(tt.args.svc)(ctx)
			if err != nil {
				t.Error(err)
			}

			if rec.Code != tt.wantCode {
				t.Errorf("Add() = %v, wantCode %v", rec.Code, tt.wantCode)
			}

			resp := rec.Result()
			buf, _ := ioutil.ReadAll(resp.Body)

			rsp := &dto.ResponseMessage{}
			_ = json.Unmarshal(buf, rsp)


			if rsp.Message != tt.want {
				t.Errorf("Add() = %v, want %v", rsp.Message, tt.want)
			}
		})
	}
}

func (m mockTagService) Get(string) (dto.GetTagsPayload, error) {
	if m.fail {
		return dto.GetTagsPayload{}, errors.New("")
	}
	return dto.GetTagsPayload{
		Tags: []string{"tag1", "tag2"},
	}, nil
}

func TestGet(t *testing.T) {
	type args struct {
		svc  TagService
	}
	tests := []struct {
		name     string
		args     args
		wantCode int
		want     []string
	}{
		{
			name: "success",
			args: args{
				svc: mockTagService{},
			},
			want:    []string{"tag1", "tag2"},
			wantCode: http.StatusOK,
		}, {
			name: "fail service",
			args: args{
				svc: mockTagService{
					fail: true,
				},
			},
			wantCode: http.StatusBadRequest,
		},
	}

	e := echo.New()
	e.Validator = middleware.NewCustomValidator(validator.New())

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			h := http.Header{}

			h.Add("Content-Type", "application/json")

			rec := httptest.NewRecorder()

			ctx := e.NewContext(&http.Request{
				Method:        echo.GET,
				URL:           &url.URL{},
				Header:        h,
			}, rec)

			ctx.Set("email", "user@test.com")

			err := Get(tt.args.svc)(ctx)
			if err != nil {
				t.Error(err)
			}

			if rec.Code != tt.wantCode {
				t.Errorf("Get() = %v, wantCode %v", rec.Code, tt.wantCode)
			}

			resp := rec.Result()
			buf, _ := ioutil.ReadAll(resp.Body)

			rsp := &dto.GetTagsPayload{}
			_ = json.Unmarshal(buf, rsp)


			if !reflect.DeepEqual(rsp.Tags, tt.want) {
				t.Errorf("Get() = %v, want %v", rsp.Tags, tt.want)
			}
		})
	}
}
