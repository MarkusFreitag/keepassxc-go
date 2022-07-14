# Major Release v1.0.0 (2022-07-14)
  * **Markus Freitag**
    * initial release
      * pkg/keystore: handle keepassxc access tokens
      * pkg/client: handle communication via the unix socket
        - change-public-keys
        - get-database-hash
        - associate
        - test-associate
        - get-logins
      * cmd: CLI tool using the client
        - search for credentials

*Released by Markus Freitag <fmarkus@mailbox.org>*
