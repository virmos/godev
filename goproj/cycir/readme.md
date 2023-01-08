![Unit Tests](https://github.com/virmos/godev/blob/main/cycir/.github/workflow/badge.svg)

![image](https://user-images.githubusercontent.com/30485720/197343342-60366eb5-d726-4f8d-9764-ac0abdedebbb.png)

--------------------------------------------------------------------

# Installation
~~~
cycir-api -dbuser='postgres' -dbpass='qwerqwer' -esAddress="http://localhost:9200" -esUsername="elastic" -esPassword="EWAq+EaS8dyQV_82TSQd" -esIndex="my-index-000001" -pusherHost="localhost" -pusherPort="4001" -pusherSecret="123abc" -pusherKey="abc123" -pusherSecure=false -pusherApp="1" -db="temp"
~~~

Change the following:
+ dbuser, dbpass, db(database name): for Postgres
+ esPassword, esIndex: for Elasticsearch

Runs
soda migrate

On Windows:

~~~
run.bat
r.bat
~~~

On Linux:

~~~
env GOOS=linux GOARCH=amd64 go build -o cycir cmd/web/*.go
env GOOS=linux GOARCH=amd64 go build -o cycir-api cmd/api/*.go
~~~

# Overview
Login
![image](https://user-images.githubusercontent.com/30485720/197343623-63c71208-82ef-4559-9405-966458d7332c.png)

Dashboard
![image](https://user-images.githubusercontent.com/30485720/197343634-bda43fef-6109-4ece-bcf9-cff13a4926a6.png)

Manage Services
![image](https://user-images.githubusercontent.com/30485720/197343668-be0d4d49-1155-42ae-8b20-f794a12d08ba.png)

![image](https://user-images.githubusercontent.com/30485720/197343685-a21cfa6f-2eb1-4659-abe3-b91c36cd79a4.png)

![image](https://user-images.githubusercontent.com/30485720/197343998-21c3ccfe-9d29-46d6-9fb0-96d7046f3e32.png)

![image](https://user-images.githubusercontent.com/30485720/197344006-355918df-86ae-469b-be7f-d4859bb5985a.png)

Setup for Mail sending
![image](https://user-images.githubusercontent.com/30485720/197344074-1adc14f5-0990-4978-952b-a33acea480e8.png)

Other pages:
![image](https://user-images.githubusercontent.com/30485720/197346971-92cdac0a-14ac-49c7-8f1d-024f32ec87cf.png)

![image](https://user-images.githubusercontent.com/30485720/197346978-a96ff70b-3bcf-45ed-883d-5870e208d3c6.png)


# System Design
## ER Diagram
![image](https://user-images.githubusercontent.com/30485720/197344473-f1329a42-3062-47db-9362-e7e60ff157b3.png)

Table preferences: holds all app states, for example:
+ preferences[monitoring_live] = "1" means monitoring is on. When monitoring is on, it automatically check hosts' availability @every 3m. 

Table host_services:
+ field service_id references Table services
+ field host_id refences Table hosts
+ field active if set true -> that service is on for that host -> When monitoring is on, automatically check service's availability @every 3m

Table hosts:
+ field active if set false -> Even if monitoring is on, no automatical check will be done

Table services: HTTP, HTTPS, SSL

Table events: Every hosts check, one event is added to the table

Table users: For Login

Table tokens: When logging in, api creates new token and sends back to frontend to store in session. Without token, frontend cannot communicate with the api. Each token has an expiry date

Table rememberme_tokens: If set true, every login without token(in table tokens) but with a rememberme_token, it will automatically extend the duration of the token(in table tokens!)

Table sessions: Not needed, since [scs]("github.com/alexedwards/scs/v2") will be used to store session

## Code
The App will be separated in to cycir/cmd/api(backend) and cycir/cmd/web(frontend)
+ Backend will handle post requests
+ Frontend will handle get requests
+ Both create a new database connection by importing cycir/internal/driver

Setup files are in
+ Backend: cycir/cmd/api/api.go
+ Frontend: cycir/cmd/web/main.go

Requests are redirected to routes -> handlers

### Files structure
Frontend:
+ all-services-status-pages.go: counts healthy/problems/pending services to render at dashboard
+ authentication-handlers.go: for Frontend login/logout
+ cache-handlers.go: stores uptime report data in redis cache (cycir/internal/cache)
+ render.go: add default data(PreferenceMap, IsAuthenticated) and render page

Backend:
+ file-handlers-api.go: import, export excel
+ handlers-api.go: handles post requests, send range uptime reports, which gets a range of dates, displays report in format of an array len == 31 (days) (this will be explained further in Send Uptime Report part)
+ start-monitoring.go: add automatically check host_services function && automatically send yesterday uptime report function to [cron]("github.com/robfig/cron/v3")
+ perform-report.go: creates yesterday uptime report function, which displays yesterday's uptime report in format of an array len == 24 (hours)
+ schedule-check.go: creates automatically check host_services function that sends get requets to host'url
+ pusher.go: probably not needed, but pushed is used by api to notify frontend for the changes without frontend having to reload

# Run
## Server Management
### Server Creatation

![image](https://user-images.githubusercontent.com/30485720/197344602-fdc955ce-4081-423a-b19c-f6a8d9427a9a.png)

![image](https://user-images.githubusercontent.com/30485720/197344898-ea1f9aae-8d76-4ebe-9c44-943933c10990.png)

![image](https://user-images.githubusercontent.com/30485720/197344905-11631985-f327-4e5a-8b4f-8d9e60ff69b1.png)

![image](https://user-images.githubusercontent.com/30485720/197345118-136e9e9b-e327-47cf-9535-fd0beebf94ba.png)

### View Server
![image](https://user-images.githubusercontent.com/30485720/197345214-086c3f70-8336-4e87-b32e-52eaf81ba0f3.png)

### Update Server
Just Update It

### Delete Server
Foreign key is deleted then readded. Potential place for issues?
![image](https://user-images.githubusercontent.com/30485720/197345238-8c453e36-08c8-440e-aae6-4c3812a9c3e9.png)

![image](https://user-images.githubusercontent.com/30485720/197345250-33286c6e-db2e-496b-982e-c3b27783e744.png)

### Import Server

![image](https://user-images.githubusercontent.com/30485720/197345132-3616ee39-72f1-4f21-953c-d95f33616a3b.png)

![image](https://user-images.githubusercontent.com/30485720/197345153-d9c56b76-c24e-49ed-9e8b-87924cc8fbcf.png)


### Export Server

![image](https://user-images.githubusercontent.com/30485720/197345183-99ea26f0-b0fa-4ea0-b4d7-ae8b3580c6e0.png)

![image](https://user-images.githubusercontent.com/30485720/197345190-17a8e89b-29be-475d-a508-7711851a07b3.png)

## Uptime Report

Input: StartDate, EndDate (admin's mail config is in preference_map)
![image](https://user-images.githubusercontent.com/30485720/197351648-2b76677b-d573-4290-a978-2dfad181da58.png)

Yesterday Uptime Report
+ First line displays uptime report in percentage for 24 hours from 0 - 23
+ Second line displays the number of uptime checks are made to that host(ONLY!!!!! counts schedule check, not when check now is pressed) 
![image](https://user-images.githubusercontent.com/30485720/197346322-f269a25d-f42e-4d66-95d8-927d7abe52e0.png)

Uptime Range Report
+ First line displays uptime report in percentage for maximum 31 days(because the array has 31 elements, arr[0] = 10%, means uptime report of day 0 is 10%, day 0 is StartDate), from the StartDate specified - EndDate
+ Second line displays the number of uptime checks are made to that host(ONLY!!!!! counts schedule check, not when check now is pressed) 
![image](https://user-images.githubusercontent.com/30485720/197346314-ac5cb367-2dde-4fda-b907-9c8d23ea0837.png)

## Redis
Uptime report is saved the first time request is sent. Next time, data is fetched from cache, no need to call the api. However after a schedule check is made, data is cleared from the cache
![image](https://user-images.githubusercontent.com/30485720/197346812-9471bac3-6636-4007-9738-66356e09645c.png)

## Logrotate
![image](https://user-images.githubusercontent.com/30485720/197346841-85c89f38-7618-4ac5-8167-329ed4bd0c72.png)

infolog
![image](https://user-images.githubusercontent.com/30485720/197346852-31d9d145-ffed-4d4d-a7e8-5ded0c33fef7.png)


## Privacy
Front end:
![image](https://user-images.githubusercontent.com/30485720/197347861-2f02d593-4008-4396-a113-8cc0f2373385.png)

![image](https://user-images.githubusercontent.com/30485720/197347868-161e39bd-3d83-4ccd-a939-bb839b421fa4.png)

Backend:

![image](https://user-images.githubusercontent.com/30485720/197347896-01216ee8-a611-4965-9e15-ff32a31f8287.png)

![image](https://user-images.githubusercontent.com/30485720/197347887-8aac7cef-fb95-4d19-ae15-d3faf2203286.png)

Token is generated by CreateAuthToken

# Unit Test

### Front end
cd cmd/web

go test -coverprofile=coverage.out && go tool cover -html=coverage.out

Handlers
![image](https://user-images.githubusercontent.com/30485720/197347461-ccc9ae33-807f-4a44-9ee5-6e87634e4001.png)

Majority of code not being covered are these functions. All route handler functions are covered
![image](https://user-images.githubusercontent.com/30485720/197347626-882d7751-557e-4010-ab6f-58d13fb8b353.png)

Same for Render, code coverage is 58.3% but it actually covers all neccessary functions like add data and render page

Render
![image](https://user-images.githubusercontent.com/30485720/197347606-7bb1f219-5fb4-4713-bfff-20806d276748.png)

### Back end
cd cmd/api

go test -coverprofile=coverage.out && go tool cover -html=coverage.out

![image](https://user-images.githubusercontent.com/30485720/197347727-7a7dfb16-6a97-4048-a85b-68437ff25ac6.png)

Handlers

![image](https://user-images.githubusercontent.com/30485720/197347756-55232ccd-c333-4864-8b1a-2c0bae39d4e9.png)


## Requirements

cycir requires:
- Postgres 11 or later (db is set up as a repository, so other databases are possible)
- An account with [Pusher](https://pusher.com/), or a Pusher alternative 
(like [ipê](https://github.com/dimiro1/ipe))
- [Soda CLI](https://gobuffalo.io/documentation/database/soda/)
