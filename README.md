This is a tool to convert a dump created with `mysqldump` to the format that TiDB Dumpling and MyDumper are using.

This needs a lot more testing before it is ready for production usage.

Example run:

```
$ ./mysqldump_converter test1.sql 
2022/11/15 09:24:05 Writing output to /tmp/mysqldump_converter
2022/11/15 09:24:05 Processing schema: test
2022/11/15 09:24:05 Processing table schema for test.t1
2022/11/15 09:24:05 Processing table data for test.t1
2022/11/15 09:24:05 Processing table schema for test.t2
2022/11/15 09:24:05 Processing table data for test.t2
$ ls -1 /tmp/mysqldump_converter
test.t1-schema.sql
test.t1.sql
test.t2-schema.sql
test.t2.sql
```
