package http

import (
	"clinic-vet-api/internal/shared/errors"
	"clinic-vet-api/internal/shared/page"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// ShouldBindAndValidate binds the request body to the provided object and validates it using the provided validation function.
func ShouldBindAndValidateBody(c *gin.Context, obj any, validate *validator.Validate) error {
	if err := c.ShouldBindBodyWithJSON(obj); err != nil {
		return errors.RequestBodyDataError(err)
	}

	if err := validate.Struct(obj); err != nil {
		validationErrors, ok := err.(validator.ValidationErrors)
		if !ok {
			return errors.InvalidDataError(err)
		}

		return errors.InvalidDataError(validationErrors)
	}

	return nil
}

func ShouldBindAndValidateQuery(c *gin.Context, obj any, validate *validator.Validate) error {
	if err := c.ShouldBindQuery(obj); err != nil {
		return errors.InvalidDataError(err)
	}

	if err := validate.Struct(obj); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return errors.InvalidDataError(validationErrors)
	}

	return nil
}

func ShouldBindPageParams(requestPageParams *page.PaginationRequest, ctx *gin.Context, validator *validator.Validate) error {
	if err := ctx.ShouldBindQuery(requestPageParams); err != nil {
		return errors.RequestURLQueryError(err, ctx.Request.URL.RawQuery)
	}

	fmt.Println(requestPageParams)
	*requestPageParams = requestPageParams.WithDefaults()

	if err := validator.Struct(requestPageParams); err != nil {
		return errors.InvalidDataError(err)
	}

	return nil
}

func ParseParamToUInt(c *gin.Context, paramName string) (uint, error) {
	idStr := c.Param(paramName)
	if idStr == "" {
		return 0, fmt.Errorf("empty id")
	}

	intValue, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, fmt.Errorf("id must be a valid number")
	}

	if intValue < 0 {
		return 0, fmt.Errorf("ID cannot be negative")
	}

	return uint(intValue), nil
}

type LoginMetadata struct {
	IP         string
	UserAgent  string
	DeviceInfo string
	Timestamp  int64
}

func NewLoginMetadata(c *gin.Context) *LoginMetadata {
	ip := extractIP(c)
	userAgent := c.Request.UserAgent()

	return &LoginMetadata{
		IP:         ip,
		UserAgent:  userAgent,
		DeviceInfo: extractDeviceInfo(userAgent),
		Timestamp:  time.Now().Unix(),
	}
}

func extractIP(c *gin.Context) string {
	if ip := c.GetHeader("X-Forwarded-For"); ip != "" {
		return strings.Split(ip, ",")[0]
	}
	if ip := c.GetHeader("X-Real-IP"); ip != "" {
		return ip
	}

	return c.ClientIP()
}

func extractDeviceInfo(userAgent string) string {
	ua := strings.ToLower(userAgent)

	switch {
	case strings.Contains(ua, "mobile"):
		return "mobile"
	case strings.Contains(ua, "tablet"):
		return "tablet"
	default:
		return "desktop"
	}
}

func (lm *LoginMetadata) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"ip":         lm.IP,
		"user_agent": lm.UserAgent,
		"device":     lm.DeviceInfo,
		"timestamp":  lm.Timestamp,
	}
}
