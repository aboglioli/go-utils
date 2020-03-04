package validator

import (
	"testing"

	"github.com/aboglioli/go-utils/errors"
	"github.com/aboglioli/go-utils/types"
	"github.com/stretchr/testify/assert"
)

type data struct {
	Name             *string  `validate:"required,min=4,max=12,alpha-space"`
	Username         string   `validate:"required,min=4,max=6,alpha-num-dash"`
	Password         string   `validate:"required,min=6"`
	Email            string   `validate:"required,email"`
	LongNameField    float64  `validate:"required"`
	PtrLongNameField *float64 `validate:"required"`
}

func TestValidatorCheckFields(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		data      data
		errFields []errors.Field
	}{{
		data{},
		[]errors.Field{{
			Field: "name",
			Code:  "required",
		}, {
			Field: "username",
			Code:  "required",
		}, {
			Field: "password",
			Code:  "required",
		}, {
			Field: "email",
			Code:  "required",
		}, {
			Field: "long_name_field",
			Code:  "required",
		}, {
			Field: "ptr_long_name_field",
			Code:  "required",
		}},
	}, {
		data{
			Name:             types.NewString(""),
			Username:         "aa",
			Password:         "aa",
			Email:            "not-an-email.com",
			LongNameField:    0,
			PtrLongNameField: types.NewFloat(0),
		},
		[]errors.Field{{
			Field: "name",
			Code:  "min",
		}, {
			Field: "username",
			Code:  "min",
		}, {
			Field: "password",
			Code:  "min",
		}, {
			Field: "email",
			Code:  "email",
		}, {
			Field: "long_name_field",
			Code:  "required",
		}},
	}, {
		data{
			Name:             types.NewString("aaaa"),
			Username:         "aaaa",
			Password:         "aaaaaa",
			Email:            "a@a.com",
			LongNameField:    1,
			PtrLongNameField: types.NewFloat(0),
		},
		[]errors.Field{},
	}, {
		data{
			Name:             types.NewString("aaaaaaaaaaaaa"),
			Username:         "aaaaaaa",
			Password:         "aaaaaaaaaaaaaaaaaaaaaaaaa",
			Email:            "asd@qwe.io",
			LongNameField:    1.5,
			PtrLongNameField: types.NewFloat(1.5),
		},
		[]errors.Field{{
			Field: "name",
			Code:  "max",
		}, {
			Field: "username",
			Code:  "max",
		}},
	}, {
		data{
			Name:             types.NewString("Alan 1"),
			Username:         "user 1",
			Password:         "123456",
			Email:            "alan@users.io",
			LongNameField:    1.5,
			PtrLongNameField: types.NewFloat(1.5),
		},
		[]errors.Field{{
			Field: "name",
			Code:  "alpha-space",
		}, {
			Field: "username",
			Code:  "alpha-num-dash",
		}},
	}, {
		data{
			Name:             types.NewString("Alan B"),
			Username:         "user-1",
			Password:         "123456",
			Email:            "alan@users.io",
			LongNameField:    1.5,
			PtrLongNameField: types.NewFloat(1.5),
		},
		[]errors.Field{},
	}, {
		data{
			Name:             types.NewString("Iván B"),
			Username:         "iván-1",
			Password:         "123456",
			Email:            "alan@users.io",
			LongNameField:    1.5,
			PtrLongNameField: types.NewFloat(1.5),
		},
		[]errors.Field{{
			Field: "username",
			Code:  "alpha-num-dash",
		}},
	}}

	for i, test := range tests {
		v := NewValidator()
		errFields, ok := v.CheckFields(test.data)

		if len(test.errFields) > 0 {
			assert.False(ok, i)
			assert.Equal(test.errFields, errFields, i)
		} else {
			assert.True(ok, i)
			assert.Equal([]errors.Field{}, errFields, i)
		}
	}
}
