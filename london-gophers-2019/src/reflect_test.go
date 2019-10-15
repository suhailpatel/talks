package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var testData = map[string]interface{}{
	"id":       uint64(100),
	"name":     "Bingo 🐕",
	"pronouns": []string{"🐶"},
	"location": "Cardboard Heaven",
	"bio":      "Woof",

	"job": map[string]interface{}{
		"role":   "Office Dog",
		"squad":  "All",
		"joined": time.Date(2015, 1, 1, 7, 0, 0, 0, time.UTC),
	},
}

func TestHandRolled(t *testing.T) {
	u := ToUserUsingHandRolled(testData)
	assert.Equal(t, u.Name, "Bingo 🐕")
	assert.Equal(t, u.Pronouns, []string{"🐶"})
	assert.Equal(t, u.Bio, "Woof")
	assert.Equal(t, u.Job.Role, "Office Dog")
	assert.Equal(t, u.Job.Squad, "All")
}

func TestReflect(t *testing.T) {
	r, err := ToUserUsingReflect(testData)
	assert.Nil(t, err)
	assert.Equal(t, r.Name, "Bingo 🐕")
	assert.Equal(t, r.Pronouns, []string{"🐶"})
	assert.Equal(t, r.Bio, "Woof")
	assert.Equal(t, r.Job.Role, "Office Dog")
	assert.Equal(t, r.Job.Squad, "All")
}

func TestBothTheSameDecoded(t *testing.T) {
	h := ToUserUsingHandRolled(testData)
	r, err := ToUserUsingReflect(testData)
	assert.Nil(t, err)
	assert.Equal(t, h, r)
}

func BenchmarkHandRolled(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ToUserUsingHandRolled(testData)
	}
}

func BenchmarkReflect(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ToUserUsingReflect(testData)
	}
}
