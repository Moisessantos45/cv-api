package config

import (
	"log"

	"github.com/supabase-community/supabase-go"
)

var Client *supabase.Client

func Init(url, key string) error {
	neClient, err := supabase.NewClient(url, key, &supabase.ClientOptions{})
	if err != nil {
		log.Printf("Error initializing Supabase client: %v", err)
		return err
	}

	Client = neClient

	log.Println("Supabase client initialized successfully")
	return nil
}
