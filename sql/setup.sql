drop table TourneyTypes;
create table TourneyTypes (
    typeid  serial primary key,
    name    text
);

drop table Tiers;
create table Tiers (
    tierid      serial primary key,
    name        text,
    multiplier  int
);

drop table Tournaments;
create table Tournaments (
    tourneyid       serial primary key,
    type            int not null,
    name            text,
    url             text,
    numentrants     int,
    uniqueplacings  int,
    bracketreset    boolean,
    tier            int not null,
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
    player      int not null,
    name        text,
    standing    int,
    foreign key (tourney) references Tournaments(tourneyid),
    foreign key (player) references Players(playerid)
);