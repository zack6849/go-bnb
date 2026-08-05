# GoBNB

A Toy AirBNB Clone service written in golang

# Setup

You'll need to download a dataset from inside-airbnb [here](https://insideairbnb.com/get-the-data/) and insert it into
the seed_data folder before running this, as the vendor requests we don't distribute their data or scrape it. Once
you've put the correct files in there, you should be able to run the app it SHOULD automatically migrate the db and seed
it **(TODO)**

The end goal would be you run docker compose up -d, it spins up the db, api, and frontend all for you and you just view
this all in a browser

# TODO

- Dockerize the go side of the application
- Migrations
- Models & Repositories
- Finish seed pipeline
- Design API Spec
- Design basic react SPA or something for the frontend
- Bonus: Authentication, right now the goal is just to import and visualize the data and expose it over an API and
  example client, but a stretch goal might be to build an entire experience with host login, user login, etc