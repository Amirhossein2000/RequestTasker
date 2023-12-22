package handlers

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/Amirhossein2000/RequestTasker/internal/app/api"
	"github.com/google/uuid"
	. "github.com/smartystreets/goconvey/convey"
)

func TestGetTaskId(t *testing.T) {
	env, terminate, err := setUpTestEnv()
	if err != nil {
		t.Fatal(err)
	}
	defer terminate()

	route := "/task/%s"

	Convey("GetTaskId all possible responses", t, func() {
		Convey("return 400 when id is missing", func() {
			req, err := env.newReq(http.MethodGet, route, nil)
			So(err, ShouldBeNil)

			resp, err := http.DefaultClient.Do(req)
			So(err, ShouldBeNil)

			So(resp.StatusCode, ShouldEqual, http.StatusBadRequest)
		})
		
		Convey("return 400 when id is invalid", func() {
			id:="invalid_uuid"
			req, err := env.newReq(http.MethodGet, fmt.Sprintf(route, id), nil)
			So(err, ShouldBeNil)

			resp, err := http.DefaultClient.Do(req)
			So(err, ShouldBeNil)

			So(resp.StatusCode, ShouldEqual, http.StatusBadRequest)
		})

		Convey("return 404 when task is not found", func() {
			body := api.GetTaskIdRequestObject{
				Id: uuid.NewString(),
			}

			req, err := env.newReq(http.MethodGet, fmt.Sprintf(route, body.Id), body)
			So(err, ShouldBeNil)

			resp, err := http.DefaultClient.Do(req)
			So(err, ShouldBeNil)

			So(resp.StatusCode, ShouldEqual, http.StatusNotFound)
		})

		Convey("return 200 when everything is fine", func() {
			body := api.GetTaskIdRequestObject{
				Id: uuid.NewString(),
			}

			req, err := env.newReq(http.MethodGet, fmt.Sprintf(route, body.Id), body)
			So(err, ShouldBeNil)

			resp, err := http.DefaultClient.Do(req)
			So(err, ShouldBeNil)

			So(resp.StatusCode, ShouldEqual, http.StatusOK)
		})
	})
}
