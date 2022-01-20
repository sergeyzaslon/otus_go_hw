package transport

type LogNotificationTransport struct {
	logger app.Logger
}

func NewLogNotificationTransport(logger app.Logger) *LogNotificationTransport {
	return &LogNotificationTransport{logger: logger}
}

func (t *LogNotificationTransport) String() string {
	return "LogNotificationTransport"
}

func (t *LogNotificationTransport) Send(n app.Notification) error {
	t.logger.Info("[notification][transport][log] %v", n)
	return nil
}
