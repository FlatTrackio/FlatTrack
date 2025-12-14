begin;

update users
set
  names = '(Deleted User)',
  email = '',
  birthday = 0,
  phoneNumber = '',
  password = ''
where deletionTimestamp <> 0;

commit;
