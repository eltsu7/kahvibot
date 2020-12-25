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
create view juonnit_ja_nimet as
    select id, update_id, username, aika, kuvaus
    from juonnit
    inner join nimet on juonnit.user_id = nimet.user_id
    order by id desc;

create view juontilaskuri as
    select nimet.user_id, username, count(*) as kupit
    from juonnit
    inner join nimet on juonnit.user_id = nimet.user_id
    group by nimet.user_id, username
    order by count(*) desc;