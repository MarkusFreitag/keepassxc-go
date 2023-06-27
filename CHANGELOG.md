# For next release
  * **Guilhem Bonnefille**
    * feat: add support for windows

*Not released yet*

# Patch Release v1.4.1 (2023-06-22)
  * **Markus Freitag**
    * pkg/keepassxc/client_linux: fix socket-path lookup routine

*Released by Markus Freitag <fmarkus@mailbox.org>*

# Minor Release v1.4.0 (2023-06-21)
  * **robert-renk**
    * pkg/keepassxc/client_linux: add support for different socket locations
      Distributions store the keepassxc socket in different locations, e.g. with snap in $HOME/snap/...

*Released by Markus Freitag <fmarkus@mailbox.org>*

# Patch Release v1.3.1 (2023-06-01)
  * **Markus Freitag**
    * go: upgrade to 1.20
    * tests: replace original monkey with license-compatible version
      The original version of monkey by bouk has been archived in 2020 and
      its license prohibits the usage of it in any project.

*Released by Markus Freitag <fmarkus@mailbox.org>*

# Minor Release v1.3.0 (2022-09-30)
  * **Markus Freitag**
    * pkg/keystore: add convenience method for profile loading
    * pkg/keepassxc: add convenience method for client initialization

*Released by Markus Freitag <fmarkus@mailbox.org>*

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
