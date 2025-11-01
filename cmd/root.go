package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/aidankeefe2022/JsonCompare/cmd/JsonCompare"

	"github.com/spf13/cobra"
)

var rootCommand = &cobra.Command{
	Use:   "JsonCompare [path to json...] [path to json...] [flags] ",
	Short: "Cli tool to assist with Json Comparing",
	Long: ` 
    /$$$$$                                /$$$$$$                                                                 
   |__  $$                               /$$__  $$                                                                
      | $$  /$$$$$$$  /$$$$$$  /$$$$$$$ | $$  \__/  /$$$$$$  /$$$$$$/$$$$   /$$$$$$   /$$$$$$   /$$$$$$   /$$$$$$ 
      | $$ /$$_____/ /$$__  $$| $$__  $$| $$       /$$__  $$| $$_  $$_  $$ /$$__  $$ |____  $$ /$$__  $$ /$$__  $$
 /$$  | $$|  $$$$$$ | $$  \ $$| $$  \ $$| $$      | $$  \ $$| $$ \ $$ \ $$| $$  \ $$  /$$$$$$$| $$  \__/| $$$$$$$$
| $$  | $$ \____  $$| $$  | $$| $$  | $$| $$    $$| $$  | $$| $$ | $$ | $$| $$  | $$ /$$__  $$| $$      | $$_____/
|  $$$$$$/ /$$$$$$$/|  $$$$$$/| $$  | $$|  $$$$$$/|  $$$$$$/| $$ | $$ | $$| $$$$$$$/|  $$$$$$$| $$      |  $$$$$$$
 \______/ |_______/  \______/ |__/  |__/ \______/  \______/ |__/ |__/ |__/| $$____/  \_______/|__/       \_______/
                                                                          | $$                                    
                                                                          | $$                                    
                                                                          |__/ 

		- Cli tool to assist with Json Comparing
`,
	Args: cobra.MinimumNArgs(2),

	Run: func(cmd *cobra.Command, args []string) {
		jsonFlag, _ := cmd.Flags().GetBool("json")
		verbose, _ := cmd.Flags().GetBool("verbose")
		filePath1 := args[0]
		filePath2 := args[1]
		output := JsonCompare.CompareFiles(filePath1, filePath2)
		if jsonFlag && verbose {
			pretty, _ := json.MarshalIndent(output, "", "  ")
			fmt.Println(string(pretty))
		} else if jsonFlag {
			newMap := map[string]interface{}{
				"Score":           output.Score,
				"MisMatchedBytes": output.MismatchBytes,
				"TotalBytes":      output.TotalBytes,
			}
			pretty, _ := json.MarshalIndent(newMap, "", "  ")
			fmt.Println(string(pretty))
		} else {
			printOutput(output, verbose)
		}
	},
}

func init() {
	rootCommand.Flags().Bool("json", false, "output result in json format")
	rootCommand.Flags().BoolP("verbose", "v", false, "output found differences between json files and not just score")
}

func Execute() {
	if err := rootCommand.Execute(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error running the tool! %v\n", err)
		os.Exit(1)
	}
}

func printOutput(o JsonCompare.Output, verbose bool) {
	fmt.Println("=== Output Summary ===")
	fmt.Printf("Total Bytes:          %.2f\n", o.TotalBytes)
	fmt.Printf("Mismatch Bytes:       %.2f\n", o.MismatchBytes)
	fmt.Printf("Similarity Score:     %.3f\n", o.Score)

	if !verbose {
		return
	}

	fmt.Println("\n--- File1 Mismatches ---")
	printSnapshots(o.File1Mismatch)

	fmt.Println("\n--- File2 Mismatches ---")
	printSnapshots(o.File2Mismatch)

}

func printSnapshots(snaps []JsonCompare.SnapShot) {
	if len(snaps) == 0 {
		fmt.Println("  (none)")
		return
	}
	for i, s := range snaps {
		fmt.Printf("  %d. Path: %s\n", i+1, formatPath(s.Path))
		fmt.Printf("     Mismatch: %s\n", s.MisMatch)
	}
}

func formatPath(p []string) string {
	if len(p) == 0 {
		return "(empty)"
	}
	return "=>" + fmt.Sprint(p)[1:len(fmt.Sprint(p))-1] // make it look like a path
}
