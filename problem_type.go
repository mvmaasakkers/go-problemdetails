package problemdetails

// ProblemType is an interface for problem type definitions. A ProblemType also implements the error interface
type ProblemType interface {
	Error() string
}
