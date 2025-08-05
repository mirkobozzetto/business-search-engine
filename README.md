# BCE API - User Guide

## Overview

API for searching Belgian companies with 47M+ rows spread across 10 tables. System based on `entitynumber` as universal key.

## Main Tables

```
entitynumber (universal key)
├── activity (36M rows) - NACE codes, classifications
├── address (2.8M rows) - Complete addresses
├── contact (683k rows) - Email, phone, web, fax
├── denomination (3.3M rows) - Names FR/NL/DE
├── enterprise (1.9M rows) - Status, legal form, dates
└── establishment (1.6M rows) - Local units
```

## How it Works

### Simple Principle

1. **Entry point** → Retrieves entitynumbers
2. **Automatic enrichment** → Joins the 6 tables
3. **Redis cache** → Stores everything
4. **Multi-criteria** → Cache intersection

### Important Prerequisites

To use `/search/multi`, you MUST create individual caches first.

## Available Endpoints

```bash
# Search by NACE code
curl "localhost:8080/api/companies/search/nace?code=62010" | jq

# Search by company name
curl "localhost:8080/api/companies/search/denomination?q=informatique" | jq

# Search by postal code
curl "localhost:8080/api/companies/search/zipcode?q=1000" | jq

# Search by creation date
curl "localhost:8080/api/companies/search/startdate?from=01-01-2025" | jq

# Multi-criteria (after creating caches)
curl "localhost:8080/api/companies/search/multi?nace=62010&zipcode=1000" | jq
```

## Basic Workflow

### Example: Health coaches created in 2025

**Step 1: NACE Cache**

```bash
curl "localhost:8080/api/companies/search/nace?code=56309" | jq
# Creates cache for all companies with NACE 56309
```

**Step 2: Date Cache**

```bash
curl "localhost:8080/api/companies/search/startdate?from=01-01-2025" | jq
# Creates cache for all companies created in 2025+
```

**Step 3: Intersection**

```bash
curl "localhost:8080/api/companies/search/multi?nace=56309&startdate_from=01-01-2025" | jq
# Intersection of 2 caches = companies NACE 56309 created in 2025
```

## Common Errors

### Missing cache

```bash
# ❌ Error
curl "localhost:8080/api/companies/search/multi?nace=62010&zipcode=1000" | jq
# "NACE cache not found: 62010. Please search by NACE first"

# ✅ Solution
curl "localhost:8080/api/companies/search/nace?code=62010" | jq
curl "localhost:8080/api/companies/search/zipcode?q=1000" | jq
curl "localhost:8080/api/companies/search/multi?nace=62010&zipcode=1000" | jq
```

### Incorrect parameter

```bash
# ❌ Error - Parameter "nace="
curl "localhost:8080/api/companies/search/nace?nace=62010" | jq

# ✅ Correct - Parameter "code="
curl "localhost:8080/api/companies/search/nace?code=62010" | jq
```

## Response Structure

Each returned company contains:

- **Basic data**: entitynumber, denomination, zipcode, city
- **Contacts**: email, web, tel, fax
- **Details**: status, juridical_form, start_date
- **Complete data**: denominations[], addresses[], contacts[], activities[], etc.

## Simple Use Cases

### All companies in a sector

```bash
curl "localhost:8080/api/companies/search/nace?code=62010" | jq
```

### Companies in a city

```bash
curl "localhost:8080/api/companies/search/zipcode?q=1000" | jq
```

### Sector + city intersection

```bash
curl "localhost:8080/api/companies/search/nace?code=62010" | jq
curl "localhost:8080/api/companies/search/zipcode?q=1000" | jq
curl "localhost:8080/api/companies/search/multi?nace=62010&zipcode=1000" | jq
```

### New companies in a sector

```bash
curl "localhost:8080/api/companies/search/nace?code=62010" | jq
curl "localhost:8080/api/companies/search/startdate?from=01-01-2024" | jq
curl "localhost:8080/api/companies/search/multi?nace=62010&startdate_from=01-01-2024" | jq
```
