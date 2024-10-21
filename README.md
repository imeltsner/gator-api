# gator
## Overview
Gator is a command line RSS feed aggregation tool. Users can input links to RSS feeds they would like to subscribe to.
While the program is running, it will check the RSS feeds for new posts at set intervals, and add them to the database.
Users can then browse the RSS feed links and click on those they are intersted in.
## Requirements
This program requires the user to have Go and Postgres installed on their machine. To install the program run
```
go install github.com/imeltsner/gator
```
The program uses a config file to connect
to the database and store the current user. It looks for the config file at ~/.gatorconfig.json. The program will automatically
set the current user once it is run. Below is example config file. Make sure to replace user with the database user and password 
with the password for that user.
```
{
    "db_url":"postgres://user:password@localhost:5432/gator?sslmode=disable"
}
```
## Commands
* Register - adds a new user to the database
```
gator register <username>
```
* Login - sets the active user to the name provided
```
gator login <username>
```
* Add feed - adds an RSS feed to the database
```
gator addfeed <title> <url>
```
* Aggregate - check RSS feeds that user follows at given interval (10s, 1m, 5m)
```
gator agg <interval>
```
* Browse - show most recent posts from followed feeds
```
gator browse <number of posts>
```
* Follow - follow a feed added by another user
```
gator follow <url>
```
* Unfollow - unfollows a feed added by another user
```
gator unfollow <url>
```
