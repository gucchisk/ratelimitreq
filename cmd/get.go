/*
Copyright © 2024 gucchisk <gucchi_sk@yahoo.co.jp>
*/
package cmd

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/gucchisk/ratelimitreq/pkg/client"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		listFile := args[0]
		command := args[1]

		fp, err := os.Open(listFile)
		if err != nil {
			fmt.Printf("%v", err)
			return
		}
		defer fp.Close()
		client := client.NewClient(0.5, map[string]float64{"www.yahoo.co.jp": 0.1, "google.co.jp": 0.1})
		// client := client.NewClient(0.5, map[string]float64{"www.yahoo.co.jp": 0.1})
		// client := client.NewClientByRateLimit(0.5)
		scanner := bufio.NewScanner(fp)
		for scanner.Scan() {
			line := scanner.Text()
			fmt.Println(line)
			req, _ := http.NewRequest("GET", line, nil)
			resp, _ := client.Do(req)
			defer resp.Body.Close()
			fmt.Printf("%s %s\n", resp.Status, time.Now().Local())
			cmds := strings.Split(command, " ")
			c := exec.Command(cmds[0], cmds[1:]...)
			c.Stdin = resp.Body
			r, e := c.CombinedOutput()
			if e != nil {
				fmt.Printf("%v", e)
				return
			}
			fmt.Printf("%s\n", r)
		}
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
}
