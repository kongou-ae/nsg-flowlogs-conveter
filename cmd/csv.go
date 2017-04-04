package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"time"
)

type Pt1h struct {
	Records []struct {
		Time          string `json:"time"`
		Systemid      string `json:"systemid"`
		Category      string `json:"category"`
		ResourceId    string `json:"resourceId"`
		OperationName string `json:"operationName`
		Properties    struct {
			Version float64 `json:"version"`
			Flows   []struct {
				Rule  string `json:"rule"`
				Flows []struct {
					Mac        string   `json:"mac`
					FlowTuples []string `json:"flowTuples"`
				} `json:"flows"`
			} `json:"flows"`
		} `json:"properties"`
	} `json:"records"`
}

// csvCmd represents the csv command
var csvCmd = &cobra.Command{
	Use:   "csv",
	Short: "convert to csv",
	Long:  `convert to csv`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Work your own magic here
		regex := regexp.MustCompile(`^([0-9]{10})`)
		file, err := ioutil.ReadFile("./PT1H.json")
		if err != nil {
			fmt.Println("flowlogs don't exist.")
			os.Exit(-1)
		}

		var Pt1h Pt1h

		err = json.Unmarshal(file, &Pt1h)
		if err != nil {
			fmt.Println("The format of PT1H.json is invalid.")
			os.Exit(-1)
		}

		for _, record := range Pt1h.Records {
			for _, flow1 := range record.Properties.Flows {
				for _, flow2 := range flow1.Flows {
					if len(flow2.FlowTuples) > 1 {
						for _, flowTuple := range flow2.FlowTuples {
							intUnix, _ := strconv.ParseInt(regex.FindStringSubmatch(flowTuple)[0], 10, 64)
							dateJst := time.Unix(intUnix, 0).In(time.FixedZone("JST", 9*60*60))
							fmt.Println(regex.ReplaceAllString(flowTuple, dateJst.String())) // + "," + flow1.Rule)
						}
					}
				}
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(csvCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// csvCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// csvCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
