package services

import (
	"context"
	"sync"
	"time"

	"github.com/IAmFutureHokage/HL-BufferService/internal/app/model"
	"github.com/IAmFutureHokage/HL-BufferService/internal/app/services/kafka_dto"
	pb "github.com/IAmFutureHokage/HL-BufferService/internal/proto"
	"github.com/IAmFutureHokage/HL-BufferService/pkg/decoder"
	"github.com/IAmFutureHokage/HL-BufferService/pkg/decoder/encoder"
	decoder_types "github.com/IAmFutureHokage/HL-BufferService/pkg/decoder/types"
	"github.com/IAmFutureHokage/HL-BufferService/pkg/kafka"
	"github.com/Shopify/sarama"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type Strorage interface {
	AddTelegram(ctx context.Context, data []model.Telegram) error
	GetTelegramByID(ctx context.Context, id uuid.UUID) (*model.Telegram, error)
	RemoveTelegrams(ctx context.Context, ids []uuid.UUID) error
	GetAll(ctx context.Context) (*[]model.Telegram, error)
	UpdateTelegram(ctx context.Context, updatedTelegram *model.Telegram) error
	GetTelegramsById(ctx context.Context, ids []uuid.UUID) (*[]model.Telegram, error)
}

type HydrologyBufferervice struct {
	pb.UnimplementedHydrologyBufferServiceServer
	storage       Strorage
	KafkaProducer sarama.SyncProducer
	KafkaConfig   kafka.KafkaConfig
}

func NewHydrologyBufferService(storage Strorage, kafkaProducer sarama.SyncProducer) *HydrologyBufferervice {
	return &HydrologyBufferervice{
		storage:       storage,
		KafkaProducer: kafkaProducer,
	}
}

func (s *HydrologyBufferervice) AddTelegram(ctx context.Context, req *pb.AddTelegramRequest) (*pb.AddTelegramResponse, error) {

	draftTelegrams, err := decoder.NewtelegramsSlice(req.Code)
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

	if err := s.storage.AddTelegram(ctx, telegrams); err != nil {
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

	err := s.storage.RemoveTelegrams(ctx, uuids)
	if err != nil {
		return nil, err
	}

	return &pb.RemoveTelegramsResponse{Success: true}, nil
}

func (s *HydrologyBufferervice) UpdateTelegramByInfo(ctx context.Context, req *pb.UpdateTelegramByInfoRequest) (*pb.UpdateTelegramResponse, error) {

	telegramId, err := uuid.Parse(req.Telegram.Id)
	if err != nil {
		return nil, err
	}

	telegram, err := s.storage.GetTelegramByID(ctx, telegramId)
	if err != nil {
		return nil, err
	}

	draftTelegram := protoToDraft(req.Telegram)

	draftTelegram.DateAndTime.EndBlockNum = telegram.EndBlockNum

	telegramCode, err := encoder.Encoder(draftTelegram)
	if err != nil {
		return nil, err
	}

	telegram.Update(draftTelegram)
	telegram.TelegramCode = telegramCode

	err = s.storage.UpdateTelegram(ctx, telegram)
	if err != nil {
		return nil, err
	}

	response := telegramToProto(telegram)

	return &pb.UpdateTelegramResponse{
		Telegram: response,
	}, nil
}

func (s *HydrologyBufferervice) UpdateTelegramByCode(ctx context.Context, req *pb.UpdateTelegramByCodeRequest) (*pb.UpdateTelegramResponse, error) {

	draftTelegram, err := decoder.NewTelegram(req.TelegramCode)
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

	telegram, err := s.storage.GetTelegramByID(ctx, telegramId)
	if err != nil {
		return nil, err
	}

	telegram.Update(draftTelegram)
	telegram.TelegramCode = telegramCode

	err = s.storage.UpdateTelegram(ctx, telegram)
	if err != nil {
		return nil, err
	}

	response := telegramToProto(telegram)

	return &pb.UpdateTelegramResponse{
		Telegram: response,
	}, nil
}

func (s *HydrologyBufferervice) GetTelegram(ctx context.Context, req *pb.GetTelegramRequest) (*pb.GetTelegramResponse, error) {

	telegramId, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, err
	}

	telegram, err := s.storage.GetTelegramByID(ctx, telegramId)
	if err != nil {
		return nil, err
	}

	response := telegramToProto(telegram)

	return &pb.GetTelegramResponse{
		Telegram: response,
	}, nil
}

func (s *HydrologyBufferervice) GetTelegrams(ctx context.Context, req *pb.GetTelegramsRequest) (*pb.GetTelegramsResponse, error) {
	telegrams, err := s.storage.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	response := make([]*pb.Telegram, len(*telegrams))

	for i := 0; i < len(response); i++ {
		telegram := telegramToProto(&(*telegrams)[i]) // передаем указатель на элемент среза
		response[i] = telegram
	}

	return &pb.GetTelegramsResponse{
		Telegrams: response,
	}, nil
}

func (s *HydrologyBufferervice) TransferToSystem(ctx context.Context, req *pb.TransferToSystemRequest) (*pb.TransferToSystemResponse, error) {

	uuids := make([]uuid.UUID, len(req.Id))

	for i := 0; i < len(uuids); i++ {
		id, err := uuid.Parse(req.Id[i])
		if err != nil {
			return nil, err
		}
		uuids[i] = id
	}

	telegrams, err := s.storage.GetTelegramsById(ctx, uuids)
	if err != nil {
		return nil, err
	}

	const maxBatchSize = 100 // Максимальное количество элементов в батче

	numBatches := (len(*telegrams)*2 + maxBatchSize - 1) / maxBatchSize
	batches := make([]kafka_dto.WaterLevelRecords, numBatches)

	for i := 0; i < len(batches); i++ {
		batches[i] = *kafka_dto.NewWaterLevelRecords(maxBatchSize)
	}

	addToBatch := func(wl kafka_dto.WaterLevel, idx int) {
		batchIdx := idx / maxBatchSize
		batches[batchIdx].Waterlevels = append(batches[batchIdx].Waterlevels, wl)
	}

	for idx, telegram := range *telegrams {
		if telegram.WaterLevelOnTime.Valid && telegram.WaterLevelOnTime.Int32 != decoder_types.CouldNotMeasure {
			addToBatch(kafka_dto.WaterLevel{
				Date:       telegram.DateTime,
				WaterLevel: telegram.WaterLevelOnTime.Int32,
				PostCode:   telegram.PostCode,
			}, idx)
			idx++
		}

		if telegram.WaterLevelOn20h.Valid && telegram.WaterLevelOn20h.Int32 != decoder_types.CouldNotMeasure {
			settime := time.Date(
				telegram.DateTime.Year(),
				telegram.DateTime.Month(),
				telegram.DateTime.Day(),
				20, 0, 0, 0,
				telegram.DateTime.Location(),
			)
			addToBatch(kafka_dto.WaterLevel{
				Date:       settime,
				WaterLevel: telegram.WaterLevelOn20h.Int32,
				PostCode:   telegram.PostCode,
			}, idx)
			idx++
		}
	}

	var wg sync.WaitGroup
	errCh := make(chan error, len(batches))

	for _, batch := range batches {
		wg.Add(1)
		go func(b kafka_dto.WaterLevelRecords) {
			defer wg.Done()
			if err := kafka.SendMessageToKafka(s.KafkaProducer, s.KafkaConfig.Topic, &b); err != nil {
				errCh <- err
			}
		}(batch)
	}
	go func() {
		wg.Wait()
		close(errCh)
	}()

	for err := range errCh {
		return nil, err
	}

	if err := s.storage.RemoveTelegrams(ctx, uuids); err != nil {
		return nil, err
	}

	return &pb.TransferToSystemResponse{
		Success: true,
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
	if req.Inflow.Valid {
		res.Inflow = &wrapperspb.DoubleValue{
			Value: req.Inflow.Float64,
		}
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

func protoToDraft(req *pb.Telegram) (res *decoder.Telegram) {
	res = &decoder.Telegram{}

	res.PostCode = decoder_types.PostCode(req.PostCode)

	res.DateAndTime.Date = byte(req.Datetime.AsTime().Day())
	res.DateAndTime.Time = byte(req.Datetime.AsTime().Hour())
	res.DateAndTime.EndBlockNum = 1

	res.IsDangerous = decoder_types.IsDangerous(req.IsDangerous)

	if req.WaterLevelOnTime != nil {
		buffer := decoder_types.WaterLevelOnTime(req.WaterLevelOnTime.Value)
		res.WaterLevelOnTime = &buffer
	}

	if req.DeltaWaterLevel != nil {
		buffer := decoder_types.DeltaWaterLevel(req.DeltaWaterLevel.Value)
		res.DeltaWaterLevel = &buffer
	}

	if req.WaterLevelOn20H != nil {
		buffer := decoder_types.WaterLevelOn20h(req.WaterLevelOn20H.Value)
		res.WaterLevelOn20h = &buffer
	}

	if req.WaterTemperature != nil || req.AirTemperature != nil {
		res.Temperature = &decoder_types.Temperature{}
	}

	if req.WaterTemperature != nil {
		buffer := req.WaterTemperature.Value
		res.Temperature.WaterTemperature = &buffer
	}

	if req.AirTemperature != nil {
		buffer := req.AirTemperature.Value
		res.Temperature.AirTemperature = &buffer
	}

	if req.IcePhenomeniaState != nil {
		buffer := decoder_types.IcePhenomeniaState(req.IcePhenomeniaState.Value)
		res.IcePhenomeniaState = &buffer
	}

	if req.IceHeight != nil || req.SnowHeight != nil {
		res.IceInfo = &decoder_types.IceInfo{}
	}

	if req.IceHeight != nil {
		buffer := req.IceHeight.Value
		res.IceInfo.Ice = &buffer
	}

	if req.SnowHeight != nil {
		buffer := decoder_types.SnowHeight(req.SnowHeight.Value)
		res.IceInfo.Snow = &buffer
	}

	if req.WaterFlow != nil {
		buffer := decoder_types.Waterflow(req.WaterFlow.Value)
		res.Waterflow = &buffer
	}

	if req.PrecipitationValue != nil || req.PrecipitationDuration != nil {
		res.Precipitation = &decoder_types.Precipitation{}
	}

	if req.PrecipitationValue != nil {
		buffer := req.PrecipitationValue.Value
		res.Precipitation.Value = &buffer
	}

	if req.PrecipitationDuration != nil {
		buffer := decoder_types.PrecipitationDuration(req.PrecipitationDuration.Value)
		res.Precipitation.Duration = &buffer
	}

	if req.ReservoirDate != nil {
		buffer := decoder_types.IsReservoirDate(req.ReservoirDate.AsTime().Day())
		res.IsReservoirDate = &buffer
		res.Reservoir = &decoder_types.Reservoir{}

		if req.HeadwaterLevel != nil {
			buffer := decoder_types.HeadwaterLevel(req.HeadwaterLevel.Value)
			res.Reservoir.HeadwaterLevel = &buffer
		}

		if req.AverageReservoirLevel != nil {
			buffer := decoder_types.AverageReservoirLevel(req.AverageReservoirLevel.Value)
			res.Reservoir.AverageReservoirLevel = &buffer
		}

		if req.DownstreamLevel != nil {
			buffer := decoder_types.DownstreamLevel(req.DownstreamLevel.Value)
			res.Reservoir.DownstreamLevel = &buffer
		}

		if req.ReservoirVolume != nil {
			buffer := decoder_types.ReservoirVolume(req.ReservoirVolume.Value)
			res.Reservoir.ReservoirVolume = &buffer
		}
	}

	if req.ReservoirWaterInflowDate != nil {
		buffer := decoder_types.IsReservoirWaterInflowDate(req.ReservoirWaterInflowDate.AsTime().Day())
		res.IsReservoirWaterInflowDate = &buffer
		res.ReservoirWaterInflow = &decoder_types.ReservoirWaterInflow{}

		if req.Inflow != nil {
			buffer := decoder_types.Inflow(req.Inflow.Value)
			res.ReservoirWaterInflow.Inflow = &buffer
		}

		if req.Reset_ != nil {
			buffer := decoder_types.Reset(req.Reset_.Value)
			res.ReservoirWaterInflow.Reset = &buffer
		}
	}

	if len(req.IcePhenomenias) != 0 {
		res.IcePhenomenia = make([]*decoder_types.Phenomenia, len(req.IcePhenomenias))

		for i := 0; i < len(res.IcePhenomenia); i++ {
			res.IcePhenomenia[i] = &decoder_types.Phenomenia{Phenomen: byte(req.IcePhenomenias[i].Phenomen)}

			if req.IcePhenomenias[i].Intensity != nil {
				res.IcePhenomenia[i].IsUntensity = true
				buffer := byte(req.IcePhenomenias[i].Intensity.Value)
				res.IcePhenomenia[i].Intensity = &buffer
			}
		}
	}

	return
}

func (s *HydrologyBufferervice) SetKafkaConfig(config kafka.KafkaConfig) {
	s.KafkaConfig = config
}
