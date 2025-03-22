package grpc

import (
	"context"
	"errors"
	"github.com/ashkan-maleki/kiwi-eats/pkg/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc/metadata"
)

const (
	HeaderAuthorization = "authorization"
	HeaderUserAgent     = "user-agent"
	HeaderXForwardedFor = "x-forwarded-for"
	HeaderXRealIP       = "x-real-ip"
)

type MD struct {
	ipAddress string
	userAgent string
	token     string
}

func (md *MD) Token() string {
	return md.token
}

func (md *MD) UserAgent() string {
	return md.userAgent
}

func (md *MD) IpAddress() string {
	return md.ipAddress
}

func Metadata(ctx context.Context, extractorFunc ...func(meta *MD, md metadata.MD) error) (*MD, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		logger.Logger.Error("No metadata found in context")
		return nil, errors.New("no metadata")
	}
	meta := &MD{}

	for _, fn := range extractorFunc {
		err := fn(meta, md)
		if err != nil {
			logger.Logger.Error("Failed to extract metadata", zap.Error(err))
			return nil, err
		}
	}

	return meta, nil
}

func IP(meta *MD, md metadata.MD) error {
	ipAddress := getFirstValue(md, HeaderXForwardedFor)
	if ipAddress == "" {
		ipAddress = getFirstValue(md, HeaderXRealIP)
	}

	if ipAddress == "" {
		return errors.New("no IP address found in metadata")
	}

	meta.ipAddress = ipAddress
	return nil
}

func UserAgent(meta *MD, md metadata.MD) error {
	userAgent := getFirstValue(md, HeaderUserAgent)
	if userAgent == "" {
		return errors.New("no user-agent found in metadata")
	}

	meta.userAgent = userAgent
	return nil
}

func Token(meta *MD, md metadata.MD) error {
	token := getFirstValue(md, HeaderAuthorization)
	if token == "" {
		return errors.New("no authorization token found in metadata")
	}

	meta.token = token
	return nil
}

func Invoke2(ctx context.Context) {
	meta, err := Metadata(ctx, IP, UserAgent, Token)
	if err != nil {
		logger.Logger.Error("Failed to extract metadata", zap.Error(err))
		return
	}

	token := meta.Token()
	ipAddress := meta.IpAddress()
	userAgent := meta.UserAgent()
	logger.Logger.Info("Extracted metadata",
		zap.String("ip_address", ipAddress),
		zap.String("user_agent", userAgent),
		zap.String("token", token),
	)
}

func getFirstValue(md metadata.MD, key string) string {
	if values := md[key]; len(values) > 0 {
		return values[0]
	}
	return ""
}
