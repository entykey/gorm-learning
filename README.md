## project required structure
```
.
├── Dockerfile
├── docker-compose.yml
├── go.mod
├── go.sum
├── main.go
├── mysql/
│   └── init.sql
└── .env
```

## Running the application:
### Build and start the services:
```sh
docker compose up -d --build
```

View logs (or we can view logs of 2 services in Docker Desktop Container tab, but laggy as shit):
```sh
docker compose logs -f
```

Verify the network was created:
```sh
docker network ls
```

Stop the services:
```sh
# Stop and remove containers
docker compose down

# Remove all volumes (including persistent data, eg. mysql db)
docker volume prune

docker network prune
```

## To inspect the db container:
1. Exec into the db container
```sh
docker exec -it gorm-learning-db-1 bash
```
2. login to mysql
```sh
mysql -u root -p
```
Enter password: "rootpassword"
NOTE: And our database name is "gorm-learning"

**Output, with a simple query**:
```sh
bash-5.1# mysql -u root -p
Enter password: 
Welcome to the MySQL monitor.  Commands end with ; or \g.
Your MySQL connection id is 88
Server version: 8.0.41 MySQL Community Server - GPL

Copyright (c) 2000, 2025, Oracle and/or its affiliates.

Oracle is a registered trademark of Oracle Corporation and/or its
affiliates. Other names may be trademarks of their respective
owners.

Type 'help;' or '\h' for help. Type '\c' to clear the current input statement.

mysql> USE gorm_learning
Reading table information for completion of table and column names
You can turn off this feature to get a quicker startup with -A

Database changed
mysql> SELECT * FROM customers;
Empty set (0.00 sec)

mysql> 
```
**TIPS:** in "mysql" -> type "quit" to exit mysql, then in bash -> Hit "Ctrl + D" to exit container terminal access


