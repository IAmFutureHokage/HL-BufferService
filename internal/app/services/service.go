package services

import (
	"context"
	"fmt"

	"github.com/IAmFutureHokage/HL-BufferService/internal/app/model"
	pb "github.com/IAmFutureHokage/HL-BufferService/internal/proto"
	"github.com/IAmFutureHokage/HL-BufferService/pkg/decoder"
	"github.com/IAmFutureHokage/HL-BufferService/pkg/encoder"
	"github.com/google/uuid"
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
	groupId := uuid.New()

	for i := 0; i < len(telegrams); i++ {

		codeTg, err := encoder.Encoder(draftTelegrams[i])

		telegrams[i] = &model.Telegram{}
		telegrams[i].Id = uuid.New()
		telegrams[i].GroupId = groupId
		telegrams[i].TelegramCode = codeTg

		if err != nil {
			return nil, err
		}

		err = telegrams[i].Update(draftTelegrams[i])

		if err != nil {
			return nil, err
		}

		fmt.Printf("Telegram %d: %+v\n", i+1, telegrams[i])
	}

	return &pb.AddTelegramResponse{}, nil
}
