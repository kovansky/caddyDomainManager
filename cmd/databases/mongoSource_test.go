package databases

import "testing"

func TestMongoSource_BuildUri(t *testing.T) {
	tables := []struct {
		input    MongoSource
		expected string
	}{
		{MongoSource{
			User:     "test",
			Password: "test",
			Host:     "127.0.0.1",
			Port:     27017,
			AuthDb:   "admin",
		}, "mongodb://test:test@127.0.0.1:27017/?authDatabase=admin"},
		{MongoSource{
			User: "test",
			Host: "127.0.0.1",
		}, "mongodb://test@127.0.0.1"},
		{MongoSource{
			Host: "127.0.0.1",
		}, "mongodb://127.0.0.1"},
		{MongoSource{
			Host:   "127.0.0.1",
			AuthDb: "admin",
		}, "mongodb://127.0.0.1/?authDatabase=admin"},
		{MongoSource{}, "mongodb://"},
	}

	for _, table := range tables {
		table.input.BuildUri()

		if table.input.connectionUri != table.expected {
			t.Errorf("MongoDB connection URI built incorrectly, expected %s, got %s", table.expected, table.input.connectionUri)
		}
	}
}
