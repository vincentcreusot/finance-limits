# finance-limits
## Context
In finance, it's common for accounts to have so-called "velocity limits". 
Each attempt to load funds will come as a single-line JSON payload, structured as follows:

```json
{
  "id": "1234",
  "customer_id": "1234",
  "load_amount": "$123.45",
  "time": "2018-01-01T00:00:00Z"
}
```

Each customer is subject to three limits:

- A maximum of $5,000 can be loaded per day
- A maximum of $20,000 can be loaded per week
- A maximum of 3 loads can be performed per day, regardless of amount

The return is a json string for each load telling it's accepted or not.
```json
{ "id": "1234", "customer_id": "1234", "accepted": true }
```


## Usage
The binary takes 2 flags :
- -i or -inputFile with the file containing the list of loads to validate
- -o or -outputFile representing the file where to write the lines of validation
## Design
Reading the file uses channels, which help decouple logic from the utilities of reading the file itself. The logic package 
then takes a channel as parameter and reads that channel to look for lines to parse.
## Building
### Makefile
A Makefile is available to simplify building and development with the following targets :
- *build*: builds the package and creates a binary in the current directory
- *test*: runs test
- *coverage*: runs test with coverage report
- *vet*: runs go vet to find suspicious constructs
- *lint*: runs the linter to find some coding styles mistakes
- *format*: formats the code
### Dependencies
The project uses [go module](https://blog.golang.org/using-go-modules) to manage the dependencies. You can type
```bash
go get
```
to get missing dependencies.
