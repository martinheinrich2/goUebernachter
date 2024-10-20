# goUebernachter

Is a web app written in Go to register short term stays of clients.
Thanks to Alex Edwards for writing the fantastic book [*"Let's Go"*](https://lets-go.alexedwards.net/).

Use makefile commands, e.g. make run.

`make build` will create the binary in /tmp/bin  
`make migrations/create name=SOMENAME` will create a new migrations file with SOMENAME.
`make migrations/up` will apply migrations

You can move the binary to another folder, but make sure to copy the /tls folder and the sqlite.db
to the same folder, before you execute the binary. Otherwise it will not find the certificates and
database. Run the binary with filename(.exe) >> logfile.log to redirect messages and append to file.

It uses the modernc.org/sqlite package, a cgo-free port of SQLite, for easy building the app in Windows.

Procedure:

1. `make migrations/up` to create an empty database or use provided empty database
2. build/run application
3. create user accounts for social worker and other staff
4. enter clients and stays

Note: When adding a new stay of a client, a social worker is required. An appointment date/time can be added.
