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

**Quite heavy size for our services**
docker_images_size.png
![docker_images_size](docker_images_size.png)


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

### or use phpmyadmin:
credentials:
- username: root
- password: rootpassword

**phpmyadmin Preview**
![phpmyadmin_gorm_learning_preview](phpmyadmin_gorm_learning_preview.png)


### testing the go webapi app:
**visit `http://localhost:3000/Customer`**:
Response:
```json
[
  {
    "CustomerId": 1,
    "Name": "John Smith",
    "Email": "john.smith@example.com",
    "Orders": null
  },
  {
    "CustomerId": 2,
    "Name": "Jane Doe",
    "Email": "jane.doe@example.com",
    "Orders": null
  },
  {
    "CustomerId": 3,
    "Name": "Bob Johnson",
    "Email": "bob.johnson@example.com",
    "Orders": null
  }
]
```

**visit `http://localhost:3000/Product`**:
Response:
```json
[
  {
    "ProductId": 1,
    "name": "Product 1",
    "price": 9.99
  },
  {
    "ProductId": 2,
    "name": "Product 2",
    "price": 14.99
  },
  {
    "ProductId": 3,
    "name": "Product 3",
    "price": 19.99
  }
]
```

**visit `http://localhost:3000/orders-miss-product`**:
Response:
```json
[
  {
    "OrderId": 2,
    "OrderDate": "2025-04-05T23:41:32Z",
    "CustomerId": 2,
    "Customer": {
      "CustomerId": 2,
      "Name": "Jane Doe",
      "Email": "jane.doe@example.com"
    },
    "OrderItems": [
      {
        "OrderItemId": 4,
        "ProductId": 3,
        "Quantity": 1,
        "Product": {
          "ProductId": 0,
          "Name": "",
          "Price": 0
        }
      },
      {
        "OrderItemId": 3,
        "ProductId": 1,
        "Quantity": 3,
        "Product": {
          "ProductId": 0,
          "Name": "",
          "Price": 0
        }
      }
    ],
    "Total": 0
  },
  {
    "OrderId": 3,
    "OrderDate": "2025-04-05T23:41:32Z",
    "CustomerId": 3,
    "Customer": {
      "CustomerId": 3,
      "Name": "Bob Johnson",
      "Email": "bob.johnson@example.com"
    },
    "OrderItems": [
      {
        "OrderItemId": 6,
        "ProductId": 3,
        "Quantity": 2,
        "Product": {
          "ProductId": 0,
          "Name": "",
          "Price": 0
        }
      },
      {
        "OrderItemId": 5,
        "ProductId": 2,
        "Quantity": 2,
        "Product": {
          "ProductId": 0,
          "Name": "",
          "Price": 0
        }
      }
    ],
    "Total": 0
  },
  {
    "OrderId": 1,
    "OrderDate": "2025-04-05T23:41:32Z",
    "CustomerId": 1,
    "Customer": {
      "CustomerId": 1,
      "Name": "John Smith",
      "Email": "john.smith@example.com"
    },
    "OrderItems": [
      {
        "OrderItemId": 2,
        "ProductId": 2,
        "Quantity": 1,
        "Product": {
          "ProductId": 0,
          "Name": "",
          "Price": 0
        }
      },
      {
        "OrderItemId": 1,
        "ProductId": 1,
        "Quantity": 2,
        "Product": {
          "ProductId": 0,
          "Name": "",
          "Price": 0
        }
      }
    ],
    "Total": 0
  }
]
```
**visit `http://localhost:3000/orders-worked`**:
Response:
```json
[
  {
    "OrderId": 1,
    "OrderDate": "2025-04-05T23:41:32Z",
    "CustomerId": 1,
    "Customer": {
      "CustomerId": 1,
      "Name": "John Smith",
      "Email": "john.smith@example.com"
    },
    "OrderItems": [
      {
        "OrderItemId": 2,
        "ProductId": 2,
        "Quantity": 1,
        "Product": {
          "ProductId": 2,
          "Name": "Product 2",
          "Price": 14.99
        }
      },
      {
        "OrderItemId": 1,
        "ProductId": 1,
        "Quantity": 2,
        "Product": {
          "ProductId": 1,
          "Name": "Product 1",
          "Price": 9.99
        }
      }
    ],
    "Total": 34.97
  },
  {
    "OrderId": 2,
    "OrderDate": "2025-04-05T23:41:32Z",
    "CustomerId": 2,
    "Customer": {
      "CustomerId": 2,
      "Name": "Jane Doe",
      "Email": "jane.doe@example.com"
    },
    "OrderItems": [
      {
        "OrderItemId": 4,
        "ProductId": 3,
        "Quantity": 1,
        "Product": {
          "ProductId": 3,
          "Name": "Product 3",
          "Price": 19.99
        }
      },
      {
        "OrderItemId": 3,
        "ProductId": 1,
        "Quantity": 3,
        "Product": {
          "ProductId": 1,
          "Name": "Product 1",
          "Price": 9.99
        }
      }
    ],
    "Total": 49.96
  },
  {
    "OrderId": 3,
    "OrderDate": "2025-04-05T23:41:32Z",
    "CustomerId": 3,
    "Customer": {
      "CustomerId": 3,
      "Name": "Bob Johnson",
      "Email": "bob.johnson@example.com"
    },
    "OrderItems": [
      {
        "OrderItemId": 6,
        "ProductId": 3,
        "Quantity": 2,
        "Product": {
          "ProductId": 3,
          "Name": "Product 3",
          "Price": 19.99
        }
      },
      {
        "OrderItemId": 5,
        "ProductId": 2,
        "Quantity": 2,
        "Product": {
          "ProductId": 2,
          "Name": "Product 2",
          "Price": 14.99
        }
      }
    ],
    "Total": 69.96
  }
]
```

**visit `http://localhost:3000/orders-not-grouped`**:
Response:
```json
[
  {
    "OrderId": 1,
    "OrderDate": "2025-04-05T23:41:32Z",
    "CustomerId": 1,
    "Customer": {
      "CustomerId": 1,
      "Name": "John Smith",
      "Email": "john.smith@example.com"
    },
    "OrderItems": [
      {
        "OrderItemId": 2,
        "ProductId": 2,
        "Quantity": 1,
        "Product": {
          "ProductId": 2,
          "Name": "Product 2",
          "Price": 14.99
        }
      }
    ],
    "Total": 14.99
  },
  {
    "OrderId": 1,
    "OrderDate": "2025-04-05T23:41:32Z",
    "CustomerId": 1,
    "Customer": {
      "CustomerId": 1,
      "Name": "John Smith",
      "Email": "john.smith@example.com"
    },
    "OrderItems": [
      {
        "OrderItemId": 1,
        "ProductId": 1,
        "Quantity": 2,
        "Product": {
          "ProductId": 1,
          "Name": "Product 1",
          "Price": 9.99
        }
      }
    ],
    "Total": 19.98
  },
  {
    "OrderId": 2,
    "OrderDate": "2025-04-05T23:41:32Z",
    "CustomerId": 2,
    "Customer": {
      "CustomerId": 2,
      "Name": "Jane Doe",
      "Email": "jane.doe@example.com"
    },
    "OrderItems": [
      {
        "OrderItemId": 4,
        "ProductId": 3,
        "Quantity": 1,
        "Product": {
          "ProductId": 3,
          "Name": "Product 3",
          "Price": 19.99
        }
      }
    ],
    "Total": 19.99
  },
  {
    "OrderId": 2,
    "OrderDate": "2025-04-05T23:41:32Z",
    "CustomerId": 2,
    "Customer": {
      "CustomerId": 2,
      "Name": "Jane Doe",
      "Email": "jane.doe@example.com"
    },
    "OrderItems": [
      {
        "OrderItemId": 3,
        "ProductId": 1,
        "Quantity": 3,
        "Product": {
          "ProductId": 1,
          "Name": "Product 1",
          "Price": 9.99
        }
      }
    ],
    "Total": 29.97
  },
  {
    "OrderId": 3,
    "OrderDate": "2025-04-05T23:41:32Z",
    "CustomerId": 3,
    "Customer": {
      "CustomerId": 3,
      "Name": "Bob Johnson",
      "Email": "bob.johnson@example.com"
    },
    "OrderItems": [
      {
        "OrderItemId": 6,
        "ProductId": 3,
        "Quantity": 2,
        "Product": {
          "ProductId": 3,
          "Name": "Product 3",
          "Price": 19.99
        }
      }
    ],
    "Total": 39.98
  },
  {
    "OrderId": 3,
    "OrderDate": "2025-04-05T23:41:32Z",
    "CustomerId": 3,
    "Customer": {
      "CustomerId": 3,
      "Name": "Bob Johnson",
      "Email": "bob.johnson@example.com"
    },
    "OrderItems": [
      {
        "OrderItemId": 5,
        "ProductId": 2,
        "Quantity": 2,
        "Product": {
          "ProductId": 2,
          "Name": "Product 2",
          "Price": 14.99
        }
      }
    ],
    "Total": 29.98
  }
]
```

**visit `http://localhost:3000/orders-raw-not-grouped`**:
Response:
```
error parsing order date: parsing time "2025-04-05 23:41:32" as "2006-01-02 15:04:05.000": cannot parse "" as ".000"
```

**visit `http://localhost:3000/orders-simple-orm-anonymous`**:
Response:
```json
[
  {
    "orderId": 1,
    "orderDate": "2025-04-05T23:41:32Z",
    "customerId": 1
  },
  {
    "orderId": 2,
    "orderDate": "2025-04-05T23:41:32Z",
    "customerId": 2
  },
  {
    "orderId": 3,
    "orderDate": "2025-04-05T23:41:32Z",
    "customerId": 3
  }
]
```