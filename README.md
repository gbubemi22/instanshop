# Project instashop

One Paragraph of project description goes here

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See deployment for notes on how to deploy the project on a live system.

## MakeFile

run all make commands with clean tests
```bash
make all build
```

build the application
```bash
make build
```

run the application
```bash
make run
```

Create DB container
```bash
make docker-run
```

Shutdown DB container
```bash
make docker-down
```

live reload the application
```bash
make watch
```

run the test suite
```bash
make test
```

clean up binary from the last build
```bash
make clean
```


postman Documentation
```bash
https://documenter.getpostman.com/view/26055191/2sAYJ6BzR4
```




Project Documentation: E-Commerce API

Table of Contents

Overview

Technologies Used

Database Schema

Endpoints

How to Test

Project Setup

Future Improvements

Overview

The E-Commerce API provides a backend for managing products, orders, and users. It includes user authentication, product creation, order placement, and many-to-many relationships between orders and products. The project follows best practices such as modular architecture, transaction management, and role-based access control.

Technologies Used

Programming Language: Go (Golang)

Framework: Gin

Database: PostgreSQL

ORM: GORM

Authentication: JWT (JSON Web Tokens)

Testing Tool: Postman

Database Schema

Users Table

Column

Type

Constraints

id

SERIAL

PRIMARY KEY

email

VARCHAR

NOT NULL, UNIQUE

password

VARCHAR

NOT NULL

role

VARCHAR

DEFAULT 'user'

Products Table

Column

Type

Constraints

id

SERIAL

PRIMARY KEY

user_id

INTEGER

NOT NULL, FOREIGN KEY

name

VARCHAR(255)

NOT NULL

description

TEXT

NOT NULL

price

DECIMAL

NOT NULL

status

VARCHAR(10)

DEFAULT 'pending'

created_at

TIMESTAMP

DEFAULT CURRENT_TIMESTAMP

updated_at

TIMESTAMP

DEFAULT CURRENT_TIMESTAMP

Orders Table

Column

Type

Constraints

id

SERIAL

PRIMARY KEY

user_id

INTEGER

NOT NULL, FOREIGN KEY

status

VARCHAR(10)

DEFAULT 'pending'

created_at

TIMESTAMP

DEFAULT CURRENT_TIMESTAMP

updated_at

TIMESTAMP

DEFAULT CURRENT_TIMESTAMP

Order_Products Table

Column

Type

Constraints

order_id

INTEGER

FOREIGN KEY, PRIMARY KEY

product_id

INTEGER

FOREIGN KEY, PRIMARY KEY

