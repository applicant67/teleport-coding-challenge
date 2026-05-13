teleport-challenge server
=========================

## Initial setup

1. Install prerequisites. Versions other than the ones listed will likely work
   since this project doesn't use particularly cutting-edge features of any
   tools, but I'll denote the versions I'm currently using in development:

    - node 16.0.0
    - yarn 1.22.10
    - GNU Make 3.81
    - openssl 1.1.1l (24 Aug 2021)

1. Install dependencies

    ```
    $ make
    ```

    `make` will pull yarn dependencies and generate keys/certificates.

1. Trust the generated CA certificate

    - For macOS Chrome and Safari this can be done with the following command:

        ```
        $ security add-trusted-cert -d -r trustRoot -k ~/Library/Keychains/login.keychain-db .keys/ca-cert.pem
        ```

    - For Firefox, navigate to `about:preferences#privacy`, scroll down to the
      "Certificates" section, and click "View Certificates". In the
      "Authorities" tab, click "Import", then choose this project's
      ".keys/ca-cert.pem".

## Run

```
$ yarn run start
```
