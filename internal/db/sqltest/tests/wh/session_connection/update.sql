-- Copyright (c) HashiCorp, Inc.
-- SPDX-License-Identifier: MPL-2.0

begin;
  select plan(8);

  select is(count(*),               1::bigint)                from wh_session_connection_accumulating_fact where connection_id = 'sc1_____clare';
  select is(bytes_up,               null)                     from wh_session_connection_accumulating_fact where connection_id = 'sc1_____clare';
  select is(bytes_down,             null)                     from wh_session_connection_accumulating_fact where connection_id = 'sc1_____clare';
  select is(connection_closed_time, 'infinity'::wh_timestamp) from wh_session_connection_accumulating_fact where connection_id = 'sc1_____clare';

  update session_connection set
    bytes_up = 10,
    bytes_down = 5,
    closed_reason = 'closed by end-user'
  where public_id = 'sc1_____clare';

  select is(count(*),               1::bigint)                from wh_session_connection_accumulating_fact;
  select is(bytes_up,               10::wh_bytes_transmitted) from wh_session_connection_accumulating_fact where connection_id = 'sc1_____clare';
  select is(bytes_down,             5::wh_bytes_transmitted)  from wh_session_connection_accumulating_fact where connection_id = 'sc1_____clare';
  select is(connection_closed_time, now()::wh_timestamp)      from wh_session_connection_accumulating_fact where connection_id = 'sc1_____clare';

  select * from finish();
rollback;
