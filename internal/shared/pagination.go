package shared

type Specification interface {
	IsSatisfiedBy(any) bool
	ToSQL() (string, []any)
}
