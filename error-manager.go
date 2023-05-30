package framework

import (
	"errors"

	"github.com/gofiber/fiber/v2"
)

type Error interface {
	Status(status int) Error
	Check(err error)
	Throw(message string)
	Internal() Error
	BadRequest() Error
	Forbidden() Error
	Unauthorized() Error
	NotFound() Error
}

type errorManager struct {
	control *control
	error   *managedError
}

func newErrorManager(control *control) *errorManager {
	return &errorManager{
		control: control,
		error:   new(managedError),
	}
}

func (m *errorManager) Status(status int) Error {
	m.error.status = status
	return m
}

func (m *errorManager) Check(err error) {
	if err != nil {
		panic(
			managedError{
				control: m.control,
				error:   err,
				message: m.error.message,
				status:  m.error.status,
			},
		)
	}
}

func (m *errorManager) Throw(message string) {
	panic(
		managedError{
			control: m.control,
			error:   errors.New(message),
			message: m.error.message,
			status:  m.error.status,
		},
	)
}

func (m *errorManager) Internal() Error {
	m.error.status = fiber.StatusInternalServerError
	m.error.message = "internal"
	return m
}

func (m *errorManager) BadRequest() Error {
	m.error.status = fiber.StatusBadRequest
	m.error.message = "bad-request"
	return m
}

func (m *errorManager) Forbidden() Error {
	m.error.status = fiber.StatusForbidden
	m.error.message = "forbidden"
	return m
}

func (m *errorManager) Unauthorized() Error {
	m.error.status = fiber.StatusUnauthorized
	m.error.message = "unauthorized"
	return m
}

func (m *errorManager) NotFound() Error {
	m.error.status = fiber.StatusNotFound
	m.error.message = "not-found"
	return m
}
