// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen/v2 version v2.0.0 DO NOT EDIT.
package api

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/labstack/echo/v4"
	"github.com/oapi-codegen/runtime"
	strictecho "github.com/oapi-codegen/runtime/strictmiddleware/echo"
)

// Defines values for HttpMethod.
const (
	DELETE  HttpMethod = "DELETE"
	GET     HttpMethod = "GET"
	HEAD    HttpMethod = "HEAD"
	OPTIONS HttpMethod = "OPTIONS"
	PATCH   HttpMethod = "PATCH"
	POST    HttpMethod = "POST"
	PUT     HttpMethod = "PUT"
)

// Defines values for TaskStatus.
const (
	Done      TaskStatus = "done"
	Error     TaskStatus = "error"
	InProcess TaskStatus = "in_process"
	New       TaskStatus = "new"
)

// GetTaskResponse defines model for GetTaskResponse.
type GetTaskResponse struct {
	// Headers Headers array from 3rd-party service response
	Headers *map[string]interface{} `json:"headers,omitempty"`

	// HttpStatusCode HTTP status code of the 3rd-party service response
	HttpStatusCode *int `json:"httpStatusCode,omitempty"`

	// Id Unique ID of the task
	Id string `json:"id"`

	// Length Content length of 3rd-party service response
	Length *int       `json:"length,omitempty"`
	Status TaskStatus `json:"status"`
}

// HttpMethod defines model for HttpMethod.
type HttpMethod string

// TaskRequest defines model for TaskRequest.
type TaskRequest struct {
	// Body Request body payload
	Body *string `json:"body,omitempty"`

	// Headers Headers for the HTTP request
	Headers *map[string]interface{} `json:"headers,omitempty"`
	Method  HttpMethod              `json:"method"`

	// Url URL of the 3rd-party service
	Url string `json:"url"`
}

// TaskStatus defines model for TaskStatus.
type TaskStatus string

// PostTaskJSONRequestBody defines body for PostTask for application/json ContentType.
type PostTaskJSONRequestBody = TaskRequest

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Create a new task for HTTP request to 3rd-party service
	// (POST /task)
	PostTask(ctx echo.Context) error
	// Get the status of a task
	// (GET /task/{id})
	GetTaskId(ctx echo.Context, id string) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// PostTask converts echo context to params.
func (w *ServerInterfaceWrapper) PostTask(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.PostTask(ctx)
	return err
}

// GetTaskId converts echo context to params.
func (w *ServerInterfaceWrapper) GetTaskId(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "id" -------------
	var id string

	err = runtime.BindStyledParameterWithLocation("simple", false, "id", runtime.ParamLocationPath, ctx.Param("id"), &id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter id: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetTaskId(ctx, id)
	return err
}

// This is a simple interface which specifies echo.Route addition functions which
// are present on both echo.Echo and echo.Group, since we want to allow using
// either of them for path registration
type EchoRouter interface {
	CONNECT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	HEAD(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	OPTIONS(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	TRACE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
}

// RegisterHandlers adds each server route to the EchoRouter.
func RegisterHandlers(router EchoRouter, si ServerInterface) {
	RegisterHandlersWithBaseURL(router, si, "")
}

// Registers handlers, and prepends BaseURL to the paths, so that the paths
// can be served under a prefix.
func RegisterHandlersWithBaseURL(router EchoRouter, si ServerInterface, baseURL string) {

	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	router.POST(baseURL+"/task", wrapper.PostTask)
	router.GET(baseURL+"/task/:id", wrapper.GetTaskId)

}

type PostTaskRequestObject struct {
	Body *PostTaskJSONRequestBody
}

type PostTaskResponseObject interface {
	VisitPostTaskResponse(w http.ResponseWriter) error
}

type PostTask201JSONResponse struct {
	// Id Generated unique ID for the task
	Id *string `json:"id,omitempty"`
}

func (response PostTask201JSONResponse) VisitPostTaskResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)

	return json.NewEncoder(w).Encode(response)
}

type PostTask400Response struct {
}

func (response PostTask400Response) VisitPostTaskResponse(w http.ResponseWriter) error {
	w.WriteHeader(400)
	return nil
}

type PostTask401Response struct {
}

func (response PostTask401Response) VisitPostTaskResponse(w http.ResponseWriter) error {
	w.WriteHeader(401)
	return nil
}

type PostTask500Response struct {
}

func (response PostTask500Response) VisitPostTaskResponse(w http.ResponseWriter) error {
	w.WriteHeader(500)
	return nil
}

type GetTaskIdRequestObject struct {
	Id string `json:"id"`
}

type GetTaskIdResponseObject interface {
	VisitGetTaskIdResponse(w http.ResponseWriter) error
}

type GetTaskId200JSONResponse GetTaskResponse

func (response GetTaskId200JSONResponse) VisitGetTaskIdResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type GetTaskId400Response struct {
}

func (response GetTaskId400Response) VisitGetTaskIdResponse(w http.ResponseWriter) error {
	w.WriteHeader(400)
	return nil
}

type GetTaskId401Response struct {
}

func (response GetTaskId401Response) VisitGetTaskIdResponse(w http.ResponseWriter) error {
	w.WriteHeader(401)
	return nil
}

type GetTaskId404Response struct {
}

func (response GetTaskId404Response) VisitGetTaskIdResponse(w http.ResponseWriter) error {
	w.WriteHeader(404)
	return nil
}

type GetTaskId500Response struct {
}

func (response GetTaskId500Response) VisitGetTaskIdResponse(w http.ResponseWriter) error {
	w.WriteHeader(500)
	return nil
}

// StrictServerInterface represents all server handlers.
type StrictServerInterface interface {
	// Create a new task for HTTP request to 3rd-party service
	// (POST /task)
	PostTask(ctx context.Context, request PostTaskRequestObject) (PostTaskResponseObject, error)
	// Get the status of a task
	// (GET /task/{id})
	GetTaskId(ctx context.Context, request GetTaskIdRequestObject) (GetTaskIdResponseObject, error)
}

type StrictHandlerFunc = strictecho.StrictEchoHandlerFunc
type StrictMiddlewareFunc = strictecho.StrictEchoMiddlewareFunc

func NewStrictHandler(ssi StrictServerInterface, middlewares []StrictMiddlewareFunc) ServerInterface {
	return &strictHandler{ssi: ssi, middlewares: middlewares}
}

type strictHandler struct {
	ssi         StrictServerInterface
	middlewares []StrictMiddlewareFunc
}

// PostTask operation middleware
func (sh *strictHandler) PostTask(ctx echo.Context) error {
	var request PostTaskRequestObject

	var body PostTaskJSONRequestBody
	if err := ctx.Bind(&body); err != nil {
		return err
	}
	request.Body = &body

	handler := func(ctx echo.Context, request interface{}) (interface{}, error) {
		return sh.ssi.PostTask(ctx.Request().Context(), request.(PostTaskRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "PostTask")
	}

	response, err := handler(ctx, request)

	if err != nil {
		return err
	} else if validResponse, ok := response.(PostTaskResponseObject); ok {
		return validResponse.VisitPostTaskResponse(ctx.Response())
	} else if response != nil {
		return fmt.Errorf("unexpected response type: %T", response)
	}
	return nil
}

// GetTaskId operation middleware
func (sh *strictHandler) GetTaskId(ctx echo.Context, id string) error {
	var request GetTaskIdRequestObject

	request.Id = id

	handler := func(ctx echo.Context, request interface{}) (interface{}, error) {
		return sh.ssi.GetTaskId(ctx.Request().Context(), request.(GetTaskIdRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "GetTaskId")
	}

	response, err := handler(ctx, request)

	if err != nil {
		return err
	} else if validResponse, ok := response.(GetTaskIdResponseObject); ok {
		return validResponse.VisitGetTaskIdResponse(ctx.Response())
	} else if response != nil {
		return fmt.Errorf("unexpected response type: %T", response)
	}
	return nil
}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/6xVTXPbNhD9K5htj2wkN+6Ft8RWJU3TSBMpp0ymAxMrEbYIwIulPaqH/72zoL4ssq09",
	"8cWm8PF29723iycofBW8Q8cR8ieIRYmVTp9j5KWOd18wBu8iylIgH5DYYjpQojZI6dNgLMgGtt5BDpN2",
	"Q2kivVUr8pV6T+aXoIm3KiI92AIV7XEz4G1AyMHf3GLB0GRQMocFa67jlTfYE2C5nKuYDqjCG1R+pbjE",
	"F0WxjnGNJGGs6UJ/dfa+RjW93mOyjnfH25HJurVc3qBbc9kFuPKO0bFq9wXllVm1dQnwz4QryOGnwVGj",
	"wU6ggUjTUgRNkwHhfW0JDeTfpKwDyvcedifM4U/k0qfy0dWV3BqPlpDBfLZI/77K3+vRp9FyJD8/LK8m",
	"kMFk9OEaMpjNl9PZ58UJ+JGW1jL3NUbuOubGm22Xsd1xJbsq6O3Ga9PH+P/6beUpSZbsQbskegioDsX/",
	"F8MnNDUZ1LTpMcuXT/9qPchg5anSDDnUZLsFnam2S6qN1CfbieInshnvJJR1fwXyBcYIGSCRJ8jA4WOP",
	"RhLYupUXFLa8wX1DLZAekBKLx2oWbTUC+4AU27ov3g3fDSUnH9DpYCGH92kpg6C5TAkOUuOIBXxrBTGC",
	"FuKmBnKY+5jGC7QsYOSPO3MUbQPJpw5hY4t0aXAbJfR+Qr2kPfY2bJ5TzVRjWmh7MCX76/DiVaGf27pv",
	"jozRSbloVH2YKHt/9o+UpqO5LJ0NF0LBFOovh8Nu1I/aqEPZcuaib8LpmktP9u8W6Lc+oKljJKc3e1OM",
	"kqUkoVhXlabtIRmllcPHVFMq8LT5FPuexhCU5I7BkzWNxF5jj0N278/UJFeRrpBT93976ci2sil+lFbQ",
	"VRqzBs6tkJ3Ieq7I945Nhm/m0PPntUft2R9vKvTl8LJ76LNn9buv3Y9ZYYycmN+9yX6ldKtC0zTNPwEA",
	"AP//q9dF2GQIAAA=",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %w", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	res := make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	resolvePath := PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		pathToFile := url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}
