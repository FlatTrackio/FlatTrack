#!/bin/bash

DEFAULT_PASSWORD="$(echo 'P@ssw0rd123!' | sha512sum | awk '{print $1}')"

# wait until API is ready
echo "waiting for FlatTrack API to be ready"
until curl -H "Accept: application/json" http://localhost:8080/api 2>&1 > /dev/null; do
    echo "FlatTrack API not ready yet"
    sleep 1
done
echo "FlatTrack API is now ready"

# insert collaborators into database with a default password
while IFS= read -r collaborator; do
    NAME=$(echo $collaborator | jq -r .name)
    EMAIL=$(echo $collaborator | jq -r .email)
    echo "Inserting '$NAME' into users"
    psql "postgres://flattrack:flattrack@localhost/flattrack" <<EOF
-- insert a new user account
with usercreate (insertedid) as (
  insert into users
    (names, email, password, registered)
    values
      ('$NAME', '$EMAIL', '$DEFAULT_PASSWORD', true)
    on conflict do nothing
    returning id
), groupflatmemberid as (
   select id from groups where name = 'flatmember'
), groupadminid as (
   select id from groups where name = 'admin'
)
insert into user_to_groups (userid, groupid) values (usercreate, groupflatmemberid);
EOF
done < <(echo $SHARINGIO_PAIR_GUESTS | jq -c .[])

echo "Complete"



# insert into user_to_groups (userid, groupid) values ()

# -- assign groups
# with inserteduserid (uid) as (
#     select id from users where email = '$EMAIL' limit 1
# ), existinggroupids (gid) as (
#     select id from groups where name = 'flatmember' or name = 'admin'
# )
# insert into user_to_groups (userid, groupid) from (inserteduserid, existinggroupids)
