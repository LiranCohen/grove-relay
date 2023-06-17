# Bitcoin Grove NOSTR Relay

This is the software that runs the whitelisted Bitcoin Grove Relay.
It's using [https://github.com/fiatjaf/relayer](https://github.com/fiatjaf/relayer) as the base framework.

Use at your own discretion.


### Usage

#### Install
```
git clone https://github.com/Bitcoin-Grove/grove-relay.git
cd grove-relay

go install ./cmd/grove-relay
```
#### Run
```
grove-relay --name "My Awesome Relay" --pgconn postgres://name:pass@localhost:5432/dbname
```

#### Optional Parameters
```
   --host - service host (default: 0.0.0.0)
   --port - service port (default: 7447)
   --description - shows up on nip11
   --pubkey - shows up on nip11 
   --contact - shows up on nip11
   --maxeventsz - Max Event Size in bytes (default: 100000)
   --maxcachesz - Max Whitelist Cache Entry Size (default: 100)
```