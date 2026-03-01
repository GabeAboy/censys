Why did I choose Postgresql?

PostgreSQL was chosen for this asset management system for several key reasons:
1. The data is inherently relational (assets → ports, assets → tags). PostgreSQL excels at modeling these one-to-many relationships with foreign keys, joins, and referential integrity. CASCADE deletes ensure that when an asset is removed, all associated ports and tags are automatically cleaned up.
2. Native support for arrays (used for ports and tags aggregation), JSONB for flexible metadata, and custom ENUM types (risk_level)
3. Critical for maintaining data integrity when creating assets with related ports and tags in transactions
4. Support for composite indexes, partial indexes, and GIN indexes for array operations and full-text search
5. Built-in array_agg() and COALESCE() functions simplify complex queries for joining assets with their ports and tags
6. Robust constraint system (UNIQUE, CHECK, FOREIGN KEY) prevents duplicate ports/tags and enforces data quality
7. Proven performance at scale with proper indexing and query optimization

Improvements For Prod
- Rate limiting (per IP/user)
- Cache layer (Redis for frequently accessed assets)
- Metrics and monitoring (Prometheus/Grafana)
- Add proper migration tool (golang-migrate)
- Add unit & integration tests
- Standardized logging (structured logging with levels)
- API documentation (Swagger/OpenAPI)
- Input sanitization and SQL injection prevention
- Graceful shutdown handling
- Background job processing for heavy operations
- Multi-tenancy support
