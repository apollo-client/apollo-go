package apollo

var (
	listeners = make([]ChangeListener, 0)
)

func AddListener(listener ChangeListener) {
	listeners = append(listeners, listener)
}

func PushChangeEvent(event *ChangeEvent) {
	if len(listeners) == 0 {
		return
	}
	for _, l := range listeners {
		if l == nil {
			continue
		}
		go l.OnChange(event)
	}
}
