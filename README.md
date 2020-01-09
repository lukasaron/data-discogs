# Data Discogs Parser

[![GoDoc](https://godoc.org/github.com/lukasaron/data-discogs?status.svg)](https://godoc.org/github.com/lukasaron/data-discogs)
[![Build Status](https://travis-ci.com/lukasaron/data-discogs.svg?branch=master)](https://travis-ci.com/lukasaron/data-discogs)

Data Discogs Parser introduces a way to categorize XML data dumps from Discogs: https://data.discogs.com. 
The library can parse provided XML files from Discogs and then save results into JSON or SQL file and lastly 
directly into SQL database. For this purpose three different writers were created.

The intention of using this parser is a library, which means there is not an executable part provided. The project has no other dependencies than the Golang language itself. 

## Writers

There are three existing writers supported by default: **JSON**, **SQL** and **DB**. 
More writers can be created by implementing the `writer` interface

### JSON Writer
As the name prompts this writer transforms input XML into JSON format. This JSON writer can be used as a solution that converts data into any NoSQL database.

### SQL Writer
The second supported writer creates SQL file with all necessary data from input in the form of insert commands within transactions. This file can be executed in any SQL database and the result will be populated table with the proper information.

### DB Writer
The last writer is for direct communication with the SQL database. All input data will be saved into appropriate tables immediately.
Before using this writer all required data tables need to be created. For that purpose please run the SQL script `sql_scripts/tables.sql`. 

You can also set up indexes on any column you want, to facilitate this process there is also a script for that `sql_scripts/indexes.sql`. 

To speed up a data transformation I would rather recommend creating indexes after the whole processing is completed.

## Examples of usage

Decoding of artists without using a writer can be done as easy like this: 
```go
package main
 
import (
    "fmt"
    "github.com/lukasaron/data-discogs"
    "log"
    "os"
)

func main() {
    f, err := os.Open("./data_samples/artists.xml")
    if err != nil {
        log.Fatal(err)
    }
    defer f.Close()

    d := discogs.NewXMLDecoder(f, nil)
    // decodes 10 artists by default, Block ItemSize can be changed via Options
    num, artists, err := d.Artists()
    fmt.Println(num, err, artists)
}
```
For the need of saving a result into output, there are writers. For instance, the SQL writer creates insert statements
and saves them into the output.
```go
package main

import (
    "fmt"
    "github.com/lukasaron/data-discogs"
    "github.com/lukasaron/data-discogs/write"
    "log"
    "os"
)

func main() {
    f, err := os.Open("./data_samples/artists.xml")
    if err != nil {
        log.Fatal(err)
    }
    defer f.Close()

    d := discogs.NewXMLDecoder(f,
        &discogs.Options{
            FileType: discogs.Artists,
            Block: discogs.Block{
            ItemSize: 20, // number of items processed at once
        },
    })

    o, _ := os.Create("./data_samples/artists.sql")
    defer o.Close()

    // for instance the SQL writer
    w := write.NewSQLWriter(o, nil)
    defer w.Close()

    err = d.Decode(w)
    fmt.Println(err)
}
```