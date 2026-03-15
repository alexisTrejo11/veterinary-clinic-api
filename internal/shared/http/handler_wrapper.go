package http

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// BusinessLogic is the core handler logic: bind request T, return result R or error.
type BusinessLogic[T any, R any] func(ctx *gin.Context, req T) (R, error)

// BusinessLogicNoBody is handler logic for requests without a body (e.g. GET, DELETE).
type BusinessLogicNoBody[R any] func(ctx *gin.Context) (R, error)

// Responder writes the successful result to the response.
type Responder[R any] func(c *gin.Context, res R)

// HandleRequest runs logic and responds with Success(c, res, "Operation successful").
func HandleRequest[T any, R any](
	v *validator.Validate,
	logic BusinessLogic[T, R],
) gin.HandlerFunc {
	return HandleRequestWithResponder(v, logic, func(c *gin.Context, res R) {
		Success(c, res, "Operation successful")
	})
}

// HandleRequestWithResponder runs logic and uses the provided responder for success.
func HandleRequestWithResponder[T any, R any](
	v *validator.Validate,
	logic BusinessLogic[T, R],
	respond Responder[R],
) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req T
		if any(req) != nil {
			if err := ShouldBindAndValidateBody(c, &req, v); err != nil {
				BadRequest(c, err)
				return
			}
		}
		res, err := logic(c, req)
		if err != nil {
			ApplicationError(c, err)
			return
		}
		respond(c, res)
	}
}

// HandleCreateRequest binds body, runs logic, then responds with Created(c, entityID, entityName).
func HandleCreateRequest[T any, R any](
	v *validator.Validate,
	entityName string,
	logic BusinessLogic[T, R],
) gin.HandlerFunc {
	return HandleRequestWithResponder(v, logic, func(c *gin.Context, res R) {
		Created(c, res, entityName)
	})
}

// HandleUpdateRequest binds body, runs logic, then responds with Updated(c, res, entityName).
func HandleUpdateRequest[T any, R any](
	v *validator.Validate,
	entityName string,
	logic BusinessLogic[T, R],
) gin.HandlerFunc {
	return HandleRequestWithResponder(v, logic, func(c *gin.Context, res R) {
		Updated(c, res, entityName)
	})
}

// HandleRequestNoBodyStruct runs logic (no body) and responds with Success(c, res, "Operation successful").
func HandleRequestNoBodyStruct[R any](
	v *validator.Validate,
	logic BusinessLogicNoBody[R],
) gin.HandlerFunc {
	return HandleRequestNoBodyWithResponder(v, logic, func(c *gin.Context, res R) {
		Success(c, res, "Operation successful")
	})
}

// HandleRequestNoBodyWithResponder runs logic (no body) and uses the provided responder.
func HandleRequestNoBodyWithResponder[R any](
	v *validator.Validate,
	logic BusinessLogicNoBody[R],
	respond Responder[R],
) gin.HandlerFunc {
	return func(c *gin.Context) {
		res, err := logic(c)
		if err != nil {
			ApplicationError(c, err)
			return
		}
		respond(c, res)
	}
}

// HandleGetRequest runs logic and responds with Found(c, res, entityName).
func HandleGetRequest[R any](
	v *validator.Validate,
	entityName string,
	logic BusinessLogicNoBody[R],
) gin.HandlerFunc {
	return HandleRequestNoBodyWithResponder(v, logic, func(c *gin.Context, res R) {
		Found(c, res, entityName)
	})
}

// HandleDeleteRequest runs logic (no body), then responds with Success(c, nil, entityName+" deleted successfully").
func HandleDeleteRequest(
	v *validator.Validate,
	entityName string,
	logic BusinessLogicNoBody[any],
) gin.HandlerFunc {
	return HandleRequestNoBodyWithResponder(v, logic, func(c *gin.Context, _ any) {
		Success(c, nil, entityName+" deleted successfully")
	})
}
