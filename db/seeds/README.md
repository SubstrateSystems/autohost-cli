# Database Seeding

This directory contains seeding files to populate the database with initial data.

## Structure

- Seeding files must follow the pattern `XXX_name.sql` where `XXX` is an order number (e.g., `001_catalog_apps.sql`)
- Seeds are automatically executed after migrations when initializing the application
- The system tracks which seeds have already been executed using the `_seeds` table
- Seeds are idempotent - they can be executed multiple times without causing errors

## Current Seeds

### 001_catalog_apps.sql
Inserts available applications into the catalog:
- bookstack: Platform for organizing and storing information
- nextcloud: Software suite for file hosting services
- redis: In-memory database for caching and messaging

### 002_catalog_apps.sql
Additional catalog applications:
- mysql: Relational database management system
- postgres: Advanced open-source relational database

## Adding New Seeds

1. Create a file with format `XXX_name.sql` in this directory
2. Use `INSERT OR IGNORE` to avoid duplicates
3. Seeds will be automatically executed on the next application startup

## Notes

- Seeds are executed in alphabetical order by filename
- Once executed, a seed is not re-executed
- If you need to modify existing data, create a new seed with a higher number
