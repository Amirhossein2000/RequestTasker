package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/Amirhossein2000/RequestTasker/internal/app/api"
	"github.com/google/uuid"
	"github.com/samber/lo"
	. "github.com/smartystreets/goconvey/convey"
)

func TestPostTask(t *testing.T) {
	env, terminate, err := setUpTestEnv()
	if err != nil {
		t.Fatal(err)
	}
	defer terminate()
	ctx := context.Background()

	Convey("PostTask all possible responses", t, func() {
		Convey("return 400 when url is invalid", func() {
			reqBody := api.PostTaskJSONRequestBody{
				Body: lo.ToPtr(`{"test":"test"}`),
				Headers: &map[string]interface{}{
					"test": "test",
				},
				Method: "GET",
				Url:    "InvalidURL",
			}

			req, err := env.newReq(http.MethodPost, "/task", reqBody)
			So(err, ShouldBeNil)

			resp, err := http.DefaultClient.Do(req)
			So(err, ShouldBeNil)

			So(resp.StatusCode, ShouldEqual, http.StatusBadRequest)
		})

		Convey("return 400 when method is invalid", func() {
			reqBody := api.PostTaskJSONRequestBody{
				Body: lo.ToPtr(`{"test":"test"}`),
				Headers: &map[string]interface{}{
					"test": "test",
				},
				Method: "InvalidMethod",
				Url:    "www.google.com",
			}

			req, err := env.newReq(http.MethodPost, "/task", reqBody)
			So(err, ShouldBeNil)

			resp, err := http.DefaultClient.Do(req)
			So(err, ShouldBeNil)

			So(resp.StatusCode, ShouldEqual, http.StatusBadRequest)
		})

		Convey("return 401 when api-key is missing", func() {
			reqBody := api.PostTaskJSONRequestBody{
				Body: lo.ToPtr(`{"test":"test"}`),
				Headers: &map[string]interface{}{
					"test": "test",
				},
				Method: "GET",
				Url:    "www.google.com",
			}

			byteBody, err := json.Marshal(reqBody)
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
			reqBody := api.PostTaskJSONRequestBody{
				Body: lo.ToPtr(`{"test":"test"}`),
				Headers: &map[string]interface{}{
					"test": "test",
				},
				Method: "GET",
				Url:    "https://www.google.com",
			}

			req, err := env.newReq(http.MethodPost, "/task", reqBody)
			So(err, ShouldBeNil)

			resp, err := http.DefaultClient.Do(req)
			So(err, ShouldBeNil)
			So(resp.StatusCode, ShouldEqual, http.StatusCreated)

			respByteBody, err := io.ReadAll(resp.Body)
			So(err, ShouldBeNil)

			respBody := api.PostTask201JSONResponse{}
			err = json.Unmarshal(respByteBody, &respBody)
			So(err, ShouldBeNil)

			So(respBody.Id, ShouldNotBeEmpty)

			publicID, err := uuid.Parse(respBody.Id)
			So(err, ShouldBeNil)

			createdTask, err := env.taskRepository.GetByPublicID(ctx, publicID)
			So(err, ShouldBeNil)

			So(createdTask.Body(), ShouldEqual, *reqBody.Body)
			So(createdTask.Headers(), ShouldEqual, convertHeadersForRequest(*reqBody.Headers))
			So(createdTask.Method(), ShouldEqual, string(reqBody.Method))
			So(createdTask.Url(), ShouldEqual, reqBody.Url)
		})
	})
}
