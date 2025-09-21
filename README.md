# E-Wallet

**"The user has a balance in the application wallet and wants to disburse the balance."**

## Features

Two main functions were implemented: **TopUp** and **Withdraw**, with the following use cases:

### **TopUp Use Cases**
- Can only top up a maximum of **10 million IDR per day**.
- Can only top up if the **amount > 0**.
- Can only top up if the source is from **"e-wallet"**, **"bank_transfer"**, or **"credit_card"**.

### **Withdraw Use Cases**
- Can only withdraw a maximum of **10 million IDR per day**.
- Can only withdraw if the **PIN is correct**.
- Can only withdraw if the **amount > 0**.
- Can only withdraw if the **balance is sufficient**.
- Can only withdraw during **operational hours**.
- Can only withdraw if the **destination bank is available**.
- Can only withdraw if there are no **pending withdrawals**.
- Withdrawals outside operational hours will be marked as **pending** and verified during the next operational hours using a **cron job**.

### **Admin Endpoint for Pending Withdrawals**
A new endpoint was added to **update the status of pending withdrawals to "approved"**, which can only be executed by an **admin**.

### **Login Endpoint**
A **login endpoint** was implemented to generate a **token**. This token is used as a **Bearer Token** to authorize access to the following endpoints:
- **TopUp**
- **Withdraw**
- **Update Pending Withdraw**

## How to Run the Application

Follow these steps to set up and run the application:

1. **Set up the database**:
   - Use the SQL script located in the folder: `db/db-money_movement.sql` as dummy data.

2. **Configure the database and environment variables**:
   - Update the `.env` file with your database credentials and configuration.

3. **Install dependencies**:
   ```bash
   go mod tidy

4. **Run**
    ```bash
   go run main.go

5. **Import postman collection**
  - Use this postman collection located in folder: `postman_collection/Money Movement DB.postman_collection.json` as api testing.
