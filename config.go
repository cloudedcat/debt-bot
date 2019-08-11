package main

type Config struct {
	BotToken string
	DBName   string
}

var config = Config{
	BotToken: "", // set up Bot token
	DBName:   "", // set up database mame
}
