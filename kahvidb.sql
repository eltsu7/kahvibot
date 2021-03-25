-- taulut
create table juonnit (
    id          serial          primary key,
    update_id   bigint          not null,
    user_id     bigint          not null,
    aika        timestamp       not null    default CURRENT_DATE,
    kuvaus      varchar(255)    not null    default ''
);

create table nimet (
    user_id     bigint          primary key,
    username    varchar(255)    not null
);
-- näkymät
create view kahvit as
    select id, update_id, username, aika, kuvaus
    from juonnit
    inner join nimet on juonnit.user_id = nimet.user_id
    order by id desc;

create view kuppilaskuri as
    select nimet.user_id, username, count(*) as kupit
    from juonnit
    inner join nimet on juonnit.user_id = nimet.user_id
    group by nimet.user_id, username
    order by count(*) desc;

create view uniikitkahvit as
    SELECT juonnit.kuvaus,
        max(juonnit.aika) AS aika,
        max(juonnit.user_id) AS user_id
    FROM juonnit
    GROUP BY juonnit.kuvaus
    ORDER BY (max(juonnit.aika)) DESC;