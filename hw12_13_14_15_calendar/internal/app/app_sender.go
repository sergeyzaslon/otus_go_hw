package app

type NotificationSource interface {
	GetNotificationChannel() (<-chan Notification, error)
}

type NotificationTransport interface {
	String() string
	Send(Notification) error
}

type NotificationSender struct {
	source     NotificationSource
	logger     Logger
	transports []NotificationTransport
}

func NewNotificationSender(
	source NotificationSource,
	logger Logger,
	transports []NotificationTransport,
) *NotificationSender {
	return &NotificationSender{source, logger, transports}
}

func (s *NotificationSender) Run() {
	s.logger.Info("[notification] start")
	ch, err := s.source.GetNotificationChannel()
	if err != nil {
		s.logger.Error("[notification] failed to start consume channel: %s", err)
		return
	}

	for notification := range ch {
		for _, t := range s.transports {
			if err := t.Send(notification); err != nil {
				s.logger.Error("[notification] Failed to send event notification %s to %s", notification.EventID, t.String())
				continue
			}

			s.logger.Info("[notification] notification %s sent via %s", notification.EventID, t.String())
		}
	}
}
