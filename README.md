# Golang CRUD API with Gin and GORM

## 📌 Overview
This is a simple **CRUD (Create, Read, Update, Delete)** API built using **Golang**, **Gin (web framework)**, and **GORM (ORM for PostgreSQL)**. The API allows users to add, retrieve, update, and delete items from a PostgreSQL database.

## 🚀 Features
- ✅ Create an item
- ✅ Retrieve all items
- ✅ Retrieve a single item by ID
- ✅ Update an item
- ✅ Delete an item
- ✅ Modular project structure

## 🏗️ Project Structure
```
my-go-app/
│── main.go           # Entry point of the application
│── config/           # Configuration (DB connection)
│   ├── database.go   # DB initialization
│── routes/           # Routes handlers
│   ├── item_routes.go
│── controllers/      # Business logic (CRUD functions)
│   ├── item_controller.go  
│── models/           # Database models
│   ├── item.go       
│── repository/       # Database operations layer
│   ├── item_repo.go  
│── services/         # Service layer (business logic)
│   ├── item_service.go  
│── .env              # Environment variables (DB creds)
│── go.mod            # Golang module file
│── go.sum            # Dependencies lock file
```

## ⚙️ Installation & Setup
### 1️⃣ Clone the Repository
```bash
git clone https://github.com/yourusername/my-go-app.git
cd my-go-app
```

### 2️⃣ Install Dependencies
```bash
go mod tidy
```

### 3️⃣ Set Up PostgreSQL Database
1. Open PostgreSQL shell:
   ```bash
   psql -U postgres -h localhost -p 5432
   ```
2. Create a new database:
   ```sql
   CREATE DATABASE mygoitems;
   ```
3. Exit PostgreSQL:
   ```sql
   \q
   ```

### 4️⃣ Configure `.env` File
Create a `.env` file in the root directory and add:
```env
DB_HOST=localhost
DB_USER=postgres
DB_PASSWORD=yourpassword
DB_NAME=mygoitems
DB_PORT=5432
```

### 5️⃣ Run the Application
```bash
go run main.go
```
The server will start on **`http://localhost:8080`**.

## 🔥 API Endpoints
### ✅ Create an Item
```http
POST /items/
```
**Request Body:**
```json
{
  "name": "Laptop",
  "price": 1200
}
```

### ✅ Get All Items
```http
GET /items/
```

### ✅ Get Item by ID
```http
GET /items/{id}
```

### ✅ Update an Item
```http
PUT /items/{id}
```
**Request Body:**
```json
{
  "name": "Gaming Laptop",
  "price": 1500
}
```

### ✅ Delete an Item
```http
DELETE /items/{id}
```

## 🐳 Run with Docker (Optional)
1. **Build the Docker image:**
   ```bash
   docker build -t my-go-app .
   ```
2. **Run the container:**
   ```bash
   docker run -p 8080:8080 --env-file .env my-go-app
   ```

## 🛠️ Next Steps
- [ ] Add authentication (JWT)
- [ ] Improve error handling
- [ ] Add unit tests
- [ ] Deploy to a cloud platform

---
Made with ❤️ using Golang & PostgreSQL 🚀

