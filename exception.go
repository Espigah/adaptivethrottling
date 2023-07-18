package adaptivethrottling

// ThrottledException é uma estrutura de erro personalizada para representar exceções de throttling.
type ThrottledException struct{}

// Implementação da interface Error para ThrottledException.
func (e ThrottledException) Error() string {
	return "Request throttled"
}
