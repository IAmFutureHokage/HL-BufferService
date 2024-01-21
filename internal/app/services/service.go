package services

import (
	"context"

	"github.com/IAmFutureHokage/HL-BufferService/internal/app/model"
	"github.com/IAmFutureHokage/HL-BufferService/internal/app/repository"
	pb "github.com/IAmFutureHokage/HL-BufferService/internal/proto"
	"github.com/IAmFutureHokage/HL-BufferService/pkg/decoder"
	"github.com/IAmFutureHokage/HL-BufferService/pkg/encoder"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type HydrologyBufferervice struct {
	pb.UnimplementedHydrologyBufferServiceServer
	repository *repository.HydrologyBufferRepository
}

func NewHydrologyBufferService(repo *repository.HydrologyBufferRepository) *HydrologyBufferervice {
	return &HydrologyBufferervice{repository: repo}
}

func (s *HydrologyBufferervice) AddTelegram(ctx context.Context, req *pb.AddTelegramRequest) (*pb.AddTelegramResponse, error) {

	draftTelegrams, err := decoder.FullDecoder(req.Code)
	if err != nil {
		return nil, err
	}

	telegrams := make([]model.Telegram, len(draftTelegrams))
	respose := make([]*pb.Telegram, len(draftTelegrams))
	groupId := uuid.New()

	for i := 0; i < len(telegrams); i++ {

		codeTg, err := encoder.Encoder(draftTelegrams[i])

		telegrams[i] = model.Telegram{}
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

		respose[i] = telegramToProto(&telegrams[i])
	}

	if err := s.repository.AddTelegram(ctx, telegrams); err != nil {
		return nil, err
	}

	return &pb.AddTelegramResponse{
		Telegrams: respose,
	}, nil
}

func (s *HydrologyBufferervice) RemoveTelegrams(ctx context.Context, req *pb.RemoveTelegramsRequest) (*pb.RemoveTelegramsResponse, error) {

	uuids := make([]uuid.UUID, len(req.Id))

	for i := 0; i < len(uuids); i++ {
		id, err := uuid.Parse(req.Id[i])
		if err != nil {
			return nil, err
		}
		uuids[i] = id
	}

	err := s.repository.RemoveTelegrams(ctx, uuids)
	if err != nil {
		return nil, err
	}

	return &pb.RemoveTelegramsResponse{Success: true}, nil
}

func (s *HydrologyBufferervice) UpdateTelegramByCode(ctx context.Context, req *pb.UpdateTelegramByCodeRequest) (*pb.UpdateTelegramResponse, error) {

	draftTelegram, err := decoder.Decoder(req.TelegramCode)
	if err != nil {
		return nil, err
	}

	telegramCode, err := encoder.Encoder(draftTelegram)
	if err != nil {
		return nil, err
	}

	telegramId, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, err
	}

	telegram, err := s.repository.GetTelegramByID(ctx, telegramId)
	if err != nil {
		return nil, err
	}

	telegram.Update(draftTelegram)
	telegram.TelegramCode = telegramCode

	err = s.repository.UpdateTelegram(ctx, telegram)
	if err != nil {
		return nil, err
	}

	response := telegramToProto(&telegram)

	return &pb.UpdateTelegramResponse{
		Telegram: response,
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
		res.WaterLevelOnTime = &wrapperspb.Int32Value{
			Value: req.WaterLevelOnTime.Int32,
		}
	}
	if req.DeltaWaterLevel.Valid {
		res.DeltaWaterLevel = &wrapperspb.Int32Value{
			Value: req.DeltaWaterLevel.Int32,
		}
	}
	if req.WaterLevelOn20h.Valid {
		res.WaterLevelOn20H = &wrapperspb.Int32Value{
			Value: req.WaterLevelOn20h.Int32,
		}
	}
	if req.WaterTemperature.Valid {
		res.WaterTemperature = &wrapperspb.DoubleValue{
			Value: req.WaterTemperature.Float64,
		}
	}
	if req.AirTemperature.Valid {
		res.AirTemperature = &wrapperspb.Int32Value{
			Value: req.AirTemperature.Int32,
		}
	}
	if req.IcePhenomeniaState.Valid {
		res.IcePhenomeniaState = &wrapperspb.Int32Value{
			Value: int32(req.IcePhenomeniaState.Byte),
		}
	}
	if req.Ice.Valid {
		res.IceHeight = &wrapperspb.Int32Value{
			Value: req.Ice.Int32,
		}
	}
	if req.Snow.Valid {
		res.SnowHeight = &wrapperspb.Int32Value{
			Value: int32(req.Snow.Byte),
		}
	}
	if req.Waterflow.Valid {
		res.WaterFlow = &wrapperspb.DoubleValue{
			Value: req.Waterflow.Float64,
		}
	}
	if req.PrecipitationValue.Valid {
		res.PrecipitationValue = &wrapperspb.DoubleValue{
			Value: req.PrecipitationValue.Float64,
		}
	}
	if req.PrecipitationDuration.Valid {
		res.PrecipitationDuration = &wrapperspb.Int32Value{
			Value: int32(req.PrecipitationDuration.Byte),
		}
	}
	if req.ReservoirDate.Valid {
		res.ReservoirDate = timestamppb.New(req.ReservoirDate.Time)
	}
	if req.HeadwaterLevel.Valid {
		res.HeadwaterLevel = &wrapperspb.Int32Value{
			Value: req.HeadwaterLevel.Int32,
		}
	}
	if req.AverageReservoirLevel.Valid {
		res.AverageReservoirLevel = &wrapperspb.Int32Value{
			Value: req.AverageReservoirLevel.Int32,
		}
	}
	if req.DownstreamLevel.Valid {
		res.DownstreamLevel = &wrapperspb.Int32Value{
			Value: req.DownstreamLevel.Int32,
		}
	}
	if req.ReservoirVolume.Valid {
		res.ReservoirVolume = &wrapperspb.DoubleValue{
			Value: req.ReservoirVolume.Float64,
		}
	}
	if req.IsReservoirWaterInflowDate.Valid {
		res.ReservoirWaterInflowDate = timestamppb.New(req.IsReservoirWaterInflowDate.Time)
	}
	if req.Reset.Valid {
		res.Reset_ = &wrapperspb.DoubleValue{
			Value: req.Reset.Float64,
		}
	}

	if len(req.IcePhenomenia) != 0 {
		res.IcePhenomenias = make([]*pb.IcePhenomenia, len(req.IcePhenomenia))

		for i := 0; i < len(res.IcePhenomenias); i++ {
			res.IcePhenomenias[i] = &pb.IcePhenomenia{Phenomen: int32(req.IcePhenomenia[i].Phenomen)}

			if req.IcePhenomenia[i].IsUntensity && req.IcePhenomenia[i].Intensity.Valid {
				res.IcePhenomenias[i].Intensity = &wrapperspb.Int32Value{
					Value: int32(req.IcePhenomenia[i].Intensity.Byte),
				}
			}
		}
	}

	return
}
