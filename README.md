# Golang CRUD API with Gin & GORM

## 📌 API Overview  
A simple **CRUD API** built with **Golang**, **Gin**, and **GORM** for managing items in a PostgreSQL database.  

## 🚀 API Endpoints  

### **1️⃣ Create an Item**  
**POST** `/items/`  
**Request Body:**  
```json
{
  "name": "Laptop",
  "price": 1200
}
```

### **2️⃣ Get All Items**  
**GET** `/items/`  

### **3️⃣ Get Item by ID**  
**GET** `/items/{id}`  

### **4️⃣ Update an Item**  
**PUT** `/items/{id}`  
**Request Body:**  
```json
{
  "name": "Gaming Laptop",
  "price": 1500
}
```

### **5️⃣ Delete an Item (Soft Delete)**  
**DELETE** `/items/{id}`  

## 🏗️ Setup & Run  
1. **Install dependencies:**  
   ```bash
   go mod tidy
   ```  
2. **Set up PostgreSQL and create the database:**  
   ```sql
   CREATE DATABASE mygoitems;
   ```  
3. **Run the server:**  
   ```bash
   go run main.go
   ```  
   Server runs on **`http://localhost:8080`**  
