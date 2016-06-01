// Copyright © 2016 Bernardo Heynemann <heynemann@gmail.com>
// This file is part of hk.

package cmd

import (
	"fmt"
	"math"
	"os"
	"time"

	"github.com/aybabtme/uniplot/histogram"
	"github.com/gosuri/uiprogress"
	"github.com/gosuri/uiprogress/util/strutil"
	"github.com/heynemann/hk/core"
	"github.com/spf13/cobra"
)

var command string
var producers int
var scripts int
var workers int

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "hk",
	Short: "hk is a benchmark tool that executes producers",
	Long: `Hollywood Killer (HK) is a tool that takes a
producer and executes it many times.

A producer is just any script that can output a json output to StdOut.
`,
	Run: func(cmd *cobra.Command, args []string) {
		startTime := time.Now().UnixNano()
		uiprogress.Start()                              // start rendering
		bar := uiprogress.AddBar(producers*scripts - 1) // Add a new bar

		// optionally, append and prepend completion and elapsed time
		bar.AppendCompleted()
		bar.PrependElapsed()

		// prepend the deploy step to the bar
		bar.PrependFunc(func(b *uiprogress.Bar) string {
			text := fmt.Sprintf("%d/%d", b.Current()+1, producers*scripts)
			return strutil.Resize(text, uint(len(text)))
		})

		results, _ := core.Run(command, producers, scripts, workers, bar.Incr)

		fmt.Println()
		fmt.Printf("%d scripts executed in %.2fs\n", producers*scripts, float64(time.Now().UnixNano()-startTime)/1000000/1000)

		fmt.Println()

		fmt.Println("Histogram for executed producers")
		fmt.Println("--------------------------------")

		bins := 10
		data := make([]float64, producers*scripts)

		for i := 0; i < producers; i++ {
			for j := 0; j < scripts; j++ {
				index := i*scripts + j
				res := results[index]
				data[index] = (res.EndDate - res.StartDate)
			}
		}

		hist := histogram.Hist(bins, data)

		maxWidth := 80
		histogram.Fprintf(os.Stdout, hist, histogram.Linear(maxWidth), func(v float64) string {
			return fmt.Sprintf("%.4gms", v)
		})

		fmt.Println()

		fmt.Println("Percentiles")
		fmt.Println("-----------")

		fmt.Printf("90th percentile: %.2f\n", getPerc(hist, 90, producers*scripts))
		fmt.Printf("99th percentile: %.2f\n", getPerc(hist, 99, producers*scripts))
		fmt.Printf("99.9th percentile: %.2f\n", getPerc(hist, 99.9, producers*scripts))
		fmt.Println()
	},
}

func getPerc(hist histogram.Histogram, perc float64, total int) float64 {
	numberOfItems := int(math.Floor(float64(total) * perc / 100))
	count := 0
	for _, bkt := range hist.Buckets {
		count += bkt.Count
		if count >= numberOfItems {
			return bkt.Max
		}
	}
	return hist.Buckets[len(hist.Buckets)-1].Max
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags, which, if defined here,
	// will be global for your application.

	//RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.hk.yaml)")
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	RootCmd.Flags().StringVarP(&command, "command", "c", "", "Producer command to execute")
	RootCmd.Flags().IntVarP(&producers, "producers", "p", 100, "Number of producers to execute")
	RootCmd.Flags().IntVarP(&scripts, "scripts", "s", 20, "Number of scripts to burn")
	RootCmd.Flags().IntVarP(&workers, "workers", "w", 50, "Number of workers to execute")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	//if cfgFile != "" { // enable ability to specify config file via flag
	//viper.SetConfigFile(cfgFile)
	//}

	//viper.SetConfigName(".hk") // name of config file (without extension)
	//viper.AddConfigPath("$HOME")  // adding home directory as first search path
	//viper.AutomaticEnv()          // read in environment variables that match

	// If a config file is found, read it in.
	//if err := viper.ReadInConfig(); err == nil {
	//fmt.Println("Using config file:", viper.ConfigFileUsed())
	//}
}
