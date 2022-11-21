package apollo

import (
	"encoding/json"
	"time"
)

func StartNotifications() {
	t := time.NewTicker(2 * time.Second)
	for range t.C {
		AsyncNotifications()
		t.Reset(2 * time.Second)
	}
}
func AsyncNotifications() {
	opts := []RequestOption{
		Timeout(10 * time.Minute),
		Success(SuccessNotifications),
		NotModified(NotModifiedNotifications),
	}
	GetNotify(conf, nil, opts...)
}

func SuccessNotifications(body []byte) error {
	ns := make([]*Notification, 0)
	err := json.Unmarshal(body, &ns)
	if err != nil {
		return err
	}
	// TODO set notifications

	return nil
}

func NotModifiedNotifications(body []byte) error {
	return nil
}

func AsyncConfigs() {
	opts := []RequestOption{
		Success(SuccessConfigs),
		NotModified(NotModifiedConfigs),
	}
	GetConfigs(conf, "", opts...)
}

func SuccessConfigs(body []byte) error {
	var apol Apollo
	err := json.Unmarshal(body, &apol)
	if err != nil {
		return err
	}
	// TODO set apol
	return nil
}

func NotModifiedConfigs(body []byte) error {
	return nil
}
