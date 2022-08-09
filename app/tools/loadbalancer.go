package tools

type ILoadBalancer interface {
	GetNextServer() string // Returns an address to the next available server
}

type LoadBalancer struct {
	srvs       []string // Servers to balaance load across
	roundRobin int
}

// Creates a new instance of load balancer
func NewLoadBalancer(srvs []string) ILoadBalancer {
	return &LoadBalancer{
		srvs:       srvs,
		roundRobin: 0,
	}
}

// Retrieve the next server
func (lb *LoadBalancer) GetNextServer() string {
	srv := lb.srvs[lb.roundRobin]
	lb.roundRobin++
	if lb.roundRobin >= len(lb.srvs) {
		lb.roundRobin = 0
	}
	return srv
}
