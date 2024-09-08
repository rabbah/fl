package examples

import (
	"fl/utils"
	"fmt"

	markdown "github.com/MichaelMure/go-term-markdown"
)

type commandDescription struct {
	Description string
	Command     string
	Output      string
}

func formatExamples(commands []commandDescription) string {
	examples := ""
	for _, c := range commands {
		examples += "\n"
		examples += fmt.Sprintf("**Description**: %s\n\n", c.Description)
		examples += fmt.Sprintf("    %s\n\n", c.Command)
		examples += fmt.Sprintf("    Sample output:\n\n    %s\n", c.Output)
	}
	return examples
}

var commands = []commandDescription{
	{
		Description: "Search for files containing a specific keyword in the src directory.",
		Command:     "fl search for files containing keyword in src directory",
		Output:      "grep -r \"keyword\" src",
	},
	{
		Description: "Process a CSV file to extract a column and count unique occurrences of a value.",
		Command:     "fl count the number of unique values that appear in the second column of a csv file, make sure the count is case insensitive, report the total count only",
		Output:      "awk -F, '{print tolower($2)}' file.csv | sort -u | wc -l",
	},
	{
		Description: "Change contents of a file to uppercase.",
		Command:     "fl change the contents of a file to uppercase and save the results back to the same file",
		Output:      "tr '[:lower:]' '[:upper:]' < file.txt > temp.txt && mv temp.txt file.txt",
	},
	{
		Description: "Generate a command to find all files with a specific extension in a directory.",
		Command:     "fl find all files ending in .ext",
		Output:      "find . -type f -name '*.ext'",
	},
	{
		Description: "Generate a command to find and replace a string in multiple files.",
		Command:     "fl find and replace \"old\" with \"new\" in multiple files",
		Output:      "find . -type f -exec sed -i 's/old/new/g' {} +",
	},
	{
		Description: "Call an authenticated API and pass in some JSON data.",
		Command:     "fl call an api that returns JSON and sends some data {\"foo\":\"bar\"} as json where the api uses basic auth and the secret is an environment variable called API_KEY",
		Output:      "curl -X POST -H 'Content-Type: application/json' -H 'Authorization: Basic $API_KEY' -d '{\"foo\":\"bar\"}' https://api.example.com/endpoint",
	},
}

func Show() {
	source := "`fl` is a command line tool that converts natural language descriptions of tasks you want to complete in your terminal into valid Unix commands."
	source += " Here are some examples.\n"
	result := markdown.Render(source, 80, 0)
	fmt.Println(string(result))

	examples := formatExamples(commands[:3])
	result = markdown.Render(examples, 80, 0)
	fmt.Println(string(result))

	showMoreExamples := utils.PromptYesNo("Do you want to see more examples?")
	if showMoreExamples {
		fmt.Println()
		examples := formatExamples(commands[3:])
		result = markdown.Render(examples, 160, 0)
		fmt.Println(string(result))
	}
}
