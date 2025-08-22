package store

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/kamogelosekhukhune777/go-chat-app/business/chat"
)

// chatKey generates a unique Redis key for storing a chat.
// Format: chat#{current_unix_milliseconds}
func chatKey() string {
	return fmt.Sprintf("chat#%d", time.Now().UnixMilli())
}

// chatIndex returns the RedisSearch index name for chats.
func chatIndex() string {
	return "idx#chats"
}

// contactListZKey returns the key for a user’s contact list ZSET.
// Format: contacts:{username}
func contactListZKey(username string) string {
	return "contacts:" + username
}

func userSetKey() string {
	return "users"
}

func sessionKey(client string) string {
	return "session#" + client
}

// =========================================================================================================================
type Document struct {
	ID      string `json:"_id"`
	Payload []byte `json:"payload"`
	Total   int64  `json:"total"`
}

func Deserialise(res any) []Document {
	switch v := res.(type) {
	case []any:
		if len(v) > 1 {
			total := len(v) - 1
			var docs = make([]Document, 0, total/2)

			for i := 1; i <= total; i = i + 2 {
				arrOfValues := v[i+1].([]any)
				value := arrOfValues[len(arrOfValues)-1].(string)

				// add _id in the response
				doc := Document{
					ID:      v[i].(string),
					Payload: []byte(value),
					Total:   v[0].(int64),
				}

				docs = append(docs, doc)
			}
			return docs
		}
	default:
		log.Printf("different response type otherthan []interface{}. type: %T", res)
		return nil
	}

	return nil
}

func DeserialiseChat(docs []Document) []chat.Chat {
	chats := []chat.Chat{}
	for _, doc := range docs {
		var c chat.Chat
		json.Unmarshal(doc.Payload, &c)

		c.ID = doc.ID
		chats = append(chats, c)
	}

	return chats
}

func DeserialiseContactList(contacts []redis.Z) []chat.ContactList {
	contactList := make([]chat.ContactList, 0, len(contacts))

	// improvement tip: use switch to get type of contact.Member
	// handle unknown type accordingly
	for _, contact := range contacts {
		contactList = append(contactList, chat.ContactList{
			Username:     contact.Member.(string),
			LastActivity: int64(contact.Score),
		})
	}

	return contactList
}
