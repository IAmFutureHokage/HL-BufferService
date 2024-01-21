package services

import (
	"context"

	"github.com/IAmFutureHokage/HL-BufferService/internal/app/model"
	pb "github.com/IAmFutureHokage/HL-BufferService/internal/proto"
	"github.com/IAmFutureHokage/HL-BufferService/pkg/decoder"
	"github.com/IAmFutureHokage/HL-BufferService/pkg/encoder"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
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
	respose := make([]*pb.Telegram, len(draftTelegrams))
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

		respose[i] = telegramToProto(telegrams[i])
	}

	return &pb.AddTelegramResponse{
		Telegrams: respose,
	}, nil
}

func telegramToProto(req *model.Telegram) (res *pb.Telegram) {
	res = &pb.Telegram{}

	res.Id = req.Id.String()
	res.GroupId = req.GroupId.String()
	res.TelegramCode = req.TelegramCode
	res.PostCode = req.PostCode
	res.Datetime = timestamppb.New(req.DateTime)
	res.IsDangerous = req.IsDangerous

	if req.WaterLevelOnTime.Valid {
		res.WaterLevelOnTime.Value = req.WaterLevelOnTime.Int32
	}
	if req.DeltaWaterLevel.Valid {
		res.DeltaWaterLevel.Value = req.DeltaWaterLevel.Int32
	}
	if req.WaterLevelOn20h.Valid {
		res.WaterLevelOn20H.Value = req.WaterLevelOn20h.Int32
	}
	if req.WaterTemperature.Valid {
		res.WaterTemperature.Value = float32(req.WaterTemperature.Float64)
	}
	if req.AirTemperature.Valid {
		res.AirTemperature.Value = req.AirTemperature.Int32
	}
	if req.IcePhenomeniaState.Valid {
		res.IcePhenomeniaState.Value = int32(req.IcePhenomeniaState.Byte)
	}
	if req.Ice.Valid {
		res.IceHeight.Value = req.Ice.Int32
	}
	if req.Snow.Valid {
		res.SnowHeight.Value = int32(req.Snow.Byte)
	}
	if req.Waterflow.Valid {
		res.WaterFlow.Value = float32(req.Waterflow.Float64)
	}
	if req.PrecipitationValue.Valid {
		res.PrecipitationValue.Value = float32(req.PrecipitationValue.Float64)
	}
	if req.PrecipitationDuration.Valid {
		res.PrecipitationDuration.Value = int32(req.PrecipitationDuration.Byte)
	}
	if req.ReservoirDate.Valid {
		res.ReservoirDate = timestamppb.New(req.ReservoirDate.Time)
	}
	if req.HeadwaterLevel.Valid {
		res.HeadwaterLevel.Value = req.HeadwaterLevel.Int32
	}
	if req.AverageReservoirLevel.Valid {
		res.AverageReservoirLevel.Value = req.AverageReservoirLevel.Int32
	}
	if req.DownstreamLevel.Valid {
		res.DownstreamLevel.Value = req.DownstreamLevel.Int32
	}
	if req.ReservoirVolume.Valid {
		res.ReservoirVolume.Value = float32(req.ReservoirVolume.Float64)
	}
	if req.IsReservoirWaterInflowDate.Valid {
		res.ReservoirWaterInflowDate = timestamppb.New(req.IsReservoirWaterInflowDate.Time)
	}
	if req.Reset.Valid {
		res.Reset_.Value = float32(req.Reset.Float64)
	}

	if len(req.IcePhenomenia) != 0 {
		res.IcePhenomenias = make([]*pb.IcePhenomenia, len(req.IcePhenomenia))

		for i := 0; i < len(res.IcePhenomenias); i++ {
			res.IcePhenomenias[i] = &pb.IcePhenomenia{Phenomen: int32(req.IcePhenomenia[i].Phenomen)}
			if req.IcePhenomenia[i].IsUntensity && req.IcePhenomenia[i].Intensity.Valid {
				res.IcePhenomenias[i].Intensity.Value = int32(req.IcePhenomenia[i].Intensity.Byte)
			}
		}
	}

	return
}
