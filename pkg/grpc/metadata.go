package grpc

import (
	"context"
	"google.golang.org/grpc/metadata"
)

func ExtractIPAndUserAgent(ctx context.Context) (ipAddress, userAgent string) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", ""
	}

	// Extract IP address (e.g., from X-Forwarded-For or X-Real-Ip)
	ipAddress = getFirstValue(md, "x-forwarded-for")
	if ipAddress == "" {
		ipAddress = getFirstValue(md, "x-real-ip")
	}

	// Extract User-Agent
	userAgent = getFirstValue(md, "user-agent")
	return ipAddress, userAgent
}

func getFirstValue(md metadata.MD, key string) string {
	if values := md[key]; len(values) > 0 {
		return values[0]
	}
	return ""
}
