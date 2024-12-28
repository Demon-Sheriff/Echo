package cmd

import (
    "fmt"
    "github.com/spf13/cobra"
    "os"
	"Echo/chat"
)

var rootCmd = &cobra.Command{
    Use:   "echo",
    Short: "echo away!",
    Long:  "echo is a cli tool just echo and recieve echoes",
    Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("The CLI is starting !")
		for {		
			cI := chat.InitiateChatInterface()
			quitting, _ := cI.Run();

			if quitting {
				break
			}
		}
    }, 
}

func Execute() {
    if err := rootCmd.Execute(); err != nil {
        fmt.Fprintf(os.Stderr, "Oops. An error while executing Echo '%s'\n", err)
        os.Exit(1)
    }
}