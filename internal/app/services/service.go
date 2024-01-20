package services

import (
	"context"

	"github.com/IAmFutureHokage/HL-BufferService/internal/app/model"
	pb "github.com/IAmFutureHokage/HL-BufferService/internal/proto"
	"github.com/IAmFutureHokage/HL-BufferService/pkg/decoder"
)

type HydrologyBufferervice struct {
	pb.UnimplementedHydrologyBufferServiceServer
}

func NewHydrologyBufferService() *HydrologyBufferervice {
	return &HydrologyBufferervice{}
}

func (s *HydrologyBufferervice) AddTelegram(ctx context.Context, req *pb.AddTelegramRequest) (*pb.AddTelegramResponse, error) {

	draftTelegrams, err := decoder.FullDecoder(req.Code)
	if err != nil {
		return nil, err
	}

	telegrams := make([]*model.Telegram, len(draftTelegrams))

	for i := 0; i < len(telegrams); i++ {
		telegrams[i] = &model.Telegram{}
		err := telegrams[i].Update(draftTelegrams[i])
		if err != nil {
			return nil, err
		}

	}

	return &pb.AddTelegramResponse{}, nil
}
