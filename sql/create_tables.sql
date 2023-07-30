-- PostgreSQL database dump

-- Dumped from database version 14.5 (Debian 14.5-1.pgdg110+1)
-- Dumped by pg_dump version 14.5 (Homebrew)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_table_access_method = heap;

-- Create the Users table
CREATE TABLE public.Users (
  id SERIAL PRIMARY KEY,
  username VARCHAR(255) NOT NULL UNIQUE,
  password VARCHAR(255) NOT NULL,
  is_admin BOOLEAN NOT NULL
);

-- Create the RequestedBooking table
CREATE TABLE public.RequestedBookings (
  id SERIAL PRIMARY KEY,
  username VARCHAR(255) NOT NULL,
  name VARCHAR(255) NOT NULL,
  start_date DATE NOT NULL,
  end_date DATE NOT NULL,
  unit_number VARCHAR(255) NOT NULL,
  start_time TIMESTAMPTZ NOT NULL,
  end_time TIMESTAMPTZ NOT NULL,
  purpose TEXT,
  facility TEXT,
  is_recurring BOOLEAN DEFAULT FALSE,  -- New is_recurring column
  recurring_weeks INT,
  FOREIGN KEY (username) REFERENCES public.Users (username)
);

-- Create the ApprovedBookings table
CREATE TABLE public.ApprovedBookings (
  id SERIAL PRIMARY KEY,
  username VARCHAR(255) NOT NULL,
  name VARCHAR(255) NOT NULL,
  start_date DATE NOT NULL,
  end_date DATE NOT NULL,
  unit_number VARCHAR(255) NOT NULL,
  start_time TIMESTAMPTZ NOT NULL,
  end_time TIMESTAMPTZ NOT NULL,
  purpose TEXT,
  facility TEXT,
  FOREIGN KEY (username) REFERENCES public.Users (username)
);

-- Create the RecurringBookings table
CREATE TABLE public.RecurringBookings (
  id SERIAL PRIMARY KEY,
  username VARCHAR(255) NOT NULL,
  name VARCHAR(255) NOT NULL,
  start_date DATE NOT NULL,
  end_date DATE NOT NULL,
  unit_number VARCHAR(255) NOT NULL,
  start_time TIMESTAMPTZ NOT NULL,
  end_time TIMESTAMPTZ NOT NULL,
  purpose TEXT,
  facility TEXT,
  FOREIGN KEY (username) REFERENCES public.Users (username)
);

-- PostgreSQL database dump complete
