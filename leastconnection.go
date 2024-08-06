package main

import (
	"fmt"
	"sync"
	"time"
)

// Server represents a server in the load balancer.
type Server struct {
	Address     string
	Status      bool
	Connections int
}

type LoadBalancer struct {
	Servers []*Server
	mutex   sync.Mutex
}

// NewLoadBalancer creates a new load balancer.
func NewLoadBalancer() *LoadBalancer {
	return &LoadBalancer{
		Servers: make([]*Server, 0),
		mutex:   sync.Mutex{},
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

// ChooseServer selects the server with the least active connections.
func (lb *LoadBalancer) ChooseServer() *Server {
	lb.mutex.Lock()
	defer lb.mutex.Unlock()

	var minConnectionsServer *Server
	minConnections := -1

	// Find the server with the least active connections
	for _, server := range lb.Servers {
		if server.Status && (minConnections == -1 || server.Connections < minConnections) {
			minConnections = server.Connections
			minConnectionsServer = server
		}
	}

	if minConnectionsServer != nil {
		minConnectionsServer.Connections++
	}

	return minConnectionsServer
}

// ReleaseConnection releases a connection from a server.
func (lb *LoadBalancer) ReleaseConnection(server *Server) {
	lb.mutex.Lock()
	defer lb.mutex.Unlock()

	// Decrement the active connections count for the server
	if server.Connections > 0 {
		server.Connections--
	}
}

// AddServer adds a new server to the load balancer.
func (lb *LoadBalancer) AddServer(server *Server) {
	lb.mutex.Lock()
	defer lb.mutex.Unlock()

	lb.Servers = append(lb.Servers, server)
}

func main() {
	fmt.Println("LoadBalancer to route request")
	// Create a new load balancer
	lb := NewLoadBalancer()

	// Add servers to the load balancer
	server1 := &Server{Address: "Server1", Status: true, Connections: 4}
	server2 := &Server{Address: "Server2", Status: true, Connections: 3}
	server3 := &Server{Address: "Server3", Status: true, Connections: 5}
	lb.AddServer(server1)
	lb.AddServer(server2)
	lb.AddServer(server3)

	// Start health checks in a goroutine
	go lb.HealthCheck()

	// Simulate incoming requests
	for i := 0; i < 10; i++ {
		go func(requestNum int) {
			chosenServer := lb.ChooseServer()
			if chosenServer != nil {
				fmt.Printf("Request %d routed to %s\n", requestNum, chosenServer.Address)
				time.Sleep(1 * time.Second) // Simulate processing time
				lb.ReleaseConnection(chosenServer)
			} else {
				fmt.Printf("Request %d: No healthy servers available\n", requestNum)
			}
		}(i)
	}

	// Allow time for health checks and requests to complete
	time.Sleep(10 * time.Second)
}
