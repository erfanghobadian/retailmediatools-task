package validator

import (
	"github.com/go-playground/validator/v10"
	"sweng-task/internal/utils"
)

var validate = validator.New()

func ValidateStruct(s interface{}) (*utils.FieldError, error) {
	err := validate.Struct(s)
	if err == nil {
		return nil, nil
	}

	if errs, ok := err.(validator.ValidationErrors); ok && len(errs) > 0 {
		ve := errs[0]

		reason := ""
		switch ve.Tag() {
		case "required":
			reason = "is required"
		case "oneof":
			reason = "must be one of: " + ve.Param()
		case "min":
			reason = "must be at least " + ve.Param()
		case "max":
			reason = "must be at most " + ve.Param()
		default:
			reason = "is invalid"
		}

		return &utils.FieldError{
			Field:  ve.Field(),
			Reason: reason,
		}, nil
	}

	return nil, err
}

type AdQueryParams struct {
	Placement string `query:"placement" validate:"required"`
	Category  string `query:"category"`
	Keyword   string `query:"keyword"`
	Limit     int    `query:"limit" validate:"omitempty,min=1,max=10"`
}

type IDParam struct {
	ID string `params:"id" validate:"required"`
}

type LineItemQueryParams struct {
	AdvertiserID string `query:"advertiser_id"`
	Placement    string `query:"placement"`
}
