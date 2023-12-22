package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/Amirhossein2000/RequestTasker/internal/app/api"
	"github.com/samber/lo"
	. "github.com/smartystreets/goconvey/convey"
)

func TestPostTask(t *testing.T) {
	env, terminate, err := setUpTestEnv()
	if err != nil {
		t.Fatal(err)
	}
	defer terminate()

	Convey("PostTask all possible responses", t, func() {
		Convey("return 400 when url is invalid", func() {
			body := api.PostTaskJSONRequestBody{
				Body: lo.ToPtr(`{"test":"test"}`),
				Headers: &map[string]interface{}{
					"test": "test",
				},
				Method: "GET",
				Url:    "InvalidURL",
			}

			byteBody, err := json.Marshal(body)
			So(err, ShouldBeNil)

			rBody := bytes.NewBuffer(byteBody)
			req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://%s/task", env.addr), rBody)
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", env.apiKey)
			So(err, ShouldBeNil)

			resp, err := http.DefaultClient.Do(req)
			So(err, ShouldBeNil)

			So(resp.StatusCode, ShouldEqual, http.StatusBadRequest)
		})

		Convey("return 400 when method is invalid", func() {
			body := api.PostTaskJSONRequestBody{
				Body: lo.ToPtr(`{"test":"test"}`),
				Headers: &map[string]interface{}{
					"test": "test",
				},
				Method: "InvalidMethod",
				Url:    "www.google.com",
			}

			byteBody, err := json.Marshal(body)
			So(err, ShouldBeNil)

			rBody := bytes.NewBuffer(byteBody)
			req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://%s/task", env.addr), rBody)
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", env.apiKey)
			So(err, ShouldBeNil)

			resp, err := http.DefaultClient.Do(req)
			So(err, ShouldBeNil)

			So(resp.StatusCode, ShouldEqual, http.StatusBadRequest)
		})

		Convey("return 401 when api-key is missing", func() {
			body := api.PostTaskJSONRequestBody{
				Body: lo.ToPtr(`{"test":"test"}`),
				Headers: &map[string]interface{}{
					"test": "test",
				},
				Method: "GET",
				Url:    "www.google.com",
			}

			byteBody, err := json.Marshal(body)
			So(err, ShouldBeNil)

			rBody := bytes.NewBuffer(byteBody)
			req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://%s/task", env.addr), rBody)
			req.Header.Set("Content-Type", "application/json")
			So(err, ShouldBeNil)

			resp, err := http.DefaultClient.Do(req)
			So(err, ShouldBeNil)

			So(resp.StatusCode, ShouldEqual, http.StatusUnauthorized)
		})

		Convey("return 201 when everything is fine", func() {
			body := api.PostTaskJSONRequestBody{
				Body: lo.ToPtr(`{"test":"test"}`),
				Headers: &map[string]interface{}{
					"test": "test",
				},
				Method: "GET",
				Url:    "https://www.google.com",
			}

			byteBody, err := json.Marshal(body)
			So(err, ShouldBeNil)

			rBody := bytes.NewBuffer(byteBody)
			req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://%s/task", env.addr), rBody)
			So(err, ShouldBeNil)

			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", env.apiKey)

			resp, err := http.DefaultClient.Do(req)
			So(err, ShouldBeNil)

			So(resp.StatusCode, ShouldEqual, http.StatusCreated)
			// TODO: check the response body and the created public-id
		})
	})
}
