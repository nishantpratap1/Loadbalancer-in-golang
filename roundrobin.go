package main

import (
	"fmt"
	"sync"
	"time"
)

// Server represents a server in the load balancer.
type Server struct {
	Address string
	Status  bool
}
type LoadBalancer struct {
	Servers []*Server
	mutex   sync.Mutex
	index   int // To keep track of the next server to choose
}

func NewLoadBalancer() *LoadBalancer {
	return &LoadBalancer{
		Servers: make([]*Server, 0),
		mutex:   sync.Mutex{},
		index:   0,
	}
}

// HealthCheck periodically checks the health of each server.
func (lb *LoadBalancer) HealthCheck() {
	for {
		time.Sleep(5 * time.Second)

		lb.mutex.Lock()
		for _, server := range lb.Servers {
			// Implement your health check logic here
			// For simplicity, we assume all servers are healthy
			server.Status = true
		}
		lb.mutex.Unlock()
	}
}

// ChooseServer selects the next server in the round-robin sequence.
func (lb *LoadBalancer) ChooseServer() *Server {
	lb.mutex.Lock()
	defer lb.mutex.Unlock()

	numServers := len(lb.Servers)
	if numServers == 0 {
		return nil
	}

	// Find the next healthy server
	for i := 0; i < numServers; i++ {
		server := lb.Servers[lb.index]
		lb.index = (lb.index + 1) % numServers
		if server.Status {
			return server
		}
	}

	return nil
}

// AddServer adds a new server to the load balancer.
func (lb *LoadBalancer) AddServer(server *Server) {
	lb.mutex.Lock()
	defer lb.mutex.Unlock()

	lb.Servers = append(lb.Servers, server)
}

// RemoveServer removes a server from the load balancer.
func (lb *LoadBalancer) RemoveServer(address string) {
	lb.mutex.Lock()
	defer lb.mutex.Unlock()

	for i, server := range lb.Servers {
		if server.Address == address {
			// Remove the server from the slice
			lb.Servers = append(lb.Servers[:i], lb.Servers[i+1:]...)
			return
		}
	}
}

func main() {
	fmt.Println("round-robin loadbalancer")

}
