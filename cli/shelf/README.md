# Coral Project - Shelf CLI
 The Shelf CLI provides a command line tooling for interfacing with the shelf database
 api.


## Install

    go get github.com/CoralProject/shelf/cli/...

## Usage
  In your terminal, once you've installed the Shelf CLI tools using `go get`, Simply


  ```bash
    > shelf
      # Shelf provides the central cli housing of various cli tools that interface with the API
      #
      # Usage:
      #   shelf [command]
      #
      # Available Commands:
      #   user        user provides a shelf CLI for managing user records.
      #   query       query provides a shelf CLI for managing and executing queries.
      #
      # Flags:
      #   -h, --help[=false]: help for shelf
      #
      # Use "shelf [command] --help" for more information about a command.  
  ```


## Tools

  - Query
    The `query` command provides CRUD based subcommands for interfacing with the
    query database API. These commands are:

    - `create`
        Create allows the creation of a new query record into to the system, by
        loading the query records content from a json file.

        Example:

          ```bash

          	> query create -f ./queries/user_advice.json

          ```

    - `get`
        Create allows the retrieval of a query record from the system using the
        giving name of the record. It returns a json response of the contents of
        the record.

        Example:

          ```bash

          	> query get -n user_advice

          ```

    - `update`
        Update provides the means of updating a giving query record in the query
        database using a file.

        Example:

          ```bash

          	> query update -n user_advice -f ./queries/user_advice_update.json

          ```

    - `delete`
        Delete provides the means of remove a giving query record in the query
        database using its giving name.

        Example:

          ```bash

          	> query delete -n user_advice

          ```

    - `execute`
        Execute provides the means of remove a executing a query record in the query
        database using its name and a optional parameter of map contain key value pairs.

        Example:

          ```bash

            > query execute -n user_advice

            >	query execute -n user_advice -p {"name":"john"}

          ```


  - Users
    The `user` command provides CRUD based subcommands for interfacing with the
    user database and authentication API. These commands are:

    - `create`
        Create allows the creation of a new user record into to the system, with the
        provided record data.

        Example:

          ```bash

          	> shelf user create -n "Alex Boulder" -e alex.boulder@gmail.com -p yefc*7fdf92

          ```

    - `get`
        Create allows the retrieval of a user record from the system using the
        giving name of the record. It returns a json response of the contents of
        the record.

        Example:

          ```bash

          	# 1. To get a user using it's name:

          	> shelf user get -n "Alex Boulder"

          	# 2. To get a user using it's email address:

          	> shelf	user get -e alex.boulder@gmail.com

          	# 3. To get a user using it's public id number:

          	> shelf user get -p 199550d7-484d-4440-801f-390d44911ade

          ```

    - `update`
        Update provides the means of updating a giving user record in the user
        database. It provides update to basic information(FullName, Email) and
        allows password change updates.


        Example:

          ```bash

          	# 1. To update the 'name' of a giving record. Simple set the update type to "name",
          	# and supply the email and new name.

          	> shelf	user update -t name -e shou.lou@gmail.com -n "Shou Lou FengZhu"

          	# 2. To update the 'email' of a giving record. Simple set the update type to "email",
          	# and supply the current email and new email.

          	> shelf user update -t email -e shou.lou@gmail.com -n shou.lou.fengzhu@gmail.com

          	# 3. To update the 'password' of a giving record. Simple set the update type to "auth",
          	# and supply the current email of the record, the current password of the record and
          	# the new password

          	> shelf user update -t auth -e shou.lou@gmail.com -o oldPassword -n newPassword

          ```

    - `delete`
        Delete provides the means of remove a giving user record in the user
        database using specific identification information such as the users Email,
        PublicID or FullName.

        Example:

          ```bash

          	# 1. To delete a user using it's name:

          	>	shelf user delete -n "Alex Boulder"

          	# 2. To delete a user using it's email address:

          	>	shelf user delete -e alex.boulder@gmail.com

          	# 3. To delete a user using it's public id number:

          	>	shelf user delete -p 199550d7-484d-4440-801f-390d44911ade

          ```


    - `auth`
        Auth provides the means of authenticating user credentials from the CLI.
        It provides authentication either using the `Username and Password` combination or
        using the records access `Token and PublicID` information.

        Example:

          ```bash

            # 1. To authenticate using the user's Public Id and Token,set the type to 'token':

          	> shelf user auth -t token -k {User PublicID} -p {User Token}

            # 2. To authenticate using the user's Email and Password, set the type to 'pass':

          	> shelf	user auth -t password -k shou.lou@gmail.com -p Shen5A43*2f3e

          ```
