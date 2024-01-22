# GO-HTMX Todo List

## Overview

This is a todo list application built using Go with HTMX and PostgreSQL. It allows users to add, mark tasks as completed or pending, and delete tasks.

## Features

- Add new tasks
- Toggle tasks between completed or pending status
- Delete tasks from the list

## Prerequisites

Before you begin, ensure you have the following installed on your system:

- [Go](https://golang.org/dl/)
- [Docker](https://www.docker.com/products/docker-desktop)

## Installation

1. Clone the repository to your local machine:

```sh
git clone https://github.com/marcelofeitoza/go-htmx-todo-list
```

2. Navigate to the cloned project directory:

```sh
cd go-htmx-todo-list
```

3. Set up the PostgreSQL database. You can do this by running a PostgreSQL instance on your machine or using Docker with the provided `docker-compose.yml`:

```sh
docker-compose up -d
```

## Running the Application

To run the application, execute the following command in the project directory:

```sh
go run main.go
```

After the server starts, you can access the todo list application by navigating to `http://localhost:8080` in your web browser.
