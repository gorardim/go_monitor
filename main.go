package main

import (
	h5 "monitor/h5_health"
	"monitor/navigation"
	// healthCheck "monitor/server_health"
	// serviceCheck "monitor/service_health"
)

func main() {
	// healthCheck.CheckServerHealth()
	// serviceCheck.GetServiceHealth()
	h5.GetH5Health()
	navigation.GetNavigationHealth()
	// var wg sync.WaitGroup
	// wg.Add(1)
	// CronStart()
	// wg.Wait()
}

func CronStart() {
	// do something
}
