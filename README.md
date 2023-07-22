<p align="center">
  <img src="./docs/adaptive-throttling-go-golang.png" width="350" title="Golang Adaptive throttling" alt="Golang Adaptive throttling">

</p>

# Adaptive throttling

Is a library that implements adaptive throttling. It is based on the sre-book + rafaelcapucho.

## Installation

```bash
go get adaptivethrottling
```


## Docs

- https://sre.google/sre-book/handling-overload/
- https://rafaelcapucho.github.io/2016/10/enhance-the-quality-of-your-api-calls-with-client-side-throttling/

## Usage


### Example

```golang
opts := adaptivethrottling.Options{
    HistoryTimeMinute:    2,
    K:                    2,
    UpperLimitToReject:   0.9,
    MaxRequestDurationMs: 300,
}
throttling := adaptivethrottling.New(opts)
func exampleFunc() (interface{}, error) {
	time.Sleep(100 * time.Millisecond)
	return "Result", nil
}
result, err := throttling(exampleFunc)
```

# Prameters

### historyTime

Each client task keeps the following information for the last N minutes of its history. In case of "Out of quota" means time to wait for server recovery.

### k

Clients can continue to issue requests to the backend until requests is K times as large as accepts. Google services and they suggest k = 2

Indicates how many minimum failures will be needed to start counting

### upperLimiteToReject

if the server goes down for more than {historyTime} minutes, the P0 value will stand in 1, rejecting locally every new request to the server, so the client app won1t be able to set up a new conection. As the result of it, the client app will never have another request reaching the server. 0.9 allowing the client to recover even in that worst scenario, when the service is down more than {historyTime} minutes.


# Client-Side Throttling

When a customer is out of quota, a backend task should reject requests quickly with the expectation that returning a "customer is out of quota" error consumes significantly fewer resources than actually processing the request and serving back a correct response. However, this logic doesn't hold true for all services. For example, it's almost equally expensive to reject a request that requires a simple RAM lookup (where the overhead of the request/response protocol handling is significantly larger than the overhead of producing the response) as it is to accept and run that request. And even in the case where rejecting requests saves significant resources, those requests still consume some resources. If the amount of rejected requests is significant, these numbers add up quickly. In such cases, the backend can become overloaded even though the vast majority of its CPU is spent just rejecting requests!

Client-side throttling addresses this problem.106 When a client detects that a significant portion of its recent requests have been rejected due to "out of quota" errors, it starts self-regulating and caps the amount of outgoing traffic it generates. Requests above the cap fail locally without even reaching the network.

We implemented client-side throttling through a technique we call adaptive throttling. Specifically, each client task keeps the following information for the last two minutes of its history:

requests
The number of requests attempted by the application layer(at the client, on top of the adaptive throttling system)
accepts
The number of requests accepted by the backend
Under normal conditions, the two values are equal. As the backend starts rejecting traffic, the number of accepts becomes smaller than the number of requests. Clients can continue to issue requests to the backend until requests is K times as large as accepts. Once that cutoff is reached, the client begins to self-regulate and new requests are rejected locally (i.e., at the client) with the probability calculated in Client request rejection probability.

# Client request rejection probability

As the client itself starts rejecting requests, requests will continue to exceed accepts. While it may seem counterintuitive, given that locally rejected requests aren't actually propagated to the backend, this is the preferred behavior. As the rate at which the application attempts requests to the client grows (relative to the rate at which the backend accepts them), we want to increase the probability of dropping new requests.

For services where the cost of processing a request is very close to the cost of rejecting that request, allowing roughly half of the backend resources to be consumed by rejected requests can be unacceptable. In this case, the solution is simple: modify the accepts multiplier K (e.g., 2) in the client request rejection probability (Client request rejection probability). In this way:

Reducing the multiplier will make adaptive throttling behave more aggressively
Increasing the multiplier will make adaptive throttling behave less aggressively
For example, instead of having the client self-regulate when requests = 2 _ accepts, have it self-regulate when requests = 1.1 _ accepts. Reducing the modifier to 1.1 means only one request will be rejected by the backend for every 10 requests accepted.

We generally prefer the 2x multiplier. By allowing more requests to reach the backend than are expected to actually be allowed, we waste more resources at the backend, but we also speed up the propagation of state from the backend to the clients. For example, if the backend decides to stop rejecting traffic from the client tasks, the delay until all client tasks have detected this change in state is shorter.

We've found adaptive throttling to work well in practice, leading to stable rates of requests overall. Even in large overload situations, backends end up rejecting one request for each request they actually process. One large advantage of this approach is that the decision is made by the client task based entirely on local information and using a relatively simple implementation: there are no additional dependencies or latency penalties.

One additional consideration is that client-side throttling may not work well with clients that only very sporadically send requests to their backends. In this case, the view that each client has of the state of the backend is reduced drastically, and approaches to increment this visibility tend to be expensive.
