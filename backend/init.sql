-- Create the database
CREATE DATABASE cfbox;

-- Connect to the newly created database
\c cfbox;

-- Create a role with login and password
CREATE ROLE cfbox_admin WITH LOGIN PASSWORD 'pa55w0rd';

-- Create the citext extension if not already installed
CREATE EXTENSION IF NOT EXISTS citext;

