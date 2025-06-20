-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE SCHEMA IF NOT EXISTS public;
CREATE SCHEMA IF NOT EXISTS core;

-- Create Tables
CREATE TABLE core.organizations (
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  name text NOT NULL,
  domain text UNIQUE,
  is_active boolean DEFAULT false,
  created_at timestamp DEFAULT now(),
  updated_at timestamp DEFAULT now()
);

CREATE TABLE core.users (
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  organization_id uuid REFERENCES core.organizations(id) ON DELETE CASCADE,
  roles text NOT NULL DEFAULT 'user',
  name text NOT NULL,
  email text UNIQUE NOT NULL,
  email_verified boolean DEFAULT false,
  image text,
  is_active boolean DEFAULT true,
  created_at timestamp DEFAULT now(),
  updated_at timestamp DEFAULT now()
);

-- This table is used for users sign-up inside an organization, the org must create a CODE and share it with the user
CREATE TABLE core.invitations (
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  organization_id uuid REFERENCES core.organizations(id) ON DELETE CASCADE,
  invited_by uuid REFERENCES core.users(id) ON DELETE CASCADE,
  expires_at timestamp NOT NULL,
  code text NOT NULL UNIQUE,
  invited_email  text NOT NULL unique,
  created_at timestamp DEFAULT now()
);

-- ACCOUNTS TABLE
CREATE TABLE core.accounts (
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  user_id uuid  NOT NULL REFERENCES core.users(id) ON DELETE CASCADE,
  provider_id text NOT NULL DEFAULT 'credential', -- It can be "google", "tiktok" etc...
  access_token text,
  refresh_token text,
  access_token_expires_at text,
  refresh_token_expires_at text,
  id_token text,
  scope text,
  password text,
  created_at timestamp DEFAULT now(),
  updated_at timestamp DEFAULT now()
);

CREATE TABLE core.sessions (
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  user_id uuid NOT NULL REFERENCES core.users(id) ON DELETE CASCADE,
  access_token text NOT NULL UNIQUE,
  refresh_token text NOT NULL UNIQUE,
  access_token_expires_at timestamp NOT NULL,
  refresh_token_expires_at timestamp NOT NULL,
  created_at timestamp DEFAULT now(),
  updated_at timestamp DEFAULT now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP SCHEMA IF EXISTS core CASCADE;
CREATE SCHEMA IF NOT EXISTS public;
-- +goose StatementEnd
