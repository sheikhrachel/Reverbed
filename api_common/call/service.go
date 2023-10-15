package call

type Service interface {
	InfoF(stringRaw string, a ...any)
	TraceF(stringRaw string, a ...any)
	ErrorF(stringRaw string, a ...any) error
}
