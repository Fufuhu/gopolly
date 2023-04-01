/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"github.com/Fufuhu/gopolly/util/client/polly"
	"github.com/Fufuhu/gopolly/util/logging"
	polly2 "github.com/aws/aws-sdk-go-v2/service/polly"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"io"
	"os"
)

// synthesizeCmd represents the synthesize command
var synthesizeCmd = &cobra.Command{
	Use:   "synthesize",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: synthesize,
}

func init() {
	rootCmd.AddCommand(synthesizeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// synthesizeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// synthesizeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	synthesizeCmd.Flags().StringVarP(&filePath, "filePath", "f", "", "Transcript file to synthesize voice audio")
	err := synthesizeCmd.MarkFlagRequired("filePath")
	if err != nil {
		return
	}
	synthesizeCmd.Flags().StringVarP(&outputFilePath, "outputFilePath", "o", "", "Output file path")
	err = synthesizeCmd.MarkFlagRequired("outputFilePath")
	if err != nil {
		return
	}
}

var filePath string
var outputFilePath string

func synthesize(cmd *cobra.Command, args []string) {
	logger := logging.GetLogger()
	defer func(logger *zap.Logger) {
		_ = logger.Sync()
	}(logger)

	logger.Info(fmt.Sprintf("file path is %s", filePath))

	content, err := os.ReadFile(filePath)
	if err != nil {
		logger.Warn(err.Error())
		os.Exit(1)
	}
	contentString := string(content)

	logger.Info(contentString)

	client, err := polly.GetPollyClient(&polly.ClientConfig{})
	if err != nil {
		logger.Warn(err.Error())
		os.Exit(1)
	}

	ctx := context.Background()

	synthesizeSpeechInput := &polly2.SynthesizeSpeechInput{
		Text:         &contentString,
		OutputFormat: "mp3",
		Engine:       "standard",
		LanguageCode: "ja-JP",
		VoiceId:      "Mizuki",
	}

	synthesizeSpeechOutput, err := client.SynthesizeSpeech(ctx, synthesizeSpeechInput)
	if err != nil {
		logger.Warn(err.Error())
		os.Exit(1)
	}
	defer func(AudioStream io.ReadCloser) {
		_ = AudioStream.Close()
	}(synthesizeSpeechOutput.AudioStream)

	outputFile, err := os.Create(outputFilePath)
	defer func(outputFile *os.File) {
		_ = outputFile.Close()
	}(outputFile)

	if err != nil {
		logger.Warn(err.Error())
		os.Exit(1)
	}

	logger.Info(fmt.Sprintf("%d characters were synthesized", synthesizeSpeechOutput.RequestCharacters))
	written, err := io.Copy(outputFile, synthesizeSpeechOutput.AudioStream)
	if err != nil {
		logger.Warn(err.Error())
		return
	}

	logger.Info(fmt.Sprintf("%d bytes were written", written))
}
