# httpshell
Shell via http protocol

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