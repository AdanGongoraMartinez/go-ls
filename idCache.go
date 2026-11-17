package main

import (
	"os/user"
	"strconv"
)

// Cache de identidades: 1 usuario y N grupos por directorio, no 2 lookups por archivo.
type idCache struct {
	users, groups map[uint32]string
}

func (c *idCache) userName(uid uint32) string {
	if name, ok := c.users[uid]; ok {
		return name
	}
	name := "?"
	if u, err := user.LookupId(strconv.FormatUint(uint64(uid), 10)); err == nil {
		name = u.Username
	}
	c.users[uid] = name
	return name
}

func (c *idCache) groupName(gid uint32) string {
	if name, ok := c.groups[gid]; ok {
		return name
	}
	name := "?"
	if g, err := user.LookupGroupId(strconv.FormatUint(uint64(gid), 10)); err == nil {
		name = g.Name
	}
	c.groups[gid] = name
	return name
}
