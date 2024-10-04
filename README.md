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

## Deployment

To build the application and deploy it to ECR (pre-Github Actions).

1. Store the AWS account ID in the `AWS_ACCOUNT_ID` in your terminal environment, or prepend all of the following commands with `AWS_ACCOUNT_ID=xxx`.
2. Log into ECR via Docker by running `aws-vault exec <profile> -- make docker-login`.
3. Build the image with `make docker-build`.
4. Tag the image with `make docker-tag`.
5. Push the image with `make docker-push`.

### Dependencies
* [Docker](https://www.docker.com/products/docker-desktop/) - Allows you to run Docker.
* [AWS Vault](https://github.com/99designs/aws-vault/tree/master) - AWS Authentication Helper, will also need corresponding profile with ECR access.

## Documentation

The details about each route and requests that can be made are in [Postman](https://kitchens-app.postman.co/workspace/89fa7a36-50b1-40e3-ae6b-e3ba8c5f9b9e/documentation/36191591-8aa45609-bc7d-43e8-852c-820feb94999e).

## Technical Debt

The following are areas that need to be revisited after initial release into beta.

- [ ] Add validation to mysql DB inputs, in particular empty string values being accepted when they should not be.
- [x] Middleware to pull the current user up into context to reduce repeated code across endpoints.
- [ ] Middleware for URL parameter validation to reduce code across endpoints.
- [ ] Automate the deployment of ECR images via GHA.
- [ ] True up process for image uploads during recipe creation.
- [ ] Handle recipe, kitchen, and account clean ups after soft deletes.
- [ ] Investigate Yoast Schema Graph on pages to see if it improves results or lowers token count.