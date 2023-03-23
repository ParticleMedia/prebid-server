package ab

import "strings"

type UserService struct {
	users map[string][]*Version
}

func NewUserService(layers []*Layer) *UserService {
	users := map[string][]*Version{}
	userZones := map[string]bool{}
	containsZone := func(user string, zone string) bool {
		if zone == "" {
			return false
		}
		zones := strings.Split(zone, ".")
		for i := len(zones) - 1; i > 0; i-- {
			parentZone := strings.Join(zones[:i], ".")
			key := user + "@" + parentZone
			if userZones[key] {
				return true
			}
			userZones[key] = true
		}
		return false
	}
	for _, layer := range layers {
		for _, exp := range layer.Exps {
			for _, version := range exp.Versions {
				for _, user := range version.Users {
					if !containsZone(user, layer.Zone) {
						users[user] = append(users[user], version)
					}
				}
			}
		}
	}
	return &UserService{users: users}
}

func (this *UserService) AB(ctx *ABContext, result *ABResult) {
	result.MergeVersions(this.users[ctx.Factor])
}
