package employee

import (
	"encoding/json"
	"errors"
	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"idm/inner/common"
	"idm/inner/web"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

type MockService struct {
	mock.Mock
}

func (svc *MockService) Create(request CreateRequest) (int64, error) {
	args := svc.Called(request)
	return args.Get(0).(int64), args.Error(1)
}

func (svc *MockService) FindById(request IdRequest) (Response, error) {
	args := svc.Called(request)
	return args.Get(0).(Response), args.Error(1)
}

func (svc *MockService) FindAll() ([]Response, error) {
	args := svc.Called()
	return args.Get(0).([]Response), args.Error(1)
}

func (svc *MockService) FindAllByIds(request IdsRequest) ([]Response, error) {
	args := svc.Called(request)
	return args.Get(0).([]Response), args.Error(1)
}

func (svc *MockService) DeleteById(request IdRequest) error {
	args := svc.Called(request)
	return args.Error(0)
}

func (svc *MockService) DeleteAllByIds(request IdsRequest) error {
	args := svc.Called(request)
	return args.Error(0)
}

func TestControllerCreateEmployee(t *testing.T) {
	a := assert.New(t)

	url := "/api/v1/employees"

	t.Run("should return created employee id", func(t *testing.T) {
		server := web.NewServer()
		svc := new(MockService)
		controller := NewController(server, svc)
		controller.RegisterRoutes()

		body := strings.NewReader(`{"name":"john doe"}`)
		req := httptest.NewRequest(fiber.MethodPost, url, body)
		req.Header.Set("Content-Type", "application/json")

		svc.On("Create", mock.AnythingOfType("CreateRequest")).Return(int64(123), nil)

		resp, err := server.App.Test(req)
		a.Nil(err)
		a.NotEmpty(resp)
		a.Equal(http.StatusOK, resp.StatusCode)

		bytesData, err := io.ReadAll(resp.Body)
		a.Nil(err)

		var responseBody common.Response[int64]
		err = json.Unmarshal(bytesData, &responseBody)
		a.Nil(err)
		a.Equal(int64(123), responseBody.Data)
		a.True(responseBody.Success)
		a.Empty(responseBody.Message)
	})

	t.Run("should return bad request on invalid json", func(t *testing.T) {
		server := web.NewServer()
		svc := new(MockService)
		controller := NewController(server, svc)
		controller.RegisterRoutes()

		body := strings.NewReader(`{invalid json`)
		req := httptest.NewRequest(fiber.MethodPost, url, body)
		req.Header.Set("Content-Type", "application/json")

		resp, err := server.App.Test(req)
		a.Nil(err)
		a.Equal(http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should return bad request on validation error", func(t *testing.T) {
		server := web.NewServer()
		svc := new(MockService)
		controller := NewController(server, svc)
		controller.RegisterRoutes()

		request := CreateRequest{Name: ""}
		body := strings.NewReader(`{"name": ""}`)
		req := httptest.NewRequest(fiber.MethodPost, url, body)
		createErr := common.RequestValidationError{Message: "name is required"}
		req.Header.Set("Content-Type", "application/json")

		svc.On("Create", request).Return(int64(0), createErr)

		resp, err := server.App.Test(req)
		a.Nil(err)
		a.Equal(http.StatusBadRequest, resp.StatusCode)
		bytesData, err := io.ReadAll(resp.Body)
		a.Nil(err)

		var responseBody common.Response[any]
		err = json.Unmarshal(bytesData, &responseBody)
		a.Nil(err)
		a.False(responseBody.Success)
		a.Equal(responseBody.Message, "name is required")
	})

	t.Run("should return bad request on already exists error", func(t *testing.T) {
		server := web.NewServer()
		svc := new(MockService)
		controller := NewController(server, svc)
		controller.RegisterRoutes()

		request := CreateRequest{Name: "john doe"}
		body := strings.NewReader(`{"name":"john doe"}`)
		req := httptest.NewRequest(fiber.MethodPost, url, body)
		req.Header.Set("Content-Type", "application/json")

		errAlready := common.AlreadyExistsError{Message: "employee already exists"}
		svc.On("Create", request).Return(int64(0), errAlready)

		resp, err := server.App.Test(req)
		a.Nil(err)
		a.Equal(http.StatusBadRequest, resp.StatusCode)
		bytesData, _ := io.ReadAll(resp.Body)

		var responseBody common.Response[any]
		_ = json.Unmarshal(bytesData, &responseBody)
		a.False(responseBody.Success)
		a.Equal(responseBody.Message, "employee already exists")
	})

	t.Run("should return internal server error on generic error", func(t *testing.T) {
		server := web.NewServer()
		svc := new(MockService)
		controller := NewController(server, svc)
		controller.RegisterRoutes()

		request := CreateRequest{Name: "john doe"}
		body := strings.NewReader(`{"name":"john doe"}`)
		req := httptest.NewRequest(fiber.MethodPost, url, body)
		req.Header.Set("Content-Type", "application/json")

		svc.On("Create", request).Return(int64(0), errors.New("unexpected server error"))

		resp, err := server.App.Test(req)
		a.Nil(err)
		a.Equal(http.StatusInternalServerError, resp.StatusCode)
		bytesData, _ := io.ReadAll(resp.Body)

		var responseBody common.Response[any]
		_ = json.Unmarshal(bytesData, &responseBody)
		a.False(responseBody.Success)
		a.Equal(responseBody.Message, "unexpected server error")
	})
}

func TestControllerFindById(t *testing.T) {
	a := assert.New(t)

	url := "/api/v1/employees/1"

	t.Run("should return employee by id", func(t *testing.T) {
		server := web.NewServer()
		svc := new(MockService)
		controller := NewController(server, svc)
		controller.RegisterRoutes()

		employee := Response{Id: 1, Name: "john doe", CreatedAt: time.Now(), UpdatedAt: time.Now()}
		svc.On("FindById", IdRequest{Id: 1}).Return(employee, nil)

		req := httptest.NewRequest(fiber.MethodGet, url, nil)
		resp, err := server.App.Test(req)

		a.Nil(err)
		a.NotEmpty(resp)
		a.Equal(http.StatusOK, resp.StatusCode)

		bytesData, err := io.ReadAll(resp.Body)
		a.Nil(err)

		var responseBody common.Response[Response]
		err = json.Unmarshal(bytesData, &responseBody)
		a.Nil(err)
		a.True(responseBody.Success)
		a.Empty(responseBody.Message)
		a.Equal(employee.Id, responseBody.Data.Id)
		a.Equal(employee.Name, responseBody.Data.Name)
		a.WithinDuration(employee.CreatedAt, responseBody.Data.CreatedAt, time.Second)
		a.WithinDuration(employee.UpdatedAt, responseBody.Data.UpdatedAt, time.Second)
	})

	t.Run("should return not found error", func(t *testing.T) {
		server := web.NewServer()
		svc := new(MockService)
		controller := NewController(server, svc)
		controller.RegisterRoutes()

		notFoundErr := common.NotFoundError{Message: "employee with id 1 not found"}
		svc.On("FindById", IdRequest{Id: 1}).Return(Response{}, notFoundErr)

		req := httptest.NewRequest(fiber.MethodGet, url, nil)
		resp, err := server.App.Test(req)
		a.Nil(err)
		a.Equal(http.StatusNotFound, resp.StatusCode)

		bytesData, err := io.ReadAll(resp.Body)
		a.Nil(err)

		var responseBody common.Response[any]
		err = json.Unmarshal(bytesData, &responseBody)
		a.Nil(err)
		a.False(responseBody.Success)
		a.Equal("employee with id 1 not found", responseBody.Message)
	})

	t.Run("should return bad request on invalid id", func(t *testing.T) {
		server := web.NewServer()
		svc := new(MockService)
		controller := NewController(server, svc)
		controller.RegisterRoutes()

		req := httptest.NewRequest(fiber.MethodGet, "/api/v1/employees/abc", nil)
		resp, err := server.App.Test(req)

		a.Nil(err)
		a.Equal(http.StatusBadRequest, resp.StatusCode)

		bytesData, err := io.ReadAll(resp.Body)
		a.Nil(err)

		var responseBody common.Response[any]
		err = json.Unmarshal(bytesData, &responseBody)
		a.Nil(err)
		a.False(responseBody.Success)
		a.Equal("invalid id parameter", responseBody.Message)
	})

	t.Run("should return internal server error on generic error", func(t *testing.T) {
		server := web.NewServer()
		svc := new(MockService)
		controller := NewController(server, svc)
		controller.RegisterRoutes()

		svc.On("FindById", IdRequest{Id: 1}).Return(Response{}, errors.New("unexpected server error"))

		req := httptest.NewRequest(fiber.MethodGet, url, nil)
		resp, err := server.App.Test(req)
		a.Nil(err)
		a.Equal(http.StatusInternalServerError, resp.StatusCode)

		bytesData, err := io.ReadAll(resp.Body)
		a.Nil(err)

		var responseBody common.Response[any]
		err = json.Unmarshal(bytesData, &responseBody)
		a.Nil(err)
		a.False(responseBody.Success)
		a.Equal("unexpected server error", responseBody.Message)
	})
}

func TestControllerFindAll(t *testing.T) {
	a := assert.New(t)
	url := "/api/v1/employees"

	t.Run("should return all employees", func(t *testing.T) {
		server := web.NewServer()
		svc := new(MockService)
		controller := NewController(server, svc)
		controller.RegisterRoutes()

		employees := []Response{
			{Id: 1, Name: "Alice", CreatedAt: time.Now(), UpdatedAt: time.Now()},
			{Id: 2, Name: "Bob", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		}
		svc.On("FindAll").Return(employees, nil)

		req := httptest.NewRequest(fiber.MethodGet, url, nil)
		resp, err := server.App.Test(req)
		a.Nil(err)
		a.NotEmpty(resp)
		a.Equal(http.StatusOK, resp.StatusCode)

		bytesData, err := io.ReadAll(resp.Body)
		a.Nil(err)

		var responseBody common.Response[[]Response]
		err = json.Unmarshal(bytesData, &responseBody)
		a.Nil(err)
		a.True(responseBody.Success)
		a.Empty(responseBody.Message)
		a.Equal(len(employees), len(responseBody.Data))
		for i, emp := range employees {
			got := responseBody.Data[i]
			a.Equal(emp.Id, got.Id)
			a.Equal(emp.Name, got.Name)
			a.WithinDuration(emp.CreatedAt, got.CreatedAt, time.Second)
			a.WithinDuration(emp.UpdatedAt, got.UpdatedAt, time.Second)
		}
	})

	t.Run("should return not found error", func(t *testing.T) {
		server := web.NewServer()
		svc := new(MockService)
		controller := NewController(server, svc)
		controller.RegisterRoutes()

		errNotFound := common.NotFoundError{Message: "no employees found"}
		svc.On("FindAll").Return([]Response(nil), errNotFound)

		req := httptest.NewRequest(fiber.MethodGet, url, nil)
		resp, err := server.App.Test(req)
		a.Nil(err)
		a.Equal(http.StatusNotFound, resp.StatusCode)

		bytesData, err := io.ReadAll(resp.Body)
		a.Nil(err)

		var responseBody common.Response[any]
		err = json.Unmarshal(bytesData, &responseBody)
		a.Nil(err)
		a.False(responseBody.Success)
		a.Equal("no employees found", responseBody.Message)
	})

	t.Run("should return internal server error on generic error", func(t *testing.T) {
		server := web.NewServer()
		svc := new(MockService)
		controller := NewController(server, svc)
		controller.RegisterRoutes()

		svc.On("FindAll").Return([]Response(nil), errors.New("unexpected server error"))

		req := httptest.NewRequest(fiber.MethodGet, url, nil)
		resp, err := server.App.Test(req)
		a.Nil(err)
		a.Equal(http.StatusInternalServerError, resp.StatusCode)

		bytesData, err := io.ReadAll(resp.Body)
		a.Nil(err)

		var responseBody common.Response[any]
		err = json.Unmarshal(bytesData, &responseBody)
		a.Nil(err)
		a.False(responseBody.Success)
		a.Equal("unexpected server error", responseBody.Message)
	})
}

func TestControllerFindAllByIds(t *testing.T) {
	a := assert.New(t)
	url := "/api/v1/employees/ids"
	validBody := `{"ids":[1,2]}`
	invalidBody := `{"ids":[]}`

	t.Run("should return all employees by ids", func(t *testing.T) {
		server := web.NewServer()
		svc := new(MockService)
		controller := NewController(server, svc)
		controller.RegisterRoutes()

		employees := []Response{
			{Id: 1, Name: "Alice", CreatedAt: time.Now(), UpdatedAt: time.Now()},
			{Id: 2, Name: "Bob", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		}
		request := IdsRequest{Ids: []int64{1, 2}}
		svc.On("FindAllByIds", request).Return(employees, nil)

		req := httptest.NewRequest(fiber.MethodPost, url, strings.NewReader(validBody))
		req.Header.Set("Content-Type", "application/json")
		resp, err := server.App.Test(req)

		a.Nil(err)
		a.NotEmpty(resp)
		a.Equal(http.StatusOK, resp.StatusCode)

		bytesData, err := io.ReadAll(resp.Body)
		a.Nil(err)

		var responseBody common.Response[[]Response]
		err = json.Unmarshal(bytesData, &responseBody)
		a.Nil(err)
		a.True(responseBody.Success)
		a.Empty(responseBody.Message)
		a.Len(responseBody.Data, len(employees))

		for i, emp := range employees {
			got := responseBody.Data[i]
			a.Equal(emp.Id, got.Id)
			a.Equal(emp.Name, got.Name)
			a.WithinDuration(emp.CreatedAt, got.CreatedAt, time.Second)
			a.WithinDuration(emp.UpdatedAt, got.UpdatedAt, time.Second)
		}
	})

	t.Run("should return validation error", func(t *testing.T) {
		server := web.NewServer()
		svc := new(MockService)
		controller := NewController(server, svc)
		controller.RegisterRoutes()

		validationErr := common.RequestValidationError{Message: "ids must not be empty"}
		request := IdsRequest{Ids: []int64{}}
		svc.On("FindAllByIds", request).Return([]Response(nil), validationErr)

		req := httptest.NewRequest(fiber.MethodPost, url, strings.NewReader(invalidBody))
		req.Header.Set("Content-Type", "application/json")
		resp, err := server.App.Test(req)

		a.Nil(err)
		a.Equal(http.StatusBadRequest, resp.StatusCode)

		bytesData, err := io.ReadAll(resp.Body)
		a.Nil(err)

		var responseBody common.Response[any]
		err = json.Unmarshal(bytesData, &responseBody)
		a.Nil(err)
		a.False(responseBody.Success)
		a.Equal("ids must not be empty", responseBody.Message)
	})

	t.Run("should return not found error", func(t *testing.T) {
		server := web.NewServer()
		svc := new(MockService)
		controller := NewController(server, svc)
		controller.RegisterRoutes()

		request := IdsRequest{Ids: []int64{1, 2}}
		notFoundErr := common.NotFoundError{Message: "employees not found"}
		svc.On("FindAllByIds", request).Return([]Response(nil), notFoundErr)

		req := httptest.NewRequest(fiber.MethodPost, url, strings.NewReader(validBody))
		req.Header.Set("Content-Type", "application/json")
		resp, err := server.App.Test(req)

		a.Nil(err)
		a.Equal(http.StatusNotFound, resp.StatusCode)

		bytesData, err := io.ReadAll(resp.Body)
		a.Nil(err)

		var responseBody common.Response[any]
		err = json.Unmarshal(bytesData, &responseBody)
		a.Nil(err)
		a.False(responseBody.Success)
		a.Equal("employees not found", responseBody.Message)
	})

	t.Run("should return internal server error on generic error", func(t *testing.T) {
		server := web.NewServer()
		svc := new(MockService)
		controller := NewController(server, svc)
		controller.RegisterRoutes()

		request := IdsRequest{Ids: []int64{1, 2}}
		svc.On("FindAllByIds", request).Return([]Response(nil), errors.New("unexpected server error"))

		req := httptest.NewRequest(fiber.MethodPost, url, strings.NewReader(validBody))
		req.Header.Set("Content-Type", "application/json")
		resp, err := server.App.Test(req)

		a.Nil(err)
		a.Equal(http.StatusInternalServerError, resp.StatusCode)

		bytesData, err := io.ReadAll(resp.Body)
		a.Nil(err)

		var responseBody common.Response[any]
		err = json.Unmarshal(bytesData, &responseBody)
		a.Nil(err)
		a.False(responseBody.Success)
		a.Equal("unexpected server error", responseBody.Message)
	})
}

func TestControllerDeleteById(t *testing.T) {
	a := assert.New(t)

	url := "/api/v1/employees/1"

	t.Run("should delete employee by id", func(t *testing.T) {
		server := web.NewServer()
		svc := new(MockService)
		controller := NewController(server, svc)
		controller.RegisterRoutes()

		svc.On("DeleteById", IdRequest{Id: 1}).Return(nil)

		req := httptest.NewRequest(fiber.MethodDelete, url, nil)
		resp, err := server.App.Test(req)

		a.Nil(err)
		a.Equal(http.StatusOK, resp.StatusCode)

		bytesData, err := io.ReadAll(resp.Body)
		a.Nil(err)

		var responseBody common.Response[any]
		err = json.Unmarshal(bytesData, &responseBody)
		a.Nil(err)
		a.True(responseBody.Success)
		a.Empty(responseBody.Message)
		a.Nil(responseBody.Data)
	})

	t.Run("should return bad request on invalid id", func(t *testing.T) {
		server := web.NewServer()
		svc := new(MockService)
		controller := NewController(server, svc)
		controller.RegisterRoutes()

		req := httptest.NewRequest(fiber.MethodDelete, "/api/v1/employees/abc", nil)
		resp, err := server.App.Test(req)

		a.Nil(err)
		a.Equal(http.StatusBadRequest, resp.StatusCode)

		bytesData, err := io.ReadAll(resp.Body)
		a.Nil(err)

		var responseBody common.Response[any]
		err = json.Unmarshal(bytesData, &responseBody)
		a.Nil(err)
		a.False(responseBody.Success)
		a.Equal("invalid id parameter", responseBody.Message)
	})

	t.Run("should return not found error", func(t *testing.T) {
		server := web.NewServer()
		svc := new(MockService)
		controller := NewController(server, svc)
		controller.RegisterRoutes()

		errNotFound := common.NotFoundError{Message: "employee not found"}
		svc.On("DeleteById", IdRequest{Id: 1}).Return(errNotFound)

		req := httptest.NewRequest(fiber.MethodDelete, url, nil)
		resp, err := server.App.Test(req)

		a.Nil(err)
		a.Equal(http.StatusNotFound, resp.StatusCode)

		bytesData, err := io.ReadAll(resp.Body)
		a.Nil(err)

		var responseBody common.Response[any]
		err = json.Unmarshal(bytesData, &responseBody)
		a.Nil(err)
		a.False(responseBody.Success)
		a.Equal("employee not found", responseBody.Message)
	})

	t.Run("should return internal server error on generic error", func(t *testing.T) {
		server := web.NewServer()
		svc := new(MockService)
		controller := NewController(server, svc)
		controller.RegisterRoutes()

		svc.On("DeleteById", IdRequest{Id: 1}).Return(errors.New("unexpected server error"))

		req := httptest.NewRequest(fiber.MethodDelete, url, nil)
		resp, err := server.App.Test(req)

		a.Nil(err)
		a.Equal(http.StatusInternalServerError, resp.StatusCode)

		bytesData, err := io.ReadAll(resp.Body)
		a.Nil(err)

		var responseBody common.Response[any]
		err = json.Unmarshal(bytesData, &responseBody)
		a.Nil(err)
		a.False(responseBody.Success)
		a.Equal("unexpected server error", responseBody.Message)
	})
}

func TestControllerDeleteAllByIds(t *testing.T) {
	a := assert.New(t)
	url := "/api/v1/employees"
	validBody := `{"ids":[1,2,3]}`
	invalidBody := `{"ids":[]}`

	t.Run("should delete all employees by ids", func(t *testing.T) {
		server := web.NewServer()
		svc := new(MockService)
		controller := NewController(server, svc)
		controller.RegisterRoutes()

		request := IdsRequest{Ids: []int64{1, 2, 3}}
		svc.On("DeleteAllByIds", request).Return(nil)

		req := httptest.NewRequest(fiber.MethodDelete, url, strings.NewReader(validBody))
		req.Header.Set("Content-Type", "application/json")
		resp, err := server.App.Test(req)

		a.Nil(err)
		a.Equal(http.StatusOK, resp.StatusCode)

		bytesData, err := io.ReadAll(resp.Body)
		a.Nil(err)

		var responseBody common.Response[any]
		err = json.Unmarshal(bytesData, &responseBody)
		a.Nil(err)
		a.True(responseBody.Success)
		a.Empty(responseBody.Message)
		a.Nil(responseBody.Data)
	})

	t.Run("should return validation error", func(t *testing.T) {
		server := web.NewServer()
		svc := new(MockService)
		controller := NewController(server, svc)
		controller.RegisterRoutes()

		request := IdsRequest{Ids: []int64{}}
		validationErr := common.RequestValidationError{Message: "ids must not be empty"}
		svc.On("DeleteAllByIds", request).Return(validationErr)

		req := httptest.NewRequest(fiber.MethodDelete, url, strings.NewReader(invalidBody))
		req.Header.Set("Content-Type", "application/json")
		resp, err := server.App.Test(req)

		a.Nil(err)
		a.Equal(http.StatusBadRequest, resp.StatusCode)

		bytesData, err := io.ReadAll(resp.Body)
		a.Nil(err)

		var responseBody common.Response[any]
		err = json.Unmarshal(bytesData, &responseBody)
		a.Nil(err)
		a.False(responseBody.Success)
		a.Equal("ids must not be empty", responseBody.Message)
	})

	t.Run("should return not found error", func(t *testing.T) {
		server := web.NewServer()
		svc := new(MockService)
		controller := NewController(server, svc)
		controller.RegisterRoutes()

		request := IdsRequest{Ids: []int64{1, 2, 3}}
		notFoundErr := common.NotFoundError{Message: "employees not found"}
		svc.On("DeleteAllByIds", request).Return(notFoundErr)

		req := httptest.NewRequest(fiber.MethodDelete, url, strings.NewReader(validBody))
		req.Header.Set("Content-Type", "application/json")
		resp, err := server.App.Test(req)

		a.Nil(err)
		a.Equal(http.StatusNotFound, resp.StatusCode)

		bytesData, err := io.ReadAll(resp.Body)
		a.Nil(err)

		var responseBody common.Response[any]
		err = json.Unmarshal(bytesData, &responseBody)
		a.Nil(err)
		a.False(responseBody.Success)
		a.Equal("employees not found", responseBody.Message)
	})

	t.Run("should return internal server error on generic error", func(t *testing.T) {
		server := web.NewServer()
		svc := new(MockService)
		controller := NewController(server, svc)
		controller.RegisterRoutes()

		request := IdsRequest{Ids: []int64{1, 2, 3}}
		svc.On("DeleteAllByIds", request).Return(errors.New("unexpected server error"))

		req := httptest.NewRequest(fiber.MethodDelete, url, strings.NewReader(validBody))
		req.Header.Set("Content-Type", "application/json")
		resp, err := server.App.Test(req)

		a.Nil(err)
		a.Equal(http.StatusInternalServerError, resp.StatusCode)

		bytesData, err := io.ReadAll(resp.Body)
		a.Nil(err)

		var responseBody common.Response[any]
		err = json.Unmarshal(bytesData, &responseBody)
		a.Nil(err)
		a.False(responseBody.Success)
		a.Equal("unexpected server error", responseBody.Message)
	})
}
