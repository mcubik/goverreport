# goverreport 

Command line tool for coverage reporting and validation.

![travis-c](https://travis-ci.org/mcubik/goverreport.svg?branch=master)
[![Coverage Status](https://coveralls.io/repos/github/mcubik/goverreport/badge.svg?branch=master)](https://coveralls.io/github/mcubik/goverreport?branch=master)
[![Maintainability](https://api.codeclimate.com/v1/badges/957f7cab3e9d45db631d/maintainability)](https://codeclimate.com/github/mcubik/goverreport/maintainability)

## Installation

```
go get -u github.com/mcubik/goverreport
```

## Usage

`goverreport` reads a coverage profile and prints a report on the terminal. Optionally, it can also validate a coverage threshold.

```
Usage: goverreport [flags] -coverprofile=coverprofile.out

Flags:
  -coverprofile string
    	Coverage output file (default "coverage.out")
  -order string
    	Sort order: asc, desc (default "asc")
  -sort string
    	Column to sort by: filename, block, stmt, missing-blocks, missing-stmts (default "filename")
  -threshold float
    	Return an error code of 1 if the coverage is below a threshold
  -metric string
    	Use a specific metric for the threshold: block, stmt (default "block")
```

## Example

```
$ goverreport -sort=block -order=desc -threshold=85

+------------------+--------+---------+-------+---------+---------------+--------------+
|       FILE       | BLOCKS | MISSING | STMTS | MISSING | BLOCK COVER % | STMT COVER % |
+------------------+--------+---------+-------+---------+---------------+--------------+
| report/view.go   |      4 |       0 |     7 |       0 |        100.00 |       100.00 |
| report/report.go |     47 |       5 |    60 |       5 |         89.36 |        91.67 |
| main.go          |     30 |      10 |    44 |      15 |         66.67 |        65.91 |
+------------------+--------+---------+-------+---------+---------------+--------------+
|      TOTAL       |   81   |   15    |  111  |   20    |     81 48     |    81 98     |
+------------------+--------+---------+-------+---------+---------------+--------------+
exit status 1

```


## Configuration

You can use a fixed threshold by configuring it in the `.goverreport.yml` configuration file. This file also
lets you configure the root path of the project, so that it gets stripped from the names of the files, and a set of paths to be excluded from the report.

Here's an example:

```
threshold: 85
metric: stmt
root: "github.com/mcubik/goverreport"
exclusions: [test/it]
```
