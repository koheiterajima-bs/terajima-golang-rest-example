package error

type ErrReviewNotFound struct{ Stack error }

func (e ErrReviewNotFound) Error() string {
	return e.Stack.Error()
}
