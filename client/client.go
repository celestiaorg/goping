package client

import (
	"fmt"
	"net"
	"time"

	"go.uber.org/zap"
)

func MeasurePacketloss(address string, numPackets int,
	packetTimeout time.Duration, logger *zap.Logger) (float64, error) {
	serverAddr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return 0, fmt.Errorf("failed to resolve UDP address: %w", err)
	}

	conn, err := net.DialUDP("udp", nil, serverAddr)
	if err != nil {
		return 0, fmt.Errorf("failed to dial UDP: %w", err)
	}
	defer conn.Close()

	var receivedPackets int
	for i := 0; i < numPackets; i++ {
		_, err := conn.Write([]byte("ping"))
		if err != nil {
			logger.Error("Error sending:", zap.Error(err))
			continue
		}

		err = conn.SetReadDeadline(time.Now().Add(packetTimeout))
		if err != nil {
			logger.Error("Failed to set read deadline:", zap.Error(err))
			continue
		}
		buffer := make([]byte, 1024)

		_, _, err = conn.ReadFromUDP(buffer)
		if err != nil {
			// Packet loss occurred (no response received)
			logger.Info(fmt.Sprintf("Packet %d lost", i+1))
			continue
		}

		receivedPackets++
		time.Sleep(10 * time.Millisecond)
		logger.Info(fmt.Sprintf("Packet %d sent successfully", i+1))
	}

	packetLossPercentage := 100 * (1 - float64(receivedPackets)/float64(numPackets))
	return packetLossPercentage, nil
}
