begin;

create table if not exists leader_election (
    name text not null,
    holderIdentity text not null,
    leaseDurationSeconds int not null,
    acquireTime bigint not null,
    renewTime bigint not null,
    leaderTransitions int not null,

    primary key(holderIdentity)
);

commit;
