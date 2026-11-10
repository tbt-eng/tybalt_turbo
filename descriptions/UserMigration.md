# Migrating Users from MySQL

## The Problem

Users authenticate via OAuth2.0 with Microsoft. We would ideally like to migrate the users and their corresponding profiles and admin_profiles to a new database.

## The Plan

This plan is built. The import tool does steps 2 to 4 in its users phase (`import_data/tool.go`, the `--users` flag). The phase first deletes every row of `users`, `profiles`, `admin_profiles`, `_externalAuths`, `user_claims`, `po_approver_props` and `mileage_reset_dates`, then loads them again.

1. We must first setup Microsoft Authentication via OAuth 2.0 in the target database.

2. The MySQL Profiles table, which is dumped as `Profiles.parquet`, contains an `azureId` column. The tool populates the `_externalAuths` table as follows:

    - collectionRef: '_pb_users_auth_'
    - created: not set by the tool
    - id: not set by the tool (the database generates it)
    - provider: 'microsoft'
    - providerId: Profiles.parquet.azureId
    - recordRef: Profiles.parquet.pocketbase_uid
    - updated: not set by the tool

3. The tool populates the `users` table as follows:

    - created: not set by the tool
    - email: Profiles.parquet.email
    - emailVisibility: 0
    - id: Profiles.parquet.pocketbase_uid
    - name: join Profiles.parquet.givenName and Profiles.parquet.surname with a space
    - password: empty
    - tokenKey: a new random UUID, as hex
    - updated: not set by the tool
    - username: Profiles.parquet.email part before @ symbol
    - verified: 1

4. The tool populates the `profiles` table as follows:

    - created: not set by the tool
    - given_name: Profiles.parquet.givenName
    - id: not set by the tool (the database generates it)
    - surname: Profiles.parquet.surname
    - updated: not set by the tool
    - manager: Profiles.parquet.pocketbase_manager
    - alternate_manager: Profiles.parquet.pocketbase_alternateManager
    - default_division: Profiles.parquet.pocketbase_defaultDivision
    - uid: Profiles.parquet.pocketbase_uid
    - notification_type: 'email_text'
    - do_not_accept_submissions: Profiles.parquet.doNotAcceptSubmissions
    - _imported: true
