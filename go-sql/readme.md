# Cycir

VCS Traning PostgreSQL exercise, using template created by Trevor Sawler

## Setup

Build in the normal way on Windows:
- From psql, create new database name cycir
- In database.yml, change your database name, password
- In run.bat, change your -dbuser, -dbpass, -db
For example:

- database.yml:
development:
  dialect: postgres
  database: cycir
  user: postgres
  password: qwerqwer
  host: localhost
  pool: 5

- run.bat:
cycir -dbuser='postgres' -dbpass='qwerqwer' -db="cycir"

## Run
- run.bat
- access localhost:4000, query employee by id with query param ?department_id=1

## Requirements

Cycir requires:
- Postgres 11 or later (Postgres 14 is the version I'm running on my machine)

## Files for Marking
All files for marking are followed by **********

```bash
├── cmd
│   └── web
│       ├── jobs-mail.go
│       ├── main.go
│       ├── middleware.go
│       ├── routes.go
│       └── setup-app.go
├── database.yml
├── go.mod
├── go.sum
├── internal
│   ├── channeldata
│   │   └── maildata.go
│   ├── config
│   │   └── config.go
│   ├── driver
│   │   └── driver.go
│   ├── handlers
│   │   ├── all-services-status-pages.go
│   │   ├── authentication-handlers.go
│   │   ├── handlers.go ********** function AdminDashboard() **********
│   │   └── schedule.go
│   ├── helpers
│   │   ├── helpers.go
│   │   ├── send-mail.go
│   │   └── template-functions.go
│   ├── models
│   │   └── models.go
│   ├── repository
│   │   ├── dbrepo
│   │   │   ├── alerts.go **********
│   │   │   ├── dbrepo.go
│   │   │   ├── departments.go **********
│   │   │   ├── employees.go **********
│   │   │   ├── preferences.go
│   │   │   └── users_postgresql.go
│   │   └── repository.go
│   └── templates
│       └── templates.go
├── ipe
│   ├── config.yml
│   ├── ipe
│   └── ipe.exe
├── migrations ----------CREATE, INSERT TO POSTGRESQL TABLE----------
│   ├── 20220825091247_create_departments_table.up.fizz  **********
│   ├── 20220825091256_create_employees_table.up.fizz ********** 
│   ├── 20220825091739_seed_departments_table.up.fizz **********
│   ├── 20220825091746_seed_employees_table.up.fizz  **********
│   ├── 20220826024714_create_alerts_table.up.fizz **********
│   ├── 20220826024723_seed_alerts_table.up.fizz **********
│   └── schema.sql
└── views
    ├── dashboard.jet **********
    ├── events.jet
    ├── healthy.jet
    ├── host.jet
    ├── hosts.jet
    ├── layouts
    │   └── layout.jet
    ├── login.jet
    ├── partials
    │   └── js.jet
    ├── pending.jet
    ├── problems.jet
    ├── schedule.jet
    ├── settings.jet
    ├── user.jet
    ├── users.jet
    └── warning.jet
```
