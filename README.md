# JSON COMPARE QUIQ TAKE HOME

## Usage

`JsonCompare -h` will provide the basics of how to use the tool.

Here is the output of `JsonCompare -h`

```
		- Cli tool to assist with Json Comparing

Usage:
  JsonCompare [path to json...] [path to json...] [flags]

Flags:
  -h, --help      help for JsonCompare
      --json      output result in json format
  -v, --verbose   output found differences between json files and not just score
```

## Json Output Format

```json
{
  "TotalBytes": 2207,
  "MismatchBytes": 844,
  "Score": 0.6175804,
  "File1Mismatch": [
    {
      "Path": [
        "root",
        "storeName"
      ],
      "MisMatch": "The Literary No...  :Infinite Pages ...  ",
      "Description": "String Values are not equal",
      "Error": "ValueError"
    },
    ...
  ],

  "File2Mismatch": [
    {
      "Path": [
        "root",
        "storeName"
      ],
      "MisMatch": "Infinite Pages ...  :The Literary No...  ",
      "Description": "String Values are not equal",
      "Error": "ValueError"
    },
    ...

  ],
{
```

- TotalBytes - is the sum of non-structural bytes in each file
- MismatchBytes - is the sum of Mismatched non-structural bytes in each file
- Score - is `(1 - (MismatchBytes / TotalBytes))`

File1Mismatch: 
This shows everything that is in file 1 but not in file 2 (file numbers are based on argument location)
File2Mismatch: This shows everything in file 2 and not in file 1

## How is JSON Compared
This tool takes two json files and returns a score based on the number of matching non-structural bytes in the JSON.
A non-structural byte are the strings, numbers and booleans that make up the data in the JSON. Mismatched byte are any
bytes that are apart of or below an initial disagreement:

Example

```json
//file 1
{
  "Hello" : "bye"
}

//file 2
{
  "Ciao" : "bye"
}
```

Ciao does not exist in file 1 so "Ciao" and "bye" are counted as mismatched bytes.

Hello does not exist in file 2 so "Hello" and "bye" are counted as mismatched bytes.

Score in this case would be 0.

### Json Arrays Special Behavior

Unlike Keys Json arrays are compared positionally. This was decided because [1,2] != [2,1] in most programming situations
and it was simpler to understand.

## Requirements

- go 1.25+
## Install 

### Project Root Install

====  Consideration ====

- Must have GOPATH added to Path env var. (https://go.dev/wiki/SettingGOPATH)

To install the tool to work like a classic Cli tool use:
```
go install .
```
In the project root and then you can run this in the command line!

```
JsonCompare -h
```

## Run

### Run from built binary

In project root run
```
go build .
```

Then run

```
./main -h
```

### Run with `go run`

In project root directory

```
go run main.go -h
```