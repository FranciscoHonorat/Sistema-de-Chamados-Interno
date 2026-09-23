-- Several replicas of the app may run the outbox relay at the same time. Each
-- one claims a batch by leasing it (locked_until) with FOR UPDATE SKIP LOCKED,
-- so an event is never delivered by two relays at once, and a relay that dies
-- mid-batch only holds its events until the lease expires.
ALTER TABLE outbox_events ADD COLUMN IF NOT EXISTS locked_until TIMESTAMPTZ;
