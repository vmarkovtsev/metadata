--
-- PostgreSQL database dump
--

-- Dumped from database version 10.11

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

--
-- Name: github_issue_comments_versioned; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.github_issue_comments_versioned (
    sum256 character varying(64) NOT NULL,
    versions integer[],
    author_association text,
    body text,
    created_at timestamp with time zone,
    htmlurl text,
    id bigint,
    issue_number bigint NOT NULL,
    node_id text,
    repository_name text NOT NULL,
    repository_owner text NOT NULL,
    repository_fullname text NOT NULL,
    updated_at timestamp with time zone,
    user_id bigint NOT NULL,
    user_login text NOT NULL
);


--
-- Name: github_issue_comments; Type: VIEW; Schema: public; Owner: -
--

CREATE VIEW public.github_issue_comments AS
 SELECT github_issue_comments_versioned.author_association,
    github_issue_comments_versioned.body,
    github_issue_comments_versioned.created_at,
    github_issue_comments_versioned.htmlurl,
    github_issue_comments_versioned.id,
    github_issue_comments_versioned.issue_number,
    github_issue_comments_versioned.node_id,
    github_issue_comments_versioned.repository_name,
    github_issue_comments_versioned.repository_owner,
    github_issue_comments_versioned.repository_fullname,
    github_issue_comments_versioned.updated_at,
    github_issue_comments_versioned.user_id,
    github_issue_comments_versioned.user_login
   FROM public.github_issue_comments_versioned;


--
-- Name: github_issues_versioned; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.github_issues_versioned (
    sum256 character varying(64) NOT NULL,
    versions integer[],
    assignees text[] NOT NULL,
    body text,
    closed_at timestamp with time zone,
    closed_by_id bigint NOT NULL,
    closed_by_login text NOT NULL,
    comments bigint,
    created_at timestamp with time zone,
    htmlurl text,
    id bigint,
    labels text[] NOT NULL,
    locked boolean,
    milestone_id text NOT NULL,
    milestone_title text NOT NULL,
    node_id text,
    number bigint,
    repository_name text NOT NULL,
    repository_owner text NOT NULL,
    repository_fullname text NOT NULL,
    state text,
    title text,
    updated_at timestamp with time zone,
    user_id bigint NOT NULL,
    user_login text NOT NULL
);


--
-- Name: github_issues; Type: VIEW; Schema: public; Owner: -
--

CREATE VIEW public.github_issues AS
 SELECT github_issues_versioned.assignees,
    github_issues_versioned.body,
    github_issues_versioned.closed_at,
    github_issues_versioned.closed_by_id,
    github_issues_versioned.closed_by_login,
    github_issues_versioned.comments,
    github_issues_versioned.created_at,
    github_issues_versioned.htmlurl,
    github_issues_versioned.id,
    github_issues_versioned.labels,
    github_issues_versioned.locked,
    github_issues_versioned.milestone_id,
    github_issues_versioned.milestone_title,
    github_issues_versioned.node_id,
    github_issues_versioned.number,
    github_issues_versioned.repository_name,
    github_issues_versioned.repository_owner,
    github_issues_versioned.repository_fullname,
    github_issues_versioned.state,
    github_issues_versioned.title,
    github_issues_versioned.updated_at,
    github_issues_versioned.user_id,
    github_issues_versioned.user_login
   FROM public.github_issues_versioned;


--
-- Name: github_organizations_versioned; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.github_organizations_versioned (
    sum256 character varying(64) NOT NULL,
    versions integer[],
    avatar_url text,
    collaborators bigint,
    created_at timestamp with time zone,
    description text,
    email text,
    htmlurl text,
    id bigint,
    login text,
    name text,
    node_id text,
    owned_private_repos bigint,
    public_repos bigint,
    total_private_repos bigint,
    updated_at timestamp with time zone
);


--
-- Name: github_organizations; Type: VIEW; Schema: public; Owner: -
--

CREATE VIEW public.github_organizations AS
 SELECT github_organizations_versioned.avatar_url,
    github_organizations_versioned.collaborators,
    github_organizations_versioned.created_at,
    github_organizations_versioned.description,
    github_organizations_versioned.email,
    github_organizations_versioned.htmlurl,
    github_organizations_versioned.id,
    github_organizations_versioned.login,
    github_organizations_versioned.name,
    github_organizations_versioned.node_id,
    github_organizations_versioned.owned_private_repos,
    github_organizations_versioned.public_repos,
    github_organizations_versioned.total_private_repos,
    github_organizations_versioned.updated_at
   FROM public.github_organizations_versioned;


--
-- Name: github_pull_request_comments_versioned; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.github_pull_request_comments_versioned (
    sum256 character varying(64) NOT NULL,
    versions integer[],
    author_association text,
    body text,
    commit_id text,
    created_at timestamp with time zone,
    diff_hunk text,
    htmlurl text,
    id bigint,
    in_reply_to bigint,
    node_id text,
    original_commit_id text,
    original_position bigint,
    path text,
    "position" bigint,
    pull_request_number bigint NOT NULL,
    pull_request_review_id bigint,
    repository_name text NOT NULL,
    repository_owner text NOT NULL,
    repository_fullname text NOT NULL,
    updated_at timestamp with time zone,
    user_id bigint NOT NULL,
    user_login text NOT NULL
);


--
-- Name: github_pull_request_comments; Type: VIEW; Schema: public; Owner: -
--

CREATE VIEW public.github_pull_request_comments AS
 SELECT github_pull_request_comments_versioned.author_association,
    github_pull_request_comments_versioned.body,
    github_pull_request_comments_versioned.commit_id,
    github_pull_request_comments_versioned.created_at,
    github_pull_request_comments_versioned.diff_hunk,
    github_pull_request_comments_versioned.htmlurl,
    github_pull_request_comments_versioned.id,
    github_pull_request_comments_versioned.in_reply_to,
    github_pull_request_comments_versioned.node_id,
    github_pull_request_comments_versioned.original_commit_id,
    github_pull_request_comments_versioned.original_position,
    github_pull_request_comments_versioned.path,
    github_pull_request_comments_versioned."position",
    github_pull_request_comments_versioned.pull_request_number,
    github_pull_request_comments_versioned.pull_request_review_id,
    github_pull_request_comments_versioned.repository_name,
    github_pull_request_comments_versioned.repository_owner,
    github_pull_request_comments_versioned.repository_fullname,
    github_pull_request_comments_versioned.updated_at,
    github_pull_request_comments_versioned.user_id,
    github_pull_request_comments_versioned.user_login
   FROM public.github_pull_request_comments_versioned;


--
-- Name: github_pull_request_reviews_versioned; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.github_pull_request_reviews_versioned (
    sum256 character varying(64) NOT NULL,
    versions integer[],
    body text,
    commit_id text,
    htmlurl text,
    id bigint,
    node_id text,
    pull_request_number bigint NOT NULL,
    repository_name text NOT NULL,
    repository_owner text NOT NULL,
    repository_fullname text NOT NULL,
    state text,
    submitted_at timestamp with time zone,
    user_id bigint NOT NULL,
    user_login text NOT NULL
);


--
-- Name: github_pull_request_reviews; Type: VIEW; Schema: public; Owner: -
--

CREATE VIEW public.github_pull_request_reviews AS
 SELECT github_pull_request_reviews_versioned.body,
    github_pull_request_reviews_versioned.commit_id,
    github_pull_request_reviews_versioned.htmlurl,
    github_pull_request_reviews_versioned.id,
    github_pull_request_reviews_versioned.node_id,
    github_pull_request_reviews_versioned.pull_request_number,
    github_pull_request_reviews_versioned.repository_name,
    github_pull_request_reviews_versioned.repository_owner,
    github_pull_request_reviews_versioned.repository_fullname,
    github_pull_request_reviews_versioned.state,
    github_pull_request_reviews_versioned.submitted_at,
    github_pull_request_reviews_versioned.user_id,
    github_pull_request_reviews_versioned.user_login
   FROM public.github_pull_request_reviews_versioned;


--
-- Name: github_pull_requests_versioned; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.github_pull_requests_versioned (
    sum256 character varying(64) NOT NULL,
    versions integer[],
    additions bigint,
    assignees text[] NOT NULL,
    author_association text,
    base_ref text NOT NULL,
    base_repository_name text NOT NULL,
    base_repository_owner text NOT NULL,
    base_repository_fullname text NOT NULL,
    base_sha text NOT NULL,
    base_user text NOT NULL,
    body text,
    changed_files bigint,
    closed_at timestamp with time zone,
    comments bigint,
    commits bigint,
    created_at timestamp with time zone,
    deletions bigint,
    head_ref text NOT NULL,
    head_repository_name text NOT NULL,
    head_repository_owner text NOT NULL,
    head_repository_fullname text NOT NULL,
    head_sha text NOT NULL,
    head_user text NOT NULL,
    htmlurl text,
    id bigint,
    labels text[] NOT NULL,
    maintainer_can_modify boolean,
    merge_commit_sha text,
    mergeable boolean,
    merged boolean,
    merged_at timestamp with time zone,
    merged_by_id bigint NOT NULL,
    merged_by_login text NOT NULL,
    milestone_id text NOT NULL,
    milestone_title text NOT NULL,
    node_id text,
    number bigint,
    repository_name text NOT NULL,
    repository_owner text NOT NULL,
    repository_fullname text NOT NULL,
    review_comments bigint,
    state text,
    title text,
    updated_at timestamp with time zone,
    user_id bigint NOT NULL,
    user_login text NOT NULL
);


--
-- Name: github_pull_requests; Type: VIEW; Schema: public; Owner: -
--

CREATE VIEW public.github_pull_requests AS
 SELECT github_pull_requests_versioned.additions,
    github_pull_requests_versioned.assignees,
    github_pull_requests_versioned.author_association,
    github_pull_requests_versioned.base_ref,
    github_pull_requests_versioned.base_repository_name,
    github_pull_requests_versioned.base_repository_owner,
    github_pull_requests_versioned.base_repository_fullname,
    github_pull_requests_versioned.base_sha,
    github_pull_requests_versioned.base_user,
    github_pull_requests_versioned.body,
    github_pull_requests_versioned.changed_files,
    github_pull_requests_versioned.closed_at,
    github_pull_requests_versioned.comments,
    github_pull_requests_versioned.commits,
    github_pull_requests_versioned.created_at,
    github_pull_requests_versioned.deletions,
    github_pull_requests_versioned.head_ref,
    github_pull_requests_versioned.head_repository_name,
    github_pull_requests_versioned.head_repository_owner,
    github_pull_requests_versioned.head_repository_fullname,
    github_pull_requests_versioned.head_sha,
    github_pull_requests_versioned.head_user,
    github_pull_requests_versioned.htmlurl,
    github_pull_requests_versioned.id,
    github_pull_requests_versioned.labels,
    github_pull_requests_versioned.maintainer_can_modify,
    github_pull_requests_versioned.merge_commit_sha,
    github_pull_requests_versioned.mergeable,
    github_pull_requests_versioned.merged,
    github_pull_requests_versioned.merged_at,
    github_pull_requests_versioned.merged_by_id,
    github_pull_requests_versioned.merged_by_login,
    github_pull_requests_versioned.milestone_id,
    github_pull_requests_versioned.milestone_title,
    github_pull_requests_versioned.node_id,
    github_pull_requests_versioned.number,
    github_pull_requests_versioned.repository_name,
    github_pull_requests_versioned.repository_owner,
    github_pull_requests_versioned.repository_fullname,
    github_pull_requests_versioned.review_comments,
    github_pull_requests_versioned.state,
    github_pull_requests_versioned.title,
    github_pull_requests_versioned.updated_at,
    github_pull_requests_versioned.user_id,
    github_pull_requests_versioned.user_login
   FROM public.github_pull_requests_versioned;


--
-- Name: github_repositories_versioned; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.github_repositories_versioned (
    sum256 character varying(64) NOT NULL,
    versions integer[],
    allow_merge_commit boolean,
    allow_rebase_merge boolean,
    allow_squash_merge boolean,
    archived boolean,
    clone_url text,
    created_at timestamp with time zone,
    default_branch text,
    description text,
    disabled boolean,
    fork boolean,
    forks_count bigint,
    fullname text,
    has_issues boolean,
    has_wiki boolean,
    homepage text,
    htmlurl text,
    id bigint,
    language text,
    name text,
    node_id text,
    open_issues_count bigint,
    owner_id bigint NOT NULL,
    owner_login text NOT NULL,
    owner_type text NOT NULL,
    private boolean,
    pushed_at timestamp with time zone,
    sshurl text,
    stargazers_count bigint,
    topics text[] NOT NULL,
    updated_at timestamp with time zone,
    watchers_count bigint
);


--
-- Name: github_repositories; Type: VIEW; Schema: public; Owner: -
--

CREATE VIEW public.github_repositories AS
 SELECT github_repositories_versioned.allow_merge_commit,
    github_repositories_versioned.allow_rebase_merge,
    github_repositories_versioned.allow_squash_merge,
    github_repositories_versioned.archived,
    github_repositories_versioned.clone_url,
    github_repositories_versioned.created_at,
    github_repositories_versioned.default_branch,
    github_repositories_versioned.description,
    github_repositories_versioned.disabled,
    github_repositories_versioned.fork,
    github_repositories_versioned.forks_count,
    github_repositories_versioned.fullname,
    github_repositories_versioned.has_issues,
    github_repositories_versioned.has_wiki,
    github_repositories_versioned.homepage,
    github_repositories_versioned.htmlurl,
    github_repositories_versioned.id,
    github_repositories_versioned.language,
    github_repositories_versioned.name,
    github_repositories_versioned.node_id,
    github_repositories_versioned.open_issues_count,
    github_repositories_versioned.owner_id,
    github_repositories_versioned.owner_login,
    github_repositories_versioned.owner_type,
    github_repositories_versioned.private,
    github_repositories_versioned.pushed_at,
    github_repositories_versioned.sshurl,
    github_repositories_versioned.stargazers_count,
    github_repositories_versioned.topics,
    github_repositories_versioned.updated_at,
    github_repositories_versioned.watchers_count
   FROM public.github_repositories_versioned;


--
-- Name: github_users_versioned; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.github_users_versioned (
    sum256 character varying(64) NOT NULL,
    versions integer[],
    avatar_url text,
    bio text,
    company text,
    created_at timestamp with time zone,
    email text,
    followers bigint,
    following bigint,
    hireable boolean,
    htmlurl text,
    id bigint,
    location text,
    login text,
    name text,
    node_id text,
    organization_id bigint NOT NULL,
    organization_login text NOT NULL,
    owned_private_repos bigint,
    private_gists bigint,
    public_gists bigint,
    public_repos bigint,
    total_private_repos bigint,
    updated_at timestamp with time zone
);


--
-- Name: github_users; Type: VIEW; Schema: public; Owner: -
--

CREATE VIEW public.github_users AS
 SELECT github_users_versioned.avatar_url,
    github_users_versioned.bio,
    github_users_versioned.company,
    github_users_versioned.created_at,
    github_users_versioned.email,
    github_users_versioned.followers,
    github_users_versioned.following,
    github_users_versioned.hireable,
    github_users_versioned.htmlurl,
    github_users_versioned.id,
    github_users_versioned.location,
    github_users_versioned.login,
    github_users_versioned.name,
    github_users_versioned.node_id,
    github_users_versioned.organization_id,
    github_users_versioned.organization_login,
    github_users_versioned.owned_private_repos,
    github_users_versioned.private_gists,
    github_users_versioned.public_gists,
    github_users_versioned.public_repos,
    github_users_versioned.total_private_repos,
    github_users_versioned.updated_at
   FROM public.github_users_versioned;


--
-- Name: issue_comments; Type: MATERIALIZED VIEW; Schema: public; Owner: -
--

CREATE MATERIALIZED VIEW public.issue_comments AS
 SELECT c.repository_owner,
    c.repository_name,
    c.repository_fullname,
    c.issue_number,
    c.created_at,
    c.body,
    c.user_id,
    c.user_login,
    c.htmlurl AS html_url
   FROM (public.github_issue_comments_versioned c
     JOIN public.github_issues_versioned i ON (((i.repository_owner = c.repository_owner) AND (i.repository_name = c.repository_name) AND (i.number = c.issue_number))))
  WITH NO DATA;


--
-- Name: issues; Type: MATERIALIZED VIEW; Schema: public; Owner: -
--

CREATE MATERIALIZED VIEW public.issues AS
 SELECT github_issues_versioned.repository_owner,
    github_issues_versioned.repository_name,
    github_issues_versioned.repository_fullname,
    github_issues_versioned.number,
    github_issues_versioned.state,
    github_issues_versioned.title,
    github_issues_versioned.body,
    github_issues_versioned.created_at,
    github_issues_versioned.closed_at,
    github_issues_versioned.updated_at,
    github_issues_versioned.comments,
    github_issues_versioned.user_id,
    github_issues_versioned.user_login,
    github_issues_versioned.htmlurl AS html_url,
    github_issues_versioned.labels
   FROM public.github_issues_versioned
  WITH NO DATA;


--
-- Name: owners; Type: MATERIALIZED VIEW; Schema: public; Owner: -
--

CREATE MATERIALIZED VIEW public.owners AS
 SELECT github_organizations_versioned.login,
    github_organizations_versioned.name
   FROM public.github_organizations_versioned
  WITH NO DATA;


--
-- Name: pull_request_comments; Type: MATERIALIZED VIEW; Schema: public; Owner: -
--

CREATE MATERIALIZED VIEW public.pull_request_comments AS
 SELECT c.repository_owner,
    c.repository_name,
    c.repository_fullname,
    c.issue_number AS pull_request_number,
    c.created_at,
    c.body,
    c.user_id,
    c.user_login,
    c.htmlurl AS html_url
   FROM (public.github_issue_comments_versioned c
     JOIN public.github_pull_requests_versioned p ON (((p.repository_owner = c.repository_owner) AND (p.repository_name = c.repository_name) AND (p.number = c.issue_number))))
UNION
 SELECT c.repository_owner,
    c.repository_name,
    c.repository_fullname,
    c.pull_request_number,
    c.created_at,
    c.body,
    c.user_id,
    c.user_login,
    c.htmlurl AS html_url
   FROM public.github_pull_request_comments_versioned c
UNION
 SELECT c.repository_owner,
    c.repository_name,
    c.repository_fullname,
    c.pull_request_number,
    c.submitted_at AS created_at,
    c.body,
    c.user_id,
    c.user_login,
    c.htmlurl AS html_url
   FROM public.github_pull_request_reviews_versioned c
  WHERE (c.body <> ''::text)
  WITH NO DATA;


--
-- Name: pull_request_reviews; Type: MATERIALIZED VIEW; Schema: public; Owner: -
--

CREATE MATERIALIZED VIEW public.pull_request_reviews AS
 SELECT github_pull_request_reviews_versioned.repository_owner,
    github_pull_request_reviews_versioned.repository_name,
    github_pull_request_reviews_versioned.repository_fullname,
    github_pull_request_reviews_versioned.pull_request_number,
    github_pull_request_reviews_versioned.submitted_at AS created_at,
    github_pull_request_reviews_versioned.user_id,
    github_pull_request_reviews_versioned.user_login,
    github_pull_request_reviews_versioned.htmlurl AS html_url,
        CASE
            WHEN (github_pull_request_reviews_versioned.state = 'CHANGES_REQUESTED'::text) THEN 'COMMENTED'::text
            ELSE github_pull_request_reviews_versioned.state
        END AS state
   FROM public.github_pull_request_reviews_versioned
  WITH NO DATA;


--
-- Name: pull_requests; Type: MATERIALIZED VIEW; Schema: public; Owner: -
--

CREATE MATERIALIZED VIEW public.pull_requests AS
 SELECT github_pull_requests_versioned.repository_owner,
    github_pull_requests_versioned.repository_name,
    github_pull_requests_versioned.repository_fullname,
    github_pull_requests_versioned.number,
    github_pull_requests_versioned.state,
    github_pull_requests_versioned.title,
    github_pull_requests_versioned.body,
    github_pull_requests_versioned.created_at,
    github_pull_requests_versioned.closed_at,
    github_pull_requests_versioned.merged_at,
    github_pull_requests_versioned.updated_at,
    github_pull_requests_versioned.commits,
    github_pull_requests_versioned.comments,
    github_pull_requests_versioned.changed_files,
    github_pull_requests_versioned.additions,
    github_pull_requests_versioned.deletions,
    github_pull_requests_versioned.review_comments AS reviews,
    github_pull_requests_versioned.user_id,
    github_pull_requests_versioned.user_login,
    github_pull_requests_versioned.base_repository_name,
    github_pull_requests_versioned.base_repository_owner,
    github_pull_requests_versioned.base_repository_fullname,
    github_pull_requests_versioned.head_ref,
    github_pull_requests_versioned.head_sha,
    github_pull_requests_versioned.merge_commit_sha,
    github_pull_requests_versioned.htmlurl AS html_url,
    github_pull_requests_versioned.labels
   FROM public.github_pull_requests_versioned
  WITH NO DATA;


--
-- Name: repositories; Type: MATERIALIZED VIEW; Schema: public; Owner: -
--

CREATE MATERIALIZED VIEW public.repositories AS
 SELECT github_repositories_versioned.owner_login AS owner,
    github_repositories_versioned.name,
    github_repositories_versioned.fullname,
    github_repositories_versioned.private,
    github_repositories_versioned.description
   FROM public.github_repositories_versioned
  WITH NO DATA;


--
-- Name: schema_migrations; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.schema_migrations (
    version bigint NOT NULL,
    dirty boolean NOT NULL
);


--
-- Name: users; Type: MATERIALIZED VIEW; Schema: public; Owner: -
--

CREATE MATERIALIZED VIEW public.users AS
 SELECT github_users_versioned.login,
    github_users_versioned.name
   FROM public.github_users_versioned
  WITH NO DATA;


--
-- Name: github_issue_comments_versioned issue_comments_versioned_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.github_issue_comments_versioned
    ADD CONSTRAINT issue_comments_versioned_pkey PRIMARY KEY (sum256);


--
-- Name: github_issues_versioned issues_versioned_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.github_issues_versioned
    ADD CONSTRAINT issues_versioned_pkey PRIMARY KEY (sum256);


--
-- Name: github_organizations_versioned organizations_versioned_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.github_organizations_versioned
    ADD CONSTRAINT organizations_versioned_pkey PRIMARY KEY (sum256);


--
-- Name: github_pull_request_comments_versioned pull_request_comments_versioned_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.github_pull_request_comments_versioned
    ADD CONSTRAINT pull_request_comments_versioned_pkey PRIMARY KEY (sum256);


--
-- Name: github_pull_request_reviews_versioned pull_request_reviews_versioned_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.github_pull_request_reviews_versioned
    ADD CONSTRAINT pull_request_reviews_versioned_pkey PRIMARY KEY (sum256);


--
-- Name: github_pull_requests_versioned pull_requests_versioned_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.github_pull_requests_versioned
    ADD CONSTRAINT pull_requests_versioned_pkey PRIMARY KEY (sum256);


--
-- Name: github_repositories_versioned repositories_versioned_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.github_repositories_versioned
    ADD CONSTRAINT repositories_versioned_pkey PRIMARY KEY (sum256);


--
-- Name: schema_migrations schema_migrations_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.schema_migrations
    ADD CONSTRAINT schema_migrations_pkey PRIMARY KEY (version);


--
-- Name: github_users_versioned users_versioned_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.github_users_versioned
    ADD CONSTRAINT users_versioned_pkey PRIMARY KEY (sum256);


--
-- Name: issue_comments_versions; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX issue_comments_versions ON public.github_issue_comments_versioned USING btree (versions);


--
-- Name: issues_versions; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX issues_versions ON public.github_issues_versioned USING btree (versions);


--
-- Name: organizations_versions; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX organizations_versions ON public.github_organizations_versioned USING btree (versions);


--
-- Name: pull_request_comments_versions; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX pull_request_comments_versions ON public.github_pull_request_comments_versioned USING btree (versions);


--
-- Name: pull_request_reviews_versions; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX pull_request_reviews_versions ON public.github_pull_request_reviews_versioned USING btree (versions);


--
-- Name: pull_requests_versions; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX pull_requests_versions ON public.github_pull_requests_versioned USING btree (versions);


--
-- Name: repositories_versions; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX repositories_versions ON public.github_repositories_versioned USING btree (versions);


--
-- Name: users_versions; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX users_versions ON public.github_users_versioned USING btree (versions);


--
-- PostgreSQL database dump complete
--

