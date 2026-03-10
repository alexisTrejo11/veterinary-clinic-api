package employees

import (
	"clinic-vet-api/internal/shared/errors"
	"context"
	"fmt"
)

func InvalidEnumParserError(enumName, enumValue string) error {
	return errors.InvalidEnumValue(context.Background(), enumValue, enumValue, fmt.Sprintf("invalid %s value: %s", enumName, enumValue), "enum parse")

}
