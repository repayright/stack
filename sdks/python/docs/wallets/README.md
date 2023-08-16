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
    hold_id='aspernatur',
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
        expires_at=dateutil.parser.isoparse('2022-09-14T03:02:47.808Z'),
        name='Miriam Hermann',
        priority=574325,
    ),
    id='0afa563e-2516-4fe4-88b7-11e5b7fd2ed0',
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
        "praesentium": 'natus',
    },
    name='Joan Satterfield',
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
            amount=807319,
            asset='ea',
        ),
        balance='excepturi',
        metadata={
            "ea": 'accusantium',
        },
        reference='ab',
        sources=[
            shared.WalletSubject(
                balance='ipsam',
                identifier='voluptate',
                type='autem',
            ),
            shared.WalletSubject(
                balance='eaque',
                identifier='pariatur',
                type='nemo',
            ),
            shared.WalletSubject(
                balance='perferendis',
                identifier='fugiat',
                type='amet',
            ),
            shared.LedgerAccountSubject(
                identifier='cumque',
                type='corporis',
            ),
        ],
    ),
    id='fbb25870-5320-42c7-bd5f-e9b90c28909b',
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
            amount=237893,
            asset='asperiores',
        ),
        balances=[
            'modi',
            'iste',
            'dolorum',
            'deleniti',
        ],
        description='pariatur',
        destination=shared.WalletSubject(
            balance='nobis',
            identifier='libero',
            type='delectus',
        ),
        metadata={
            "quos": 'aliquid',
            "dolorem": 'dolorem',
        },
        pending=False,
    ),
    id='323f9b77-f3a4-4100-a74e-bf69280d1ba7',
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
    balance_name='voluptate',
    id='a89ebf73-7ae4-4203-8e5e-6a95d8a0d446',
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
    hold_id='quod',
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
        "qui": 'dolorum',
        "a": 'esse',
        "harum": 'iusto',
        "ipsum": 'quisquam',
    },
    page_size=947371,
    wallet_id='amet',
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
    page_size=730856,
    wallet_id='accusamus',
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
    id='453f870b-326b-45a7-b429-cdb1a8422bb6',
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
    id='79d23227-15bf-40cb-b1e3-1b8b90f3443a',
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
    id='1108e0ad-cf4b-4921-879f-ce953f73ef7f',
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
        "quod": 'odio',
        "similique": 'facilis',
        "vero": 'ducimus',
    },
    name='Rosemarie Stanton',
    page_size=773326,
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
            "voluptatibus": 'exercitationem',
        },
    ),
    id='d2cff7c7-0a45-4626-9436-813f16d9f5fc',
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
    hold_id='saepe',
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
