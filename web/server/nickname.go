package server

import (
	"math/rand"
	"strings"
)

// Adjectives for random nickname generation - fun and friendly
var adjectives = []string{
	"Adventurous", "Brave", "Cheerful", "Daring", "Eager",
	"Fearless", "Gentle", "Happy", "Inventive", "Jolly",
	"Keen", "Lively", "Mighty", "Noble", "Optimistic",
	"Playful", "Quick", "Radiant", "Speedy", "Trusty",
	"Unique", "Valiant", "Witty", "Zany", "Zippy",
	"Cosmic", "Dazzling", "Electric", "Fluffy", "Groovy",
	"Hyper", "Jazzy", "Luminous", "Majestic", "Nimble",
	"Plucky", "Quirky", "Rustic", "Snappy", "Turbo",
	"Wacky", "Bouncy", "Crafty", "Dizzy", "Fancy",
	"Goofy", "Hasty", "Icy", "Jumpy", "Kooky",
}

// Animals for random nickname generation
var animals = []string{
	"Aardvark", "Badger", "Capybara", "Dolphin", "Elephant",
	"Falcon", "Giraffe", "Hedgehog", "Iguana", "Jaguar",
	"Koala", "Lemur", "Mongoose", "Narwhal", "Octopus",
	"Pangolin", "Quokka", "Raccoon", "Sloth", "Tapir",
	"Unicorn", "Vulture", "Wombat", "Xerus", "Yak",
	"Zebra", "Axolotl", "Bison", "Chinchilla", "Dingo",
	"Emu", "Flamingo", "Gecko", "Hamster", "Ibex",
	"Jackal", "Kiwi", "Llama", "Manatee", "Newt",
	"Ocelot", "Panda", "Quail", "Rhino", "Starfish",
	"Toucan", "Urchin", "Viper", "Walrus", "Yeti",
}

// GenerateRandomNickname creates a fun random nickname like "HappyPanda" or "SwiftFalcon"
func GenerateRandomNickname() string {
	adjective := adjectives[rand.Intn(len(adjectives))]
	animal := animals[rand.Intn(len(animals))]
	return adjective + animal
}

// GenerateRandomNicknameWithSpace creates a nickname with space like "Happy Panda"
func GenerateRandomNicknameWithSpace() string {
	adjective := adjectives[rand.Intn(len(adjectives))]
	animal := animals[rand.Intn(len(animals))]
	return adjective + " " + animal
}

// IsValidNickname checks if a nickname is valid (not empty, reasonable length)
func IsValidNickname(nickname string) bool {
	trimmed := strings.TrimSpace(nickname)
	return len(trimmed) >= 2 && len(trimmed) <= 30
}
