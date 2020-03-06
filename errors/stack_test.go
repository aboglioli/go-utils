package errors

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildStack(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		err      error
		opts     *StackOptions
		expected []Stack
	}{{
		Internal.New("I").Wrap(errors.New("err")),
		FullStack,
		[]Stack{{
			Error{
				Type: Internal,
				Code: "I",
			},
			[]Stack{{
				Error: Error{
					Type:    Unknown,
					Message: "err",
				},
			}},
		}},
	}, {
		Validation.New("V").S(404).Wrap(
			Internal.New("I").M("not found err db").C("connectionString", "%s:%d", "localhost", 27017).Wrap(
				errors.New("db: not found"),
			),
		),
		FullStack,
		[]Stack{{
			Error: Error{
				Type:   Validation,
				Code:   "V",
				Status: 404,
			},
			Stack: []Stack{{
				Error: Error{
					Type:    Internal,
					Code:    "I",
					Message: "not found err db",
					Context: Context{
						"connectionString": "localhost:27017",
					},
				},
				Stack: []Stack{{
					Error: Error{
						Type:    Unknown,
						Message: "db: not found",
					},
				}},
			}},
		}},
	}, {
		Errors{
			Internal.New("I1").Wrap(errors.New("err1")),
			Internal.New("I2").P("p").M("m").Wrap(errors.New("err2")),
		},
		FullStack,
		[]Stack{{
			Error: Error{
				Type: Internal,
				Code: "I1",
			},
			Stack: []Stack{{
				Error: Error{
					Type:    Unknown,
					Message: "err1",
				},
			}},
		}, {
			Error: Error{
				Type:    Internal,
				Code:    "I2",
				Path:    "p",
				Message: "m",
			},
			Stack: []Stack{{
				Error: Error{
					Type:    Unknown,
					Message: "err2",
				},
			}},
		}},
	}, {
		Status.New("S").Wrap(Errors{Internal.New("I"), errors.New("raw")}),
		FullStack,
		[]Stack{{
			Error: Error{
				Type: Status,
				Code: "S",
			},
			Stack: []Stack{{
				Error: Error{
					Type: Internal,
					Code: "I",
				},
			}, {
				Error: Error{
					Type:    Unknown,
					Message: "raw",
				},
			}},
		}},
	}, {
		Status.New("S").P("p").Wrap(Errors{Internal.New("I").P("p"), Errors{errors.New("raw1"), errors.New("raw2")}}),
		FullStack,
		[]Stack{{
			Error: Error{
				Type: Status,
				Code: "S",
				Path: "p",
			},
			Stack: []Stack{{
				Error: Error{
					Type: Internal,
					Code: "I",
					Path: "p",
				},
			}, {
				Error: Error{
					Type:    Unknown,
					Message: "raw1",
				},
			}, {
				Error: Error{
					Type:    Unknown,
					Message: "raw2",
				},
			}},
		}},
	}, {
		errors.New("raw"),
		FullStack,
		[]Stack{{
			Error: Error{
				Type:    Unknown,
				Message: "raw",
			},
		}},
	}, {
		errors.New("raw"),
		InfoStack,
		[]Stack{},
	}, {
		Validation.New("V").Wrap(Errors{errors.New("raw"), Status.New("S")}),
		InfoStack,
		[]Stack{{
			Error: Error{
				Type: Validation,
				Code: "V",
			},
			Stack: []Stack{{
				Error: Error{
					Type: Status,
					Code: "S",
				},
			}},
		}},
	}, {
		Validation.New("V").S(404).P("p").Wrap(
			Internal.New("I").P("p").M("not found err db").C("connectionString", "%s:%d", "localhost", 27017).Wrap(
				errors.New("db: not found"),
			),
		),
		InfoStack,
		[]Stack{{
			Error: Error{
				Type:   Validation,
				Code:   "V",
				Status: 404,
			},
		}},
	}, {
		Status.New("S").P("p").Wrap(Errors{Internal.New("I"), Errors{errors.New("raw1"), errors.New("raw2")}}),
		InfoStack,
		[]Stack{{
			Error: Error{
				Type: Status,
				Code: "S",
			},
		}},
	}, {
		Status.New("S").P("p").Wrap(Errors{Internal.New("I").P("p").Wrap(errors.New("raw")), Errors{Internal.New("I").Wrap(errors.New("raw1")), errors.New("raw2")}, errors.New("raw3"), Status.New("s_code").Wrap(Validation.New("v_code"))}),
		InfoStack,
		[]Stack{{
			Error: Error{
				Type: Status,
				Code: "S",
			},
			Stack: []Stack{{
				Error: Error{
					Type: Status,
					Code: "s_code",
				},
				Stack: []Stack{{
					Error: Error{
						Type: Validation,
						Code: "v_code",
					},
				}},
			}},
		}},
	}}

	for i, test := range tests {
		stack := BuildStack(test.err, test.opts)
		assert.Equal(test.expected, stack, i)
	}
}
