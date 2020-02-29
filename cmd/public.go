package cmd

import (
	"github.com/gardod/shorty-api/internal/driver/http"
	"github.com/gardod/shorty-api/internal/handler/public"

	"github.com/spf13/cobra"
)

var publicCmd = &cobra.Command{
	Use:   "public",
	Short: "Serve public facing API",
	Run: func(cmd *cobra.Command, args []string) {
		http.Serve(public.GetRouter())
	},
}

func init() {
	rootCmd.AddCommand(publicCmd)
}
