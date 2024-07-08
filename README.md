# Command Line Tool: fl

This command line tool, `fl`, is designed to convert natural language descriptions of tasks into valid Unix commands.
It simplifies the process of executing tasks on a Unix command line by allowing users to describe their actions in plain language.

## Usage

To use `fl`, simply provide a natural language description of the task you want to perform, and the tool will generate the corresponding Unix command.
Here are some sample calls for using `fl`:

### Sample Calls

1. **Description**: List all files in the current directory.
   ```bash
   fl list all files in current directory
   ```

2. **Description**: Create a new directory named documents.
   ```bash
   fl create a new directory named documents
   ```

3. **Description**: Search for files containing a specific keyword in the src directory.
   ```bash
   fl search for files containing keyword in src directory
   ```

4. **Description**: Call an authenticated API and pass in some JSON data.
   ```bash
   fl call an api that returns JSON and sends some data {"foo":"bar"} as json where the api uses basic auth and the secret is an environment variable called API_KEY
   ```
