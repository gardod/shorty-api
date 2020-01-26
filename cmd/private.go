package cmd

import (
	"github.com/gardod/shorty-api/internal/driver/http"
	"github.com/gardod/shorty-api/internal/handler/private"
	"github.com/spf13/cobra"
)

var privateCmd = &cobra.Command{
	Use:   "private",
	Short: "Serve private facing API",
	Run: func(cmd *cobra.Command, args []string) {
		http.Serve(private.GetRouter())
	},
}

func init() {
	rootCmd.AddCommand(privateCmd)
}
