package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var registerCmd = &cobra.Command{
	Use: "register",
	Short: "Create an account for mini-paas",
	RunE: func(cmd *cobra.Command, args []string) error {
		email, _ := cmd.Flags().GetString("email")
		username, _ := cmd.Flags().GetString("username")
		password, _ := cmd.Flags().GetString("password")

		payload := map[string]string{
			"email": email,
			"username": username,
			"password": password,
		}
		
		resp, err := SendHttpRequest("POST", "/api/user/create", payload, false)

		if err != nil {
			return err
		}

		fmt.Println("Registration Response: ")
		fmt.Println(string(resp))

		return nil
	},
}

var loginCmd = &cobra.Command{
	Use: "login",
	Short: "Login Into mini-paas",
	RunE: func(cmd *cobra.Command, args []string) error {
		email, _ := cmd.Flags().GetString("email")
		password, _ := cmd.Flags().GetString("password")

		payload := map[string]string{
			"email": email,
			"password": password,
		}

		resp, err := SendHttpRequest("POST", "/api/user/login", payload, false)

		if err != nil {
			return err
		}

		fmt.Println("Login Response: ")
		fmt.Println(string(resp))

		return nil
	},
}

func init() {
	rootCmd.AddCommand(registerCmd)
	rootCmd.AddCommand(loginCmd)

	registerCmd.Flags().StringP("email", "e", "", "User Email")
	registerCmd.Flags().StringP("username", "u", "", "Username")
	registerCmd.Flags().StringP("password", "p", "", "password")
	registerCmd.MarkFlagRequired("email")
	registerCmd.MarkFlagRequired("username")
	registerCmd.MarkFlagRequired("password")

	loginCmd.Flags().StringP("email", "e", "", "User Email")
	loginCmd.Flags().StringP("password", "p", "", "password")
	loginCmd.MarkFlagRequired("email")
	loginCmd.MarkFlagRequired("password")
}