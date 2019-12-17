# Discogs Parser

Discogs Parser introduces the way to categorize data dumps from Discogs [https://data.discogs.com/]

## JSON Writer
It's a default writer which needs to know what is a name of the output file, therefore argument `output` needs to be specified.
- `output` file name that will be created with the JSON output. When this flag is not specified the new file with is created in the same folder as the input file name is saved. The extension of the file name is JSON instead of XML.

## Postgres Writer
This writer can write results in a PostgreSQL database. To do this the specified environmental variables need to be defined:
- DB_HOST (default: `localhost`)
- DB_PORT (default: `5432`)
- DB_NAME (default: `discogs`)
- DB_USERNAME (default: `user`)
- DB_PASSWORD (default: `password`)
- DB_SSL_MODE (default: `disable`), possible values defined in the PG driver's documentation [https://godoc.org/github.com/lib/pq]

Before processing,  make sure you have created all necessary database tables. For this, I would recommend running the `tables.sql` script from the sql_scripts folder. For the performance purposes index creation is recommended to perform after the processing of a data. Therefore the creation of indexes for all tables is separated into `indexes.sql`.  

## Common command line arguments
- `filename` this is the only one mandatory argument! It's a file name as input (usually `discogs_`[`date of dump`]`_`[`artists `|` labels `|` masters `|` releases`]`.xml`)
- `writer-type` type of a output writer (default: `json`), possible values [`json` | `postgres`]

### Quality Filter
-`quality` filter output based on the input data quality field defined by Discogs (default: `All`), all possible values: [`All` | `EntirelyIncorrect` | `NeedsVote` | `NeedsMajorChanges` | `NeedsMinorChanges` | `Correct` | `CompleteAndCorrect`]

### Processing blocks
- `block-size` specifies the processing block size (the number of items processed at once). Default value: `1000`
- `block-skip` the number of blocks from the beginning that will  be skipped, default: `0`
- `block-limit` limit amount of blocks to be processed. Default: `2147483647`

## Example
 - Have a look into the `start.sh` script for inspiration how to run the Discogs-parser  

# TODO
- Comments
- Tests