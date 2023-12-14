package client

import (
	"fmt"
	"net"
	"time"

	"go.uber.org/zap"
)

// MeasurePacketloss measures the packet loss percentage for a given address
// by sending numPackets packets and waiting for a response for each packet
// for packetTimeout duration.
// Returns the packet loss percentage and any error that occurred.
func MeasurePacketloss(address string, numPackets int, packetTimeout time.Duration, logger *zap.Logger) (float64, error) {
	metricFunc := func(packetNum int, startTime time.Time) {
		logger.Info(fmt.Sprintf("Packet %d sent successfully", packetNum))
	}

	return MeasurePacketMetric(address, numPackets, packetTimeout, logger, metricFunc)
}

// MeasureLatency measures the average latency for a given address
// by sending numPackets packets and waiting for a response for each packet
// for packetTimeout duration.
// Returns the average latency and any error that occurred.
func MeasureLatency(address string, numPackets int, packetTimeout time.Duration, logger *zap.Logger) (time.Duration, error) {
	var totalLatency time.Duration

	metricFunc := func(packetNum int, startTime time.Time) {
		elapsed := time.Since(startTime)
		logger.Info(fmt.Sprintf("Packet %d received in %s", packetNum, elapsed.String()))
		totalLatency += elapsed
	}

	_, err := MeasurePacketMetric(address, numPackets, packetTimeout, logger, metricFunc)
	if err != nil {
		return 0, err
	}

	averageLatency := totalLatency / time.Duration(numPackets)
	return averageLatency, nil
}

func MeasureJitter(address string, numPackets int, packetTimeout time.Duration, logger *zap.Logger) (time.Duration, error) {
	var totalJitter time.Duration
	prevLatency := time.Duration(0)

	for i := 0; i < numPackets; i++ {
		latency, err := MeasureLatency(address, 1, packetTimeout, logger)
		if err != nil {
			logger.Error("Error measuring latency:", zap.Error(err))
			continue
		}

		if prevLatency != 0 {
			jitter := (latency - prevLatency).Abs()
			totalJitter += jitter
		}

		prevLatency = latency
	}

	averageJitter := totalJitter / time.Duration(numPackets-1)
	return averageJitter, nil
}

// MeasurePacketMetric measures a custom packet metric for a given address
// by sending numPackets packets and waiting for a response for each packet
// for packetTimeout duration.
func MeasurePacketMetric(address string, numPackets int, packetTimeout time.Duration, logger *zap.Logger, metricFunc func(int, time.Time)) (float64, error) {
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
		startTime := time.Now()

		if err := sendMessage(conn, "ping", packetTimeout, logger); err != nil {
			continue
		}

		buffer := make([]byte, 1024)

		_, _, err = conn.ReadFromUDP(buffer)
		if err != nil {
			logger.Info(fmt.Sprintf("Packet %d lost", i+1))
			continue
		}

		receivedPackets++
		metricFunc(i+1, startTime)
		time.Sleep(10 * time.Millisecond)
	}

	packetLossPercentage := 100 * (1 - float64(receivedPackets)/float64(numPackets))
	return packetLossPercentage, nil
}

func sendMessage(conn *net.UDPConn, msg string, packetTimeout time.Duration, logger *zap.Logger) error {
	_, err := conn.Write([]byte(msg))
	if err != nil {
		logger.Error("Error sending:", zap.Error(err))
		return err
	}

	err = conn.SetReadDeadline(time.Now().Add(packetTimeout))
	if err != nil {
		logger.Error("Failed to set read deadline:", zap.Error(err))
		return err
	}

	return nil
}
