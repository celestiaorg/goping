package server

import (
	"fmt"
	"net"

	"go.uber.org/zap"
)

func Serve(serveAddr string, logger *zap.Logger) (func(), error) {
	udpAddr, err := net.ResolveUDPAddr("udp", serveAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve UDP address: %w", err)
	}

	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to listen UDP: %w", err)
	}

	logger.Info("Server listening", zap.String("address", serveAddr))

	buf := make([]byte, 1024)
	cancel := make(chan struct{})

	go func() {
		defer conn.Close()

		for {
			select {
			case <-cancel:
				logger.Info("Server stopped")
				return
			default:
				n, addr, err := conn.ReadFromUDP(buf)
				if err != nil {
					logger.Error("Error reading", zap.Error(err))
					continue
				}

				logger.Info("Received packet", zap.String("from", addr.String()), zap.String("data", string(buf[:n])))

				// Echo back the received packet to the client
				_, err = conn.WriteToUDP(buf[:n], addr)
				if err != nil {
					logger.Error("Error sending response", zap.Error(err))
					continue
				}
			}
		}
	}()

	cancelFunc := func() {
		close(cancel)
	}

	return cancelFunc, nil
}
