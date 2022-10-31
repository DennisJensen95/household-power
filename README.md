# household-power

![example workflow](https://github.com/DennisJensen95/household-power/actions/workflows/build.yml/badge.svg)
![Code coverage](https://img.shields.io/endpoint?url=https://gist.githubusercontent.com/DennisJensen95/2b7862c80c14d562c8659e1283543190/raw/household-power-coverage.json)

Power consumption prices are crazy. Therefore, I made this tool to identify if
you would be better off using a variable price arrangement paying hour by hour
prices for your energy. It is hard to figure out, when looking at Energy vendors
apps if it is financilly a benefit to have a variable paying arrangement. 

This tool calculates what you would pay if you had a variable arrangement vs
fixed price. You must provide a period of interest, a fixed price and transport
cost of the energy per kWh in DKK. 

## Usage 

The example below will show you your consumption and how much it would cost with
variable power cost vs fixed price.

First you need a token from https://eloverblik.dk/welcome. You can generate one
there by logging in with MitID and then replacing `TokenElOverblik` with the
generated token and also changing the `meteringPointId` to your metering point.

```bash
go build && ./household-power -token TokenForElOverblik -meteringPointId 132813343 -start-date 2022-10-01 -en
d-date 2022-10-31 -log-level i -fixed-price 6.0 -transport-cost 1.35
INFO[0000] Power used                                    PowerKWH=146                                                 
INFO[0000] Power used                                    Money spend with fixed agreement [DKK]=877.1 Money spend with
 variable agreement [DKK]=355.6
```

This can save you a lot of money if you are running with a high fixed price and
not using power in the busy hours. As you can see above I can save 500 DKK per
month with a variable price arrangement.