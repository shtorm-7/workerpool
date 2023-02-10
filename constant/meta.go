package constant

type (
	Tag  string
	Meta interface {
		Name() string
		Tags() []Tag
	}
)
