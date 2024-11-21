# Instagram Stats

An application that fetches Instagram stats of public accounts and periodically updates the profile information and videos.
For now the account username must be provided manually.

## Migrate pg models

Update the constant Dsn in const.go with your postgres configuration.
Then run the migration.
```
go run migrate.go
```

## Run the server

```
go run server.go
```

## Run the cron job

```
go run cron.go
```
