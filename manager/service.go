package manager

import "github.com/cloudedcat/finance-bot/model"

// Service manages groups and participants
type Service interface {
	RegisterGroup(group model.Group) error
	RegisterParticipant(model.GroupID, model.Participant) error
	ListParticipant(model.GroupID) (model.Participants, error)
}

type service struct {
	groups       model.GroupRepository
	participants model.ParticipantRepository
}

// NewService creates a Manager Service
func NewService(groups model.GroupRepository, partics model.ParticipantRepository) Service {
	return &service{
		groups:       groups,
		participants: partics,
	}
}

func (s *service) RegisterGroup(group model.Group) error {
	if err := group.Validate(); err != nil {
		return err
	}
	return s.groups.Store(&group)
}

func (s *service) RegisterParticipant(gID model.GroupID, partic model.Participant) error {
	if err := partic.Validate(); err != nil {
		return err
	}
	return s.participants.Store(gID, &partic)
}

func (s *service) ListParticipant(gID model.GroupID) (model.Participants, error) {
	return s.participants.FindAll(gID)
}
