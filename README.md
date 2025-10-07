# gatorcli
Blog Aggregator Project - Boot.dev

Prerequisites:
For this program, you will need Go and PostgreSQL installed.

Go installation terminal command: curl -sS https://webi.sh/golang | sh

PostgreSQL installation: 
    MacOS with brew: 
     - brew install postgresql@15
    Linux/WSL(Debian):
     - sudo apt update
     - sudo apt install postgresql postgresql-contrib
     - sudo passwd postgres //Linus only. Sets password.
        - enter a memorable password
PostgreSQL documentation: https://docs.sqlc.dev/en/latest/tutorials/getting-started-postgresql.html

Gator CLI requires a .json config file to run. In your home directory, manually create the following file: ~/.gatorconfig.json

In the config file, put the following content:

{
    "db_url": "postgres://example"
}

This will be your connection string. The application will add a current_user_name field.


Intalling the Gator CLI:

After downloading the repo, navigate to the root directory of the application. This should be the gatorcli directory. Type "go install" in your terminal from the root directory of the project.

Once installed, type gatorcli followed by a command and any required arguments. Available commands are as follows:

	- "login" //logs in an existing user
	- "register" //adds a user to the user table. takes user name as argument
	- "reset" //deletes all tables and recreates as empty
	- "users" //lists users in the user table
	- "agg" //records posts from followed feeds to the posts table. takes time duration in seconds as argument
	- "addfeed" //adds a new feed from a URL and follows for the user who added it
	- "feeds" //shows feeds in the feeds table
	- "follow" //follows and existing feed for the active user
	- "following" //shows all feeds followed by the active user
	- "unfollow" //unfollows a feed based on URL
	- "browse" //shows posts that have been logged from the user's followed feeds.  takes a limit argument (defaults to 2)

