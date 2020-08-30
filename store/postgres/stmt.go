package postgres

// schema name
const prefix = `goblog`

// This stmt is run on every startup to ensure that database structures exist, and if they don't, it creates them
const stmtStartup = `
create schema ` + prefix + `
    create table if not exists users
    (
        id           bigserial not null
            constraint users_pk
                primary key,
        display_name text,
        login        text unique,
        password     text
    )

	create unique index if not exists users_id_uindex
    on users (id)

    create table if not exists authors
    (
        id      bigserial not null
            constraint authors_pk
                primary key,
        user_id bigint
            constraint authors_users_id_fk
                references users
				unique,
        name    text
    )

	create unique index if not exists authors_id_uindex
    on authors (id)

    create table if not exists articles
    (
        title        text,
        article_id   bigserial not null
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
        id          bigint not null
            constraint sessions_pk
                primary key
				unique,
        user_id     bigint not null
            constraint sessions_users_id_fk
                references users
				on delete cascade,
        valid_until timestamp
    )
    
    create unique index if not exists sessions_id_uindex
        on sessions (id)
    
    create table if not exists comments
    (
        comment_id     bigint  not null
            constraint comments_pk
                primary key
				unique,
        user_id        bigint not null
            constraint comments_users_id_fk
                references users
                on update cascade,
        unsafe_content text
    )
    
    create unique index if not exists comments_comment_id_uindex
        on comments (comment_id)
    
    create table if not exists admins
    (
        user_id bigint not null
            constraint admins_pk
                primary key
				unique
            constraint admins_users_id_fk
                references users
                on delete cascade
    );
`

// articles
const stmtLoadArticlesSortedByNewest = `select title, article_id, author_id, html_preview, timestamp from 
` + prefix + `.articles order by timestamp desc offset $1 limit $2;`

const stmtNewArticle = `insert into ` + prefix + `.articles (title, author_id, html_content, 
html_preview, timestamp) values ($1,$2,$3,$4,$5);`

const stmtEditArticle = `update ` + prefix + `.articles set title = $1, author_id = $2, 
html_content = $3, html_preview = $4, timestamp = $5 where article_id = $6;`

const stmtRemoveArticle = `delete from ` + prefix + `.articles where article_id = $1;`

const stmtGetArticleByID = `select title, author_id, html_content, timestamp from ` + prefix + `.articles where article_id = $1;`

const stmtArticleNumber = `select count(article_id) from ` + prefix + `.articles;`

// users
const stmtListUsers = `select id, display_name, login from ` + prefix + `.users order by 
id offset $1 limit $2;`

const stmtAddUser = `insert into ` + prefix + `.users (display_name, login, password) values 
($1,$2,$3);`

const stmtEditUser = `update ` + prefix + `.users set display_name = $1, login = $2, password = $3 
where id = $4;`

const stmtRemoveUser = `delete from ` + prefix + `.users where id = $1;`

const stmtGetUserID = `select id from ` + prefix + `.users where login = $1;`

const stmtGetUser = `select id, display_name, login, password from ` + prefix + `.users where id = $1;`

// authors
const stmtListAuthors = `select ` + prefix + `.users.id as user_id, display_name as user_display_name, 
login,` + prefix + `.authors.id as author_id, name as author_name from ` + prefix +
	`.users inner join ` + prefix + `.authors on ` + prefix + `.authors.user_id=` +
	prefix + `.users.id order by ` + prefix + `.users.id offset $1 limit $2;`

const stmtGetAuthor = `select id, name from ` + prefix + `.authors where user_id = $1;`

const stmtAddAuthor = `insert into ` + prefix + `.authors (user_id, name) values ($1, $2);`

const stmtLinkAuthor = `update ` + prefix + `.authors set user_id = $1 where id = $2;`

const stmtRemoveAuthor = `delete from ` + prefix + `.authors where id = $1;`

// admins
const stmtPromoteToAdmin = `insert into ` + prefix + `.admins values ($1);`

const stmtDemoteFromAdmin = `delete from ` + prefix + `.admins where user_id = $1;`

const stmtIsAdmin = `select count(user_id) from ` + prefix + `.admins where user_id = $1;`

const stmtListAdmins = `select id, display_name, login from ` + prefix + `.users 
inner join ` + prefix + `.admins on user_id=id order by ` + prefix +
	`.users.id offset $1 limit $2;`
