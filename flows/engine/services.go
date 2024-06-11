package engine

import (
	"github.com/nyaruka/goflow/flows"

	"github.com/pkg/errors"
)

// EmailServiceFactory resolves a session to a email service
type EmailServiceFactory func(flows.Session) (flows.EmailService, error)

// WebhookServiceFactory resolves a session to a webhook service
type WebhookServiceFactory func(flows.Session) (flows.WebhookService, error)

// ClassificationServiceFactory resolves a session and classifier to an NLU service
type ClassificationServiceFactory func(flows.Session, *flows.Classifier) (flows.ClassificationService, error)

// TicketServiceFactory resolves a session to a ticket service
type TicketServiceFactory func(flows.Session, *flows.Ticketer) (flows.TicketService, error)

// AirtimeServiceFactory resolves a session to an airtime service
type AirtimeServiceFactory func(flows.Session) (flows.AirtimeService, error)

// ExternalServiceServiceFactory resolves a session to an external service service
type ExternalServiceServiceFactory func(flows.Session, *flows.ExternalService) (flows.ExternalServiceService, error)

// WeniGPTServiceFactory resolves a session to a weniGPT service
type WeniGPTServiceFactory func(flows.Session) (flows.WeniGPTService, error)

// BrainServiceFactory resolves a session to a brain service
type BrainServiceFactory func(flows.Session) (flows.BrainService, error)

type MsgCatalogServiceFactory func(flows.Session, *flows.MsgCatalog) (flows.MsgCatalogService, error)

type OrgContextServiceFactory func(flows.Session, *flows.OrgContext) (flows.OrgContextService, error)

type services struct {
	email           EmailServiceFactory
	webhook         WebhookServiceFactory
	classification  ClassificationServiceFactory
	ticket          TicketServiceFactory
	airtime         AirtimeServiceFactory
	externalService ExternalServiceServiceFactory
	wenigpt         WeniGPTServiceFactory
	msgCatalog      MsgCatalogServiceFactory
	orgContext      OrgContextServiceFactory
	brain           BrainServiceFactory
}

func newEmptyServices() *services {
	return &services{
		email: func(flows.Session) (flows.EmailService, error) {
			return nil, errors.New("no email service factory configured")
		},
		webhook: func(flows.Session) (flows.WebhookService, error) {
			return nil, errors.New("no webhook service factory configured")
		},
		wenigpt: func(flows.Session) (flows.WeniGPTService, error) {
			return nil, errors.New("no wenigpt service factory configured")
		},
		classification: func(flows.Session, *flows.Classifier) (flows.ClassificationService, error) {
			return nil, errors.New("no classification service factory configured")
		},
		ticket: func(flows.Session, *flows.Ticketer) (flows.TicketService, error) {
			return nil, errors.New("no ticket service factory configured")
		},
		airtime: func(flows.Session) (flows.AirtimeService, error) {
			return nil, errors.New("no airtime service factory configured")
		},
		externalService: func(flows.Session, *flows.ExternalService) (flows.ExternalServiceService, error) {
			return nil, errors.New("no external service factory configured")
		},
		msgCatalog: func(flows.Session, *flows.MsgCatalog) (flows.MsgCatalogService, error) {
			return nil, errors.New("no msg catalog service factory configured")
		},
		orgContext: func(flows.Session, *flows.OrgContext) (flows.OrgContextService, error) {
			return nil, errors.New("no org context service factory configured")
		},
		brain: func(flows.Session) (flows.BrainService, error) {
			return nil, errors.New("no brain service factory configured")
		},
	}
}

func (s *services) Email(session flows.Session) (flows.EmailService, error) {
	return s.email(session)
}

func (s *services) Webhook(session flows.Session) (flows.WebhookService, error) {
	return s.webhook(session)
}

func (s *services) WeniGPT(session flows.Session) (flows.WeniGPTService, error) {
	return s.wenigpt(session)
}

func (s *services) Classification(session flows.Session, classifier *flows.Classifier) (flows.ClassificationService, error) {
	return s.classification(session, classifier)
}

func (s *services) Ticket(session flows.Session, ticketer *flows.Ticketer) (flows.TicketService, error) {
	return s.ticket(session, ticketer)
}

func (s *services) Airtime(session flows.Session) (flows.AirtimeService, error) {
	return s.airtime(session)
}

func (s *services) ExternalService(session flows.Session, externalService *flows.ExternalService) (flows.ExternalServiceService, error) {
	return s.externalService(session, externalService)
}

func (s *services) MsgCatalog(session flows.Session, msgCatalog *flows.MsgCatalog) (flows.MsgCatalogService, error) {
	return s.msgCatalog(session, msgCatalog)
}

func (s *services) OrgContext(session flows.Session, context *flows.OrgContext) (flows.OrgContextService, error) {
	return s.orgContext(session, context)
}

func (s *services) Brain(session flows.Session) (flows.BrainService, error) {
	return s.brain(session)
}
