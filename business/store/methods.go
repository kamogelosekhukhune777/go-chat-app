package store

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/kamogelosekhukhune777/go-chat-app/business/chat"
)

// Cache defines the contract for any chat storage or caching layer.
//
// This allows the business logic to remain agnostic of the underlying
// persistence mechanism (Redis, DynamoDB, in-memory, etc.).
//
// Any implementation of Cache should provide methods to:
//
//   - Maintain a user's contact list (UpdateContactList, FetchContactList)
//   - Create and store chat messages (CreateChat)
//   - Fetch chat messages between two users in a given time window (FetchChatBetween)
type Cache interface {
	// UpdateContactList adds a contact to a user's contact list if it doesn't
	// exist, or updates the timestamp of last interaction if it does.
	UpdateContactList(username, contact string) error

	// CreateChat stores a new chat message.
	// It should also ensure both participants have each other added to their contact list.
	//
	// Returns:
	//   - a unique identifier for the stored chat (string)
	//   - error if the operation fails
	CreateChat(c *chat.Chat) (string, error)

	// FetchChatBetween retrieves all chats exchanged between username1 and username2
	// within the given timestamp range [fromTS, toTS].
	//
	// Parameters:
	//   - username1, username2: identifiers of chat participants
	//   - fromTS, toTS: timestamp range (inclusive) in string form
	//
	// Returns:
	//   - slice of chat.Chat objects
	//   - error if retrieval fails
	FetchChatBetween(username1, username2, fromTS, toTS string) ([]chat.Chat, error)

	// FetchContactList retrieves the contact list for the given user,
	// sorted by most recent activity.
	//
	// Returns:
	//   - slice of chat.ContactList (username + last activity timestamp)
	//   - error if retrieval fails
	FetchContactList(username string) ([]chat.ContactList, error)
}

// UpdateContactList adds or updates the timestamp of a contact in a user’s contact list.
//
// Internally uses a Redis ZSET where:
//   - Key: "contacts:{username}"
//   - Member: contact's username
//   - Score: last interaction timestamp (Unix time)
func (r *RedisCache) UpdateContactList(username, contact string) error {
	zs := &redis.Z{Score: float64(time.Now().Unix()), Member: contact}
	return r.client.ZAdd(context.Background(), contactListZKey(username), zs).Err()
}

// CreateChat stores a chat message in RedisJSON and updates both users' contact lists.
//
// Storage pattern:
//   - Key: "chat#{timestamp}"
//   - Value: JSON representation of chat.Chat
//
// It also calls UpdateContactList for both participants.
func (r *RedisCache) CreateChat(c *chat.Chat) (string, error) {
	chatKey := chatKey()
	by, _ := json.Marshal(c)

	_, err := r.client.Do(context.Background(), "JSON.SET", chatKey, "$", string(by)).Result()
	if err != nil {
		log.Println("error while setting chat json", err)
		return "", err
	}

	_ = r.UpdateContactList(c.From, c.To)
	_ = r.UpdateContactList(c.To, c.From)

	return chatKey, nil
}

// CreateFetchChatBetweenIndex sets up a RediSearch index for querying chats.
//
// This method is Redis-specific and not part of the generic store.Cache interface.
// It should be called once (e.g. during application startup) to create the index
// used in FetchChatBetween queries.
//
// Redis CLI equivalent:
//
//	FT.CREATE idx#chats ON JSON PREFIX 1 chat#
//	  SCHEMA $.from AS from TAG
//	         $.to AS to TAG
//	         $.timestamp AS timestamp NUMERIC SORTABLE
func (r *RedisCache) CreateFetchChatBetweenIndex() error {
	res, err := r.client.Do(context.Background(),
		"FT.CREATE",
		chatIndex(),
		"ON", "JSON",
		"PREFIX", "1", "chat#",
		"SCHEMA", "$.from", "AS", "from", "TAG",
		"$.to", "AS", "to", "TAG",
		"$.timestamp", "AS", "timestamp", "NUMERIC", "SORTABLE",
	).Result()

	if err != nil {
		log.Println("error creating chat index:", err)
		return err
	}

	log.Println("chat index created:", res)
	return nil
}

// FetchChatBetween retrieves all chats exchanged between two users in a given timestamp range.
//
// Uses RediSearch to query by `from`, `to`, and `timestamp` fields.
// Results are sorted in descending order by timestamp.
func (r *RedisCache) FetchChatBetween(username1, username2, fromTS, toTS string) ([]chat.Chat, error) {
	query := fmt.Sprintf("@from:{%s|%s} @to:{%s|%s} @timestamp:[%s %s]",
		username1, username2, username1, username2, fromTS, toTS)

	res, err := r.client.Do(context.Background(),
		"FT.SEARCH", chatIndex(), query, "SORTBY", "timestamp", "DESC").Result()
	if err != nil {
		return nil, err
	}

	data := Deserialise(res)

	return DeserialiseChat(data), nil
}

// FetchContactList retrieves a user’s contact list, sorted by most recent activity.
//
// Internally uses a Redis ZSET and returns a list of ContactList structs
// (username + last interaction timestamp).
func (r *RedisCache) FetchContactList(username string) ([]chat.ContactList, error) {
	zRangeArg := redis.ZRangeArgs{
		Key:   contactListZKey(username),
		Start: 0,
		Stop:  -1,
		Rev:   true, // newest first
	}

	res, err := r.client.ZRangeArgsWithScores(context.Background(), zRangeArg).Result()
	if err != nil {
		return nil, err
	}

	return DeserialiseContactList(res), nil
}
