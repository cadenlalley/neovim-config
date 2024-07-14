# kitchens-api
Kitchens API

## Development

To work on the Kitchens API you'll need to follow the below steps.

1. Duplicate the `.env.example` file to `.env` and fill out the required environment variables.
2. Spin up the Kitchens API dependencies via `docker compose up`.
3. Execute `make fixtures` which will spin up a few dummy accounts and the test-service@kitchens-app.com which can be logged into.
3. Execute `make serve` and the API should be available on `localhost:1313`.

### Dependencies

* [Golang](https://formulae.brew.sh/formula/go) - Allows you to run Golang.
* [Docker](https://www.docker.com/products/docker-desktop/) - Allows you to run Docker.
* [Auth0 Login](https://github.com/auth0-samples/auth0-vue-samples/tree/master/01-Login) - Simple way to get the `id_token` needed for `v1` requests.

## Documentation

The details about each route and requests that can be made are in [Postman](https://kitchens-app.postman.co/workspace/89fa7a36-50b1-40e3-ae6b-e3ba8c5f9b9e/documentation/36191591-8aa45609-bc7d-43e8-852c-820feb94999e).