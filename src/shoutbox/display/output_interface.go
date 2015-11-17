package display

type Output interface {
	Buffer() *Buffer
	Flush()
	Clear()
}
