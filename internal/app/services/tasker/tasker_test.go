package tasker

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/Amirhossein2000/RequestTasker/internal/app/services/logger"
	"github.com/Amirhossein2000/RequestTasker/internal/domain/common"
	"github.com/Amirhossein2000/RequestTasker/internal/domain/entities"
	"github.com/Amirhossein2000/RequestTasker/internal/infrastructures/kafka"
	"github.com/Amirhossein2000/RequestTasker/internal/infrastructures/mysql"
	"github.com/Amirhossein2000/RequestTasker/internal/pkg/integration"
	"github.com/Amirhossein2000/RequestTasker/internal/pkg/test"
	. "github.com/smartystreets/goconvey/convey"
)

func TestTasker(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	addr, cleanup, err := integration.SetupKafkaContainer(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	defer cleanup()

	taskEventRepository, err := kafka.NewTaskEventRepository(ctx,
		kafka.Config{
			Brokers:           []string{addr},
			Topic:             "test-topic",
			GroupID:           "test-group-id",
			Timeout:           time.Second * 10,
			NumPartitions:     1,
			ReplicationFactor: 1,
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	conn, tearDown, err := integration.SetupMySQLContainer()
	if err != nil {
		t.Fatal(err)
	}
	defer tearDown()

	logger, err := logger.NewLogger(false)
	if err != nil {
		panic(err)
	}

	taskRepository := mysql.NewTaskRepository(conn, common.TaskTable)
	taskStatusRepository := mysql.NewTaskStatusRepository(conn, common.TaskStatusTable)
	taskResultRepository := mysql.NewTaskResultRepository(conn, common.TaskResultTable)

	tasker := NewTasker(
		logger,
		taskEventRepository,
		taskRepository,
		taskStatusRepository,
		taskResultRepository,
		http.DefaultClient,
	)
	tasker.Start(ctx)

	port, url := getPortAndUrl()
	newTask := entities.NewTask(
		url,
		http.MethodPost,
		map[string]string{
			"1": "1",
			"2": "2",
		},
		`
			{
				"1":1,
				"2":"2",
				"test":"test"
			}
		`,
	)
	expectedTask, err := taskRepository.Create(ctx, newTask)
	if err != nil {
		t.Fatal(err)
	}
	newTestThirdPartyServer(port, thirdPartyHandler)

	Convey("Register and send a request", t, func() {
		Convey("When status is done", func() {
			err := tasker.Process(ctx, *expectedTask)
			So(err, ShouldBeNil)

			time.Sleep(time.Second * 5)

			returnedStatus, err := taskStatusRepository.GetLatestByTaskID(ctx, expectedTask.ID())
			So(err, ShouldBeNil)
			So(returnedStatus.Status(), ShouldEqual, common.StatusDone)

			result, err := taskResultRepository.GetByTaskID(ctx, expectedTask.ID())
			So(err, ShouldBeNil)
			So(result.StatusCode(), ShouldEqual, http.StatusOK)
			So(result.Length(), ShouldEqual, 23)
			So(result.Headers()["Testheaderkey"], ShouldEqual, "testHeaderValue")

		})
	})
}

var route = "/ThirdPartyServer"

func getPortAndUrl() (int, string) {
	port, err := test.GetAvailablePort()
	if err != nil {
		panic(err)
	}
	return port, fmt.Sprintf("http://localhost:%d/%s", port, route)
}

func newTestThirdPartyServer(port int, h http.HandlerFunc) {
	go func() {
		http.HandleFunc(route, h)

		address := fmt.Sprintf(":%d", port)
		err := http.ListenAndServe(address, nil)
		if err != nil {
			panic(err)
		}
	}()

	time.Sleep(100 * time.Millisecond)
}

func thirdPartyHandler(w http.ResponseWriter, r *http.Request) {
	resp := map[string]string{
		"testResp": "testResp",
	}
	bResp, err := json.Marshal(resp)
	if err != nil {
		panic(err)
	}

	w.Header().Add("Content-Type", "application/json")
	w.Header().Add("Testheaderkey", "testHeaderValue")
	w.WriteHeader(http.StatusOK)
	w.Write(bResp)
}
