package repo

type MessageRepository interface {
	Add(string) uint
	Fetch(uint) (string, bool)
}

type RuntimeMessageRepository struct {
	messages []string
}

func (r *RuntimeMessageRepository) Add(message string) uint {
	r.messages = append(r.messages, message)
	return uint(len(r.messages) - 1)
}

func (r *RuntimeMessageRepository) Fetch(id uint) (string, bool) {
	if id >= uint(len(r.messages)) {
		return "", false
	}
	return r.messages[id], true
}
