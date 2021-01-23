package mongodb

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/mongo/readpref"

	"github.com/sharmarajdaksh/yorpoll-api/config"
	"github.com/sharmarajdaksh/yorpoll-api/internal/log"
	"github.com/sharmarajdaksh/yorpoll-api/internal/poll"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Connection represents a new MySQL database connection
type Connection struct {
	client *mongo.Client
	dbName string
}

// GetConnectionString builds and returns an appropriate connection string
func (c *Connection) GetConnectionString(config *config.Config) string {
	connectionString := fmt.Sprintf("mongodb://%s:%s@%s:%d/%s",
		config.Database.Username,
		config.Database.Password(),
		config.Database.Hostname,
		config.Database.Port,
		config.Database.Name,
	)
	c.dbName = config.Database.Name
	return connectionString
}

// Connect will establish a database connection
func (c *Connection) Connect(config *config.Config) error {
	cstr := c.GetConnectionString(config)

	clientOpts := options.Client().ApplyURI(cstr)
	client, err := mongo.NewClient(clientOpts)
	if err != nil {
		log.Error(err)
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Error(err)
		return err
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Error(err)
		return err
	}

	c.client = client
	return nil
}

// Close will close the database connection
func (c *Connection) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return c.client.Disconnect(ctx)
}

// Ping pings the database to verify that the database connection is alive
func (c *Connection) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return c.client.Ping(ctx, readpref.Primary())
}

// GetPollByID fetches a poll from the database by ID
func (c *Connection) GetPollByID(pollID string) (*poll.Poll, error) {
	col := c.client.Database(c.dbName).Collection(pollCollection)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var p poll.Poll
	if err := col.FindOne(ctx, bson.M{"_id": pollID}).Decode(&p); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		log.Error(err)
		return nil, err
	}

	return &p, nil
}

// SavePoll saves a Poll and corresponding Options
func (c *Connection) SavePoll(p *poll.Poll) error {
	col := c.client.Database(c.dbName).Collection(pollCollection)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	data, err := bson.Marshal(p)
	if err != nil {
		log.Error(err)
		return err
	}

	_, err = col.InsertOne(ctx, data)
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}

// DeletePoll will delete a poll from the database
func (c *Connection) DeletePoll(p *poll.Poll) (bool, error) {
	return c.DeletePollByID(p.ID)
}

// DeletePollByID deletes a poll from the database
func (c *Connection) DeletePollByID(pollID string) (bool, error) {
	col := c.client.Database(c.dbName).Collection(pollCollection)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	r, err := col.DeleteOne(ctx, bson.M{"_id": pollID})
	if err != nil {
		log.Error(err)
		return false, err
	}

	if r.DeletedCount == 0 {
		return false, nil
	}

	return true, nil
}

// AddPollVote adds a vote to a particular option for a poll
func (c *Connection) AddPollVote(pollID string, optionID string) (bool, error) {
	col := c.client.Database(c.dbName).Collection(pollCollection)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	mod := time.Now().Unix()

	r, err := col.UpdateOne(
		ctx,
		bson.M{"_id": pollID, "options._id": optionID},
		bson.D{
			{Key: "$inc", Value: bson.D{{Key: "votes", Value: 1}, {Key: "options.$.votes", Value: 1}}},
			{Key: "$set", Value: bson.D{{Key: "modified", Value: mod}}}})
	if err != nil {
		log.Error(err)
		return false, err
	}

	if r.ModifiedCount == 0 {
		return false, nil
	}
	return true, nil
}
