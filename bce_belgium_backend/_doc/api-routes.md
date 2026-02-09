# API Routes Documentation

## Base URL

`http://localhost:8080/api`

## Health Check

- **GET** `/health` - API status check

## Tables Routes

- **GET** `/tables` - List all tables with row counts
- **GET** `/tables/structure` - Complete database structure with all columns
- **GET** `/tables/:name/info` - Specific table information
- **GET** `/tables/:name/columns` - Column details for a table

## Data Routes

- **GET** `/data/:table/preview?limit=5` - Preview table data
- **GET** `/data/:table/values/:column?limit=20` - Column value statistics

## Search Routes

- **GET** `/search/:table/:column?q=term&limit=50` - Search in table column
- **GET** `/count/:table/:column?q=term` - Count matching rows

## Export Routes

- **GET** `/export/:table` - Export full table as CSV
- **GET** `/export/:table?column=col&search=term&limit=10000` - Export filtered data

## Examples

### Get all tables

```http
GET /api/tables
```

### Search companies by NACE code

```http
GET /api/search/activity/nace_code?q=62020&limit=100
```

### Export IT companies

```http
GET /api/export/activity?column=nacecode&search=62020&limit=5000
```

### Preview enterprise table

```http
GET /api/data/enterprise/preview?limit=10
```
