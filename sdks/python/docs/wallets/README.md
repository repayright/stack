# wallets

### Available Operations

* [confirm_hold](#confirm_hold) - Confirm a hold
* [create_balance](#create_balance) - Create a balance
* [create_wallet](#create_wallet) - Create a new wallet
* [credit_wallet](#credit_wallet) - Credit a wallet
* [debit_wallet](#debit_wallet) - Debit a wallet
* [get_balance](#get_balance) - Get detailed balance
* [get_hold](#get_hold) - Get a hold
* [get_holds](#get_holds) - Get all holds for a wallet
* [get_transactions](#get_transactions)
* [get_wallet](#get_wallet) - Get a wallet
* [get_wallet_summary](#get_wallet_summary) - Get wallet summary
* [list_balances](#list_balances) - List balances of a wallet
* [list_wallets](#list_wallets) - List all wallets
* [update_wallet](#update_wallet) - Update a wallet
* [void_hold](#void_hold) - Cancel a hold
* [walletsget_server_info](#walletsget_server_info) - Get server info

## confirm_hold

Confirm a hold

### Example Usage

```python
import sdk
from sdk.models import operations, shared

s = sdk.SDK(
    security=shared.Security(
        authorization="Bearer YOUR_ACCESS_TOKEN_HERE",
    ),
)

req = operations.ConfirmHoldRequest(
    confirm_hold_request=shared.ConfirmHoldRequest(
        amount=100,
        final=True,
    ),
    hold_id='nisi',
)

res = s.wallets.confirm_hold(req)

if res.status_code == 200:
    # handle response
```

## create_balance

Create a balance

### Example Usage

```python
import sdk
import dateutil.parser
from sdk.models import operations, shared

s = sdk.SDK(
    security=shared.Security(
        authorization="Bearer YOUR_ACCESS_TOKEN_HERE",
    ),
)

req = operations.CreateBalanceRequest(
    create_balance_request=shared.CreateBalanceRequest(
        expires_at=dateutil.parser.isoparse('2022-05-20T03:14:12.839Z'),
        name='Fernando Aufderhar',
        priority=716075,
    ),
    id='a4469b6e-2141-4959-890a-fa563e2516fe',
)

res = s.wallets.create_balance(req)

if res.create_balance_response is not None:
    # handle response
```

## create_wallet

Create a new wallet

### Example Usage

```python
import sdk
from sdk.models import shared

s = sdk.SDK(
    security=shared.Security(
        authorization="Bearer YOUR_ACCESS_TOKEN_HERE",
    ),
)

req = shared.CreateWalletRequest(
    metadata={
        "maxime": 'deleniti',
        "facilis": 'in',
    },
    name='Diane VonRueden',
)

res = s.wallets.create_wallet(req)

if res.create_wallet_response is not None:
    # handle response
```

## credit_wallet

Credit a wallet

### Example Usage

```python
import sdk
from sdk.models import operations, shared

s = sdk.SDK(
    security=shared.Security(
        authorization="Bearer YOUR_ACCESS_TOKEN_HERE",
    ),
)

req = operations.CreditWalletRequest(
    credit_wallet_request=shared.CreditWalletRequest(
        amount=shared.Monetary(
            amount=469249,
            asset='repellat',
        ),
        balance='quibusdam',
        metadata={
            "saepe": 'pariatur',
        },
        reference='accusantium',
        sources=[
            shared.WalletSubject(
                balance='natus',
                identifier='magni',
                type='sunt',
            ),
        ],
    ),
    id='cddc6926-01fb-4576-b0d5-f0d30c5fbb25',
)

res = s.wallets.credit_wallet(req)

if res.status_code == 200:
    # handle response
```

## debit_wallet

Debit a wallet

### Example Usage

```python
import sdk
from sdk.models import operations, shared

s = sdk.SDK(
    security=shared.Security(
        authorization="Bearer YOUR_ACCESS_TOKEN_HERE",
    ),
)

req = operations.DebitWalletRequest(
    debit_wallet_request=shared.DebitWalletRequest(
        amount=shared.Monetary(
            amount=521037,
            asset='dignissimos',
        ),
        balances=[
            'quis',
        ],
        description='nesciunt',
        destination=shared.LedgerAccountSubject(
            identifier='perferendis',
            type='dolores',
        ),
        metadata={
            "quam": 'dolor',
            "vero": 'nostrum',
            "hic": 'recusandae',
            "omnis": 'facilis',
        },
        pending=False,
    ),
    id='90c28909-b3fe-449a-8d9c-bf48633323f9',
)

res = s.wallets.debit_wallet(req)

if res.debit_wallet_response is not None:
    # handle response
```

## get_balance

Get detailed balance

### Example Usage

```python
import sdk
from sdk.models import operations

s = sdk.SDK(
    security=shared.Security(
        authorization="Bearer YOUR_ACCESS_TOKEN_HERE",
    ),
)

req = operations.GetBalanceRequest(
    balance_name='cum',
    id='77f3a410-0674-4ebf-a928-0d1ba77a89eb',
)

res = s.wallets.get_balance(req)

if res.get_balance_response is not None:
    # handle response
```

## get_hold

Get a hold

### Example Usage

```python
import sdk
from sdk.models import operations

s = sdk.SDK(
    security=shared.Security(
        authorization="Bearer YOUR_ACCESS_TOKEN_HERE",
    ),
)

req = operations.GetHoldRequest(
    hold_id='asperiores',
)

res = s.wallets.get_hold(req)

if res.get_hold_response is not None:
    # handle response
```

## get_holds

Get all holds for a wallet

### Example Usage

```python
import sdk
from sdk.models import operations

s = sdk.SDK(
    security=shared.Security(
        authorization="Bearer YOUR_ACCESS_TOKEN_HERE",
    ),
)

req = operations.GetHoldsRequest(
    cursor='aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ==',
    metadata={
        "ipsum": 'voluptate',
        "id": 'saepe',
    },
    page_size=263322,
    wallet_id='aspernatur',
)

res = s.wallets.get_holds(req)

if res.get_holds_response is not None:
    # handle response
```

## get_transactions

### Example Usage

```python
import sdk
from sdk.models import operations

s = sdk.SDK(
    security=shared.Security(
        authorization="Bearer YOUR_ACCESS_TOKEN_HERE",
    ),
)

req = operations.GetTransactionsRequest(
    cursor='aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ==',
    page_size=20651,
    wallet_id='amet',
)

res = s.wallets.get_transactions(req)

if res.get_transactions_response is not None:
    # handle response
```

## get_wallet

Get a wallet

### Example Usage

```python
import sdk
from sdk.models import operations

s = sdk.SDK(
    security=shared.Security(
        authorization="Bearer YOUR_ACCESS_TOKEN_HERE",
    ),
)

req = operations.GetWalletRequest(
    id='ce5e6a95-d8a0-4d44-ace2-af7a73cf3be4',
)

res = s.wallets.get_wallet(req)

if res.get_wallet_response is not None:
    # handle response
```

## get_wallet_summary

Get wallet summary

### Example Usage

```python
import sdk
from sdk.models import operations

s = sdk.SDK(
    security=shared.Security(
        authorization="Bearer YOUR_ACCESS_TOKEN_HERE",
    ),
)

req = operations.GetWalletSummaryRequest(
    id='53f870b3-26b5-4a73-829c-db1a8422bb67',
)

res = s.wallets.get_wallet_summary(req)

if res.get_wallet_summary_response is not None:
    # handle response
```

## list_balances

List balances of a wallet

### Example Usage

```python
import sdk
from sdk.models import operations

s = sdk.SDK(
    security=shared.Security(
        authorization="Bearer YOUR_ACCESS_TOKEN_HERE",
    ),
)

req = operations.ListBalancesRequest(
    id='9d232271-5bf0-4cbb-9e31-b8b90f3443a1',
)

res = s.wallets.list_balances(req)

if res.list_balances_response is not None:
    # handle response
```

## list_wallets

List all wallets

### Example Usage

```python
import sdk
from sdk.models import operations

s = sdk.SDK(
    security=shared.Security(
        authorization="Bearer YOUR_ACCESS_TOKEN_HERE",
    ),
)

req = operations.ListWalletsRequest(
    cursor='aHR0cHM6Ly9nLnBhZ2UvTmVrby1SYW1lbj9zaGFyZQ==',
    metadata={
        "aut": 'quas',
    },
    name='William Ortiz',
    page_size=984330,
)

res = s.wallets.list_wallets(req)

if res.list_wallets_response is not None:
    # handle response
```

## update_wallet

Update a wallet

### Example Usage

```python
import sdk
from sdk.models import operations

s = sdk.SDK(
    security=shared.Security(
        authorization="Bearer YOUR_ACCESS_TOKEN_HERE",
    ),
)

req = operations.UpdateWalletRequest(
    request_body=operations.UpdateWalletRequestBody(
        metadata={
            "facilis": 'cupiditate',
            "qui": 'quae',
        },
    ),
    id='879fce95-3f73-4ef7-bbc7-abd74dd39c0f',
)

res = s.wallets.update_wallet(req)

if res.status_code == 200:
    # handle response
```

## void_hold

Cancel a hold

### Example Usage

```python
import sdk
from sdk.models import operations

s = sdk.SDK(
    security=shared.Security(
        authorization="Bearer YOUR_ACCESS_TOKEN_HERE",
    ),
)

req = operations.VoidHoldRequest(
    hold_id='exercitationem',
)

res = s.wallets.void_hold(req)

if res.status_code == 200:
    # handle response
```

## walletsget_server_info

Get server info

### Example Usage

```python
import sdk


s = sdk.SDK(
    security=shared.Security(
        authorization="Bearer YOUR_ACCESS_TOKEN_HERE",
    ),
)


res = s.wallets.walletsget_server_info()

if res.server_info is not None:
    # handle response
```
