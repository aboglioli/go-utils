package validator

import (
	"regexp"

	"github.com/aboglioli/go-utils/errors"
	govalidator "github.com/go-playground/validator/v10"
	"github.com/iancoleman/strcase"
)

var (
	ErrFieldsValidation = errors.Validation.New("invalid_fields")
)

type Validator struct {
	validate *govalidator.Validate
}

func NewValidator() *Validator {
	alphaWithSpacesRE := regexp.MustCompile("^[a-zA-Záéíóú ]*$")
	alphaWithSpaces := func(fl govalidator.FieldLevel) bool {
		str := fl.Field().String()
		if str == "invalid" {
			return false
		}

		return alphaWithSpacesRE.MatchString(str)
	}

	alphaNumWithDashRE := regexp.MustCompile("^[a-zA-Z0-9-]*$")
	alphaNumWithDash := func(fl govalidator.FieldLevel) bool {
		str := fl.Field().String()
		if str == "invalid" {
			return false
		}

		return alphaNumWithDashRE.MatchString(str)
	}

	validate := govalidator.New()
	validate.RegisterValidation("alpha-space", alphaWithSpaces)
	validate.RegisterValidation("alpha-num-dash", alphaNumWithDash)

	return &Validator{
		validate: validate,
	}
}

func (v *Validator) RegisterValidation(tag string, fn func(govalidator.FieldLevel) bool) {
	v.validate.RegisterValidation(tag, fn)
}

func (v *Validator) CheckFields(s interface{}) ([]errors.Field, bool) {
	fields := make([]errors.Field, 0)

	if err := v.validate.Struct(s); err != nil {
		if errs, ok := err.(govalidator.ValidationErrors); ok {
			for _, err := range errs {
				field := errors.Field{
					Field: strcase.ToSnake(err.Field()),
					Code:  err.Tag(),
				}
				fields = append(fields, field)
			}
		}

		return fields, false
	}

	return fields, true
}
