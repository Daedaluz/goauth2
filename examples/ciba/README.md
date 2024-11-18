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

The [hard requirement of sending a hint](https://openid.net/specs/openid-client-initiated-backchannel-authentication-core-1_0.html#rfc.section.7) to the OP is ignored in this example
as it is used to allow any user to authenticate, not just the user that the hint is for.

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


eg:
```bash
./ciba register "Kalle Anka"
```

```
Challenge ID:  6ca14a4b8ca81332
Please register the key at https://idp.inits.se/authenticator?id=6ca14a4b8ca81332
```

After creating the key, you'll be redirected to a page that displays the created user id, that later can be used with
`./ciba <user_id>` to authenticate the specific user.

## Authenticating any user

To start an authentication request for an arbitrary user, run the following command:

```bash
./ciba
```

This method will skip sending any hints to the provider.

eg:
```bash
./ciba
```

```
<Terminal QR code>
https://idp.inits.se/authenticator?id=f637d537ffc7468d
Authorization Pending:  Waiting for user to view the challenge
```


## Authenticating a specific user

To start an authentication request for a specific user, run the following command:

```bash
./ciba <username>
```

This method will send a `login_hint` to the provider.

eg:
```bash
./ciba 5e09028099a65f488aed4d5bcf378a0a7404
```

```
<Terminal QR code>
https://idp.inits.se/authenticator?id=f637d537ffc7468d
Authorization Pending:  Waiting for user to view the challenge
```
