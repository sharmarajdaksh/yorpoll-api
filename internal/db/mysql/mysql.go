package mysql

import (
	"database/sql"
	gosql "database/sql"
	"fmt"
	"time"

	"github.com/sharmarajdaksh/yorpoll-api/internal/log"

	// MySQL driver
	_ "github.com/go-sql-driver/mysql"
	"github.com/sharmarajdaksh/yorpoll-api/config"
	"github.com/sharmarajdaksh/yorpoll-api/internal/poll"
)

// Connection represents a new MySQL database connection
type Connection struct {
	conn *sql.DB
}

// GetConnectionString builds and returns an appropriate connection string
func (c *Connection) GetConnectionString(config *config.Config) string {
	connectionString := fmt.Sprintf("%s:%s@%s(%s:%d)/%s",
		config.Database.Username,
		config.Database.Password(),
		config.Database.NetworkType,
		config.Database.Hostname,
		config.Database.Port,
		config.Database.Name,
	)
	return connectionString
}

// Connect will establish a database connection
func (c *Connection) Connect(config *config.Config) error {
	cstr := c.GetConnectionString(config)
	conn, err := gosql.Open("mysql", cstr)
	if err != nil {
		log.Error(err)
		return err
	}

	// run migrations
	err = doMigrations(conn)
	if err != nil {
		log.Error(err)
		return err
	}

	log.Logger.Debug().Str("event", "migrations successful").Msg("mysql database migrations successful")

	c.conn = conn
	return nil
}

// Close will close the database connection
func (c *Connection) Close() error {
	return c.conn.Close()
}

// Ping pings the database to verify that the database connection is alive
func (c *Connection) Ping() error {
	return c.conn.Ping()
}

// GetPollByID fetches a poll from the database by ID
func (c *Connection) GetPollByID(pollID string) (*poll.Poll, error) {
	ostmt, err := c.conn.Prepare("SELECT id, votes, title FROM poll_option WHERE poll_id = ?")
	if err != nil {
		log.Error(err)
		return nil, err
	}
	defer ostmt.Close()

	opts := []poll.Option{}
	rows, err := ostmt.Query(pollID)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	defer rows.Close()

	totalVotes := int64(0)
	for rows.Next() {
		var opt poll.Option
		err = rows.Scan(&opt.ID, &opt.Votes, &opt.Title)
		if err != nil {
			log.Error(err)
			return nil, err
		}

		totalVotes += opt.Votes
		opts = append(opts, opt)
	}
	err = rows.Err()
	if err != nil {
		log.Error(err)
		return nil, err
	}

	stmt, err := c.conn.Prepare("SELECT title, poll_description, created, modified, expiry FROM poll WHERE id = ?")
	if err != nil {
		log.Error(err)
		return nil, err
	}
	defer stmt.Close()

	var poll poll.Poll
	err = stmt.QueryRow(pollID).Scan(&poll.Title, &poll.Description, &poll.Created, &poll.Modified, &poll.Expiry)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		log.Error(err)
		return nil, err
	}

	poll.ID = pollID
	poll.Votes = totalVotes
	poll.Options = opts

	return &poll, nil
}

// SavePoll saves a Poll and corresponding Options
func (c *Connection) SavePoll(p *poll.Poll) error {
	// Insert the poll record itself
	pstmt, err := c.conn.Prepare("INSERT INTO poll (id, title, poll_description, created, modified, expiry) VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Error(err)
		return err
	}
	defer pstmt.Close()

	_, err = pstmt.Exec(p.ID, p.Title, p.Description, p.Created, p.Modified, p.Expiry)
	if err != nil {
		log.Error(err)
		return err
	}

	ostmt, err := c.conn.Prepare("INSERT INTO poll_option (id, votes, title, poll_id) VALUES (?, ?, ?, ?)")
	if err != nil {
		log.Error(err)
		return err
	}
	defer ostmt.Close()

	for _, o := range p.Options {
		_, err := ostmt.Exec(o.ID, o.Votes, o.Title, p.ID)
		if err != nil {
			log.Error(err)
			return err
		}
	}

	return nil
}

// DeletePoll will delete a poll from the database
func (c *Connection) DeletePoll(p *poll.Poll) (bool, error) {
	return c.DeletePollByID(p.ID)
}

// DeletePollByID deletes a poll from the database
func (c *Connection) DeletePollByID(pollID string) (bool, error) {
	stmt, err := c.conn.Prepare("DELETE FROM poll WHERE id = ?")
	if err != nil {
		log.Error(err)
		return false, err
	}
	defer stmt.Close()

	r, err := stmt.Exec(pollID)
	if err != nil {
		log.Error(err)
		return false, err
	}

	if ra, _ := r.RowsAffected(); ra == 0 {
		return false, nil
	}

	return true, nil
}

// AddPollVote adds a vote to a particular option for a poll
func (c *Connection) AddPollVote(pollID string, optionID string) (bool, error) {
	// verify that the poll indeed does exist
	estmt, err := c.conn.Prepare("SELECT expiry FROM poll WHERE id = ?")
	if err != nil {
		log.Error(err)
		return false, err
	}
	defer estmt.Close()

	var expiry int64
	err = estmt.QueryRow(pollID).Scan(&expiry)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	if expiry < time.Now().Unix() {
		// Poll exists but is expired
		return false, poll.ErrPollExpired
	}

	stmt, err := c.conn.Prepare("UPDATE poll_option SET votes = votes + 1 WHERE id = ? AND poll_id = ?")
	if err != nil {
		log.Error(err)
		return false, err
	}
	defer stmt.Close()

	r, err := stmt.Exec(optionID, pollID)
	if err != nil {
		log.Error(err)
		return false, err
	}

	_ = c.updateModifiedTimestamp(pollID)

	if ra, _ := r.RowsAffected(); ra == 0 {
		return false, nil
	}
	return true, nil
}

// updateModifiedTimestamp updates a poll's updated timestamp
func (c *Connection) updateModifiedTimestamp(pollID string) error {
	t := time.Now().Unix()

	stmt, err := c.conn.Prepare("UPDATE poll SET modified = ? WHERE id = ?")
	if err != nil {
		log.Error(err)
		return err
	}

	_, err = stmt.Exec(t, pollID)
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}
