teleport-challenge client
=========================

## Initial setup

1. Install prerequisites. Versions other than the ones listed will likely work
   since this project doesn't use particularly cutting-edge features of any
   tools, but I'll denote the versions I'm currently using in development:

    - node 16.0.0
    - yarn 1.22.10
    - GNU Make 3.81

1. Install dependencies, build

    ```
    $ make
    ```

    `make` will pull yarn dependencies, typecheck, and build the app. Outputs
    are in the "build/" directory.

## Run

1. Start the server

    ```
    $ cd ../server
    $ yarn run start
    ```

1. Navigate to https://localhost:8443 (with the default server settings)
