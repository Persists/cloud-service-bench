package log

import (
	"math/rand"
)

type LogSynthesizer struct {
	Name          string
	MessageLength int
}

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// getRandomMessage generates a random message of length MessageLength
func (g *LogSynthesizer) getRandomMessage() string {
	b := make([]byte, g.MessageLength)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// generateSeverity generates a random severity based on the weights (using a roulette wheel selection)
func (g *LogSynthesizer) getRandomSeverity() Severity {
	r := rand.Float64()
	switch {
	case r < 0.6:
		return INFO
	case r < 0.8:
		return WARN
	case r < 0.9:
		return ERROR
	default:
		return DEBUG
	}
}

// SynthesizeLog generates a log message with random severity, message, and tags
func (g *LogSynthesizer) SynthesizeLog() *LogMessage {
	severity := g.getRandomSeverity()
	message := g.getRandomMessage()
	tags := []string{"tag1", "tag2"}
	return &LogMessage{
		Severity: severity,
		Name:     g.Name,
		Message:  message,
		Tags:     tags,
	}
}

// SynthesizeLogs generates n log messages
func (g *LogSynthesizer) SynthesizeLogs(n int) []*LogMessage {
	logs := make([]*LogMessage, n)
	for i := 0; i < n; i++ {
		logs[i] = g.SynthesizeLog()
	}
	return logs
}

// NewLogSynthesizer creates a new LogSynthesizer
func NewLogSynthesizer(name string, messageLength int) *LogSynthesizer {
	return &LogSynthesizer{
		Name:          name,
		MessageLength: messageLength,
	}
}
