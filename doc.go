// Discogs Parser introduces a way to categorize XML data dumps from Discogs [https://data.discogs.com/].
//
// The intention of using this parser is a library, which means there is not an executable part provided.
// The project has no other dependencies than the Golang language itself.
//
// There are three existing writers supported by default: JSON, SQL and DB.
//
// JSON Writer - As the name prompts this writer transforms input XML into JSON format.
// This writer could be used as a solution that converts data into any NoSQL database.
//
// SQL Writer - The second supported writer creates SQL file with all necessary data from input.
// This file can be executed in any SQL database and the result will be populated table with the proper information.
//
// DB Writer - The last writer is for direct communication with the SQL database.
// All input data will be saved into appropriate tables immediately.
// Before using this writer all required data tables need to be created.
// For that purpose please run the SQL script named tables.sql in the sql_scripts folder..
// You can also set up indexes on any column you want,
// to facilitate this process there is also the script called indexes.sql situated in the sql_script folder.
// To speed up a data transformation I would rather recommend creating indexes after the whole processing is completed.
//
// By Lukas Aron
package parser

// Example of usage:
// import (
//	"fmt"
//	"github.com/lukasaron/discogs-parser/decode"
// )
//
// func main() {
//	d := decode.NewXmlDecoder("./data_samples/artists.xml", nil)
//	num, _, err := d.Artists(100)
//	fmt.Println(num, err)
//	_ = d.Close()
// }