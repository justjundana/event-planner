package event

import _models "github.com/justjundana/event-planner/models"

type EventInterface interface {
	Get() ([]_models.Event, error)
}
