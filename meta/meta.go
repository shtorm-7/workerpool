package meta

import C "github.com/shtorm-7/workerpool/constant"

type Meta struct {
	name string
	tags []C.Tag
}

func NewMeta(opts ...MetaOption) *Meta {
	meta := new(Meta)
	for _, opt := range opts {
		opt(meta)
	}
	return meta
}

func (m *Meta) Name() string {
	return m.name
}

func (m *Meta) Tags() []C.Tag {
	return m.tags
}
