createdb pr_tracker
psql -U postgres -d pr_tracker
grant all privileges on all tables in schema public to jacob;
grant all privileges on all sequences in schema public to jacob;
grant all privileges on all functions in schema public to jacob;