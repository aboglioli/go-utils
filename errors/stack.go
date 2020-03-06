package errors

type Stack struct {
	Error
	Stack []Stack `json:"stack,omitempty"`
}

type StackOptions struct {
	ExcludePath     bool
	ExcludeInternal bool
}

var (
	FullStack = &StackOptions{
		ExcludePath:     false,
		ExcludeInternal: false,
	}
	InfoStack = &StackOptions{
		ExcludePath:     true,
		ExcludeInternal: true,
	}
)

func BuildStack(err error, opts *StackOptions) []Stack {
	stack := Stack{}

	switch err := err.(type) {
	case Errors:
		stacks := make([]Stack, 0)
		for _, err := range err {
			stacks = append(stacks, BuildStack(err, opts)...)
		}
		return stacks
	case Error:
		if opts.ExcludeInternal && err.Type == Internal {
			return []Stack{}
		}

		stack.Error = err

		if opts.ExcludePath {
			stack.Path = ""
		}

		if err.Cause != nil {
			causeStack := BuildStack(err.Cause, opts)
			if len(causeStack) > 0 {
				stack.Stack = causeStack
			}
			stack.Cause = nil
		}
	case error:
		if opts.ExcludeInternal {
			return []Stack{}
		}

		stack.Error = Error{
			Type:    Unknown,
			Message: err.Error(),
		}
	}

	return []Stack{stack}
}
