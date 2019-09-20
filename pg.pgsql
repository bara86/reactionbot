--
-- PostgreSQL database dump
--

-- Dumped from database version 11.5
-- Dumped by pg_dump version 11.5

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

SET default_with_oids = false;

--
-- Name: emojis; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.emojis (
    name text NOT NULL
);


--
-- Name: groups; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.groups (
    iduser text NOT NULL,
    groupname text NOT NULL
);


--
-- Name: groupsemojis; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.groupsemojis (
    groupname text NOT NULL,
    userid text NOT NULL,
    emojiname text NOT NULL
);


--
-- Name: temporaryusertokens; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.temporaryusertokens (
    uuid text NOT NULL,
    userid text NOT NULL
);


--
-- Name: users; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.users (
    id text NOT NULL,
    token text NOT NULL
);


--
-- Name: emojis emojis_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.emojis
    ADD CONSTRAINT emojis_pkey PRIMARY KEY (name);


--
-- Name: groups groups_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.groups
    ADD CONSTRAINT groups_pkey PRIMARY KEY (iduser, groupname);


--
-- Name: groupsemojis groupsemojis_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.groupsemojis
    ADD CONSTRAINT groupsemojis_pkey PRIMARY KEY (userid, groupname, emojiname);


--
-- Name: temporaryusertokens temporaryusertokens_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.temporaryusertokens
    ADD CONSTRAINT temporaryusertokens_pkey PRIMARY KEY (uuid, userid);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: groups groups_iduser_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.groups
    ADD CONSTRAINT groups_iduser_fkey FOREIGN KEY (iduser) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: groupsemojis groupsemojis_emojiname_fkey1; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.groupsemojis
    ADD CONSTRAINT groupsemojis_emojiname_fkey1 FOREIGN KEY (emojiname) REFERENCES public.emojis(name) ON DELETE CASCADE;


--
-- Name: groupsemojis groupsemojis_groupname_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.groupsemojis
    ADD CONSTRAINT groupsemojis_groupname_fkey FOREIGN KEY (groupname, userid) REFERENCES public.groups(groupname, iduser) ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--

