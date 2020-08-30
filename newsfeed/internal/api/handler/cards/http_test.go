package cards

import (
	"articles/newsfeed/internal/dto"
	"articles/newsfeed/internal/pkg/middleware"
	"bytes"
	"encoding/json"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/magiconair/properties/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

type mockCardService struct {
	fail bool
}

func (m mockCardService) Add(dto.Card) error {
	if m.fail {
		return errors.New("internal")
	}
	return nil
}

func TestAdd(t *testing.T) {
	type args struct {
		svc  CardService
		body dto.Card
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
				svc: mockCardService{},
				body: dto.Card{
					Title:     "card1",
					Timestamp: time.Now(),
					Tags: []string{
						"tag1", "tag2",
					},
				},
			},
			want:     "card added successfully",
			wantCode: http.StatusOK,
		}, {
			name: "fail service",
			args: args{
				svc: mockCardService{true},
				body: dto.Card{
					Title:     "card1",
					Timestamp: time.Now(),
					Tags: []string{
						"tag1", "tag2",
					},
				},
			},
			want:     "could not add card: internal",
			wantCode: http.StatusBadRequest,
		}, {
			name: "fail validation",
			args: args{
				svc: mockCardService{false},
				body: dto.Card{
					Title:     "card1",
					Timestamp: time.Now(),
				},
			},
			want:     "Key: 'Card.Tags' Error:Field validation for 'Tags' failed on the 'required' tag",
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

			err := Add(tt.args.svc)(ctx)
			require.Nil(t, err, "error adding card", err)
			assert.Equal(t, tt.wantCode, rec.Code, "status code is not equal")

			resp := rec.Result()
			buf, _ := ioutil.ReadAll(resp.Body)

			rsp := &dto.ResponseMessage{}
			_ = json.Unmarshal(buf, rsp)
			assert.Equal(t, tt.want, rsp.Message, "response message is not equal")
		})
	}
}

func (m mockCardService) GetByUser(email string) (dto.GetCardsPayload, error) {
	if m.fail {
		return dto.GetCardsPayload{}, errors.New("")
	}

	return getCardsPayload, nil
}

func TestGet(t *testing.T) {
	type args struct {
		svc CardService
	}
	tests := []struct {
		name     string
		args     args
		wantCode int
		want     dto.GetCardsPayload
	}{
		{
			name: "success",
			args: args{
				svc: mockCardService{},
			},
			want:     getCardsPayload,
			wantCode: http.StatusOK,
		}, {
			name: "fail service",
			args: args{
				svc: mockCardService{
					fail: true,
				},
			},
			want:     dto.GetCardsPayload{},
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
				Method: echo.GET,
				URL:    &url.URL{},
				Header: h,
			}, rec)

			ctx.Set("email", "user@test.com")

			err := Get(tt.args.svc)(ctx)
			require.Nil(t, err, "error getting card", err)
			assert.Equal(t, tt.wantCode, rec.Code, "status code is not equal")

			resp := rec.Result()
			buf, _ := ioutil.ReadAll(resp.Body)

			rsp := dto.GetCardsPayload{}
			_ = json.Unmarshal(buf, &rsp)

			assert.Equal(t, tt.want, rsp, "response does not match expected output")
		})
	}
}

func (m mockCardService) GetByTags(tags []string) (dto.GetCardsPayload, error) {
	if m.fail {
		return dto.GetCardsPayload{}, errors.New("")
	}

	return getCardsPayload, nil
}

func TestGetByTags(t *testing.T) {
	type args struct {
		svc  CardService
		tags []string
	}
	tests := []struct {
		name     string
		args     args
		wantCode int
		want     dto.GetCardsPayload
	}{
		{
			name: "success",
			args: args{
				svc:  mockCardService{},
				tags: []string{"tag1", "tag2"},
			},
			want:     getCardsPayload,
			wantCode: http.StatusOK,
		}, {
			name: "fail service",
			args: args{
				svc:  mockCardService{fail: true},
				tags: []string{"tag1", "tag2"},
			},
			wantCode: http.StatusBadRequest,
		}, {
			name: "fail validation: not enough tags",
			args: args{
				svc:  mockCardService{},
				tags: []string{"tag1"},
			},
			wantCode: http.StatusBadRequest,
		}, {
			name: "fail validation: tags not unique",
			args: args{
				svc:  mockCardService{},
				tags: []string{"tag1", "tag1"},
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

			q := url.Values{}
			for _, tag := range tt.args.tags {
				q.Add("tags", tag)
			}

			ctx := e.NewContext(&http.Request{
				Method: echo.GET,
				URL: &url.URL{
					RawQuery: q.Encode(),
				},
				Header: h,
			}, rec)

			err := GetByTags(tt.args.svc)(ctx)
			require.Nil(t, err, "error getting card", err)
			assert.Equal(t, tt.wantCode, rec.Code, "status code is not equal")

			resp := rec.Result()
			buf, _ := ioutil.ReadAll(resp.Body)

			rsp := dto.GetCardsPayload{}
			_ = json.Unmarshal(buf, &rsp)

			assert.Equal(t, tt.want, rsp, "response does not match expected output")
		})
	}
}

var getCardsPayload = dto.GetCardsPayload{
	Cards: []dto.Card{
		{
			Title:     "card1",
			Timestamp: time.Time{},
			Tags:      []string{"tag1", "tag2"},
		}, {
			Title:     "card2",
			Timestamp: time.Time{},
			Tags:      []string{"tag3", "tag3"},
		},
	},
}
