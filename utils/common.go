package utils

import (
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

// hashing a password
func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}

	return string(bytes)
}

// compares a hashed password
func VerifyPassword(userPassword string, providedPassword string) (bool, string) {

	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))
	check := true
	msg := ""

	if err != nil {
		msg = "login or password is incorrect"
		check = false
	}
	return check, msg
}

// converts a struct to a map using BSON.
func StructToMap(input interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	bsonBytes, err := bson.Marshal(input)
	if err == nil {
		bson.Unmarshal(bsonBytes, &result)
	}
	return result
}

// to delete key from JSON
func DeleteKey(m map[string]interface{}, keyToDelete string) {
	// Check if the key exists in the map
	if _, ok := m[keyToDelete]; ok {
		newMap := make(map[string]interface{})
		for k, v := range m {
			if k != keyToDelete {
				newMap[k] = v
			}
		}

		for k := range m {
			delete(m, k)
		}
		// Copy values from the new map to the original map
		for k, v := range newMap {
			m[k] = v
		}
	}
}
