package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

var rootCmd = &cobra.Command{
	Use:   "gdtg",
	Short: "gdtg is a utility to find locally stored Discord tokens",
	Long:  `A Fast and Flexible token finder`,
}

var searchCmd = &cobra.Command{
	Use:   "search [platform]",
	Short: "Search for tokens in specified platform",
	Long: `Search for tokens in a specified platform. The platform can be one of the following:

- Discord
- Discord Canary

Browsers are coming soon.
You can also use the keyword "all" to search all platforms.

Example usage:

gdtg search Discord
gdtg search all`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		usingOS := runtime.GOOS

		if usingOS == "windows" || usingOS == "darwin" {
			fmt.Println("Windows /macOS currently not supported")
		}

		home, err := os.UserHomeDir()
		if err != nil {
			log.Fatal(err.Error())
		}

		configFolder := home + "/.config/"
		paths := map[string]string{
			"Discord":        configFolder + "discord",
			"Discord Canary": configFolder + "discordcanary",
			//TODO how to handle profiles for browsers
			/* "Google Chrome":  configFolder + "google-chrome",
			"Brave":         configFolder + "BraveSoftware/Brave-Browser",
			"Brave Nightly": configFolder + "BraveSoftware/Brave-Browser-Nightly", */
		}

		if args[0] == "all" {
			tokens, err := getTokens(paths)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println("----------------------------------------------------")
			for key, tokensForKey := range tokens {
				fmt.Printf("Found %d token for %s\n", len(tokensForKey), key)
				fmt.Println(tokensForKey)
				fmt.Println("----------------------------------------------------")
			}
		} else {
			path, exists := paths[args[0]]
			if !exists {
				fmt.Printf("No such key: %s\n", args[0])
				return
			}

			tokens, err := getTokens(map[string]string{args[0]: path})
			if err != nil {
				fmt.Println(err)
				return
			}

			fmt.Printf("%s: %v\n", args[0], tokens[args[0]])
		}
	},
}

func main() {
	rootCmd.AddCommand(searchCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func contains(slice []string, str string) bool {
	for _, v := range slice {
		if v == str {
			return true
		}
	}
	return false
}

func getTokens(paths map[string]string) (map[string][]string, error) {
	var tokens = make(map[string][]string)

	fmt.Println("Searching tokens...")

	for key, path := range paths {
		path = filepath.Join(path, "Local Storage", "leveldb")
		var tokensForPath []string

		files, err := os.ReadDir(path)
		if err != nil {
			return nil, err
		}

		for _, file := range files {
			fileName := file.Name()
			if !strings.HasSuffix(fileName, ".log") && !strings.HasSuffix(fileName, ".ldb") {
				continue
			}

			content, err := os.ReadFile(filepath.Join(path, fileName))
			if err != nil {
				return nil, err
			}

			lines := strings.Split(string(content), "\n")

			for _, line := range lines {
				line = strings.TrimSpace(line)
				if line == "" {
					continue
				}

				for _, regex := range []string{`[\w-]{24}\.[\w-]{6}\.[\w-]{27}`, `mfa\.[\w-]{84}`} {
					r := regexp.MustCompile(regex)
					matches := r.FindAllString(line, -1)
					for _, match := range matches {
						if !contains(tokensForPath, match) {
							tokensForPath = append(tokensForPath, match)
						}
					}
				}
			}
		}

		tokens[key] = tokensForPath
	}

	return tokens, nil
}
