package static_test

import (
	"testing"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/assets/static"

	"github.com/stretchr/testify/assert"
)

func TestTopic(t *testing.T) {
	topic := static.NewTopic(
		assets.TopicUUID("37657cf7-5eab-4286-9cb0-bbf270587bad"),
		"Weather",
		"1876c846-ea1f-43b4-8ffa-7330772845b6",
	)
	assert.Equal(t, assets.TopicUUID("37657cf7-5eab-4286-9cb0-bbf270587bad"), topic.UUID())
	assert.Equal(t, "Weather", topic.Name())
	assert.Equal(t, assets.QueueUUID("1876c846-ea1f-43b4-8ffa-7330772845b6"), topic.QueueUUID())
}
