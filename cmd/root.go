package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func BuildRootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "casino",
		Short: "Bot de discord degenerado",
	}

	dburl := ""
	rootCmd.PersistentFlags().StringVar(&dburl, "database", "./casino.db", "Path to database file")

	viper.BindPFlag("database", rootCmd.PersistentFlags().Lookup("database"))

	return rootCmd
}
