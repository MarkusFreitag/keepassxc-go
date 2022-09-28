# Minor Release v1.2.0 (2022-09-28)
  * **Markus Freitag**
    * cmd/getLogins: add flag to enable json output

*Released by Markus Freitag <fmarkus@mailbox.org>*

# Patch Release v1.1.1 (2022-09-02)
  * **Markus Freitag**
    * pkg/client: fix imports for darwin

*Released by Markus Freitag <fmarkus@mailbox.org>*

# Minor Release v1.1.0 (2022-09-02)
  * **Markus Freitag**
    * pkg/client: add darwin support, search socket in $TMPDIR

*Released by Markus Freitag <fmarkus@mailbox.org>*

# Patch Release v1.0.1 (2022-07-15)
  * **Markus Freitag**
    * pkg/keystore: fix keystore file permission

*Released by Markus Freitag <fmarkus@mailbox.org>*

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
