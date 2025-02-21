# Restaurant-Management
Restaurant Management API (Golang, Gin, MongoDB)  This is a RESTful API for a restaurant management system built using Golang, Gin framework, and MongoDB. It includes features for managing food items, orders, order items, invoices, tables, and menus.

## Features

### 🛠  User Authentication
- Signup & Login with secure password hashing (bcrypt).
- JWT-based authentication for secure access.
- Refresh token implementation for session management.
### 🍽 Restaurant Management
- **Food Management**: Add, update, and retrieve food items.
- **Menu Handling**: Organize food items into different categories.
- **Table Management**: Assign and track restaurant tables.
- **Order Management**: Create, update, and retrieve orders.
- **Order Item Management**: Manage individual items within an order.

### 🔍 Fetch Order Details Using Order ID
- Retrieve complete order details using order_id, including:
- Food items
- Table number
- Order details
- Number of persons at the table
- Total cost 
- Efficiently fetch structured order details using MongoDB Aggregation functions.
  
### 🧾 Invoice Generation
- Generate invoices based on orders and items.
- Include total bill amount with tax calculations.
### 🚀 Tech Stack
- Golang with Gin Framework for building REST APIs.
- MongoDB for NoSQL database storage.
- JWT Authentication for security.
- Docker for containerized deployment.
- Swagger API Documentation for easy API testing and visualization.


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

## Order Billing Process  POST http://localhost:8000/orders

  ### Request for Create table 

   {
  "number_of_guests": 2,
  "table_number": "T2"
  }
  
  ### Response
  {
    "InsertedID": "67b7dbef1c0f9ece68877069"
  } 
  ### Request for Create Order 

   
{
  "order_date": "2025-02-21T00:00:00Z",
  "table_id": "67b7dbef1c0f9ece68877069"
}


### Response 

{
    "InsertedID": "67b7d75e1c0f9ece68877066"
}

### Request for Order_items POST http://localhost:8000/OrderItems

{
  "order_items": [
    {
      "quantity": 3,
      "food_id": "67ad34a2fc209f120628890e",
      "order_id": "67b7d75e1c0f9ece68877066"
    },
    {
      "quantity": 1,
      "food_id": "67ab7289290928d729753d6b",
      "order_id": "67b7d75e1c0f9ece68877066"
    }
  ]
}

### Response 

   {
    "InsertedIDs": [
        "67b7da001c0f9ece68877067",
        "67b7da001c0f9ece68877068"
    ]
}


 ### Fetching the all Deatails GET http://localhost:8000/OrderItemsbyorder_id/67b7d75e1c0f9ece68877066

 ### Response 

    [
    {
        "OrderId": "67b7d75e1c0f9ece68877066",
        "TableId": "67b7dbef1c0f9ece68877069",
        "Table_number": "T2",
        "order_items": [
            [
                {
                    "Unit_price": 70,
                    "food_image": "http://example.com/pizza.jpg",
                    "foodname": "Pizza",
                    "quantity": 1
                }
            ],
            [
                {
                    "Unit_price": 10,
                    "food_image": "https://www.eatingwell.com/recipe/252379/classic-hamburger/",
                    "foodname": "Dosa",
                    "quantity": 3
                }
            ]
        ],
        "payment_due": 100
    }
]





