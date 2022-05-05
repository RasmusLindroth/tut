package ui

type Shared struct {
	Top    *Top
	Bottom *Bottom
}

func NewShared(tv *TutView) *Shared {
	return &Shared{
		Top:    NewTop(tv),
		Bottom: NewBottom(tv),
	}
}
