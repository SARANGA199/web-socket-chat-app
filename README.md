# 🗨️ WebSocket Chat App  

A real-time one-to-one messaging application using **WebSockets**, **Go**, **MongoDB**, and **React**. This project enables instant communication between users with secure message storage and seamless updates.  

## 🚀 Features  
✅ Real-time one-to-one messaging using WebSockets  
✅ Secure JWT authentication  
✅ MongoDB for message storage  
✅ Multi-user WebSocket connections  
✅ React frontend for an interactive chat interface  

## 🛠️ Tech Stack  
- **Backend:** Go (WebSocket, JWT, MongoDB)  
- **Frontend:** React  
- **Database:** MongoDB  

## 📌 How It Works  
1. Open two browser windows at `http://localhost:3000`.  
2. Enter `userId` and `receiverId` in the text boxes and connect as chat users.  
3. Send and receive messages in real time.  

## 📂 Setup & Installation  

### 🔹 Clone the Repository  
```bash
git clone https://github.com/SARANGA199/web-socket-chat-app.git
cd web-socket-chat-app

🏗️ Backend Setup
 1️⃣ Install Go Dependencies
   Ensure you have Go installed. If not, download Go here.

 2️⃣ Run the following command inside the project folder to install dependencies:
   go mod tidy
 
 3️⃣ Start the WebSocket Server
   Run the backend server using
   go run main.go

The server should now be running on http://localhost:8080.

💻 Frontend Setup
 1️⃣ Navigate to the Frontend Folder
   cd frontend

 2️⃣ Install Dependencies
   Now, install the required packages
   yarn install

 3️⃣ Start the React App
  Run the frontend with
  yarn start

The application should now be accessible at http://localhost:3000.
