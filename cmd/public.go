package cmd

import (
	"github.com/gardod/shorty-api/cmd/public"
	"github.com/spf13/cobra"
)

var publicCmd = &cobra.Command{
	Use:   "public",
	Short: "Serve public facing API",
	Run: func(cmd *cobra.Command, args []string) {
		public.Serve()
	},
}

func init() {
	rootCmd.AddCommand(publicCmd)
}
