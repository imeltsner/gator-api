# gator
## Overview
Gator is the backend for an RSS feed aggregation service, built with Go. Users, feeds, and posts are stored in a Postgres database.
Users can follow feeds added by other users, and browse posts made by feeds that they follow.
## Requirements
This program requires the user to have Go and Postgres installed on their machine. To install the program run
```
go install github.com/imeltsner/gator
```
Make sure to create a .env file and fill it out based on the sample.env. 

