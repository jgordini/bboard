-- Pinned posts: staff can pin/unpin; pinned posts appear first in lists
ALTER TABLE posts ADD COLUMN IF NOT EXISTS pinned_at TIMESTAMPTZ NULL;
ALTER TABLE posts ADD COLUMN IF NOT EXISTS pinned_by_id INT NULL;

DO $$
BEGIN
  IF NOT EXISTS (
    SELECT 1 FROM information_schema.table_constraints
    WHERE constraint_name = 'posts_pinned_by_id_fkey'
    AND table_name = 'posts'
  ) THEN
    ALTER TABLE posts
      ADD CONSTRAINT posts_pinned_by_id_fkey
      FOREIGN KEY (pinned_by_id, tenant_id) REFERENCES users(id, tenant_id);
  END IF;
END $$;
