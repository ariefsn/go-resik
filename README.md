# GO RESIK

- Go boilerplate with clean architecture

## Deps

- [Fiber](https://docs.gofiber.io/)
- [Zerolog](https://github.com/rs/zerolog)
- [DotEnv](https://github.com/joho/godotenv)

## Tools

- [Air](https://github.com/cosmtrek/air)
- [Mockery](https://github.com/vektra/mockery)

## Structure

- app
  - <app_name>
    - delivery
    - repository
    - service
- common
- domain
- helper
- logger

## How to

1. Copy `.env.example` to `.env` and adjust the variables
2. Start the app
   - Run

       ```shell
         make run
       ```

   - Build

       ```shell
         make build
       ```

   - Build and Run

       ```shell
         make build.run
       ```

## Reff

- [Clean Architecture by Uncle Bob](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
