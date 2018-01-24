# HttpShell [![Build Status](https://travis-ci.org/kaneg/httpshell.svg?branch=master)](https://travis-ci.org/kaneg/httpshell)
Shell via http protocol
#### Features
* Written by Go
* Https based
* Can authenticated by certificate
* Can use http proxy
* Customize command to run

### Usage

#### Server

`httpshelld [-k] -l <listen host:listen port> <command>`
* **-l**:  Listening host and listening port, e.g. 192.168.1.1:2200. Listen host can be empty, meaing listeing on all addresses.
* **-k**: Whether to validate client certificate. If it is enabled, the client must provide certificate. The authorized certificates are stored in $HOME/.httpshell/authorized.pem. The file's format is PEM certificate.
* **command**: command to be executed, usually *bash* or *login* are used.

Typical use case:

`httpshelld -k -l :2200 login`

#### Client
`httpshell https://<server host:port>`

If $HOME/.httpshell/crt.pem and key.pem are available, they will be automatically loaded as client certificate.

Typical use case:

`httpshell https://192.168.1.1:2200`

### Examples
1. Auth by built-in username/password of Linux:

    ```shell
    httpshelld login
    ```
1. Provide shell for each Docker container by dynamically parameter:
    * In server side:

    ```shell
    httpshelld docker exec -it {{.docker}} bash
    
    ```
    * In client side, access shell of Docker container named t1
    
    ```bash
    httpshell https://localhost:5000?docker=t1
    ```
