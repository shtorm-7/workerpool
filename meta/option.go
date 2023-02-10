package meta

import C "github.com/shtorm-7/workerpool/constant"

type MetaOption func(meta *Meta)

func WithName(name string) MetaOption {
	return func(meta *Meta) {
		meta.name = name
	}
}

func WithTags(tags ...C.Tag) MetaOption {
	return func(meta *Meta) {
		meta.tags = tags
	}
}
