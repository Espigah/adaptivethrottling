package adaptivethrottling

import (
	"math"
	"math/rand"
	"time"
)

type adaptiveThrottlingHistory struct {
	requestsHistory []int64
	acceptsHistory  []int64
}

// Função utilitária para filtrar registros mais antigos do histórico
func filterOldestHistory(historyTimeMinute int, value int64) bool {
	date := time.Now().Unix() - int64(historyTimeMinute*60)
	return value > date
}

// Cria uma nova instância do adaptiveThrottlingHistory
func createAdaptiveThrottlingHistory(historyTimeMinute int) *adaptiveThrottlingHistory {
	return &adaptiveThrottlingHistory{
		requestsHistory: []int64{},
		acceptsHistory:  []int64{},
	}
}

// Adiciona uma nova requisição ao histórico
func (h *adaptiveThrottlingHistory) addRequests() {
	h.requestsHistory = append(h.requestsHistory, time.Now().Unix())
}

// Adiciona uma nova requisição aceita ao histórico
func (h *adaptiveThrottlingHistory) addAccepts() {
	h.acceptsHistory = append(h.acceptsHistory, time.Now().Unix())
}

// filter é uma função auxiliar para filtrar elementos de um slice baseado em um predicado
func filter(slice []int64, predicate func(value int64) bool) []int64 {
	result := make([]int64, 0, len(slice))
	for _, value := range slice {
		if predicate(value) {
			result = append(result, value)
		}
	}
	return result
}

// Atualiza o histórico removendo registros mais antigos
func (h *adaptiveThrottlingHistory) refresh(historyTimeMinute int) {
	filterFn := func(value int64) bool {
		return filterOldestHistory(historyTimeMinute, value)
	}

	h.requestsHistory = filter(h.requestsHistory, filterFn)
	h.acceptsHistory = filter(h.acceptsHistory, filterFn)
}

// Retorna o comprimento do histórico de requisições
func (h *adaptiveThrottlingHistory) getRequestsHistoryLength() int {
	return len(h.requestsHistory)
}

// Retorna o comprimento do histórico de requisições aceitas
func (h *adaptiveThrottlingHistory) getAcceptsHistoryLength() int {
	return len(h.acceptsHistory)
}

// New é uma função que cria uma instância do adaptive throttling com base nas opções fornecidas.
func New(opts Options) func(func() (interface{}, error)) (interface{}, error) {
	opts.Fill()
	requestRejectionProbability := 0.0
	adaptiveThrottling := createAdaptiveThrottlingHistory(opts.HistoryTimeMinute)

	// Função utilitária para verificar a probabilidade de rejeição de uma nova requisição
	checkRequestRejectionProbability := func() bool {
		return rand.Float64() < requestRejectionProbability
	}

	// Função utilitária para atualizar a probabilidade de rejeição de requisições
	updateRequestRejectionProbability := func() {
		adaptiveThrottling.refresh(opts.HistoryTimeMinute)

		requests := adaptiveThrottling.getRequestsHistoryLength()
		accepts := adaptiveThrottling.getAcceptsHistoryLength()

		p0 := math.Max(0, (float64(requests)-opts.K*float64(accepts))/(float64(requests)+1))
		p1 := math.Min(p0, opts.UpperLimitToReject)

		requestRejectionProbability = p1
	}

	return func(fn func() (interface{}, error)) (interface{}, error) {
		adaptiveThrottling.addRequests()

		if checkRequestRejectionProbability() {
			updateRequestRejectionProbability()
			return nil, ThrottledException{}
		}

		startTime := time.Now()

		result, err := fn()
		if err != nil {
			updateRequestRejectionProbability()
			return nil, err
		}

		duration := time.Since(startTime)

		if duration.Seconds() < 1 {
			adaptiveThrottling.addAccepts()
		}

		updateRequestRejectionProbability()

		return result, nil
	}
}

func main() {

}
