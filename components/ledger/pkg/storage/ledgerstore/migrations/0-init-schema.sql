/**
  Some utils
 */
create aggregate aggregate_objects(jsonb) (
  sfunc = jsonb_concat,
  stype = jsonb,
  initcond = '{}'
);

create function first_agg (anyelement, anyelement)
    returns anyelement
    language sql
    immutable
    strict
    parallel safe
as $$
    select $1
$$;

create aggregate first (anyelement) (
    sfunc    = first_agg,
    stype    = anyelement,
    parallel = safe
);

create function array_distinct(anyarray)
    returns anyarray
    language sql
    immutable
as $$
    select array_agg(distinct x)
    from unnest($1) t(x);
$$ ;

/** Define tables **/
create table transactions (
    id numeric not null,
    metadata jsonb not null default '{}'::jsonb,
    date timestamp without time zone not null,
    reference varchar,
    revision numeric default 0 not null,
    last_update timestamp not null,
    reverted bool default false not null,
    postings varchar not null,

    -- all information for pre/post commit volumes
    involved_sources varchar[] not null,
    involved_destinations varchar[] not null,
    involved_assets varchar[] not null
);

create table accounts (
    address varchar,
    address_array jsonb,
    metadata jsonb default '{}'::jsonb,
    revision numeric default 0,
    last_update timestamp
);

create table moves (
    seq serial not null primary key ,
    transaction_id numeric not null,
    posting_index int8 not null,
    account_address varchar not null,
    account_address_array jsonb not null,
    asset varchar not null,
    amount numeric not null,
    date timestamp not null,
    post_commit_aggregated_input numeric,
    post_commit_aggregated_output numeric,
    is_source boolean not null,
    stale boolean not null
);

create type log_type as enum (
    'NEW_TRANSACTION',
    'REVERTED_TRANSACTION',
    'SET_METADATA'
);

create table logs (
    id numeric not null,
    type log_type not null,
    hash bytea not null,
    date timestamp not null,
    data jsonb not null,
    idempotency_key varchar(255)
);

/** Define types **/
create type account_with_volumes as (
    address varchar,
    metadata jsonb,
    volumes jsonb
);

create type volumes as (
    asset varchar,
    inputs numeric,
    outputs numeric
);

/** Define index **/

/** Index required for write part */
-- todo: if 'where not stale' is used, logs insertion is speedup (maybe by 40%), but slow down some read query
create index moves_range_dates on moves (account_address, asset, date); -- where not stale

create index moves_range_dates_not_staled on moves (account_address, asset, date) where not stale;

/** Index requires for read */
create index transactions_date on transactions (date);
create index transactions_metadata on transactions using gin (metadata);
create index transactions_involved_sources on transactions using gin (involved_sources);
create index transactions_involved_destinations on transactions using gin (involved_destinations);
create unique index transactions_revisions on transactions(id desc, revision desc);

create index moves_account_address on moves (account_address);
create index moves_account_address_array on moves using gin (account_address_array jsonb_ops);
create index moves_account_address_array_length on moves (jsonb_array_length(account_address_array));
create index moves_date on moves (date);
--create index moves_date_ranges on moves using brin(date); --todo: monitor usage
create index moves_asset on moves(asset);

create unique index accounts_revisions on accounts(address asc, revision desc);

/** Define write functions **/
create function insert_new_account(_address varchar, _date timestamp)
    returns void
    language sql
as $$
    insert into accounts (address, last_update, address_array, revision)
    values (_address, _date, to_json(string_to_array(_address, ':')), 0)
    on conflict do nothing
$$;

create function get_account(_account_address varchar, _before timestamp default null)
    returns setof accounts
    language sql
    stable
as $$
    select distinct on (address) *
    from accounts t
    where (_before is null or t.last_update <= _before)
        and t.address = _account_address
    order by address, revision desc
    limit 1;
$$;

create function get_transaction(_id numeric, _before timestamp default null)
    returns setof transactions
    language sql
    stable
as $$
    select distinct on (id) *
    from transactions t
    where (_before is null or t.last_update <= _before) and t.id = _id
    order by id desc, revision desc
    limit 1;
$$;

-- a simple 'select distinct asset from moves' would be more simple
-- but Postgres is extremely inefficient with distinct
-- so the query implementation use a "hack" to emulate skip scan feature which Postgres lack natively
-- see https://wiki.postgresql.org/wiki/Loose_indexscan for more information
create function get_all_assets()
    returns setof varchar
    language sql
as $$
    with recursive t as (
        select min(asset) as asset
        from moves
        union all
        select (
            select min(asset)
            from moves
            where asset > t.asset
        )
        from t
        where t.asset is not null
    )
    select asset from t where asset is not null
    union all
    select null where exists(select 1 from moves where asset is null)
$$;

create function get_moves(_before timestamp default null)
    returns setof moves
    language sql
    stable
as $$
    select *
    from moves s
    where _before is null or s.date <= _before
    order by date desc, seq desc
$$;

create function get_moves_for_account(_account_address varchar, _before timestamp default null)
    returns setof moves
    language sql
    stable
as $$
    select *
    from get_moves(_before) s
    where s.account_address = _account_address
$$;

create function get_moves_for_account_and_asset(_account_address varchar, _asset varchar, _before timestamp default null)
    returns setof moves
    language sql
    stable
as $$
    select *
    from get_moves_for_account(_account_address, _before) s
    where s.asset = _asset
$$;

create function get_latest_computed_move_for_account_and_asset(_account_address varchar, _asset varchar, _before timestamp default null)
    returns setof moves
    language sql
    stable
as $$
    select v.*
    from get_moves_for_account_and_asset(_account_address, _asset, _before) v
    where not v.stale
    limit 1
$$;

create function get_latest_move_for_account_and_asset(_account_address varchar, _asset varchar, _before timestamp default null)
    returns setof moves
    language sql
    stable
as $$
    select *
    from get_moves_for_account_and_asset(_account_address, _asset, _before) v
    limit 1;
$$;

create function update_account_metadata(_address varchar, _metadata jsonb, _date timestamp)
    returns void
    language sql
as $$
    insert into accounts (address, metadata, last_update, revision, address_array)
    select _address, originalAccount.metadata || _metadata, _date, originalAccount.revision + 1, to_json(string_to_array(originalAccount.address, ':'))
    from get_account(_address) originalAccount
    union all -- if account doesn't exists
    select _address, _metadata, _date, 0, to_json(string_to_array(_address, ':'))
    limit 1;
$$;

create function update_transaction_metadata(_id numeric, _metadata jsonb, _date timestamp)
    returns void
    language sql
as $$
    insert into transactions (id, metadata, date, reference, reverted, involved_sources, involved_destinations,
                              involved_assets, last_update, revision, postings)
    select originalTX.id,
           originalTX.metadata || _metadata,
           originalTX.date,
           originalTX.reference,
           originalTX.reverted,
           originalTX.involved_sources,
           originalTX.involved_destinations,
           originalTX.involved_assets,
           _date,
            originalTX.revision + 1,
            originalTX.postings
    from get_transaction(_id) originalTX
$$;

create function revert_transaction(_id numeric, _date timestamp)
    returns void
    language sql
as $$
    insert into transactions (id, metadata, date, reference, reverted, involved_sources, involved_destinations,
                              involved_assets, last_update, revision, postings)
    select originalTX.id,
        originalTX.metadata,
        originalTX.date,
        originalTX.reference,
        true,
        originalTX.involved_sources,
        originalTX.involved_destinations,
        originalTX.involved_assets,
        _date,
        originalTX.revision + 1,
        originalTX.postings
    from get_transaction(_id) originalTX
$$;

-- todo: maybe we could avoid plpgsql functions
create function insert_transaction(data jsonb)
    returns void
    language plpgsql
as $$
    declare
        posting jsonb;
        index int8 = 0;
        involved_sources varchar[];
        involved_destinations varchar[];
        involved_assets varchar[];
    begin
        index = 1;
        for posting in (select jsonb_array_elements(data->'postings')) loop
            involved_sources[index] = posting->>'source';
            involved_destinations[index] = posting->>'destination';
            involved_assets[index] = posting->>'asset';
            index = index + 1;
        end loop;

        insert into transactions (id, metadata, date, reference, involved_sources, involved_destinations, involved_assets, last_update, postings)
        values ((data->>'id')::numeric, coalesce(data->'metadata', '{}'::jsonb), (data->>'date')::timestamp without time zone, data->>'reference', involved_sources, involved_destinations, involved_assets, (data->>'date')::timestamp without time zone, jsonb_pretty(data->'postings'));

        index = 0;
        for posting in (select jsonb_array_elements(data->'postings')) loop
            -- todo: sometimes the balance is known at commit time (for sources != world), we need to forward the value to populate the pre_commit_aggregated_input and output
            insert into moves (date, account_address, asset, transaction_id, posting_index, amount, is_source, account_address_array, stale)
            values
                ((data->>'date')::timestamp without time zone, posting->>'source', posting->>'asset', (data->>'id')::numeric, index, (posting->>'amount')::numeric, true, (select to_json(string_to_array(posting->>'source', ':'))), true),
                ((data->>'date')::timestamp without time zone, posting->>'destination', posting->>'asset', (data->>'id')::numeric, index, (posting->>'amount')::numeric, false, (select to_json(string_to_array(posting->>'destination', ':'))), true);

            -- todo: we could probably avoid insertion using some kind of full join later
            perform insert_new_account(posting->>'source', (data->>'date')::timestamp without time zone);
            perform insert_new_account(posting->>'destination', (data->>'date')::timestamp without time zone);

            index = index + 1;
        end loop;

        -- invalid balances of future transaction
        -- todo: use a window?
        update moves b
        set stale = true
        where not b.stale and b.date > (data->>'date')::timestamp without time zone and (
            account_address = any(involved_sources) or account_address = any(involved_destinations)
        );
    end
$$;

-- function ensuring a specific balance, at a specific time, is properly computed
-- to compute a balance for an account and an asset, we take the last not staled value, then we add all amounts of balances
-- between the last not staled value and the balance we're actually trying to update.
-- (remember the balance is versioned and each new fund movements give a new row for the balance of an account)
create function compute_move(record_to_update moves)
    returns moves
    language sql
as $$
    with
         latest_computed_move as (
             (
                 select moves.seq, moves.date, moves.post_commit_aggregated_input, moves.post_commit_aggregated_output, moves.is_source, moves.amount
                 from moves
                 where date = (
                     select max(date)
                     from moves
                     where
                         account_address = record_to_update.account_address and
                         asset = record_to_update.asset and
                         date <= record_to_update.date and
                         not stale
                 ) and
                       account_address = record_to_update.account_address and
                       asset = record_to_update.asset and
                       not stale
                 order by seq desc
             ) union all
             (
                 select -1, '-Infinity', 0, 0, false, 0
             )
             limit 1
         ),
         new_moves_since_latest_computed_move_at_previous_date as (
             select m.*
             from moves m
             join latest_computed_move on true
             where m.account_address = record_to_update.account_address and
                 m.asset = record_to_update.asset and
                 m.date < record_to_update.date and
                 m.date > latest_computed_move.date
         ),
         new_moves_since_latest_computed_move_at_same_date as (
             select m.*
             from moves m
             join latest_computed_move on true
             where m.account_address = record_to_update.account_address and
                 m.asset = record_to_update.asset and
                 m.date = record_to_update.date and
                 m.seq > latest_computed_move.seq and
                 m.seq < record_to_update.seq
         ),
         -- We could use one query using and/or on dates, but using two queries allow the query planner to take better decision and speed up results
         new_moves_since_latest_computed_move as (
             select * from new_moves_since_latest_computed_move_at_previous_date
             union all
             select * from new_moves_since_latest_computed_move_at_same_date
         ),
         new_outputs as (
             select coalesce(sum(m.amount), 0) as amount
             from new_moves_since_latest_computed_move m
             where is_source
         ),
         new_inputs as (
             select coalesce(sum(m.amount), 0) as amount
             from new_moves_since_latest_computed_move m
             where not is_source
         )
    update moves
    set
        post_commit_aggregated_input = latest_computed_move.post_commit_aggregated_input + new_inputs.amount + case when not moves.is_source then moves.amount else 0 end,
        post_commit_aggregated_output = latest_computed_move.post_commit_aggregated_output + new_outputs.amount + case when moves.is_source then moves.amount else 0 end,
        stale = false
    from new_inputs, new_outputs, latest_computed_move
    where moves.seq = record_to_update.seq and latest_computed_move.seq <> moves.seq
    returning moves.*;
$$;

create function ensure_move_computed(m moves)
    returns moves
    language sql
as $$
select m.*
where not m.stale
union all
select *
from compute_move(m)
limit 1
$$;

-- function allowing to force update all balances
create function update_pre_post_commit_volumes(_limit numeric default 100)
    returns setof moves
    language plpgsql
as $$
    declare move moves;
    begin
        select * into move
        from moves where stale
        order by date, seq
        limit 1;

        return query with last_computed_move as (
            (
                select moves.date, moves.seq, moves.post_commit_aggregated_input, moves.post_commit_aggregated_output, moves.is_source, moves.amount
                from moves
                where date = (
                    select max(date)
                    from moves
                    where account_address = move.account_address and asset = move.asset and not stale and date < move.date
                ) and
                    account_address = move.account_address and
                    asset = move.asset and
                    not stale
                order by seq desc
            ) union all (
                select '-Infinity', -1, 0, 0, false, 0
            )
            limit 1
        ),
        computed_moves as (
            select moves.seq, moves.amount, moves.is_source,
                last_computed_move.post_commit_aggregated_output + sum(case when moves.is_source then moves.amount else 0 end) over (order by moves.date asc, moves.seq asc) as outputs,
                last_computed_move.post_commit_aggregated_input + sum(case when not moves.is_source then moves.amount else 0 end) over (order by moves.date asc, moves.seq asc) as inputs
            from moves
            join last_computed_move on true
            where moves.account_address = move.account_address and moves.asset = move.asset and (moves.date > last_computed_move.date or (moves.date = last_computed_move.date and moves.seq > last_computed_move.seq))
            order by moves.date asc, moves.seq asc
            limit _limit
        )
        update moves
        set post_commit_aggregated_input = computed_moves.inputs,
            post_commit_aggregated_output = computed_moves.outputs,
            stale = false
        from computed_moves
        where moves.seq = computed_moves.seq
        returning moves.*;
    end;
$$;

create function refresh_volumes(_account varchar, _asset varchar, _before timestamp default null)
    returns volumes
    language sql
as $$
    select computed_move.asset,
           computed_move.post_commit_aggregated_input as inputs,
           computed_move.post_commit_aggregated_output as outputs
    from get_latest_move_for_account_and_asset(_account, _asset, _before) m, compute_move(m) computed_move
$$;

create function handle_log() returns trigger
  security definer
  language plpgsql
as $$
  declare
    _key varchar;
    _value jsonb;
  begin
    if new.type = 'NEW_TRANSACTION' then
      perform insert_transaction(new.data->'transaction');
      for _key, _value in (select * from jsonb_each_text(new.data->'accountMetadata')) loop
          perform update_account_metadata(_key, _value, (new.data->'transaction'->>'date')::timestamp);
      end loop;
    end if;
    if new.type = 'REVERTED_TRANSACTION' then
        perform insert_transaction(new.data->'transaction');
        perform revert_transaction((new.data->>'revertedTransactionID')::numeric, (new.data->'transaction'->>'date')::timestamp);
    end if;
    if new.type = 'SET_METADATA' then
        if new.data->>'targetType' = 'TRANSACTION' then
            perform update_transaction_metadata((new.data->>'targetId')::numeric, new.data->'metadata', new.date);
        else
            perform update_account_metadata((new.data->>'targetId')::varchar, new.data ->'metadata', new.date);
        end if;
    end if;

    return new;
  end;
$$;

/** Define the trigger which populate table in response to new logs **/
create trigger account_insert after insert on logs
    for each row execute procedure handle_log();

create function get_account_volumes_for_asset(_account varchar, _asset varchar, _before timestamp default null)
    returns volumes
    language sql
    stable
as $$
    (
        select v.asset, v.post_commit_aggregated_input, v.post_commit_aggregated_output
        from get_latest_move_for_account_and_asset(_account, _asset, _before) v
        where not v.stale
    ) union all (
        select * from refresh_volumes(_account, _asset, _before)
    )
    limit 1
$$;

create function get_all_account_volumes(_account varchar, _before timestamp default null)
    returns setof volumes
    language sql
    stable
as $$
    with
        all_assets as (
            select v.v as asset
            from get_all_assets() v
        ),
        moves as (
            select m.*
            from all_assets, get_latest_move_for_account_and_asset(_account, all_assets.asset, _before := _before) m
        ),
        fresh_moves as (
            select moves.asset, moves.post_commit_aggregated_input, moves.post_commit_aggregated_output
            from moves
            where not moves.stale
        ),
        refreshed_moves as (
            select refreshed_move.asset, refreshed_move.post_commit_aggregated_input, refreshed_move.post_commit_aggregated_output
            from moves, compute_move(moves) as refreshed_move
            where moves.stale
        )
    select *
    from fresh_moves
    union
    select *
    from refreshed_moves
$$;

create function volumes_to_jsonb(v volumes)
    returns jsonb
    language sql
    immutable
as $$
    select ('{"' || v.asset || '": {"input": ' || v.inputs || ', "output": ' || v.outputs || '}}')::jsonb
$$;

create function get_account_aggregated_volumes(_account_address varchar, _before timestamp default null)
    returns jsonb
    language sql
    stable
as $$
    select aggregate_objects(volumes_to_jsonb(volumes))
    from get_all_account_volumes(_account_address, _before := _before) volumes
$$;

create function get_account_balance(_account varchar, _asset varchar, _before timestamp default null)
    returns numeric
    language sql
    stable
as $$
    select volumes.inputs - volumes.outputs
    from get_account_volumes_for_asset(_account, _asset, _before := _before) volumes
$$;

create function aggregate_ledger_volumes(
    _before timestamp default null,
    _accounts varchar[] default null,
    _assets varchar[] default null
)
    returns setof volumes
    language sql
    stable
as $$
    with
        moves as (
            select distinct on (m.account_address, m.asset) m.*
            from get_moves(_before := _before) m
            where (_accounts is null or account_address = any(_accounts)) and
                (_assets is null or asset = any(_assets))
            order by account_address, asset, m.seq desc
        ),
        fresh_moves as (
            select moves.asset, moves.post_commit_aggregated_input as inputs, moves.post_commit_aggregated_output as outputs
            from moves
            where not moves.stale
        ),
        refreshed_moves as (
            select refreshed_move.asset, refreshed_move.post_commit_aggregated_input as inputs, refreshed_move.post_commit_aggregated_output as outputs
            from moves, compute_move(moves) as refreshed_move
            where moves.stale
        )
    select v.asset, sum(v.inputs) as inputs, sum(v.outputs) as outputs
    from (
        select *
        from fresh_moves
        union all
        select *
        from refreshed_moves
    ) v
    group by v.asset
$$;

create function get_aggregated_volumes_for_transaction(tx transactions) returns jsonb
    stable
    language sql
as
$$
select aggregate_objects(jsonb_build_object(data.account_address, data.aggregated))
from (
    select distinct on (safe_move.account_address, safe_move.asset) safe_move.account_address,
        volumes_to_jsonb((safe_move.asset, first(safe_move.post_commit_aggregated_input), first(safe_move.post_commit_aggregated_output))) as aggregated
    from moves move
    join ensure_move_computed(move) safe_move on true
    where move.transaction_id = tx.id
    group by safe_move.account_address, safe_move.asset
) data
$$;
