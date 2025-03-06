--
-- PostgreSQL database dump
--

-- Dumped from database version 17.4 (Debian 17.4-1.pgdg120+2)
-- Dumped by pg_dump version 17.4

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET transaction_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: resource_histories; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.resource_histories (
    id bigint NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    resource_id text,
    current_token text,
    previous_token text
);


ALTER TABLE public.resource_histories OWNER TO postgres;

--
-- Name: resource_histories_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.resource_histories_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.resource_histories_id_seq OWNER TO postgres;

--
-- Name: resource_histories_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.resource_histories_id_seq OWNED BY public.resource_histories.id;


--
-- Name: resources; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.resources (
    id bigint NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    resource_id text NOT NULL,
    consistency_token text
);


ALTER TABLE public.resources OWNER TO postgres;

--
-- Name: resources_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.resources_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.resources_id_seq OWNER TO postgres;

--
-- Name: resources_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.resources_id_seq OWNED BY public.resources.id;


--
-- Name: resource_histories id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.resource_histories ALTER COLUMN id SET DEFAULT nextval('public.resource_histories_id_seq'::regclass);


--
-- Name: resources id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.resources ALTER COLUMN id SET DEFAULT nextval('public.resources_id_seq'::regclass);


--
-- Data for Name: resource_histories; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.resource_histories (id, created_at, updated_at, deleted_at, resource_id, current_token, previous_token) FROM stdin;
1	2025-03-06 17:51:27.455649+00	2025-03-06 17:51:27.455649+00	\N	my_resource_one	myrandomconsistencytoken	
2	2025-03-06 17:51:27.47713+00	2025-03-06 17:51:27.47713+00	\N	another_resource_here	myrandomconsistencytoken	
3	2025-03-06 17:51:27.501406+00	2025-03-06 17:51:27.501406+00	\N	third_resource_this_is	myrandomconsistencytoken	
\.


--
-- Data for Name: resources; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.resources (id, created_at, updated_at, deleted_at, resource_id, consistency_token) FROM stdin;
1	2025-03-06 17:51:27.446549+00	2025-03-06 17:51:27.446549+00	\N	my_resource_one	myrandomconsistencytoken
2	2025-03-06 17:51:27.467608+00	2025-03-06 17:51:27.467608+00	\N	another_resource_here	myrandomconsistencytoken
3	2025-03-06 17:51:27.485807+00	2025-03-06 17:51:27.485807+00	\N	third_resource_this_is	myrandomconsistencytoken
\.


--
-- Name: resource_histories_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.resource_histories_id_seq', 9, true);


--
-- Name: resources_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.resources_id_seq', 9, true);


--
-- Name: resource_histories resource_histories_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.resource_histories
    ADD CONSTRAINT resource_histories_pkey PRIMARY KEY (id);


--
-- Name: resources resources_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.resources
    ADD CONSTRAINT resources_pkey PRIMARY KEY (id, resource_id);


--
-- Name: idx_resource_histories_deleted_at; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_resource_histories_deleted_at ON public.resource_histories USING btree (deleted_at);


--
-- Name: idx_resources_deleted_at; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_resources_deleted_at ON public.resources USING btree (deleted_at);


--
-- PostgreSQL database dump complete
--

