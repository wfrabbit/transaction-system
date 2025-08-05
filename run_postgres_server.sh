docker run --name local-postgres \
  -e POSTGRES_USER=wangfeng \
  -e POSTGRES_PASSWORD=wangfeng \
  -e POSTGRES_DB=transaction_db \
  -v pgdata:"/var/lib/postgresql/data" \
  -p 5432:5432 \
  -d postgres:15
