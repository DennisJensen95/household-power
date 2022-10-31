# household-power

Power consumption prices are crazy. Therefore, be more knowledgable for you
consumation Calculate the price of your power usage if you are paying by the
hour vs the fixed price you might give.

## Usage 

The example below will show you your consumption and how much it would cost with
variable power cost vs fixed price.

First you need a token from https://eloverblik.dk/welcome. You can generate one
there by logging in with MitID and then replacing `TokenElOverblik` with the
generated token.

```bash
go build && ./household-power -token TokenForElOverblik -meteringPointId 132813343 -start-date 2022-10-01 -en
d-date 2022-10-31 -log-level i -fixed-price 6.0 -transport-cost 1.35
INFO[0000] Power used                                    PowerKWH=146                                                 
INFO[0000] Power used                                    Money spend with fixed agreement [DKK]=877.1 Money spend with
 variable agreement [DKK]=355.6
```

This can save you a lot of money if you are running with a high fixed price and
not using power in the busy hours. 