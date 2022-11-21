package apollo

type ChangeKind int32

const (
	ADDED ChangeKind = iota
	MODIFIED
	DELETED
)
const ()

type ChangeListener interface {
	OnChange(event *ChangeEvent)
}

type ChangeEvent struct {
	Namespace string
	Changes   map[string]*ConfigChange
}

type ConfigChange struct {
	OldValue string
	NewValue string
	Kind     ChangeKind
}

func CreateModifyChange(oldValue string, newValue string) *ConfigChange {
	return &ConfigChange{
		OldValue: oldValue,
		NewValue: newValue,
		Kind:     MODIFIED,
	}
}

func CreateAddChange(newValue string) *ConfigChange {
	return &ConfigChange{
		NewValue: newValue,
		Kind:     ADDED,
	}
}

func CreateDelChange(oldValue string) *ConfigChange {
	return &ConfigChange{
		OldValue: oldValue,
		Kind:     DELETED,
	}
}
