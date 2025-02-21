# Restaurant-Management
Restaurant Management API (Golang, Gin, MongoDB)  This is a RESTful API for a restaurant management system built using Golang, Gin framework, and MongoDB. It includes features for managing food items, orders, order items, invoices, tables, and menus.

## Features

🛠 User Authentication
- Signup & Login with secure password hashing (bcrypt).
- JWT-based authentication for secure access.
- Refresh token implementation for session management.
🍽 Restaurant Management
Food Management: Add, update, and retrieve food items.
Menu Handling: Organize food items into different categories.
Table Management: Assign and track restaurant tables.
📦 Order Management
Place orders and associate them with tables.
Retrieve order details using order_id.
Fetch order information, including:
Number of food items
Quantity of each item
Total cost calculation
Table number
📊 MongoDB Aggregation for Order Insights
Use of MongoDB aggregation functions to fetch structured order details efficiently.
🧾 Invoice Generation
Generate invoices based on orders and items.
Include total bill amount with tax calculations.
🚀 Tech Stack
Golang with Gin Framework for building REST APIs.
MongoDB for NoSQL database storage.
JWT Authentication for security.


## User Authentication (Signup & Login)
Endpoint: POST /users/signup
### Request:

{
  "first_name": "Vj",
  "last_name": "Raj",
  "password": "123456",
  "email": "vj@gmail.com",
  "phone": "9876543210"
}
### Response:


{
  "InsertedID": "67b7ceb740642997053061aa"
}

### Login
Endpoint: POST /users/login
### Request:

{
  "email": "vj@gmail.com",
  "password": "123456"
}
### Response:


{
  "ID": "67b7ceb740642997053061aa",
  "first_name": "Vj",
  "last_name": "Raj",
  "email": "vj@gmail.com",
  "phone": "9876543210",
  "token": "JWT_ACCESS_TOKEN",
  "refresh_token": "JWT_REFRESH_TOKEN",
  "created_at": "2025-02-21T00:54:15Z",
  "updated_at": "2025-02-21T00:54:15Z",
  "user_id": "67b7ceb740642997053061aa"
}
