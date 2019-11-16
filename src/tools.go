package src

import guuid "github.com/satori/go.uuid"

// GetHash we use it to get hasj=h for todo command
func GetHash() string {
	id, _ := guuid.NewV4()
	return id.String()
}
