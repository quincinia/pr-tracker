drop table TourneyTypes;
create table TourneyTypes (
    typeid  serial primary key,
    name    text not null
);

drop table Tiers;
create table Tiers (
    tierid      serial primary key,
    name        text not null,
    multiplier  int not null
);

drop table Tournaments;
create table Tournaments (
    tourneyid       serial primary key,
    type            int,
    name            text,
    url             text,
    numentrants     int,
    uniqueplacings  int,
    bracketreset    boolean,
    tier            int,
    foreign key (type) references TourneyTypes(typeid),
    foreign key (tier) references Tiers(tierid)
);

drop table Players;
create table Players (
    playerid    serial primary key,
    name        text
);

drop table Attendees;
create table Attendees (
    attendeeid  serial primary key,
    tourney     int not null,
    player      int,
    name        text,
    standing    int,
    foreign key (tourney) references Tournaments(tourneyid) on delete cascade,
    foreign key (player) references Players(playerid) on delete set null
    unique (tourney, player) -- A player represents at most 1 attendee
);