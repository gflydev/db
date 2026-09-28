package migrate

import "strings"

// upTemplate and downTemplate are the bodies New writes into a fresh migration pair. Every
// line is an SQL comment, so an untouched pair is still valid, no-op SQL. They exist to make
// the two mistakes that are easy to miss and hard to fix later visible at the moment of
// writing: an up file whose objects the down file forgets to remove (e.g. a CREATE TYPE with
// no DROP TYPE, which makes the next up fail), and an edit to an up file after it was applied
// (which trips the checksum check). Keep both dialect-neutral: New does not know the dialect.
//
// "{{name}}" is replaced with the migration's full name by renderTemplate.
const upTemplate = `-- ============================================================================
-- Migration: {{name}}
-- Direction: UP (apply)
-- ============================================================================
--
-- Before you write:
--   * One concern per migration. Unrelated changes go in separate migrations, so each one can
--     be rolled back on its own.
--   * PostgreSQL runs this whole file in one transaction: any error rolls all of it back.
--     MySQL auto-commits every DDL statement, so a failure halfway leaves the statements
--     before it applied. On MySQL, keep a file to what can be re-run or cleaned up by hand.
--   * Once this file has been applied anywhere, NEVER edit it again, not even a comment.
--     Its checksum is recorded, and a changed file stops every later migration. Fix mistakes
--     forward, in a new migration.
--   * Some statements cannot run inside a transaction on PostgreSQL, e.g.
--     CREATE INDEX CONCURRENTLY, and ALTER TYPE ... ADD VALUE whose new value is then used in
--     the same file. Give them a migration of their own and check they work with this tool.
--
-- For every object you create below, write its removal in the .down.sql file NOW, in reverse
-- order. Keep this list up to date as you go; the down file must undo each line of it:
--   - [ ] ...
--
-- Write the statements in the sections below, in this order, and delete the sections you
-- do not need.

-- ----------------------------------------------------------------------------
-- 1. Types, extensions, sequences
--    (CREATE TYPE / CREATE EXTENSION / CREATE SEQUENCE)
--    Needed before the tables and columns that use them.
-- ----------------------------------------------------------------------------


-- ----------------------------------------------------------------------------
-- 2. Tables
--    (CREATE TABLE, with primary keys and NOT NULL / DEFAULT on each column)
-- ----------------------------------------------------------------------------


-- ----------------------------------------------------------------------------
-- 3. Changes to existing tables
--    (ALTER TABLE ... ADD / ALTER / RENAME / DROP COLUMN)
--    A new NOT NULL column on a table that has rows needs a DEFAULT, or a backfill in
--    section 5 before the NOT NULL is set. DROP COLUMN loses data: the down file cannot
--    bring it back.
-- ----------------------------------------------------------------------------


-- ----------------------------------------------------------------------------
-- 4. Constraints and indexes
--    (FOREIGN KEY / UNIQUE / CHECK / CREATE INDEX)
--    Name them explicitly (e.g. fk_orders_user_id, idx_orders_created_at), so the down file
--    can drop them by name.
-- ----------------------------------------------------------------------------


-- ----------------------------------------------------------------------------
-- 5. Data: backfill and seed
--    (INSERT / UPDATE)
--    Only data the schema needs to work. Say in the down file whether it removes this data.
-- ----------------------------------------------------------------------------


-- ----------------------------------------------------------------------------
-- 6. Functions, triggers, views
--    (CREATE FUNCTION / CREATE TRIGGER / CREATE VIEW)
-- ----------------------------------------------------------------------------

`

const downTemplate = `-- ============================================================================
-- Migration: {{name}}
-- Direction: DOWN (roll back)
-- ============================================================================
--
-- Goal: after this file runs, the schema is exactly what it was before the .up.sql file ran,
-- so that running the .up.sql file again succeeds.
--
-- Checklist:
--   * Every object the .up.sql file creates is removed here, in REVERSE order: the up file's
--     sections 6 -> 1. An object must be dropped after everything that uses it, e.g. drop a
--     column before the TYPE it uses, and a FOREIGN KEY before the table it points to.
--   * Easy to forget, because nothing fails until the next up:
--       CREATE TYPE / DOMAIN      ->  DROP TYPE / DROP DOMAIN
--       CREATE EXTENSION          ->  DROP EXTENSION (only if no other migration needs it)
--       CREATE SEQUENCE           ->  DROP SEQUENCE
--       CREATE FUNCTION / TRIGGER ->  DROP TRIGGER, then DROP FUNCTION
--       ADD CONSTRAINT / INDEX    ->  DROP CONSTRAINT / DROP INDEX
--       RENAME a -> b             ->  RENAME b -> a
--       ALTER COLUMN (type, default, NOT NULL)  ->  ALTER it back to the old definition
--   * Data the up file deleted or overwrote cannot be restored here. Say so in a comment.
--   * Before you commit, run up -> down -> up on a scratch database. The second up must pass.
--
-- Unlike the .up.sql file, this file is not checksummed, so it can still be fixed after the
-- migration was applied.
--
-- Write the statements in the sections below, in this order, and delete the sections you
-- do not need.

-- ----------------------------------------------------------------------------
-- 6. Functions, triggers, views
--    (DROP VIEW / DROP TRIGGER / DROP FUNCTION)
-- ----------------------------------------------------------------------------


-- ----------------------------------------------------------------------------
-- 5. Data: backfill and seed
--    (DELETE the seeded rows, or a comment saying why they stay)
-- ----------------------------------------------------------------------------


-- ----------------------------------------------------------------------------
-- 4. Constraints and indexes
--    (DROP CONSTRAINT / DROP INDEX)
-- ----------------------------------------------------------------------------


-- ----------------------------------------------------------------------------
-- 3. Changes to existing tables
--    (ALTER TABLE ... DROP the added columns, RENAME back, ALTER back)
-- ----------------------------------------------------------------------------


-- ----------------------------------------------------------------------------
-- 2. Tables
--    (DROP TABLE)
-- ----------------------------------------------------------------------------


-- ----------------------------------------------------------------------------
-- 1. Types, extensions, sequences
--    (DROP TYPE / DROP EXTENSION / DROP SEQUENCE)
--    Last: only after every column that uses them is gone.
-- ----------------------------------------------------------------------------

`

// renderTemplate fills a migration file template with the migration's full name.
func renderTemplate(tmpl, name string) string {
	return strings.ReplaceAll(tmpl, "{{name}}", name)
}
