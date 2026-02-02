-- Post flags: any user can flag a post/idea for inappropriateness; flags visible only to moderators
CREATE TABLE IF NOT EXISTS post_flags (
    id          SERIAL PRIMARY KEY,
    tenant_id   INT NOT NULL,
    post_id     INT NOT NULL,
    user_id     INT NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    reason      TEXT,
    CONSTRAINT post_flags_tenant_id_fkey FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    CONSTRAINT post_flags_post_id_fkey FOREIGN KEY (post_id) REFERENCES posts(id),
    CONSTRAINT post_flags_user_fkey FOREIGN KEY (tenant_id, user_id) REFERENCES users(tenant_id, id),
    CONSTRAINT post_flags_unique UNIQUE (tenant_id, post_id, user_id)
);

CREATE INDEX IF NOT EXISTS post_flags_tenant_post ON post_flags (tenant_id, post_id);
