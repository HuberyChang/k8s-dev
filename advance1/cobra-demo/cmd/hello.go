/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"log"
	"time"

	"github.com/spf13/cobra"
)

var (
	Name *string
	Time *bool
)

// helloCmd represents the hello command
var helloCmd = &cobra.Command{
	Use:   "hello",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		str, _ := cmd.Flags().GetString("name")
		fmt.Printf("hello, %s!\n", str)
		t, _ := cmd.Flags().GetBool("time")
		if t {
			fmt.Println("Time:", time.Now().Format("2006-01-02 15:04:05"))
		}

		// fmt.Println("hello called")
	},
}

func init() {
	rootCmd.AddCommand(helloCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// helloCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// helloCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	// 这样的flag无法通过if判断是否要输出
	// 且对于所有子命令都有效
	helloCmd.PersistentFlags().String("name", "", "Say hello to someone")

	// 这种flag就可以通过if来判断是否输出了,默认是false
	// 在调用command的时候加上 --time 或 -t 就可以变为true
	// 这样的flag仅对当前命令有效
	helloCmd.Flags().BoolP("time", "t", false, "Add time info to hello")

	// 设置使用hello的时候后面必须加上name
	err := helloCmd.MarkPersistentFlagRequired("name")
	if err != nil {
		log.Fatal("Set flag required fail!")
	}
}
