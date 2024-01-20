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
	response := make([]*pb.Telegram, len(draftTelegrams))
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

		pbTelegram := convertToProto(telegrams[i])

		response[i] = pbTelegram
	}

	return &pb.AddTelegramResponse{
		Telegrams: response,
	}, nil
}

func convertToProto(tg *model.Telegram) *pb.Telegram {
	pbTelegram := &pb.Telegram{
		Id:           tg.Id.String(),
		GroupId:      tg.GroupId.String(),
		TelegramCode: tg.TelegramCode,
		PostCode:     tg.PostCode,
		IsDangerous:  tg.IsDangerous,
	}

	if !tg.DateTime.IsZero() {
		pbTelegram.Datetime = timestamppb.New(tg.DateTime)
	}

	if tg.WaterLevelOnTime.Valid {
		pbTelegram.WaterLevelOnTime = int32(tg.WaterLevelOnTime.Int32)
	}

	if tg.DeltaWaterLevel.Valid {
		pbTelegram.DeltaWaterLevel = int32(tg.DeltaWaterLevel.Int32)
	}

	if tg.WaterLevelOn20h.Valid {
		pbTelegram.WaterLevelOn_20H = int32(tg.WaterLevelOn20h.Int32)
	}

	if tg.WaterTemperature.Valid {
		pbTelegram.WaterTemperature = float32(tg.WaterTemperature.Float64)
	}

	if tg.AirTemperature.Valid {
		pbTelegram.AirTemperature = int32(tg.AirTemperature.Int32)
	}

	if tg.IcePhenomeniaState.Valid {
		pbTelegram.IcePhenomeniaState = pb.IcePhenomeniaState(tg.IcePhenomeniaState.Byte)
	}

	if tg.Ice.Valid {
		pbTelegram.IceHeight = int32(tg.Ice.Int32)
	}

	if tg.Snow.Valid {
		pbTelegram.SnowHeight = pb.SnowHeight(tg.Snow.Byte)
	}

	if tg.Waterflow.Valid {
		pbTelegram.WaterFlow = float32(tg.WaterTemperature.Float64)
	}

	if tg.PrecipitationValue.Valid {
		pbTelegram.PrecipitationValue = float32(tg.PrecipitationValue.Float64)
	}

	if tg.PrecipitationDuration.Valid {
		pbTelegram.PrecipitationDuration = pb.PrecipitationDuration(tg.PrecipitationDuration.Byte)
	}

	if !tg.ReservoirDate.Time.IsZero() {
		pbTelegram.ReservoirDate = timestamppb.New(tg.ReservoirDate.Time)

	}

	if tg.IsReservoirWaterInflowDate.Valid && !tg.IsReservoirWaterInflowDate.Time.IsZero() {
		pbTelegram.ReservoirWaterInflowDate = timestamppb.New(tg.IsReservoirWaterInflowDate.Time)
	}

	if tg.HeadwaterLevel.Valid {
		pbTelegram.ReservoirData.HeadwaterLevel = int32(tg.HeadwaterLevel.Int32)
	}

	if tg.AverageReservoirLevel.Valid {
		pbTelegram.ReservoirData.AverageReservoirLevel = int32(tg.AverageReservoirLevel.Int32)
	}

	if tg.DownstreamLevel.Valid {
		pbTelegram.ReservoirData.DownstreamLevel = int32(tg.DownstreamLevel.Int32)
	}

	if tg.ReservoirVolume.Valid {
		pbTelegram.ReservoirData.ReservoirVolume = float32(tg.ReservoirVolume.Float64)
	}

	if tg.Inflow.Valid {
		pbTelegram.ReservoirWaterInflowData.Inflow = float32(tg.Inflow.Float64)
	}

	if tg.Reset.Valid {
		pbTelegram.ReservoirWaterInflowData.Reset_ = float32(tg.Reset.Float64)
	}

	for _, p := range tg.IcePhenomenia {
		pbPhenomenia := &pb.IcePhenomenia{
			Phenomen: int32(p.Phenomen),
		}

		if p.Intensity.Valid {
			pbPhenomenia.Intensity = int32(p.Intensity.Byte)
		}

		pbTelegram.IcePhenomenias = append(pbTelegram.IcePhenomenias, pbPhenomenia)
	}

	return pbTelegram
}
