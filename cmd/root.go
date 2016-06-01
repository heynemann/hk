// Copyright © 2016 Bernardo Heynemann <heynemann@gmail.com>
// This file is part of hk.

package cmd

import (
	"fmt"
	"os"

	"github.com/aybabtme/uniplot/histogram"
	"github.com/heynemann/hk/core"
	"github.com/sethgrid/curse"
	"github.com/sethgrid/multibar"
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
		fmt.Println()

		progressBars, _ := multibar.New()
		bar := progressBars.MakeBar(producers*scripts, "Total Scripts")
		progressBars.Println("Executing producers...")
		go progressBars.Listen()

		results, _ := core.Run(command, producers, scripts, workers, bar)

		fmt.Println()
		fmt.Println()

		fmt.Printf("%c0m", curse.ESC)

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
		histogram.Fprint(os.Stdout, hist, histogram.Linear(maxWidth))
		//fmt.Println(results)
		//fmt.Println(errors)
	},
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
