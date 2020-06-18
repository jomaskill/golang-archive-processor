package domain

type Base interface {
	Data() Base
	Index() string
}
