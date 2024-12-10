package main

import (
	"fmt"
	"os/exec"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var ssh_mon_status = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Name: "ssh_mon_status",
	Help: "SSH Monitor Status and response time",
}, []string{"sshd_host", "status"})

func measureCommandResponseTime(command string, args ...string) (float64, error) {
	startTime := time.Now()

	cmd := exec.Command(command, args...)
	err := cmd.Run()

	endTime := time.Now()
	duration := endTime.Sub(startTime).Seconds()

	return duration, err
}

func check_ssh_status() {

	// username := "vagrant"
	// host_ip := "192.168.16.128"
	username := "testuser"
	host_ip := "8.19.55.99"
	ssh_id := username + "@" + host_ip

	duration, err := measureCommandResponseTime("ssh", "-o", "BatchMode=yes", "-o", "ConnectTimeout=30", ssh_id, "exit")
	if err != nil {
		ssh_mon_status.WithLabelValues(host_ip, "down").Set(-1)
		fmt.Printf("Command failed: %v\n", err)
	} else {
		ssh_mon_status.WithLabelValues(host_ip, "up").Set(duration)
		fmt.Printf("Command response time: %.2f seconds\n", duration)
	}
}
