# Connecting to Valentine Wall Database with DBeaver

Quick guide to view and query the Valentine Wall SQLite database using DBeaver.

## Prerequisites

- DBeaver installed ([download here](https://dbeaver.io/download/))
- Valentine Wall application

## Step 1: Start the Application

The SQLite database is created by PocketBase on first startup.

**Windows:**
```powershell
.\start-backend.sh
```

**Linux/Mac:**
```bash
./start-backend.sh
```

Wait for the message: `Server started at http://127.0.0.1:8090`

The database file will be created at: `./pb_data/data.db`

## Step 2: Locate the Database File

Navigate to your valentine-wall project folder and confirm the database exists:

```
valentine-wall/
└── pb_data/
    └── data.db  ← SQLite database file
```

Full path example: `d:\valentine-wall\pb_data\data.db`

## Step 3: Connect DBeaver to SQLite

1. Open DBeaver
2. Click **Database** → **New Database Connection** (or click the plug icon)
3. Select **SQLite** from the list
4. Click **Next**

## Step 4: Add the Valentine Wall Database

1. In the connection settings:
   - **Path**: Click **Browse** and navigate to `pb_data/data.db` in your valentine-wall folder
   - Or paste the full path: `d:\valentine-wall\pb_data\data.db`
   
2. Click **Test Connection** to verify
3. If prompted to download SQLite drivers, click **Download**
4. Click **Finish**

## Step 5: View Tables and Data

Your database will appear in the Database Navigator:

```
SQLite - data.db
└── main
    └── Tables
        ├── _collections
        ├── _params
        ├── college_departments
        ├── gifts
        ├── messages
        ├── message_replies
        ├── user_details
        ├── users
        ├── virtual_transactions
        └── virtual_wallets
```

**To view data:**
- Double-click any table to open it
- Right-click → **View Data** to see table contents
- Use the SQL editor for custom queries: **SQL Editor** → **New SQL Script**

## Common Queries

View all messages:
```sql
SELECT * FROM messages ORDER BY created DESC LIMIT 100;
```

Check user details:
```sql
SELECT * FROM user_details;
```

View virtual transactions:
```sql
SELECT * FROM virtual_transactions ORDER BY created DESC;
```

## Tips

- Keep the backend running while using DBeaver for live data
- Refresh tables (F5) to see recent changes
- Be careful with DELETE/UPDATE queries in production environments
