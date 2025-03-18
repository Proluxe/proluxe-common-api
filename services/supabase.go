package services

import (
	"fmt"

	u "github.com/scottraio/go-utils"
	supabase "github.com/supabase-community/supabase-go"
)

func SupabaseClient() *supabase.Client {
	client, err := supabase.NewClient(u.GetDotEnvVariable("SUPABASE_URL"), u.GetDotEnvVariable("SUPABASE_KEY"), &supabase.ClientOptions{
		Schema: "next_auth",
	})

	if err != nil {
		fmt.Println("cannot initalize client", err)
	}

	return client
}
