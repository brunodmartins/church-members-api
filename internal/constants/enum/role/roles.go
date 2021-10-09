package role

import (
	"github.com/BrunoDM2943/church-members-api/internal/constants/enum"
	"strings"
)

const (
	ADMIN enum.Role = iota
	USER
)

func From(value string) enum.Role {
	switch strings.ToUpper(value) {
	case "ADMIN":
		return ADMIN
	case "USER":
		return USER
	default:
		return -1
	}
}
