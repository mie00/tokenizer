package tokenizers

type UnSupportedOperation struct{}

func (UnSupportedOperation) Error() string {
	return "operation unsupported"
}

type InvalidCount struct{}

func (InvalidCount) Error() string {
	return "invalid count"
}
