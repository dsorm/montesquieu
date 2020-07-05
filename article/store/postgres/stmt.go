package postgres

// schema name
const prefix = `goblog`

// This stmt is run on every startup to ensure that database structures exist, and if they don't, it creates them
const stmtStartup = `
create schema ` + prefix + `
    create table if not exists users
    (
        id           integer not null
            constraint users_pk
                primary key,
        display_name text,
        login        text,
        password     text
    )

    create table if not exists authors
    (
        id      integer not null
            constraint authors_pk
                primary key,
        user_id integer not null
            constraint authors_users_id_fk
                references users,
        name    text
    )
    
    create table if not exists articles
    (
        title        text,
        article_id   integer not null
            constraint articles_pk
                primary key,
        author_id    integer not null
            constraint articles_authors_id_fk
                references authors,
        html_content text,
		html_preview text,
		timestamp bigint
    )
    
    create unique index if not exists articles_article_id_uindex
        on articles (article_id)
    
    create unique index if not exists authors_id_uindex
        on authors (id)
    
    create unique index if not exists users_id_uindex
        on users (id)
    
    create table if not exists sessions
    (
        id          integer not null
            constraint sessions_pk
                primary key,
        user_id     integer not null
            constraint sessions_users_id_fk
                references users,
        valid_until timestamp
    )
    
    create unique index if not exists sessions_id_uindex
        on sessions (id)
    
    create table if not exists comments
    (
        comment_id     serial  not null
            constraint comments_pk
                primary key,
        user_id        integer not null
            constraint comments_users_id_fk
                references users
                on update cascade,
        unsafe_content text
    )
    
    create unique index if not exists comments_comment_id_uindex
        on comments (comment_id)
    
    create table if not exists admins
    (
        user_id integer not null
            constraint admins_pk
                primary key
            constraint admins_users_id_fk
                references users
                on delete cascade
    );
`

const stmtIndex = `
select title, article_id, author_id, html_preview, timestamp from ` + prefix + `.articles
order by timestamp desc
offset $1
limit $2;
`
