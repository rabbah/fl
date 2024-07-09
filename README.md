# `fl`: command line aid powered by Postman Flows

This command line tool, `fl`, is designed to convert natural language descriptions of tasks into valid Unix commands.
It simplifies the process of executing tasks on a Unix command line by allowing users to describe their actions in plain language.

## Usage

To use `fl`, simply provide a natural language description of the task you want to perform, and the tool will generate the corresponding Unix command.
Here are some sample calls for using `fl`:

### Sample Calls

1. **Description**: Remove a directory and all its contents.
   ```sh
   fl remove a directory and all its contents
   ```
   Sample output:
   ```sh
   rm -r directory_name
   ```

2. **Description**: Search for files containing a specific keyword in the src directory.
   ```sh
   fl search for files containing keyword in src directory
   ```

   Sample output:
   ```sh
   grep -r "keyword" src
   ```

3. **Description**: Process a CSV file to extract a column and count unique occurences of a value.
   ```sh
   fl count the number of unique values that appear in the second column of a csv file, make sure the count is case insensitive, report the total count only
   ```
   Sample output:
   ```sh
   awk -F, '{print tolower($2)}' file.csv | sort -u | wc -l
   ```
   
4. **Description**: Call an authenticated API and pass in some JSON data.
   ```sh
   fl call an api that returns JSON and sends some data {"foo":"bar"} as json where the api uses basic auth and the secret is an environment variable called API_KEY
   ```

   Sample output:
   ```sh
   curl -X POST -H 'Content-Type: application/json' -H 'Authorization: Basic $API_KEY' -d '{"foo":"bar"}' https://api.example.com/endpoint
   ```

## Postman Flows

The entire backend for `fl` is implemented using [Postman Flows](https://learning.postman.com/docs/postman-flows/overview). Flows is a visual and low-code programming language for working with APIs and creating workflows with direct manipulation of APIs and data.

- TODO: embeddable flow

## Notes and References

This README Usage and Sample Calls sections were entirely generated with `fl` itself.
There are related command line tools for using LLMs with a CLI that I've found. There are surely others, so let me know what other interesting and related tools I should check out.

- `gencmd`
  - Generate commands with Generative AI [blog](https://sathishvj.medium.com/gencmd-generate-commands-with-generative-ai-d945b2a1d050)
  - `gencmd-bash` https://github.com/sathishvj/gencmd-bash and related Go implementation `gencmd-go` https://github.com/sathishvj/gencmd-go
- `ShellGPT` https://github.com/TheR1D/shell_gpt
