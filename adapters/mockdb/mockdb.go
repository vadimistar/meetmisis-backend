package mockdb

// import (
// 	"log"

// 	"github.com/aws/aws-sdk-go/aws/awserr"
// 	"github.com/google/uuid"
// 	"github.com/guregu/dynamo"
// 	"github.com/pkg/errors"
// 	"github.com/vadimistar/hackathon1/models"
// )

// type Client struct {
// 	users       []*models.User
// 	taggedUsers []*models.TaggedUser
// }

// func New() (*Client, error) {
// 	return &Client{}, nil
// }

// func (c *Client) GetUser(username string) (*models.User, error) {
// 	var user *models.User
// 	for _, u := range c.users {
// 		if username == u.Username {
// 			user = u
// 			break
// 		}
// 		log.Printf("get user with name = %s not found", username)
// 		return nil, nil
// 	}
// 	log.Printf("get user with username = %s, got %+v", username, user)
// 	return user, nil
// }

// func (c *Client) SaveUser(user *models.User) error {
// 	log.Printf("save user %+v", user)

// 	id, err := uuid.NewRandom()
// 	if err != nil {
// 		return errors.Wrap(err, "cannot create random")
// 	}

// 	user.ID = id.String()

// 	log.Printf("save user: give it to user: %+v", user)

// 	return nil
// }

// func (c *Client) SaveVerification(v *models.Verification) error {
// 	return nil
// }

// func (c *Client) SaveTagsForUser(userID string, tags []string) error {
// 	tu := &models.TaggedUser{
// 		UserID: userID,
// 		Tags:   tags,
// 	}

// 	log.Printf("save tags for user: %+v", tu)

// 	return nil
// }

// func (c *Client) GetTaggedUser(userID string) (*models.TaggedUser, error) {
// 	if userID == "" {
// 		return nil, errors.New("userID is empty")
// 	}

// 	err := c.createTableIfNotExists("taggedUsers", &models.TaggedUser{})
// 	if err != nil {
// 		return nil, errors.Wrap(err, "create table")
// 	}

// 	table := c.db.Table("taggedUsers")

// 	user := new(models.TaggedUser)

// 	it := table.Scan().
// 		Filter("'UserID' = ?", userID).
// 		Iter()
// 	it.Next(user)

// 	if it.Err() != nil {
// 		if errors.Is(it.Err(), dynamo.ErrNotFound) {
// 			return nil, nil
// 		}
// 		return nil, errors.Wrap(it.Err(), "scan table")
// 	}

// 	return user, nil
// }

// func (c *Client) GetTag(tagID string) (*models.Tag, error) {
// 	if tagID == "" {
// 		return nil, errors.New("tagID is empty")
// 	}

// 	err := c.createTableIfNotExists("tags", &models.Tag{})
// 	if err != nil {
// 		return nil, errors.Wrap(err, "create table")
// 	}

// 	table := c.db.Table("tags")

// 	tag := new(models.Tag)

// 	it := table.Scan().
// 		Filter("'Id' = ?", tagID).
// 		Iter()
// 	it.Next(tag)

// 	if it.Err() != nil {
// 		if errors.Is(it.Err(), dynamo.ErrNotFound) {
// 			return nil, nil
// 		}
// 		return nil, errors.Wrap(it.Err(), "scan table")
// 	}

// 	return tag, nil
// }

// func (c *Client) GetPartner(userId string) (*models.Partner, error) {
// 	if userId == "" {
// 		return nil, errors.New("partnerID are empty")
// 	}

// 	err := c.createTableIfNotExists("partners", &models.Partner{})
// 	if err != nil {
// 		return nil, errors.Wrap(err, "create table")
// 	}

// 	table := c.db.Table("partners")

// 	partner := new(models.Partner)

// 	it := table.Scan().
// 		Filter("'UserID' = ?", userId).
// 		Iter()
// 	it.Next(partner)

// 	if it.Err() != nil {
// 		if errors.Is(it.Err(), dynamo.ErrNotFound) {
// 			return nil, nil
// 		}
// 		return nil, errors.Wrap(it.Err(), "scan table")
// 	}

// 	return partner, nil
// }

// func (c *Client) SaveTag(tag *models.Tag) error {
// 	err := c.createTableIfNotExists("tags", &models.Tag{})
// 	if err != nil {
// 		return errors.Wrap(err, "create table")
// 	}

// 	table := c.db.Table("tags")

// 	if tag.Id == "" {
// 		id, err := uuid.NewRandom()
// 		if err != nil {
// 			return errors.Wrap(err, "cannot create random")
// 		}

// 		tag.Id = id.String()
// 	}

// 	err = table.Put(tag).Run()
// 	if err != nil {
// 		return errors.Wrap(err, "put into the table")
// 	}

// 	return nil
// }

// func (c *Client) SavePartner(p *models.Partner) error {
// 	err := c.createTableIfNotExists("partners", &models.Partner{})
// 	if err != nil {
// 		return errors.Wrap(err, "create table")
// 	}

// 	table := c.db.Table("partners")

// 	err = table.Put(p).Run()
// 	if err != nil {
// 		return errors.Wrap(err, "put into the table")
// 	}

// 	return nil
// }

// func (c *Client) GetUsersWithTags(tags []string) ([]string, error) {
// 	err := c.createTableIfNotExists("taggedUsers", &models.TaggedUser{})
// 	if err != nil {
// 		return nil, errors.Wrap(err, "create table")
// 	}

// 	table := c.db.Table("taggedUsers")

// 	var items []string

// 	it := table.Scan().
// 		Filter("contains('Tags', ?)", tags).
// 		Iter()

// 	for {
// 		var taggedUser models.TaggedUser

// 		next := it.Next(&taggedUser)

// 		if it.Err() != nil {
// 			if errors.Is(it.Err(), dynamo.ErrNotFound) {
// 				return nil, nil
// 			}
// 			return nil, errors.Wrap(it.Err(), "scan table")
// 		}

// 		items = append(items, taggedUser.UserID)

// 		if !next {
// 			break
// 		}
// 	}

// 	return items, nil
// }

// func (c *Client) createTableIfNotExists(name string, from interface{}) error {
// 	err := c.db.CreateTable(name, from).Run()
// 	if err != nil {
// 		if awsErr, ok := err.(awserr.Error); ok {
// 			if awsErr.Code() == "ResourceInUseException" {
// 				// exists
// 				return nil
// 			}
// 		}
// 		return err
// 	}
// 	return nil
// }
