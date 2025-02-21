# Restaurant-Management
Restaurant Management API (Golang, Gin, MongoDB)  This is a RESTful API for a restaurant management system built using Golang, Gin framework, and MongoDB. It includes features for managing food items, orders, order items, invoices, tables, and menus.

## User Authentication (Signup & Login)
Signup
Endpoint: POST /users/signup
Request:

{
  "first_name": "Vj",
  "last_name": "Raj",
  "password": "123456",
  "email": "vj@gmail.com",
  "phone": "9876543210"
}
Response:

json
Copy
Edit
{
  "InsertedID": "67b7ceb740642997053061aa"
}
Login
Endpoint: POST /users/login
Request:

json
Copy
Edit
{
  "email": "vj@gmail.com",
  "password": "123456"
}
Response:

json
Copy
Edit
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
