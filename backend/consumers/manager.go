package consumers

import (
	"dota-gsi/backend/events"
	"dota-gsi/backend/handlers"

	"github.com/sirupsen/logrus"
)

// Consumer interface for all domain consumers
type Consumer interface {
	Start()
	Stop()
}

// ConsumerManager manages multiple domain consumers
type ConsumerManager struct {
	consumers []Consumer
	logger    *logrus.Entry
}

// NewConsumerManager creates a new consumer manager
func NewConsumerManager(logger *logrus.Entry) *ConsumerManager {
	return &ConsumerManager{
		consumers: make([]Consumer, 0),
		logger:    logger,
	}
}

// AddMapConsumer adds a MapConsumer to the manager
func (cm *ConsumerManager) AddMapConsumer(eventBus *events.EventBus, handlerList []handlers.Handler) {
	mapConsumer := NewMapConsumer(eventBus, cm.logger.WithField("consumer", "map"), handlerList)
	cm.consumers = append(cm.consumers, mapConsumer)
}

// AddHeroConsumer adds a HeroConsumer to the manager
func (cm *ConsumerManager) AddHeroConsumer(eventBus *events.EventBus, handlerList []handlers.Handler) {
	heroConsumer := NewHeroConsumer(eventBus, cm.logger.WithField("consumer", "hero"), handlerList)
	cm.consumers = append(cm.consumers, heroConsumer)
}

// AddRuneConsumer adds a RuneConsumer to the manager
func (cm *ConsumerManager) AddRuneConsumer(eventBus *events.EventBus, handlerList []handlers.Handler, gameConfig interface{}) {
	runeConsumer := NewRuneConsumer(eventBus, cm.logger.WithField("consumer", "rune"), handlerList, gameConfig)
	cm.consumers = append(cm.consumers, runeConsumer)
}

// AddTimingConsumer adds a TimingConsumer to the manager
func (cm *ConsumerManager) AddTimingConsumer(eventBus *events.EventBus, handlerList []handlers.Handler, gameConfig interface{}) {
	timingConsumer := NewTimingConsumer(eventBus, cm.logger.WithField("consumer", "timing"), handlerList, gameConfig)
	cm.consumers = append(cm.consumers, timingConsumer)
}

// AddAbilitiesConsumer adds an AbilitiesConsumer to the manager (future implementation)
func (cm *ConsumerManager) AddAbilitiesConsumer(eventBus *events.EventBus, handlerList []handlers.Handler) {
	// TODO: Implement AbilitiesConsumer
	// abilitiesConsumer := NewAbilitiesConsumer(eventBus, cm.logger.WithField("consumer", "abilities"), handlerList)
	// cm.consumers = append(cm.consumers, abilitiesConsumer)
	cm.logger.Info("⚡ AbilitiesConsumer will be implemented soon")
}

// AddItemsConsumer adds an ItemsConsumer to the manager (future implementation)
func (cm *ConsumerManager) AddItemsConsumer(eventBus *events.EventBus, handlerList []handlers.Handler) {
	// TODO: Implement ItemsConsumer
	// itemsConsumer := NewItemsConsumer(eventBus, cm.logger.WithField("consumer", "items"), handlerList)
	// cm.consumers = append(cm.consumers, itemsConsumer)
	cm.logger.Info("🎒 ItemsConsumer will be implemented soon")
}

// StartAll starts all registered consumers
func (cm *ConsumerManager) StartAll() {
	cm.logger.WithField("count", len(cm.consumers)).Info("🚀 Starting all consumers")

	for _, consumer := range cm.consumers {
		consumer.Start()
	}

	cm.logger.Info("✅ All consumers started")
}

// StopAll stops all registered consumers
func (cm *ConsumerManager) StopAll() {
	cm.logger.WithField("count", len(cm.consumers)).Info("🛑 Stopping all consumers")

	for _, consumer := range cm.consumers {
		consumer.Stop()
	}

	cm.logger.Info("✅ All consumers stopped")
}

// Count returns the number of registered consumers
func (cm *ConsumerManager) Count() int {
	return len(cm.consumers)
}
