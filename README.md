# Discogs Parser


Discogs Parser introduces the way to categorize XML data dumps from Discogs [https://data.discogs.com/].

The intention of usage this parser is library, which means there is not an executable part provided. 

There are three existing writers supported by default: `json`, `sql` and `db`.


#### JSON Writer
As the name prompts this writer transforms input XML into JSON format. This writer could be used as a solution that converts data into any NoSQL database.

#### SQL Writer
The second supported writer creates SQL file with all necessary data from input. This file can be executed in any SQL database and the result will be populated table with the proper information.

#### DB Writer
The last writer is for direct communication with SQL database. All input data will be saved into appropriate tables immediately.
Before using this writer all required data tables need to be created. For that purpose please run the SQL script `sql_scripts/tables.sql`. 

You can also setup indexes on any column you want, to facilitate this process there is also a script for that `sql_scripts/indexes.sql`. 

In order to speed up a data transformation I would rather recommend to create indexes after the whole processing is completed. 