# Extended CIBA POC

## About

This example is to demonstrate how the CIBA flow could look like if the OP
returns a `qr_data` and a `qr_type` in the authentication request response.

They are used deliver oob information to the authenticator to identify the authentication request.

Example usecase:

1. The RP initiates an authentication request with the OP.
2. The RP Presents the user with a QR code that can be used by the authenticator.
3. The user scans the QR code with the authenticator.
4. The authenticator presents the user with some kind of consent screen.
5. Once the user has responded to the request, the RP can collect the result as normal.

In this example the `qr_data` is an url and the `qr_type` contains the value `url`.

Due to the fact that there now is some relevant reference to the authenticator returned from the OP, the hard requirement
of some hint being sent to the OP can be removed.

`qr_type` could potentially be `opaque` and `qr_data` could be some other kind of relevance to the authenticator.

The flow is very similar to the one used for [BankID](https://www.bankid.com/), where the user is typically presented with
an animated QR code that can be scanned by the BankID app.

## Build

```bash
go build
```

## Register a user

The id-provider used in this example exclusively uses passkeys as a form of authentication.
To create a user, and associate a key; run the following command and follow the generated link:

```bash
./ciba register <name>
```

The name is not important and can be anything, for your own reference.

## Authenticating any user

To start an authentication request for an arbitrary user, run the following command:

```bash
./ciba
```

This method will skip sending any hints to the provider.

## Authenticating a specific user

To start an authentication request for a specific user, run the following command:

```bash
./ciba <username>
```

This method will send a `login_hint` to the provider.
