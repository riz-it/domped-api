create_table:
	migrate -database "postgresql://postgres.nrtkjdcoqmhjsxdwlold:dompeddatabase123@aws-0-ap-southeast-1.pooler.supabase.com:6543/postgres?search_path=public" -path internal/database/migration -verbose up
drop_table:
	migrate -database "postgresql://postgres.nrtkjdcoqmhjsxdwlold:dompeddatabase123@aws-0-ap-southeast-1.pooler.supabase.com:6543/postgres?search_path=public" -path internal/database/migration -verbose down
