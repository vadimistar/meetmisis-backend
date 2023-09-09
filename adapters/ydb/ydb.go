package ydb

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/google/uuid"
	"github.com/guregu/dynamo"
	"github.com/pkg/errors"
	"github.com/vadimistar/hackathon1/models"
)

type Client struct {
	db *dynamo.DB
}

func New() (*Client, error) {
	sess, err := session.NewSession()
	if err != nil {
		return nil, errors.Wrap(err, "create session")
	}

	region := os.Getenv("AWS_DEFAULT_REGION")
	if region == "" {
		return nil, errors.New("no AWS_DEFAULT_REGION env variable provided")
	}

	endpoint := os.Getenv("YDB_ENDPOINT")
	if endpoint == "" {
		return nil, errors.New("no YDB_ENDPOINT env variable provided")
	}

	idKey := os.Getenv("YDB_ID_KEY")
	if idKey == "" {
		return nil, errors.New("no YDB_ID_KEY env variable provided")
	}

	idSecret := os.Getenv("YDB_SECRET_KEY")
	if idSecret == "" {
		return nil, errors.New("no YDB_SECRET_KEY env variable provided")
	}

	db := dynamo.New(sess,
		&aws.Config{
			Credentials: credentials.NewStaticCredentials(idKey, idSecret, ""),
			Endpoint:    aws.String(endpoint),
			Region:      aws.String(region),
		},
	)

	return &Client{db: db}, nil
}

func (c *Client) GetUser(username string) (*models.User, error) {
	if username == "" {
		return nil, errors.New("empty input")
	}

	err := c.createTableIfNotExists("users", &models.User{})
	if err != nil {
		return nil, errors.Wrap(err, "create table")
	}

	table := c.db.Table("users")
	user := new(models.User)

	it := table.Scan().
		Filter("'Username' = ?", username).
		Iter()
	it.Next(user)

	if it.Err() != nil {
		if errors.Is(it.Err(), dynamo.ErrNotFound) {
			return nil, nil
		}
		return nil, errors.Wrap(it.Err(), "scan table")
	}

	return user, nil
}

func (c *Client) SaveUser(user *models.User) error {
	err := c.createTableIfNotExists("users", &models.User{})
	if err != nil {
		return errors.Wrap(err, "create table")
	}

	table := c.db.Table("users")

	id, err := uuid.NewRandom()
	if err != nil {
		return errors.Wrap(err, "cannot create random")
	}

	user.ID = id.String()

	err = table.Put(user).Run()
	if err != nil {
		return errors.Wrap(err, "put into the table")
	}

	return nil
}

func (c *Client) SaveVerification(v *models.Verification) error {
	if v.Token == "" {
		return errors.New("token is empty")
	}
	if v.UserID == "" {
		return errors.New("userID is empty")
	}

	err := c.createTableIfNotExists("verification", &models.Verification{})
	if err != nil {
		return errors.Wrap(err, "create table")
	}

	table := c.db.Table("verification")

	err = table.Put(v).Run()
	if err != nil {
		return errors.Wrap(err, "put into the table")
	}

	return nil
}

func (c *Client) createTableIfNotExists(name string, from interface{}) error {
	err := c.db.CreateTable(name, from).Run()
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			if awsErr.Code() == "ResourceInUseException" {
				// exists
				return nil
			}
		}
		return err
	}
	return nil
}
