package utils

func AnyErrorNotNil(errors ...error) bool {
	for _, v := range errors {
		if v != nil {
			return true
		}
	}
	return false
}

func AllErrorsPresent(errors ...error) bool {
	nilErrorFound := false
	for _, v := range errors {
		if v == nil {
			nilErrorFound = true
		}
	}
	// if no nil error is in the slice then they all aren't errors
	return !nilErrorFound
}
