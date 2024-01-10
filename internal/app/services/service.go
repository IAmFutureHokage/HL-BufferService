package services

import (
	"context"

	pb "github.com/IAmFutureHokage/HL-BufferService/internal/proto"
	"github.com/IAmFutureHokage/HL-BufferService/pkg/decoder"
)

type HydrologyBufferervice struct {
	pb.UnimplementedHydrologyBufferServiceServer
}

func (s *HydrologyBufferervice) AddTelegram(ctx context.Context, req *pb.AddTelegramRequest) (*pb.AddTelegramResponse, error) {

	_, err := decoder.FullDecoder(req.Code)
	if err != nil {
		return nil, err
	}

	return &pb.AddTelegramResponse{}, nil
}
