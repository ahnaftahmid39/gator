# Gator CLI

Gator is a command-line tool for managing RSS feeds, users, and subscriptions. It allows users to register, login, follow RSS feeds, browse posts, and more.

## Prerequisites

Before running Gator, ensure you have the following installed:

- [Go](https://go.dev/dl/) (version 1.23.6 or higher)
- [PostgreSQL](https://www.postgresql.org/download/)

## Installation

You can install the Gator CLI using `go install`:

```sh
 go install github.com/ahnaftahmid39/gator@latest
```

## Configuration

Before running Gator, you need to set up a configuration file.

1. Create a postgres database for this application.
2. Create a configuration file at `~/.gator/config.json`:

```json
{
  "db_url":"postgres://username:password@localhost:5432/gator?sslmode=disable",
  "current_user_name":"Bob"
}
```

3. Replace `username`, `password`, and `gator` with your actual PostgreSQL credentials and database name.
4. Make sure PostgreSQL is running and that your database is set up correctly.

## Running the Program

Once installed and configured, you can run Gator commands using:

```sh
gator <command> [arguments]
```

## Available Commands

Here are some of the commands you can use:

### User Management

- **Register a new user:**

  ```sh
  gator register <username>
  ```
  Example: `gator register alice`

- **Login as a user:**

  ```sh
  gator login <username>
  ```
  Example: `gator login alice`

- **List all users:**

  ```sh
  gator users
  ```

### Feed Management

- **Add a new RSS feed:**

  ```sh
  gator addfeed <feed_name> <feed_url>
  ```
  Example: `gator addfeed TechCrunch https://techcrunch.com/feed/`

- **List all available feeds:**

  ```sh
  gator feeds
  ```

- **Follow a feed:**

  ```sh
  gator follow <feed_url>
  ```
  Example: `gator follow https://techcrunch.com/feed/`

- **List followed feeds:**

  ```sh
  gator following
  ```

- **Unfollow a feed:**

  ```sh
  gator unfollow <feed_url>
  ```
  Example: `gator unfollow https://techcrunch.com/feed/`

### Aggregation

- **Start the feed aggregator:**

  ```sh
  gator agg <time_between_requests>
  ```
  Example: `gator agg 10m` (Fetch feeds every 10 minutes)

### Browsing Posts

- **Browse recent posts from your followed feeds:**

  ```sh
  gator browse [limit]
  ```
  Example: `gator browse 5` (View 5 latest posts)

## Contributing

Feel free to contribute by submitting issues or pull requests!


