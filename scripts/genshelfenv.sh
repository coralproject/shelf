#!/bin/bash
#
# When the script is ran, the variables from the current environment will be
# used to populate the following variables:
#
# - DISABLE_GENERATE_PLATFORM_KEYS (defaults to `FALSE`)
# - DISABLE_SHELL_EXPORT (defaults to `FALSE`)
# - LOGGING_LEVEL (defaults to `1`)
# - RECAPTCHA_SECRET
# - AUTH_PUBLIC_KEY
# - MONGO_URI (defaults to `mongodb://localhost/test`)
# - CORAL_HOST (defaults to `:16180`)
# - SPONGE_HOST (defaults to `:16181`)
# - XENIA_HOST (defaults to `:16182`)
# - ASK_HOST (defaults to `:16183`)
# - XENIA_WEB_HOST (defaults to `http://localhost:16182`)
# - CORAL_SPONGED_URL (defaults to `http://localhost:16181`)
# - CORAL_XENIAD_URL (defaults to $XENIA_WEB_HOST)
#
# If you DISABLE_SHELL_EXPORT, then the outputted $CONFIG_FILE will not contain
# notations used for sourcing the file (i.e. `export ...`). If this option is
# not enabled, then you can simply run: `source $CONFIG_FILE` and the variables
# will be loaded into the environment. If not you can use a tool like dotenv
# (as seen https://github.com/bkeepers/dotenv) which will load the file into
# a running process directly.
#
# If you DISABLE_GENERATE_PLATFORM_KEYS, then the shelf platform will not
# generate new keys to use for internal platform authentication and/or it will
# prefer the keys available in the following environment variables:
#
# - PLATFORM_PRIVATE_KEY
# - PLATFORM_PUBLIC_KEY
#

##############
## settings ##
##############

CONFIG_FILE=.shelfenv
LOGGING_LEVEL=${LOGGING_LEVEL:-1}
MONGO_URI=${MONGO_URI:-mongodb://localhost/test}

SHEBANG=
EXPORT=
if [ "$DISABLE_SHELL_EXPORT" != "TRUE" ]
then

  EXPORT="export "
  SHEBANG="#!/bin/bash

"

fi

CORAL_HOST=${CORAL_HOST:-:16180}
SPONGE_HOST=${SPONGE_HOST:-:16181}
XENIA_HOST=${XENIA_HOST:-:16182}
ASK_HOST=${ASK_HOST:-:16183}
XENIA_WEB_HOST=${XENIA_WEB_HOST:-http://localhost:16182}
CORAL_SPONGED_URL=${CORAL_SPONGED_URL:-http://localhost:16181}
CORAL_XENIAD_URL=${CORAL_XENIAD_URL:-$XENIA_WEB_HOST}

################################
## generate the shelf secrets ##
################################

if [ "$DISABLE_GENERATE_PLATFORM_KEYS" != "TRUE" ]
then

  # generate private key
  PLATFORM_PRIVATE_KEY=$(openssl ecparam -genkey -name secp384r1 -noout | openssl base64 -e | tr -d '\n')

  # generate public key
  PLATFORM_PUBLIC_KEY=$(echo $PLATFORM_PRIVATE_KEY | openssl base64 -d -A | openssl ec -pubout | openssl base64 -e | tr -d '\n')

fi

#############################
## create the $CONFIG_FILE ##
#############################

cat > $CONFIG_FILE <<EOF
${SHEBANG}## CONFIG

# CORAL
${EXPORT}CORAL_LOGGING_LEVEL=$LOGGING_LEVEL
${EXPORT}CORAL_HOST=$CORAL_HOST
${EXPORT}CORAL_SPONGED_URL=$CORAL_SPONGED_URL
${EXPORT}CORAL_XENIAD_URL=$CORAL_XENIAD_URL

# XENIA
${EXPORT}XENIA_LOGGING_LEVEL=$LOGGING_LEVEL
${EXPORT}XENIA_HOST=$XENIA_HOST
${EXPORT}XENIA_WEB_HOST=$XENIA_WEB_HOST

# ASK
${EXPORT}ASK_RECAPTCHA_SECRET=$RECAPTCHA_SECRET
${EXPORT}ASK_LOGGING_LEVEL=$LOGGING_LEVEL
${EXPORT}ASK_HOST=$ASK_HOST

# SPONGE
${EXPORT}SPONGE_LOGGING_LEVEL=$LOGGING_LEVEL
${EXPORT}SPONGE_HOST=$SPONGE_HOST

## SECRETS

# CORAL
${EXPORT}CORAL_PLATFORM_PRIVATE_KEY=$PLATFORM_PRIVATE_KEY
${EXPORT}CORAL_AUTH_PUBLIC_KEY=$AUTH_PUBLIC_KEY

# XENIA
${EXPORT}XENIA_AUTH_PUBLIC_KEY=$PLATFORM_PUBLIC_KEY
${EXPORT}XENIA_MONGO_URI=$MONGO_URI

# ASK
${EXPORT}ASK_AUTH_PUBLIC_KEY=$PLATFORM_PUBLIC_KEY
${EXPORT}ASK_MONGO_URI=$MONGO_URI

# SPONGE
${EXPORT}SPONGE_AUTH_PUBLIC_KEY=$PLATFORM_PUBLIC_KEY
${EXPORT}SPONGE_MONGO_URI=$MONGO_URI

EOF

if [ "$DISABLE_SHELL_EXPORT" != "TRUE" ]
then
  chmod +x $CONFIG_FILE
fi
