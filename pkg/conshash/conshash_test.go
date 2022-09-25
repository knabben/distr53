package conshash

import (
	"fmt"
	"hash/fnv"
	"testing"
)

func newConfig() Config {
	return Config{
		PartitionCount:    23,
		ReplicationFactor: 20,
		Load:              1.25,
		Hasher:            hasher{},
	}
}

type testMember string

func (tm testMember) String() string {
	return string(tm)
}

type hasher struct{}

func (hs hasher) Sum64(data []byte) uint64 {
	h := fnv.New64()
	h.Write(data)
	return h.Sum64()
}

func TestConsistentAdd(t *testing.T) {
	cfg := newConfig()
	c := New(nil, cfg)
	members := make(map[string]struct{})

	for i := 0; i < 8; i++ {
		member := testMember(fmt.Sprintf("node%d.olric", i))
		members[member.String()] = struct{}{}
		c.Add(member)
	}
	for member := range members {
		found := false
		for _, mem := range c.GetMembers() {
			if member == mem.String() {
				found = true
			}
		}
		if !found {
			t.Fatalf("%s could not be found", member)
		}
	}
}

func TestConsistentRemove(t *testing.T) {
	members := []Member{}
	for i := 0; i < 8; i++ {
		member := testMember(fmt.Sprintf("node%d.olric", i))
		members = append(members, member)
	}
	cfg := newConfig()
	c := New(members, cfg)
	if len(c.GetMembers()) != len(members) {
		t.Fatalf("inserted member count is different")
	}
	for _, member := range members {
		c.Remove(member.String())
	}
	if len(c.GetMembers()) != 0 {
		t.Fatalf("member count should be zero")
	}
}

func TestConsistentLoad(t *testing.T) {
	members := []Member{}
	for i := 0; i < 8; i++ {
		member := testMember(fmt.Sprintf("node%d.olric", i))
		members = append(members, member)
	}
	cfg := newConfig()
	c := New(members, cfg)
	if len(c.GetMembers()) != len(members) {
		t.Fatalf("inserted member count is different")
	}
	maxLoad := c.AverageLoad()
	for member, load := range c.LoadDistribution() {
		if load > maxLoad {
			t.Fatalf("%s exceeds max load. Its load: %f, max load: %f", member, load, maxLoad)
		}
	}
}

func TestConsistentLocateKey(t *testing.T) {
	cfg := newConfig()
	c := New(nil, cfg)
	key := []byte("Olric")
	res := c.LocateKey(key)
	if res != nil {
		t.Fatalf("This should be nil: %v", res)
	}
	members := make(map[string]struct{})
	for i := 0; i < 8; i++ {
		member := testMember(fmt.Sprintf("node%d.olric", i))
		members[member.String()] = struct{}{}
		c.Add(member)
	}
	res = c.LocateKey(key)
	if res == nil {
		t.Fatalf("This shouldn't be nil: %v", res)
	}
}

func TestConsistentInsufficientMemberCount(t *testing.T) {
	members := []Member{}
	for i := 0; i < 8; i++ {
		member := testMember(fmt.Sprintf("node%d.olric", i))
		members = append(members, member)
	}
	cfg := newConfig()
	c := New(members, cfg)
	key := []byte("Olric")
	_, err := c.GetClosestN(key, 30)
	if err != ErrInsufficientMemberCount {
		t.Fatalf("Expected ErrInsufficientMemberCount(%v), Got: %v", ErrInsufficientMemberCount, err)
	}
}

func TestConsistentClosestMembers(t *testing.T) {
	members := []Member{}
	for i := 0; i < 8; i++ {
		member := testMember(fmt.Sprintf("node%d.olric", i))
		members = append(members, member)
	}
	cfg := newConfig()
	c := New(members, cfg)
	key := []byte("Olric")
	closestn, err := c.GetClosestN(key, 2)
	if err != nil {
		t.Fatalf("Expected nil, Got: %v", err)
	}
	if len(closestn) != 2 {
		t.Fatalf("Expected closest member count is 2. Got: %d", len(closestn))
	}
	partID := c.FindPartitionID(key)
	owner := c.GetPartitionOwner(partID)
	for i, cl := range closestn {
		if i != 0 && cl.String() == owner.String() {
			t.Fatalf("Backup is equal the partition owner: %s", owner.String())
		}
	}
}
