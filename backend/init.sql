-- Create the database
CREATE DATABASE cfbox;

-- Connect to the newly created database
\c cfbox;

-- Create the citext extension if not already installed
CREATE EXTENSION IF NOT EXISTS citext;

-- Create a user with administration privileges
CREATE USER cfbox_admin WITH PASSWORD 'pa55w0rd';
ALTER DATABASE cfbox OWNER TO cfbox_admin;