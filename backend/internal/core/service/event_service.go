package service

import (
	"context"
	"encoding/json"
	"fernandoglatz/home-management/internal/core/common/utils"
	"fernandoglatz/home-management/internal/core/common/utils/exceptions"
	"fernandoglatz/home-management/internal/core/common/utils/log"
	"fernandoglatz/home-management/internal/core/entity"
	"fernandoglatz/home-management/internal/core/entity/event"
	"fernandoglatz/home-management/internal/core/model/message"
	"fernandoglatz/home-management/internal/infrastructure/config"
	"fernandoglatz/home-management/internal/infrastructure/repository"
	"fmt"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

var eventServices map[string]any
var eventServiceMutex sync.Mutex
var processRfMessageMutex sync.Mutex

const REDIS_RF_EVENT_PREFIX = "last-rf-event:"

type EventService[T entity.IEvent] struct {
	service    Service[T]
	repository repository.EventRepository[T]
}

func GetEventService[T entity.IEvent]() EventService[T] {
	entity := utils.Instance[T]()
	typeName := utils.GetTypeName(entity)

	eventServiceMutex.Lock()
	defer eventServiceMutex.Unlock()

	if eventServices == nil {
		eventServices = make(map[string]any)
	}

	eventService := eventServices[typeName]

	if eventService == nil {
		service := GetGenericService[T]()
		repository := repository.GetEventRepository[T]()

		eventService = EventService[T]{
			service:    service,
			repository: repository,
		}

		eventServices[typeName] = eventService
	}

	return eventService.(EventService[T])
}

func (eventService EventService[T]) Get(ctx context.Context, id string) (T, *exceptions.WrappedError) {
	return eventService.service.Get(ctx, id)
}

func (eventService EventService[T]) GetAll(ctx context.Context) ([]T, *exceptions.WrappedError) {
	return eventService.service.GetAll(ctx)
}

func (eventService EventService[T]) Save(ctx context.Context, entity T) *exceptions.WrappedError {
	return eventService.service.Save(ctx, entity)
}

func (eventService EventService[T]) Remove(ctx context.Context, entity T) *exceptions.WrappedError {
	return eventService.service.Remove(ctx, entity)
}

func (eventService EventService[T]) ProcessMessage(ctx context.Context, body []byte) *exceptions.WrappedError {
	var eventMessage message.EventMessage

	err := json.Unmarshal(body, &eventMessage)
	if err != nil {
		return &exceptions.WrappedError{
			BaseError: exceptions.InvalidJSON,
			Error:     err,
		}
	}

	eventType := eventMessage.Type
	switch eventType {
	case event.RECEIVED_RF:
		return eventService.processRfEvent(ctx, body)

	default:
		return eventService.processDefaultMessage(ctx, eventMessage)
	}
}

func (eventService EventService[T]) GetRfEvents(ctx context.Context, code int, bits int, protocol int, frequency int, startDate time.Time, endDate time.Time) ([]entity.RfEvent, *exceptions.WrappedError) {
	return eventService.repository.GetRfEvents(ctx, code, bits, protocol, frequency, startDate, endDate)
}

func (eventService EventService[T]) populateEvent(event *entity.Event, eventMessage message.EventMessage) {
	event.Type = eventMessage.Type
	event.Version = eventMessage.Version
	event.Date = eventMessage.Date
}

func (eventService EventService[T]) processDefaultMessage(ctx context.Context, eventMessage message.EventMessage) *exceptions.WrappedError {
	log.Warn(ctx).Msg("Processing unknown event: " + string(eventMessage.Type))
	event := &entity.Event{}

	defaultEventService := GetEventService[*entity.Event]()
	defaultEventService.populateEvent(event, eventMessage)

	return defaultEventService.Save(ctx, event)
}

func (eventService EventService[T]) processRfEvent(ctx context.Context, body []byte) *exceptions.WrappedError {
	log.Info(ctx).Msg("Processing RF event")

	processRfMessageMutex.Lock()
	defer processRfMessageMutex.Unlock()

	var rfEventMessage message.RfEventMessage
	var lastRfEventMessage message.RfEventMessage

	err := json.Unmarshal(body, &rfEventMessage)
	if err != nil {
		return &exceptions.WrappedError{
			BaseError: exceptions.InvalidJSON,
			Error:     err,
		}
	}

	ttl := config.ApplicationConfig.Data.Redis.TTL.RfEvent
	redisKey := fmt.Sprintf("%s%d-%d-%d-%d", REDIS_RF_EVENT_PREFIX, rfEventMessage.Code, rfEventMessage.Bits, rfEventMessage.Protocol, rfEventMessage.Frequency)
	err = utils.RedisDatabase.GetStruct(ctx, redisKey, &lastRfEventMessage)

	if err != nil && err != redis.Nil {
		log.Error(ctx).Msg("Error on retrieving RF event from redis: " + err.Error())
	}

	err = utils.RedisDatabase.SetStruct(ctx, redisKey, rfEventMessage, ttl)
	if err != nil {
		log.Error(ctx).Msg("Error on putting RF event on redis: " + err.Error())
	}

	mergePeriod := config.ApplicationConfig.Application.Processing.RfEvents.MergePeriod
	currentDate := rfEventMessage.Date
	lastDate := lastRfEventMessage.Date

	if lastDate.IsZero() {
		startDate := currentDate.Add(-mergePeriod)
		endDate := currentDate.Add(mergePeriod)
		events, errw := eventService.GetRfEvents(ctx, rfEventMessage.Code, rfEventMessage.Bits, rfEventMessage.Protocol, rfEventMessage.Frequency, startDate, endDate)
		if errw != nil {
			log.Error(ctx).Msg("Error on retrieving old RF events: " + errw.GetMessage())
		}

		for _, event := range events {
			lastDate = event.Date
		}
	}

	diff := currentDate.Sub(lastDate)

	if diff > mergePeriod {
		event := &entity.RfEvent{
			Code:             rfEventMessage.Code,
			Bits:             rfEventMessage.Bits,
			Protocol:         rfEventMessage.Protocol,
			Frequency:        rfEventMessage.Frequency,
			ReceiveTolerance: rfEventMessage.ReceiveTolerance,
		}

		defaultEventService := GetEventService[*entity.RfEvent]()
		defaultEventService.populateEvent(&event.Event, rfEventMessage.EventMessage)

		return defaultEventService.Save(ctx, event)

	} else {
		log.Warn(ctx).Msg("Ignoring RF event, duplicated: " + redisKey)
	}

	return nil
}
