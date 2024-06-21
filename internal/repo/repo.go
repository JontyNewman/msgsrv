package repo

type MessageRepository interface {
	Add(string) (uint, error)
	Fetch(uint) (string, bool, error)
}

type RuntimeMessageRepository struct {
	messages []string
}

func (r *RuntimeMessageRepository) Add(message string) (uint, error) {
	r.messages = append(r.messages, message)
	return uint(len(r.messages) - 1), nil
}

func (r *RuntimeMessageRepository) Fetch(id uint) (string, bool, error) {
	if id >= uint(len(r.messages)) {
		return "", false, nil
	}
	return r.messages[id], true, nil
}
