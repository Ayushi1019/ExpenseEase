ExpenseEase: A PERSONAL FINANCE MANAGER

Project Objective
The objective of this project is to build a personal finance manager-EXPENSEEASE that allows users to track their income, expenses, and other financial transactions, and generate reports and insights to help them make better financial decisions.

Requirements:
Go (version 1.20.2)
Node.js (version 16.13.2)
PostgreSQL (version 15)

Installion:

1. Install the required version of Go using https://golang.org/dl/. You can verify your installation by executing this command in your terminal "go version".
2. Install the required version of Node.js using https://nodejs.org/. To verify the installation was successfull, type "node -v" and "npm -v" in your terminal window.
3. Download the latests version of PostgreSQL from https://www.postgresql.org/download/ and install it. To verify that PostgreSQL is running correctly, you can use the "psql" command-line tool. Open a command prompt or terminal window and type "psql". If PostgreSQL is running correctly, you will be prompted to enter your password and then you will see the "psql" prompt.
4. Clone the repository: <br />
   git clone https://github.com/Ayushi1019/ExpenseEase.git <br />
   cd ExpenseEase <br />
5. Setting up the backend: <br />
   2.1 Install Go dependencies:
   > $ go mod download
   > 2.2 Create a .env file in the root directory of the project and add the following:
   > DATABASE_URL=postgres://user:password@localhost:port/<database-name>
   > Replace user, password, and <database-name> with your PostgreSQL credentials
   > 2.3 Run the backend server:
   > $ go run main.go

Frontend Setup: <br />

1. Install Node.js dependencies: <br />
   $ cd client <br />
   $ npm install <br />
2. Start the development server: <br />
   $ npm start <br />

Open the application in your browser at http://localhost:3000. <br />

Usage <br />
Users can track their expenses according to the various categories mentioned.
User can specify their monthly budget with respect to different categories of expense.
User will be able to track their income history
User will be able to analyze their budget and eventually improve their financial situation.

Contributions <br />
Ayushi Mundra (801208050)
Harsh Raval (801257980)
Shivashish Naramdeo (801208044)
Urvashi Murari (801205124)
Vidit Sethi (801203308)
