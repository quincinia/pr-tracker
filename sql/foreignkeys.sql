-- https://popsql.com/learn-sql/postgresql/how-to-remove-a-not-null-constraint-in-postgresql
alter table tournaments alter column type drop not null;
alter table tournaments alter column tier drop not null;

alter table attendees alter column player drop not null;